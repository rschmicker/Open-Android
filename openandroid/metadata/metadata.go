package metadata

import (
	"crypto/sha256"
	"fmt"
	"github.com/Open-Android/openandroid/utils"
	"io"
	"os"
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
