package stringApk

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func GetStrings(ApkDir string, CodeDir string) []string {
	prog := "java"
	args := []string{"-Dfile.encoding=UTF-8",
		"-cp",
		CodeDir + "/stringApk/:./apis/Rapid.jar",
		"StringParser",
		ApkDir}
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		return []string{"ERROR:	no DEX file is found in the APK file."}
	}
	if errout.String() != "" {
		log.Printf(errout.String())
	}
	data := strings.Split(out.String(), "\n")
	data = append(data[5:], data[:5+1]...)
	return data
}
