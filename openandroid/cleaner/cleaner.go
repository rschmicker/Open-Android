package cleaner

import (
	"github.com/Open-Android/openandroid/cache"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"os"
	"strings"
	"sync"
)

func CleanDirectory(config utils.ConfigData) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 1000)
	fileList := cache.GetPaths(config.ApkDir, ".apk")
	for _, file := range fileList {
		wg.Add(1)
		go func(file string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			hash := metadata.Sha256File(file)
			newPath := ""
			if strings.Contains(file, "benign") {
				newPath = config.ApkDir + "/benign/" + hash + ".apk"
			} else {
				newPath = config.ApkDir + "/malware/" + hash + ".apk"
			}
			if file == newPath {
				log.Printf("Skipping: %v already cleaned", file)
				return
			}
			err := os.Rename(file, newPath)
			utils.Check(err)
			log.Printf("Cleaned: %v -> %v", file, newPath)
		}(file)
	}
}
