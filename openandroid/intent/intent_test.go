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
	configPath, err := filepath.Abs("../../openandroid.yaml")
	c.Assert(err, IsNil)
	config := utils.ReadConfig(configPath)
	testLoc := config.DecodedDir + "/8fc218d35790b7c363b7423f9bd6faa71b2adcc59e55444431eced0cf0e60a4d"
	intents := GetIntents(testLoc)
	c.Check(67, Equals, len(intents))
}
