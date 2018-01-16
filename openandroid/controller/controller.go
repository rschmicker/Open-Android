package controller

import (
	"github.com/Open-Android/openandroid/apis"
	. "github.com/Open-Android/openandroid/apkdata"
	"github.com/Open-Android/openandroid/cleaner"
	"github.com/Open-Android/openandroid/intent"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/stringApk"
	"github.com/Open-Android/openandroid/utils"
	"github.com/rschmicker/FileCache/cache"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var javaMutex = &sync.Mutex{}
var countMutex = &sync.Mutex{}
var wg sync.WaitGroup

func Runner(config utils.ConfigData) {
	if config.Clean {
		cleaner.CleanDirectory(config)
	}
	cacheTable := &cache.CacheTable{}
	cacheTable.RamDiskPath = config.CacheDir + "/cache/"
	toDoFiles := utils.GetPaths(config.ApkDir, ".apk")
	if config.Force {
		cacheTable.Files = toDoFiles
	} else {
		doneFiles := utils.GetPaths(config.OutputDir, ".json")
		cacheTable.Files = utils.CrossCompare(toDoFiles, doneFiles)
	}
	length := len(cacheTable.Files)
	cacheTable.Initialize()

	sigChannel := make(chan os.Signal)
	go func() {
		for sig := range sigChannel {
			switch sig {
			case syscall.SIGINT:
				log.Printf("Clearing cache...")
				cacheTable.Close()
				os.Exit(1)
			}
		}
	}()
	signal.Notify(sigChannel, syscall.SIGINT)

	go cacheTable.Runner()
	sem := make(chan struct{}, runtime.NumCPU())
	count := 0
	for !cacheTable.IsEmpty() {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			apk := cacheTable.GetFilePath()
			defer cacheTable.Completed(apk)
			defer func() { <-sem }()
			defer wg.Done()
			if apk == "" {
				return
			}
			err := extract(apk, config)
			if err != nil {
				log.Printf("Warning: " + apk + " is not a valid APK file")
				return
			}
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
	log.Printf("All files parsed... clearing cache...")
	cacheTable.Close()
}

func extract(path string, config utils.ConfigData) error {
	apkd := &ApkData{}
	err := apkd.GetMetaData(path)
	if err != nil {
		return err
	}
	err = apkd.IsMalicious(path, config.VtApiKey, config.VtApiCheck)
	if err != nil {
		log.Printf(err.Error())
	}
	apkd.Intents = intent.GetIntents(path)
	javaMutex.Lock()
	apkd.Apis = apis.GetApis(path, config.CodeDir)
	apkd.Strings = stringApk.GetStrings(path, config.CodeDir)
	javaMutex.Unlock()
	apkd.WriteJSON(config.OutputDir)
	return nil
}
