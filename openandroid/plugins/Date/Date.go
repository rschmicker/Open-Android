package main

import (
	"github.com/Open-Android/openandroid/utils"
	"os"
)

func NeedLock() bool { return false }

func GetKey() string { return "DateLastModified" }

func GetValue(path string, config utils.ConfigData) (interface{}, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	mod := info.ModTime()
	return mod.Format("2006-01-02 15:04:05"), nil
}
