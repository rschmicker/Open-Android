package apkdata

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type ApkDataTestSuite struct{}

var _ = Suite(&ApkDataTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *ApkDataTestSuite) TestIsMalicious(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	apkd := &ApkData{}
	apkd.IsMalicious(testLoc, config.VtApiKey, true)
	c.Check(apkd.Malicious, Equals, false)
}

func (s *ApkDataTestSuite) TestcheckReport(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	vt, err := NewVirusTotal(config.VtApiKey)
	c.Assert(err, IsNil)
	testHash := "8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	rr, err := vt.checkReport(testHash)
	c.Assert(err, IsNil)
	c.Check(rr.Sha256, Equals, testHash)
}

func (s *ApkDataTestSuite) TestscanApk(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	vt, err := NewVirusTotal(config.VtApiKey)
	c.Assert(err, IsNil)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	err = vt.scanApk(testLoc)
	c.Assert(err, IsNil)
}

func (s *ApkDataTestSuite) TestWriteJSON(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	apkd := &ApkData{}
	err = apkd.GetMetaData(testLoc)
	c.Assert(err, IsNil)
	apkd.WriteJSON(config.OutputDir)
}

func (s *ApkDataTestSuite) TestNewVirusTotal(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	_, err = NewVirusTotal(config.VtApiKey)
	c.Assert(err, IsNil)
}
