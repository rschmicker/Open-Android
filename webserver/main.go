package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "fmt"
    "io"
    "path/filepath"
)

type Config struct {
	TemplateDir		string
	KeysDir 		string
}

var c = Config{
	KeysDir: "/home/rschmicker/src/github.com/Open-Android/webserver/keys/",
}

func getArg(name string, aMap map[string][]string) (arg string, err error) {
    var ok bool
    var arglist []string
    if arglist, ok = aMap[name]; !ok {
        return "", fmt.Errorf("unable to obtain arg from map with key %s", name)
    }
    return arglist[0], nil
}

func Query(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	args := req.URL.Query()
    q, err := getArg("q", args)
    if err != nil {
        log.Println("Error: " + err.Error())
        http.Error(w, "Error: " + err.Error(), http.StatusBadRequest)
        return
    }
    log.Println("q: " + q)
    io.WriteString(w, "q: " + q)
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