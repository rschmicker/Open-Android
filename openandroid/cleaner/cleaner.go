package cleaner

import (
	"github.com/Open-Android/openandroid/cache"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"os"
	"strings"
)

func CleanDirectory(config utils.ConfigData) {
	fileList := cache.GetPaths(config.ApkDir, ".apk")
	for _, file := range fileList {
		hash := metadata.Sha256File(file)
		newPath := ""
		if strings.Contains(file, "benign") {
			newPath = config.ApkDir + "/benign/" + hash + ".apk"
		} else {
			newPath = config.ApkDir + "/malicious/" + hash + ".apk"
		}
		err := os.Rename(file, newPath)
		utils.Check(err)
	}
}
