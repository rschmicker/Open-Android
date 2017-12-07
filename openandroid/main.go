package main

import (
	"flag"
	"fmt"
	"github.com/Open-Android/openandroid/controller"
	"github.com/Open-Android/openandroid/utils"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configFlag := flag.String("config", "", "Location to YAML config file.")
	flag.Parse()

	if len(*configFlag) == 0 {
		printUsage()
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(*configFlag)
	utils.Check(err)
	config := utils.ConfigData{}
	err = yaml.Unmarshal(data, &config)
	utils.Check(err)

	apkDir, err := filepath.Abs(config.ApkDir)
	utils.Check(err)
	decodedDir, err := filepath.Abs(config.DecodedDir)
	utils.Check(err)
	outputDir, err := filepath.Abs(config.OutputDir)
	utils.Check(err)

	log.Printf("apkDir: " + apkDir)
	log.Printf("decodedDir: " + decodedDir)
	log.Printf("outputDir: " + outputDir)

	controller.Run(apkDir, decodedDir, outputDir)
}

func printUsage() {
	fmt.Println(`
Syntax:
	>openandroid -config <YAML config file>

Example:
	>openandroid -config ./openandroid.yaml
		`)
}
