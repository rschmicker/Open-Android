package controller

import (
	"bytes"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func Run(ApkDir string, DecodedDir string, OutputDir string) {
	paths := getPaths(ApkDir, ".apk")
	if len(paths) == 0 {
		log.Fatal("No APKs found")
	}
	pathMap := decode(paths, DecodedDir)
	extract(pathMap)
}

func getPaths(ApkDir string, Containing string) []string {
	fileList := make([]string, 0)
	err := filepath.Walk(ApkDir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, Containing) {
			fileList = append(fileList, path)
		}
		return err
	})
	utils.Check(err)
	return fileList
}

func extract(pathMap map[string]string) {
	var wg sync.WaitGroup
	for apkPath, decodedPath := range pathMap {
		wg.Add(1)
		go func(apkPath string, decodedPath string) {
			defer wg.Done()
			ApkData := &utils.ApkData{}
			ApkData.Md5 = metadata.Md5File(apkPath)
		}(apkPath, decodedPath)
	}
	wg.Wait()
}

func decode(ApkPaths []string, DecodedDir string) (map[string]string) {
	var wg sync.WaitGroup
	pathMap := make(map[string]string)
	for _, ApkPath := range ApkPaths {
		wg.Add(1)
		go func(ApkPath string, DecodedDir string) {
			defer wg.Done()
			sha256Hash := metadata.Sha256File(ApkPath)
			apktoolPath, err := filepath.Abs("./apktool.sh")
			utils.Check(err)
			apkDecodedDir := DecodedDir + "/" + sha256Hash
			if _, err := os.Stat(apkDecodedDir); os.IsNotExist(err) {
				os.Mkdir(apkDecodedDir, os.FileMode(0700))
			}
			args := []string{"d", "--quiet", "-f", "-o", apkDecodedDir, ApkPath}
			cmd := exec.Command(apktoolPath, args...)
			var out bytes.Buffer
			var errout bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &errout
			err = cmd.Run()
			utils.Check(err)
			if errout.String() != "" {
				log.Printf(errout.String())
			}
			pathMap[ApkPath] = apkDecodedDir
		}(ApkPath, DecodedDir)
	}
	wg.Wait()
	return pathMap
}
