package intent

import (
	"github.com/Open-Android/openandroid/utils"
	. "gopkg.in/check.v1"
	"path/filepath"
	"testing"
)

type IntentTestSuite struct{}

var _ = Suite(&IntentTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *IntentTestSuite) TestGetIntents(c *C) {
	configPath, err := filepath.Abs("../../test.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.ApkDir + "/Facebook Lite_v70.0.0.9.116_apkpure.com.apk"
	intents := GetIntents(testLoc)
	c.Check(60, Equals, len(intents))
}
