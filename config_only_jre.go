/*
 * Copyright 2023s the original author or authors.
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
	"github.com/buildpacks/libcnb/v2"
	"github.com/paketo-buildpacks/libjvm/v2"
	"github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/log"
)

// consts for the config only JRE
const isBuild = false
const isLaunch = true
const isCache = false
const name = "ajre" //alphabetically must come before 'helper' and 'java-security-properties'.
const ubi8JrePath = "/usr/lib/jvm/jre"

type ConfigOnlyJRE struct {
	ApplicationPath   string
	JavaVersion       string
	CertificateLoader libjvm.CertificateLoader
	DistributionType  libjvm.DistributionType
	LayerContributor  libpak.LayerContributor
	Logger            log.Logger
}

func NewConfigOnlyJRE(logger log.Logger, info libcnb.BuildpackInfo, applicationPath string, javaVersion string, certificateLoader libjvm.CertificateLoader) (ConfigOnlyJRE, error) {
	certLoader := libjvm.NewCertificateLoader()
	contributor := libpak.NewLayerContributor("Configuration for JRE", info, libcnb.LayerTypes{Launch: isLaunch, Build: isBuild, Cache: isCache}, logger)
	return ConfigOnlyJRE{
		ApplicationPath:   applicationPath,
		JavaVersion:       javaVersion,
		CertificateLoader: certLoader,
		DistributionType:  libjvm.JREType,
		LayerContributor:  contributor,
		Logger:            logger,
	}, nil
}

func (j ConfigOnlyJRE) Contribute(layer *libcnb.Layer) error {

	return j.LayerContributor.Contribute(layer, func(layer *libcnb.Layer) error {
		j.Logger.Body("Configuring installed JRE")
		return ConfigureJRE(layer, j.Logger,
			ubi8JrePath,          //java home
			j.JavaVersion,       //java version
			j.ApplicationPath,   //app path
			isBuild,             //isBuild
			isLaunch,            //isLaunch
			j.CertificateLoader, //certLoader
			j.DistributionType)  //jdk/jre
	})
}

func (j ConfigOnlyJRE) Name() string {
	return name
}
