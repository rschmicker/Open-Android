package intent

import (
	"bufio"
	"github.com/Open-Android/openandroid/utils"
	"os"
	"strings"
)

func GetIntents(decodedPath string) []string {
	ret := []string{}
	file, err := os.Open(decodedPath + "/AndroidManifest.xml")
	utils.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ".intent.") {
			lines := strings.Split(line, "\"")
			ret = append(ret, lines[1])
		}
	}
	err = scanner.Err()
	utils.Check(err)
	return ret
}
