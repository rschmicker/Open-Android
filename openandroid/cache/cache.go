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
	"syscall"
)

var CacheTableMutex = &sync.Mutex{}

type CacheTable struct {
	Table            []CacheObject
	RamDiskPath      string
	Length           int
	Size             uint64
	CurrentSize      uint64
	DirectoryToCache string
	Files            []string
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
	ct.CurrentSize = 0
	ct.DirectoryToCache = config.ApkDir
	toDoFiles := GetPaths(ct.DirectoryToCache, ".apk")
	doneFiles := GetPaths(config.OutputDir, ".json")
	ct.Files = CrossCompare(toDoFiles, doneFiles)
	if len(ct.Files) == 0 {
		log.Fatal("No Files found")
	}
	ct.Size = ct.AvailableRamSpace() / 2
	ct.Populate()
	return len(ct.Files)
}

func (ct *CacheTable) Populate() {
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	end := 0
	for i := 0; i < len(ct.Files); i++ {
		if (ct.CurrentSize + GetFileSize(ct.Files[i])) > ct.Size {
			end = i
			break
		}
		name := metadata.GetApkName(ct.Files[i])
		err := copyFileContents(ct.Files[i], ct.RamDiskPath+name)
		utils.Check(err)
		ct.CurrentSize += GetFileSize(ct.Files[i])
		co := CacheObject{ct.RamDiskPath + name, false, false}
		ct.Table = append(ct.Table, co)
		log.Println("Caching: " + name)
		log.Printf("Cache size(%v): %v/%v", len(ct.Table), ct.CurrentSize/1024/1024, ct.Size/1024/1024)
	}
	ct.Files = ct.Files[end:]
}

func (ct *CacheTable) AvailableRamSpace() uint64 {
	var stat syscall.Statfs_t
	syscall.Statfs(ct.RamDiskPath, &stat)
	return (stat.Bavail * uint64(stat.Bsize))
}

func GetFileSize(path string) uint64 {
	fi, err := os.Stat(path)
	utils.Check(err)
	return uint64(fi.Size())
}

func (ct *CacheTable) Runner() {
	for {
		if ct.IsEmpty() {
			break
		}
		ct.GarbageCollector()
		ct.Populate()
	}
}

func (ct *CacheTable) GarbageCollector() {
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	for i := 0; i < len(ct.Table); i++ {
		file := ct.Table[i]
		if file.Completed {
			fileSize := GetFileSize(file.FilePath)
			err := os.Remove(file.FilePath)
			utils.Check(err)
			ct.Table[i] = ct.Table[len(ct.Table)-1]
			ct.Table = ct.Table[:len(ct.Table)-1]
			log.Println("Removed: " + metadata.GetApkName(file.FilePath) + " from cache")
			ct.CurrentSize -= fileSize
		}
	}
}

func (ct *CacheTable) Completed(path string) {
	name := metadata.GetApkName(path)
	for i := 0; i < len(ct.Table); i++ {
		ramName := metadata.GetApkName(ct.Table[i].FilePath)
		if ramName == name {
			ct.Table[i].Completed = true
		}
	}
}

func (ct *CacheTable) IsEmpty() bool {
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	return (len(ct.Table) == 0 && len(ct.Files) == 0)
}

func (ct *CacheTable) IsNotEmpty() bool {
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
	return (len(ct.Table) == 0 && len(ct.Files) != 0)
}

func (ct *CacheTable) GetFilePath() string {
	path := ""
	CacheTableMutex.Lock()
	defer CacheTableMutex.Unlock()
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

func CrossCompare(todoFiles []string, doneFiles []string) []string {
	ret := []string{}
	found := false
	for _, todo := range todoFiles {
		name := metadata.GetApkName(todo)
		name = name[:len(name)-4]
		for _, done := range doneFiles {
			if strings.Contains(done, name) {
				found = true
				log.Printf("Skipping: %v already completed...", todo)
				break
			}
		}
		if found == false {
			ret = append(ret, todo)
		}
		found = false
	}
	return ret
}

func GetPaths(dir string, Containing string) []string {
	fileList := []string{}
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
