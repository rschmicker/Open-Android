package utils

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type ConfigData struct {
	ApkDir    string `yaml:"apkDir"`
	OutputDir string `yaml:"outputDir"`
	CodeDir   string `yaml:"codeDir"`
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

	return config
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
