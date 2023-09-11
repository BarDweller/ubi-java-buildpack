package ubi8javabuildpack

import (
	"os"
	"strings"

	libcnb "github.com/buildpacks/libcnb/v2"
	libjvm "github.com/paketo-buildpacks/libjvm/v2"
	libpak "github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/log"
)

func Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	result := libcnb.BuildResult{}

	logger := log.NewPaketoLogger(os.Stdout)
	logger.Title(context.Buildpack.Info.Name, context.Buildpack.Info.Version, context.Buildpack.Info.Homepage)

	//read the env vars set via the extension.
	versionb, err := os.ReadFile("/bpi.paketo.ubi.java.version")
	if err != nil {
		return result, err
	}
	helperb, err := os.ReadFile("/bpi.paketo.ubi.java.helpers")
	if err != nil {
		return result, err
	}

	version := string(versionb)
	helperstr := string(helperb)

	//only act if the version is set, otherwise we are a no-op.
	if version != "" {
		//recreate the various Contributable's that the extension could not use to create layers.
		logger.Body(" - Helper buildpack contributing helpers '" + helperstr + "' for version " + version)
		helpers := strings.Split(helperstr, ",")
		h := libpak.NewHelperLayerContributor(context.Buildpack, logger, helpers...)
		logger.Body(" - Helper buildpack adding java security properties")
		jsp := libjvm.NewJavaSecurityProperties(context.Buildpack.Info, logger)

		//use libpak to process the contributable's into layers, by invoking the buildfunc.
		logger.Body(" - Helper buildpack creating layers")
		return libpak.ContributableBuildFunc(func(context libcnb.BuildContext, result *libcnb.BuildResult) ([]libpak.Contributable, error) {
			return []libpak.Contributable{h, jsp}, nil
		})(context)
	} else {
		logger.Body(" - Helper buildpack did not detect config from extension. Disabling.")
	}

	return result, nil
}
