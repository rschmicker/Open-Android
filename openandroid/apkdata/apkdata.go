package apkdata

import (
	"encoding/json"
	"github.com/Open-Android/openandroid/metadata"
	"github.com/Open-Android/openandroid/utils"
	"io/ioutil"
	"strings"
)

type ApkData struct {
	Strings     []string
	Apis        []string
	Permissions []string
	Md5         string
	Sha256      string
	Sha1        string
	PackageName string
	Version     string
	Intents     []string
	Malicious   bool
}

func (apk *ApkData) IsMalicious(apkPath string) {
	if strings.Contains(apkPath, "malware") {
		apk.Malicious = true
	} else {
		apk.Malicious = false
	}
}

func (apk *ApkData) WriteJSON(OutputDir string) {
	data, err := json.Marshal(apk)
	utils.Check(err)
	outputFile := OutputDir + "/" + apk.Sha256 + ".json"
	err = ioutil.WriteFile(outputFile, data, 0644)
	utils.Check(err)
}

func (apk *ApkData) GetMetaData(apkPath string) {
	apk.Md5 = metadata.Md5File(apkPath)
	apk.Sha1 = metadata.Sha1File(apkPath)
	apk.Sha256 = metadata.Sha256File(apkPath)
	apk.Version = metadata.GetVersion(apkPath)
	apk.PackageName = metadata.GetPackageName(apkPath)
	apk.Permissions = metadata.GetPermissions(apkPath)
}
