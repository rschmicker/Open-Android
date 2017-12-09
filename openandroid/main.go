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
	flag.Parse()

	if len(*configFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	config := utils.ReadConfig(*configFlag)

	log.Printf("apkDir: " + config.ApkDir)
	log.Printf("decodedDir: " + config.DecodedDir)
	log.Printf("outputDir: " + config.OutputDir)
	log.Printf("codeDir: " + config.CodeDir)

	controller.Runner(config.ApkDir, config.DecodedDir, config.OutputDir, config.CodeDir)
}

func printUsage() {
	fmt.Println(`
Syntax:
	>openandroid -config <YAML config file>

Example:
	>openandroid -config ./openandroid.yaml
		`)
}
