package controller

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type ControllerTestSuite struct{}

var _ = Suite(&ControllerTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *ControllerTestSuite) TestgetPaths(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir
	paths := getPaths(testLoc, ".apk")
	c.Check(len(paths), Equals, 1)
}
