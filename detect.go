package ubi8javabuildpack

import (
	libcnb "github.com/buildpacks/libcnb/v2"
)

func Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Pass:  true,
		Plans: []libcnb.BuildPlan{},
	}, nil
}
