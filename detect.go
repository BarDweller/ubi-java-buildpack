package ubi8javabuildpack

import (
	libcnb "github.com/buildpacks/libcnb/v2"
)

func Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "ubi-java-helper"},
				},
			},
			{
				Requires: []libcnb.BuildPlanRequire{
					{Name: "ubi-java-helper"},
				},
			},
		},
	}, nil
}
