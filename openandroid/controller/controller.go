package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Open-Android/openandroid/apis"
	"github.com/Open-Android/openandroid/intent"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/stringApk"
	"github.com/Open-Android/openandroid/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type ApkData struct {
	Strings     []string
	Apis        []string
	Permissions []string
	Md5         string
	Sha256      string
	Sha1        string
	PackageName string
	Version     string
	Intents     []string
	Malicious   bool
}

func Runner(ApkDir string, DecodedDir string, OutputDir string, CodeDir string) {
	paths := getPaths(ApkDir, ".apk")
	if len(paths) == 0 {
		log.Fatal("No APKs found")
	}
	pathMap := decode(paths, DecodedDir)
	extract(pathMap, OutputDir, CodeDir)
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

func extract(pathMap map[string]string, OutputDir string, CodeDir string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 12)
	var mutex = &sync.Mutex{}
	count := 0
	for apkPath, decodedPath := range pathMap {
		wg.Add(1)
		go func(apkPath string, decodedPath string, outputDir string, codePath string) {
			var err error
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			ApkData := &ApkData{}
			ApkData.GetMetaData(apkPath, decodedPath)
			ApkData.IsMalicious(apkPath)
			ApkData.Intents = intent.GetIntents(decodedPath)
			mutex.Lock()
			ApkData.Apis, err = apis.GetApis(apkPath, codePath)
			if err != nil {
				log.Printf("Error extracting apis: " + apkPath)
				return
			}
			ApkData.Strings, err = stringApk.GetStrings(apkPath, codePath)
			if err != nil {
				log.Printf("Error extracting strings: " + apkPath)
				return
			}
			mutex.Unlock()
			ApkData.WriteJSON(outputDir)
			mutex.Lock()
			count++
			mutex.Unlock()
			percent := (float64(count) / float64(len(pathMap))) * float64(100)
			log.Printf("(%.2f%%) Extracted: "+metadata.GetApkName(decodedPath), percent)
		}(apkPath, decodedPath, OutputDir, CodeDir)
	}
	wg.Wait()
	close(sem)
}

func (apk *ApkData) IsMalicious(apkPath string) {
	if strings.Contains(apkPath, "malware") {
		apk.Malicious = true
	} else {
		apk.Malicious = false
	}
}

func (apk *ApkData) WriteJSON(OutputDir string) {
	data, err := json.Marshal(apk)
	utils.Check(err)
	outputFile := OutputDir + "/" + apk.Sha256 + ".json"
	err = ioutil.WriteFile(outputFile, data, 0644)
	utils.Check(err)
}

func (apk *ApkData) GetMetaData(apkPath string, decodedPath string) {
	apk.Md5 = metadata.Md5File(apkPath)
	apk.Sha1 = metadata.Sha1File(apkPath)
	apk.Sha256 = metadata.Sha256File(apkPath)
	apk.Version = metadata.GetVersion(decodedPath)
	apk.PackageName = metadata.GetPackageName(decodedPath)
	apk.Permissions = metadata.GetPermissions(decodedPath)
}

func decode(ApkPaths []string, DecodedDir string) map[string]string {
	var wg sync.WaitGroup
	pathMap := make(map[string]string)
	sem := make(chan struct{}, 12)
	count := 0
	var mutex = &sync.Mutex{}
	for _, ApkPath := range ApkPaths {
		wg.Add(1)
		go func(ApkPath string, DecodedDir string) {
			sem <- struct{}{}
			defer func() { <-sem }()
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
			if err != nil {
				log.Printf("Error decoding: " + ApkPath)
				return
			}
			if errout.String() != "" {
				log.Printf(errout.String())
			}
			pathMap[ApkPath] = apkDecodedDir
			mutex.Lock()
			count++
			mutex.Unlock()
			percent := (float64(count) / float64(len(ApkPaths))) * float64(100)
			log.Printf("(%.2f%%) Decoded: "+metadata.GetApkName(apkDecodedDir), percent)
		}(ApkPath, DecodedDir)
	}
	wg.Wait()
	close(sem)
	return pathMap
}
