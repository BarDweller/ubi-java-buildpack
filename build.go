package ubi8javabuildpack

import (
	"os"
	"strings"

	libcnb "github.com/buildpacks/libcnb/v2"
	libjvm "github.com/paketo-buildpacks/libjvm/v2"
	libpak "github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/log"
)

const javaVersionBuilderFile = "/bpi.paketo.ubi.java.version"
const javaHelpersBuilderFile = "/bpi.paketo.ubi.java.helpers"
const openSslLoaderName = "openssl-certificate-loader"

func Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	result := libcnb.BuildResult{}

	logger := log.NewPaketoLogger(os.Stdout)
	logger.Title(context.Buildpack.Info.Name, context.Buildpack.Info.Version, context.Buildpack.Info.Homepage)

	//read the env vars set via the extension.
	versionb, err := os.ReadFile(javaVersionBuilderFile)
	if err != nil {
		return result, err
	}
	helperb, err := os.ReadFile(javaHelpersBuilderFile)
	if err != nil {
		return result, err
	}

	version := strings.TrimSuffix(string(versionb), "\n")
	helperstr := strings.TrimSuffix(string(helperb), "\n")

	//only act if the version is set, otherwise we are a no-op.
	if version != "" {
		//recreate the various Contributable's that the extension could not use to create layers.
		logger.Body("Helper buildpack configuring using '" + helperstr + "' and java security properties for version " + version)

		jre, err := NewConfigOnlyJRE(logger, context.Buildpack.Info, context.ApplicationPath, version, libjvm.NewCertificateLoader(logger))
		if err != nil {
			return result, err
		}

		helpers := strings.Split(helperstr, ",")

		//remove certloader.. we don't do this.
		var temp []string
		for _, h := range helpers {
			if h != openSslLoaderName {
				temp = append(temp, h)
			}
		}
		helpers = temp

		h := libpak.NewHelperLayerContributor(context.Buildpack, logger, helpers...)
		jsp := libjvm.NewJavaSecurityProperties(context.Buildpack.Info, logger)

		//use libpak to process the contributable's into layers, by invoking the buildfunc.
		logger.Body("Initiating layer creation")
		return libpak.ContributableBuildFunc(func(context libcnb.BuildContext, result *libcnb.BuildResult) ([]libpak.Contributable, error) {
			return []libpak.Contributable{jre, h, jsp}, nil
		})(context)

	} else {
		logger.Body("Helper buildpack did not detect config from extension. Disabling.")
	}

	return result, nil
}
