package utils

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ConfigData struct {
	ApkDir    string `yaml:"apkDir"`
	OutputDir string `yaml:"outputDir"`
	CodeDir   string `yaml:"codeDir"`
	CacheDir  string `yaml:"cacheDir"`
	CacheSize int    `yaml:"cacheSize"`
	Clean     bool
	Force     bool
}

func ReadConfig(configPath string) ConfigData {
	data, err := ioutil.ReadFile(configPath)
	Check(err)
	config := ConfigData{}
	err = yaml.Unmarshal(data, &config)
	Check(err)

	config.ApkDir, err = filepath.Abs(config.ApkDir)
	Check(err)
	config.OutputDir, err = filepath.Abs(config.OutputDir)
	Check(err)
	config.CodeDir, err = filepath.Abs(config.CodeDir)
	Check(err)
	config.CacheDir, err = filepath.Abs(config.CacheDir)
	Check(err)

	return config
}

func GetPaths(dir string, Containing string) []string {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, Containing) {
			fileList = append(fileList, path)
		}
		return err
	})
	Check(err)
	return fileList
}

func CrossCompare(todoFiles []string, doneFiles []string) []string {
	ret := []string{}
	found := false
	for _, todo := range todoFiles {
		_, name := filepath.Split(todo)
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

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
