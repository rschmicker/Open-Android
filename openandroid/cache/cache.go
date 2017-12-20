package cache

import (
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	//"time"
)

var CacheTableMutex = &sync.Mutex{}

type CacheTable struct {
	Table            []CacheObject
	RamDiskPath      string
	Size             int
	DirectoryToCache string
	Files            []string
	//Location         int
}

type CacheObject struct {
	FilePath  string
	InProcess bool
	Completed bool
}

func (ct *CacheTable) Initialize(config utils.ConfigData) int {
	ct.RamDiskPath = config.CacheDir + "/cache/"
	os.RemoveAll(ct.RamDiskPath)
	err := os.Mkdir(ct.RamDiskPath, 0777)
	utils.Check(err)
	ct.Size = config.CacheSize
	ct.DirectoryToCache = config.ApkDir
	ct.Files = getPaths(ct.DirectoryToCache, ".apk")
	if len(ct.Files) == 0 {
		log.Fatal("No APKs found")
	}
	if ct.Size > len(ct.Files) {
		ct.Size = len(ct.Files)
	}
	//ct.Location = 0
	length := len(ct.Files)
	ct.Populate(ct.Size)
	return length
}

func (ct *CacheTable) Populate(end int) {
	if len(ct.Table) == 0 {
		return
	}
	CacheTableMutex.Lock()
	if end > ct.Size {
		end = ct.Size
	}
	if end > len(ct.Files) {
		end = len(ct.Files)
	}
	for i := 0; i < end; i++ {
		name := metadata.GetApkName(ct.Files[i])
		err := copyFileContents(ct.Files[i], ct.RamDiskPath+name)
		utils.Check(err)
		co := CacheObject{ct.RamDiskPath + name, false, false}
		ct.Table = append(ct.Table, co)
		log.Println("Caching: " + name)
	}
	ct.Files = ct.Files[end:]
	CacheTableMutex.Unlock()
}

func (ct *CacheTable) Runner() {
	for {
		CacheTableMutex.Lock()
		if len(ct.Table) == 0 {
			break
		}
		for i := 0; i < len(ct.Table); i++ {
			file := ct.Table[i]
			if file.Completed {
				err := os.Remove(file.FilePath)
				utils.Check(err)
				ct.Table = append(ct.Table[:i], ct.Table[i+1:]...)
				log.Println("Removed: " + metadata.GetApkName(file.FilePath) + " from cache")
				ct.Populate(1)
			}
		}
		if len(ct.Table) == 0 && len(ct.Files) == 0 {
			CacheTableMutex.Unlock()
			return
		}
		CacheTableMutex.Unlock()
		//time.Sleep(5 * time.Second)
	}
}

func (ct *CacheTable) Completed(path string) {
	name := metadata.GetApkName(path)
	for i := 0; i < len(ct.Table); i++ {
		if strings.Contains(ct.Table[i].FilePath, name) {
			ct.Table[i].Completed = true
		}
	}
}

func (ct *CacheTable) IsEmpty() bool {
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	return !(len(ct.Table) > 0)
}

func (ct *CacheTable) GetFilePath() string {
	path := ""
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	if len(ct.Table) == 0 {
		return path
	}
	for i := 0; i < len(ct.Table); i++ {
		if !ct.Table[i].InProcess {
			path = ct.Table[i].FilePath
			ct.Table[i].InProcess = true
			break
		}
	}

	return path
}

func (ct *CacheTable) Close() {
	err := os.RemoveAll(ct.RamDiskPath)
	utils.Check(err)
}

func getPaths(dir string, Containing string) []string {
	fileList := make([]string, 0)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, Containing) {
			fileList = append(fileList, path)
		}
		return err
	})
	utils.Check(err)
	return fileList
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
