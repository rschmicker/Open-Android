package apis

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type ApiTestSuite struct{}

var _ = Suite(&ApiTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *ApiTestSuite) TestGetStrings(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	apisInApk := GetApis(testLoc, config.CodeDir)
	c.Check(len(apisInApk), Equals, 5007)
}
