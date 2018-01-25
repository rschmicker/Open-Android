package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	TemplateDir string
	KeysDir     string
}

var c = Config{
	KeysDir: "/Users/rschmicker/src/github.com/Open-Android/webserver/keys/",
}

const InfoMsg string = `
No GET parameters found!

Here's some brief info on how to query:

q=
     - Query search parameter (Required)
     - If you'd like to download everything, use q=*:*
     - Search using regex or for a particular text with q=<field name>:<regex or literal>
     - Example: q=Malicious:true
        - Returns all APKs that are malicious

fl=
     - Return only these fields
     - Example: fl=Permissions+Apis+PackageName
        - Returns only the fields Permissions, Apis, and PackageName

omitHeader=
     - Remove extra query info
     - Example: omitHeader=true

rows=
     - Number of entries to return
     - Example: rows=7000
         - Returns top 7000 entries of given query

wt=
     - Writer Type
     - Type of file that is returned from the query
     - json, csv, xml, etc.
     - Default of json
     - Example: wt=json

More detailed descriptions and examples can be seen on our wiki: https://github.com/rschmicker/Open-Android/wiki

`

func getArg(name string, aMap map[string][]string) (arg string, err error) {
	var ok bool
	var arglist []string
	if arglist, ok = aMap[name]; !ok {
		return "", fmt.Errorf("unable to obtain arg from map with key %s", name)
	}
	return arglist[0], nil
}

func getSolrQuery(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Warning: " + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Warning: " + err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Warning: " + err.Error())
	}
	return string(body)
}

func Query(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	args := req.URL.Query()
	q, err := getArg("q", args)
	if err != nil {
		log.Println("Warning: " + err.Error())
		http.Error(w, InfoMsg, http.StatusBadRequest)
		return
	}
	fl, _ := getArg("fl", args)
	omitHeader, _ := getArg("omitHeader", args)
	rows, _ := getArg("rows", args)
	wt, _ := getArg("wt", args)

	solrQuery := "http://localhost:8983/solr/apks/select?"
	solrQuery += "q=" + q
	if len(fl) != 0 {
		solrQuery += "&fl=" + fl
	}
	if len(omitHeader) != 0 {
		solrQuery += "&omitHeader=" + omitHeader
	}
	if len(rows) != 0 {
		solrQuery += "&rows=" + rows
	}
	if len(wt) != 0 {
		solrQuery += "&wt=" + wt
	} else {
		solrQuery += "&wt=json"
	}
	io.WriteString(w, "solrQuery: "+solrQuery)
	io.WriteString(w, getSolrQuery(solrQuery))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Query)
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	cert, _ := filepath.Abs(c.KeysDir + "server.crt")
	key, _ := filepath.Abs(c.KeysDir + "server.key")
	log.Fatal(srv.ListenAndServeTLS(cert, key))
}
