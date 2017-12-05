package utils

import(
	"log"
)

type ApkData struct {
	Strings		[]string
	Apis		[]string
	Permissions	[]string
	Md5			string
	Sha256		string
	Sha1		string
	PackageName	string
	Version		string
	Intents		string
}

type ConfigData struct {
	ApkDir		string		`yaml:"apkDir"`
	DecodedDir	string		`yaml:"decodedDir"`
	OutputDir	string		`yaml:"outputDir"`
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}