package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	"github.com/Open-Android/openandroid/utils"
	"github.com/Open-Android/openandroid/controller"
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

	log.Printf("apkDir: " + config.ApkDir)
	log.Printf("decodedDir: " + config.DecodedDir)
	log.Printf("outputDir: " + config.OutputDir)

	controller.Run(config.ApkDir, config.DecodedDir, config.OutputDir)
}

func printUsage() {
	fmt.Println(`
Syntax:
	>openandroid -config <YAML config file>

Example:
	>openandroid -config ./openandroid.yaml
		`)
}