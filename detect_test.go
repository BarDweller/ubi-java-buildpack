package ubi8javabuildpack_test

import (
	"testing"

	libcnb "github.com/buildpacks/libcnb/v2"
	ubi8javabuildpack "github.com/paketo-community/ubi-java-buildpack"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = NewWithT(t).Expect
		dc     libcnb.DetectContext
		result libcnb.DetectResult
		err    error
	)

	context("Detect Result Check", func() {
		it.Before(func() {
			dc = libcnb.DetectContext{}
		})
		it("includes build plan options", func() {
			result, err = ubi8javabuildpack.Detect(dc)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(
				libcnb.DetectResult{
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
				}))
		})
	})
}
