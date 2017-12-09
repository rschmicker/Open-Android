package metadata

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type MetaDataTestSuite struct{}

var _ = Suite(&MetaDataTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *MetaDataTestSuite) TestSha256File(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	hash := Sha256File(testLoc)
	testHash := "8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	c.Check(hash, Equals, testHash)
}

func (s *MetaDataTestSuite) TestSha1File(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	hash := Sha1File(testLoc)
	testHash := "3ce20472e647d0194fd11518c302236934d5f605"
	c.Check(hash, Equals, testHash)
}

func (s *MetaDataTestSuite) TestMd5File(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	hash := Md5File(testLoc)
	testHash := "a1c88d70e6ffe6ed5167f75c8399af4e"
	c.Check(hash, Equals, testHash)
}

func (s *MetaDataTestSuite) TestGetPackageName(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.DecodedDir + "/8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	name := GetPackageName(testLoc)
	testName := "com.facebook.lite"
	c.Check(name, Equals, testName)
}

func (s *MetaDataTestSuite) TestGetVersion(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.DecodedDir + "/8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	version := GetVersion(testLoc)
	testVersion := "70.0.0.9.116"
	c.Check(version, Equals, testVersion)
}

func (s *MetaDataTestSuite) TestGetApkName(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.DecodedDir + "/8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	name := GetApkName(testLoc)
	testName := "Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	c.Check(name, Equals, testName)
}

func (s *MetaDataTestSuite) TestGetPermissions(c *C) {
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.DecodedDir + "/8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	permissions := GetPermissions(testLoc)
	c.Check(len(permissions), Equals, 45)
}
