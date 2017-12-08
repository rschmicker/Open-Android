package utils

import (
	"log"
)

type ConfigData struct {
	ApkDir     string `yaml:"apkDir"`
	DecodedDir string `yaml:"decodedDir"`
	OutputDir  string `yaml:"outputDir"`
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
