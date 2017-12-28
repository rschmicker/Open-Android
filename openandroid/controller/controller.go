package controller

import (
	"github.com/Open-Android/openandroid/apis"
	. "github.com/Open-Android/openandroid/apkdata"
	"github.com/Open-Android/openandroid/cache"
	"github.com/Open-Android/openandroid/cleaner"
	"github.com/Open-Android/openandroid/intent"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/stringApk"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"runtime"
	"sync"
)

var javaMutex = &sync.Mutex{}
var countMutex = &sync.Mutex{}
var wg sync.WaitGroup

func Runner(config utils.ConfigData) {
	if config.Clean == true {
		cleaner.CleanDirectory(config)
	}
	cacheTable := &cache.CacheTable{}
	length := cacheTable.Initialize(config)
	go cacheTable.Runner()
	sem := make(chan struct{}, runtime.NumCPU())
	count := 0
	for i := 0; i < length; i++ {
		if cacheTable.IsEmpty() {
			break
		}
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			apk := cacheTable.GetFilePath()
			defer cacheTable.Completed(apk)
			defer func() { <-sem }()
			defer wg.Done()
			extract(apk, config)
			countMutex.Lock()
			count++
			countMutex.Unlock()
			percent := (float64(count) / float64(length) * float64(100))
			name := metadata.GetApkName(apk)
			log.Printf("(%.2f%%) Completed: "+name, percent)
		}()
	}
	wg.Wait()
	close(sem)
	cacheTable.Close()
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
