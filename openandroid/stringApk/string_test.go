package stringApk

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type StringTestSuite struct{}

var _ = Suite(&StringTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *StringTestSuite) TestGetStrings(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	stringsInApk := GetStrings(testLoc, config.CodeDir)
	c.Check(len(stringsInApk), Equals, 14967)
}
