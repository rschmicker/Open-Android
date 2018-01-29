package main

import (
	"flag"
	"fmt"
	"github.com/Open-Android/openandroid/controller"
	"github.com/Open-Android/openandroid/utils"
	"log"
	"os"
)

func main() {
	configFlag := flag.String("config", "", "Location to YAML config file.")
	cleanFlag := flag.Bool("clean", false, "Move all apk files to their SHA256 value.")
	vtFlag := flag.Bool("vt", false, "Scan all files through Virus Total.")
	appendFlag := flag.Bool("append", false, "Append new feature extractor data.")
	flag.Parse()

	if len(*configFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	config := utils.ReadConfig(*configFlag)
	config.Clean = *cleanFlag
	config.VtApiCheck = *vtFlag
	config.Append = *appendFlag
	log.Printf("apkDir: " + config.ApkDir)
	log.Printf("outputDir: " + config.OutputDir)
	log.Printf("codeDir: " + config.CodeDir)
	log.Printf("clean: %t", config.Clean)
	log.Printf("vt: %t", config.VtApiCheck)
	log.Printf("append: %t", config.Append)

	controller.Runner(config)
}

func printUsage() {
	fmt.Println(`
Syntax:
	>openandroid -config <YAML config file> [-clean] [-force] [-vtFlag] [-append]

Example:
	>openandroid -config ./openandroid.yaml
		`)
}
