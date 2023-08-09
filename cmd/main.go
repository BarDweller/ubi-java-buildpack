package main

import (
	ubi8javabuildpack "github.com/paketo-community/ubi-java-buildpack"

	"github.com/paketo-buildpacks/libpak"
)

func main() {
	libpak.BuildpackMain(
		ubi8javabuildpack.Detect,
		ubi8javabuildpack.Build,
	)
}
