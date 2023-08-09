package ubi8javabuildpack

import (
	"fmt"
	"os"
	"strings"

	libcnb "github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libjvm"
	libpak "github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {

	result := libcnb.BuildResult{}

	//read the env vars set via the extension.
	version := os.Getenv("UBI_JAVA_EXTENSION_VERSION")
	helperstr := os.Getenv("UBI_JAVA_EXTENSION_HELPERS")

	if version != "" {
		helpers := strings.Split(helperstr, ",")
		h := libpak.NewHelperLayerContributor(context.Buildpack, helpers...)
		h.Logger = bard.NewLogger(context.Logger.DebugWriter())

		h.Logger.Title(context.Buildpack.Info.Name, context.Buildpack.Info.Version, context.Buildpack.Info.Homepage)
		h.Logger.Body(" - Helper buildpack processing helpers '" + helperstr + "' for version " + version)

		l, err := libjvm.DefaultFlattenContributorFn(h, context)
		if err != nil {
			return result, fmt.Errorf("unable to contribute helper layer\n%w", err)
		}
		result.Layers = append(result.Layers, l)

		h.Logger.Body(" - Helper buildpack adding java security properties")
		jsp := libjvm.NewJavaSecurityProperties(context.Buildpack.Info)
		jsp.Logger = h.Logger
		l, err = libjvm.DefaultFlattenContributorFn(jsp, context)
		if err != nil {
			return result, fmt.Errorf("unable to contribute jsp layer\n%w", err)
		}
		result.Layers = append(result.Layers, l)
	}

	return result, nil
}
