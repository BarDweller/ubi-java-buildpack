/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ubi8javabuildpack

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/heroku/color"

	"github.com/buildpacks/libcnb/v2"
	"github.com/magiconair/properties"
	"github.com/paketo-buildpacks/libjvm/v2"
	"github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/crush"
	"github.com/paketo-buildpacks/libpak/v2/log"

	"github.com/paketo-buildpacks/libjvm/v2/count"
)

type JRE struct {
	ApplicationPath   string
	CertificateLoader libjvm.CertificateLoader
	DistributionType  libjvm.DistributionType
	LayerContributor  libpak.DependencyLayerContributor
	Logger            log.Logger
	Metadata          map[string]interface{}
}

func NewJRE(applicationPath string, dependency libpak.BuildModuleDependency, cache libpak.DependencyCache, distributionType libjvm.DistributionType, certificateLoader libjvm.CertificateLoader, metadata map[string]interface{}) (JRE, error) {
	expected := map[string]interface{}{"dependency": dependency}

	if md, err := certificateLoader.Metadata(); err != nil {
		return JRE{}, fmt.Errorf("unable to generate certificate loader metadata\n%w", err)
	} else {
		for k, v := range md {
			expected[k] = v
		}
	}

	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Build:  libjvm.IsBuildContribution(metadata),
		Cache:  libjvm.IsBuildContribution(metadata),
		Launch: libjvm.IsLaunchContribution(metadata),
	}, cache.Logger)
	contributor.ExpectedMetadata = expected

	return JRE{
		ApplicationPath:   applicationPath,
		CertificateLoader: certificateLoader,
		DistributionType:  distributionType,
		LayerContributor:  contributor,
		Metadata:          metadata,
		Logger:            cache.Logger,
	}, nil
}

func ConfigureJRE(layer *libcnb.Layer, logger log.Logger,
	javaHome string,
	javaVersion string,
	appPath string,
	isBuild bool,
	isLaunch bool,
	certLoader libjvm.CertificateLoader,
	distType libjvm.DistributionType) error {
	var cacertsPath string
	if libjvm.IsBeforeJava9(javaVersion) && distType == libjvm.JDKType {
		cacertsPath = filepath.Join(javaHome, "jre", "lib", "security", "cacerts")
	} else {
		cacertsPath = filepath.Join(javaHome, "lib", "security", "cacerts")
	}
	if err := os.Chmod(cacertsPath, 0664); err != nil {
		return fmt.Errorf("unable to set keystore file permissions\n%w", err)
	}

	if libjvm.IsBeforeJava18(javaVersion) {
		if err := certLoader.Load(cacertsPath, "changeit"); err != nil {
			return fmt.Errorf("unable to load certificates\n%w", err)
		}
	} else {
		logger.Bodyf("%s: The JVM cacerts entries cannot be loaded with Java 18+, for more information see: https://github.com/paketo-buildpacks/libjvm/issues/158", color.YellowString("Warning"))
	}

	if isBuild {
		layer.BuildEnvironment.Default("JAVA_HOME", javaHome)
	}

	if isLaunch {
		layer.LaunchEnvironment.Default("BPI_APPLICATION_PATH", appPath)
		layer.LaunchEnvironment.Default("BPI_JVM_CACERTS", cacertsPath)

		if c, err := count.Classes(javaHome); err != nil {
			return fmt.Errorf("unable to count JVM classes\n%w", err)
		} else {
			layer.LaunchEnvironment.Default("BPI_JVM_CLASS_COUNT", c)
		}

		if libjvm.IsBeforeJava9(javaVersion) && distType == libjvm.JDKType {
			layer.LaunchEnvironment.Default("BPI_JVM_EXT_DIR", filepath.Join(javaHome, "jre", "lib", "ext"))
		} else if libjvm.IsBeforeJava9(javaVersion) && distType == libjvm.JREType {
			layer.LaunchEnvironment.Default("BPI_JVM_EXT_DIR", filepath.Join(javaHome, "lib", "ext"))
		}

		var securityFile string
		if libjvm.IsBeforeJava9(javaVersion) && distType == libjvm.JDKType {
			securityFile = filepath.Join(javaHome, "jre", "lib", "security", "java.security")
		} else if libjvm.IsBeforeJava9(javaVersion) && distType == libjvm.JREType {
			securityFile = filepath.Join(javaHome, "lib", "security", "java.security")
		} else {
			securityFile = filepath.Join(javaHome, "conf", "security", "java.security")
		}

		p, err := properties.LoadFile(securityFile, properties.UTF8)
		if err != nil {
			return fmt.Errorf("unable to read properties file %s\n%w", securityFile, err)
		}
		p = p.FilterStripPrefix("security.provider.")

		var providers []string
		for k, v := range p.Map() {
			providers = append(providers, fmt.Sprintf("%s|%s", k, v))
		}
		sort.Strings(providers)
		layer.LaunchEnvironment.Default("BPI_JVM_SECURITY_PROVIDERS", strings.Join(providers, " "))

		layer.LaunchEnvironment.Default("JAVA_HOME", javaHome)
		layer.LaunchEnvironment.Default("MALLOC_ARENA_MAX", "2")

		layer.LaunchEnvironment.Append("JAVA_TOOL_OPTIONS", " ", "-XX:+ExitOnOutOfMemoryError")
	}
	return nil
}

func (j JRE) Contribute(layer *libcnb.Layer) error {

	return j.LayerContributor.Contribute(layer, func(layer *libcnb.Layer, artifact *os.File) error {
		j.Logger.Bodyf("Expanding to %s", layer.Path)
		if err := crush.Extract(artifact, layer.Path, 1); err != nil {
			return fmt.Errorf("unable to expand JRE\n%w", err)
		}

		ConfigureJRE(layer, j.Logger,
			layer.Path,                              //java home
			j.LayerContributor.Dependency.Version,   //java version
			j.ApplicationPath,                       //app path
			libjvm.IsBuildContribution(j.Metadata),  //isBuild
			libjvm.IsLaunchContribution(j.Metadata), //isLaunch
			j.CertificateLoader,                     //certLoader
			j.DistributionType)                      //jdk/jre

		return nil
	})
}

func (j JRE) Name() string {
	return j.LayerContributor.LayerName()
}
