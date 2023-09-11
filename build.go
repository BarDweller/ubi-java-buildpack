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

	//read the env vars set via the extension.
	version := os.Getenv("BPI_UBI_JAVA_EXTENSION_VERSION")
	helperstr := os.Getenv("BPI_UBI_JAVA_EXTENSION_HELPERS")

	//only act if the version is set, otherwise we are a no-op.
	if version != "" {
		logger := log.NewPaketoLogger(context.Logger.DebugWriter())
		logger.Title(context.Buildpack.Info.Name, context.Buildpack.Info.Version, context.Buildpack.Info.Homepage)

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
	}

	return result, nil
}
