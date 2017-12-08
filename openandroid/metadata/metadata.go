package metadata

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"github.com/Open-Android/openandroid/utils"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type ApkToolConfig struct {
	Name string `yaml:"apkFileName"`
}

type Package struct {
	PackageName string `xml:"package,attr"`
}

type AndroidManifest struct {
	XMLName         xml.Name     `xml:"manifest"`
	UsesPermissions []Permission `xml:"uses-permission"`
	Permissions     []Permission `xml:"permission"`
}

type Permission struct {
	Name string `xml:"name,attr"`
}

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

func GetPackageName(decodedPath string) string {
	androidManifest := decodedPath + "/AndroidManifest.xml"
	data, err := ioutil.ReadFile(androidManifest)
	utils.Check(err)
	packageName := Package{}
	err = xml.Unmarshal(data, &packageName)
	utils.Check(err)
	return packageName.PackageName
}

func GetVersion(decodedPath string) string {
	apkToolOutputPath := decodedPath + "/apktool.yml"
	data, err := ioutil.ReadFile(apkToolOutputPath)
	utils.Check(err)
	stringData := string(data)
	stringData = strings.Split(stringData, "versionName: ")[1]
	stringData = strings.Trim(stringData, "\n")
	return strings.Trim(stringData, " ")
}

func GetApkName(decodedPath string) string {
	apkToolOutputPath := decodedPath + "/apktool.yml"
	data, err := ioutil.ReadFile(apkToolOutputPath)
	utils.Check(err)
	config := ApkToolConfig{}
	err = yaml.Unmarshal(data, &config)
	utils.Check(err)
	return config.Name
}

func GetPermissions(decodedPath string) []string {
	androidManifest := decodedPath + "/AndroidManifest.xml"
	data, err := ioutil.ReadFile(androidManifest)
	utils.Check(err)
	permissions := AndroidManifest{}
	err = xml.Unmarshal(data, &permissions)
	utils.Check(err)
	ret := []string{}
	for _, name := range permissions.UsesPermissions {
		ret = append(ret, name.Name)
	}
	for _, name := range permissions.Permissions {
		ret = append(ret, name.Name)
	}
	return ret
}
