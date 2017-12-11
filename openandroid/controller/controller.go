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
var wg sync.WaitGroup

func Runner(config utils.ConfigData) {
	paths := getPaths(config.ApkDir, ".apk")
	if len(paths) == 0 {
		log.Fatal("No APKs found")
	}
	sem := make(chan struct{}, runtime.NumCPU()/2*100)
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
			name := metadata.GetApkName(apk)
			log.Printf("(%.2f%%) Completed: "+name, percent)
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
	apkd := &ApkData{}
	apkd.GetMetaData(path)
	apkd.IsMalicious(path)
	apkd.Intents = intent.GetIntents(path)
	javaMutex.Lock()
	apkd.Apis = apis.GetApis(path, config.CodeDir)
	apkd.Strings = stringApk.GetStrings(path, config.CodeDir)
	javaMutex.Unlock()
	apkd.WriteJSON(config.OutputDir)
}
