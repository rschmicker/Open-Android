package controller

import (
	"encoding/json"
	"github.com/Open-Android/openandroid/cleaner"
	"github.com/Open-Android/openandroid/utils"
	"github.com/rschmicker/FileCache/cache"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"runtime"
	"sync"
	"syscall"
)

type WorkerData struct {
	CacheTable *cache.CacheTable
	Sem        chan struct{}
	Config     utils.ConfigData
	Length     int
	Count      *int
}

var mutex = &sync.Mutex{}
var countMutex = &sync.Mutex{}
var wg sync.WaitGroup

func Runner(config utils.ConfigData) {
	count := 0
	ct, length := initCache(config)
	wd := &WorkerData{
		CacheTable: ct,
		Sem:        make(chan struct{}, runtime.NumCPU()),
		Config:     config,
		Length:     length,
		Count:      &count,
	}
	go wd.CacheTable.Runner()
	for !wd.CacheTable.IsEmpty() {
		wg.Add(1)
		wd.Sem <- struct{}{}
		go worker(wd)
	}
	wg.Wait()
	close(wd.Sem)
	log.Printf("All files parsed... clearing cache...")
	wd.CacheTable.Close()
}

func worker(wd *WorkerData) {
	apk := wd.CacheTable.GetFilePath()
	defer wd.CacheTable.Completed(apk)
	defer func() { <-wd.Sem }()
	defer wg.Done()
	if apk == "" {
		return
	}
	err := extract(apk, wd.Config)
	if err != nil {
		log.Printf("Warning: " + apk + " is not a valid APK file")
		return
	}
	countMutex.Lock()
	*wd.Count++
	countMutex.Unlock()
	percent := (float64(*wd.Count) / float64(wd.Length) * float64(100))
	_, name := filepath.Split(apk)
	log.Printf("(%.2f%%) Completed: "+name, percent)
}

func initCache(config utils.ConfigData) (*cache.CacheTable, int) {
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
	return cacheTable, length
}

func extract(path string, config utils.ConfigData) error {
	jsonBuilder := make(map[string]interface{})
	plugins := utils.GetPaths(config.CodeDir+"/plugins/", ".so")
	for _, plug := range plugins {
		p, err := plugin.Open(plug)
		utils.Check(err)
		needLock, err := p.Lookup("NeedLock")
		utils.Check(err)
		needLockFunc, ok := needLock.(func() bool)
		if !ok {
			log.Fatal("Error: Malformed NeedLock function in " + plug)
		}

		k, err := p.Lookup("GetKey")
		utils.Check(err)
		keyfunc, ok := k.(func() string)
		if !ok {
			log.Fatal("Error: Malformed GetKey function in " + plug)
		}
		key := keyfunc()

		v, err := p.Lookup("GetValue")
		utils.Check(err)
		valuefunc, ok := v.(func(string, utils.ConfigData) (interface{}, error))
		if !ok {
			log.Fatal("Error: Malformed GetValue function in " + plug)
		}
		if needLockFunc() {
			mutex.Lock()
		}
		result, err := valuefunc(path, config)
		if needLockFunc() {
			mutex.Unlock()
		}
		if err != nil {
			return err
		}
		jsonBuilder[key] = result
	}
	WriteJSON(jsonBuilder, config.OutputDir)
	return nil
}

func WriteJSON(toWrite map[string]interface{}, OutputDir string) {
	data, err := json.Marshal(toWrite)
	utils.Check(err)
	Sha256, ok := toWrite["Sha256"].(string)
	if !ok {
		log.Fatal("Error: Count not validate Sha256 value as a string")
	}
	outputFile := OutputDir + "/" + Sha256 + ".json"
	err = ioutil.WriteFile(outputFile, data, 0644)
	utils.Check(err)
}
