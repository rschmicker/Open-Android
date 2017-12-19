package cache

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type CacheTestSuite struct{}

var _ = Suite(&CacheTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *CacheTestSuite) TestgetPaths(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir
	paths := getPaths(testLoc, ".apk")
	c.Check(len(paths), Equals, 1)
}

func (s *CacheTestSuite) TestInitialize(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	ct := &CacheTable{}
	length := ct.Initialize(config)
	c.Check(length, Equals, 1)
	c.Check(ct.RamDiskPath, Equals, "/dev/shm/cache/")
	c.Check(ct.Size, Equals, 1)
	c.Check(ct.DirectoryToCache, Equals, config.ApkDir)
	c.Check(ct.Location, Equals, 1)
}
