package metadata

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/Open-Android/openandroid/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Sha256File(fileName string) string {
	f, err := os.Open(fileName)
	utils.Check(err)
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	utils.Check(err)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5File(fileName string) string {
	f, err := os.Open(fileName)
	utils.Check(err)
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	utils.Check(err)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha1File(fileName string) string {
	f, err := os.Open(fileName)
	utils.Check(err)
	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	utils.Check(err)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetPackageName(path string) string {
	prog := "aapt"
	args := []string{
		"dump",
		"permissions",
		path}
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		return ""
	}
	if errout.String() != "" {
		log.Printf(errout.String())
	}
	data := strings.Split(out.String(), "\n")
	return strings.Split(data[0], "package: ")[1]
}

func GetVersion(path string) string {
	prog := "aapt"
	args := []string{
		"dump",
		"badging",
		path}
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		return ""
	}
	if errout.String() != "" {
		log.Printf(errout.String())
	}
	data := strings.Split(out.String(), "\n")
	version := data[0]
	version = strings.Split(version, "versionName='")[1]
	version = strings.Split(version, "'")[0]
	return version
}

func GetApkName(path string) string {
	name := strings.Split(path, "/")
	return name[len(name)-1]
}

func GetPermissions(path string) []string {
	prog := "aapt"
	args := []string{
		"dump",
		"permissions",
		path}
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	if errout.String() != "" {
		log.Printf(errout.String())
	}
	tmp := strings.Split(out.String(), "\n")
	tmp = tmp[1:]
	data := []string{}
	for _, line := range tmp {
		if !strings.Contains(line, "permission: ") {
			continue
		}
		line = strings.Trim(line, " ")
		line = strings.Split(line, "permission: ")[1]
		if strings.Contains(line, "name='") {
			line = strings.Split(line, "name='")[1]
		}
		line = strings.Trim(line, "'")
		data = append(data, line)
	}
	return data
}
