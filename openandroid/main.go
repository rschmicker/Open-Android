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
	forceFlag := flag.Bool("force", false, "Re-parse apk already completed in output folder.")
	vtFlag := flag.Bool("vt", false, "Scan all files through Virus Total.")
	flag.Parse()

	if len(*configFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	config := utils.ReadConfig(*configFlag)
	config.Clean = *cleanFlag
	config.Force = *forceFlag
	config.VtApiCheck = *vtFlag
	log.Printf("apkDir: " + config.ApkDir)
	log.Printf("outputDir: " + config.OutputDir)
	log.Printf("codeDir: " + config.CodeDir)
	log.Printf("clean: %t", config.Clean)
	log.Printf("force: %t", config.Force)
	log.Printf("vt: %t", config.VtApiCheck)

	controller.Runner(config)
}

func printUsage() {
	fmt.Println(`
Syntax:
	>openandroid -config <YAML config file> [-clean] [-force] [-vtFlag]

Example:
	>openandroid -config ./openandroid.yaml
		`)
}
