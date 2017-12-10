package controller

import (
	"github.com/Open-Android/openandroid/apis"
	. "github.com/Open-Android/openandroid/apkdata"
	"github.com/Open-Android/openandroid/intent"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/stringApk"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var javaMutex = &sync.Mutex{}
var countMutex = &sync.Mutex{}

func Runner(config utils.ConfigData) {
	paths := getPaths(config.ApkDir, ".apk")
	if len(paths) == 0 {
		log.Fatal("No APKs found")
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())
	count := 0
	for _, apk := range paths {
		wg.Add(1)
		go func(apk string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			extract(apk, config)
			countMutex.Lock()
			count++
			countMutex.Unlock()
			percent := (float64(count) / float64(len(paths))) * float64(100)
			log.Printf("(%.2f%%) Completed: "+metadata.GetApkName(apk), percent)
		}(apk)
	}
	wg.Wait()
	close(sem)
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

func extract(path string, config utils.ConfigData) {
	var err error
	apkd := &ApkData{}
	apkd.GetMetaData(path)
	apkd.IsMalicious(path)
	apkd.Intents = intent.GetIntents(path)
	javaMutex.Lock()
	apkd.Apis, err = apis.GetApis(path, config.CodeDir)
	if err != nil {
		log.Printf("Error extracting apis: " + path)
		return
	}
	apkd.Strings, err = stringApk.GetStrings(path, config.CodeDir)
	if err != nil {
		log.Printf("Error extracting strings: " + path)
		return
	}
	javaMutex.Unlock()
	apkd.WriteJSON(config.OutputDir)
}
