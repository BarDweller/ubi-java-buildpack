package ubi8javabuildpack_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitJava(t *testing.T) {
	suite := spec.New("java", spec.Report(report.Terminal{}))
	suite("Detect", testDetect)
	suite.Run(t)
}
