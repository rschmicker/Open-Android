package api

import(
	"net/http"
	"io"
	"github.com/Open-Android/webserver/utils"
	"log"
	"os/exec"
	"bytes"
)

func Query(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	args := req.URL.Query()
        toArg, _ := utils.GetArg("to", args)
        fromArg, _ := utils.GetArg("from", args)
	flds, _ := utils.GetArg("fields", args)
        //http.Error(w, InfoMsg, http.StatusBadRequest)

	log.Println("===============================")
	log.Println("From: " + fromArg)
	log.Println("To: " + toArg)
	log.Println("Fields: ")
	log.Println(flds)

	argsP := []string{"run", "query/query.go", "-from", fromArg, "-to", toArg, "-fields", flds}
	cmd := exec.Command("/usr/local/go/bin/go", argsP...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Out: " + out.String())
	log.Println("Error: " + errout.String())
	log.Println("===============================")

	io.WriteString(w, "Your query is being processed\n")
	io.WriteString(w, "Your filename will be the date time of the query\n")
	io.WriteString(w, "Please check back in an hour\n")
	io.WriteString(w, "This file will be purged in a month\n")
}
