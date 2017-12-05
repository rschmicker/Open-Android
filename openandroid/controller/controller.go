package controller

import(
	"path/filepath"
	"strings"
	"log"
	"os"
	"bytes"
	"os/exec"
	"github.com/Open-Android/openandroid/utils"
)

func Run(ApkDir string, DecodedDir string, OutputDir string) {
	paths := getApkPaths(ApkDir)
	if len(paths) == 0 {
		log.Fatal("No APKs found")
	}
	decode(paths, DecodedDir)
}

func getApkPaths(ApkDir string) ([]string) {
	fileList := make([]string, 0)
	err := filepath.Walk(ApkDir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".apk") {
			fileList = append(fileList, path)
		}
		return err
	})
    utils.Check(err)
    return fileList
}

func decode(ApkPaths []string, DecodedDir string) {
	for _, ApkPath := range ApkPaths {
		cmd := exec.Command("apktool", "d", "--quiet", "-o", DecodedDir, ApkPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		log.Printf(out.String())
		utils.Check(err)
	}
}