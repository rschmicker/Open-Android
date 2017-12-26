package cleaner

import (
	"github.com/Open-Android/openandroid/cache"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"os"
	"strings"
)

func CleanDirectory(config utils.ConfigData) {
	fileList := cache.GetPaths(config.ApkDir, ".apk")
	for _, file := range fileList {
		go func(file string) {
			hash := metadata.Sha256File(file)
			newPath := ""
			if strings.Contains(file, "benign") {
				newPath = config.ApkDir + "/benign/" + hash + ".apk"
			} else {
				newPath = config.ApkDir + "/malware/" + hash + ".apk"
			}
			if file == newPath {
				log.Printf("Skipping: %v already cleaned", file)
				continue
			}
			err := os.Rename(file, newPath)
			utils.Check(err)
			log.Printf("Cleaned: %v -> %v", file, newPath)
		}(file)
	}
}
