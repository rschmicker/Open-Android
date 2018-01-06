package cleaner

import (
	"bufio"
	"fmt"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func CleanDirectory(config utils.ConfigData) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())
	fileList := utils.GetPaths(config.ApkDir, ".apk")
	for _, file := range fileList {
		r, err := os.Open(file)
		utils.Check(err)
		var header [2]byte
		_, err = io.ReadFull(r, header[:])
		utils.Check(err)
		r.Close()
		magic := fmt.Sprintf("%s", header)
		if magic != "PK" {
			log.Printf("Unknown file: %v", file)
			out, err := exec.Command("xxd " + file + " | head -n10").Output()
			utils.Check(err)
			log.Printf(string(out))
			reader := bufio.NewReader(os.Stdin)
			log.Printf("Delete file? (y/N): ")
			text, err := reader.ReadString('\n')
			utils.Check(err)
			if text == "y" {
				err = os.Remove(file)
				utils.Check(err)
			}
		}
	}
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
	wg.Wait()
	close(sem)
}
