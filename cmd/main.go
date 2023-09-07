package main

import (
	ubi8javabuildpack "github.com/paketo-community/ubi-java-buildpack"

	libpak "github.com/paketo-buildpacks/libpak/v2"
)

func main() {
	libpak.BuildpackMain(
		ubi8javabuildpack.Detect,
		ubi8javabuildpack.Build,
	)
}
