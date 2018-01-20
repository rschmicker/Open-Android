package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "path/filepath"
    "html/template"
)

type Config struct {
	TemplateDir		string
	KeysDir 		string
}

var c = Config{
	TemplateDir: "/home/rschmicker/src/github.com/Open-Android/webserver/templates/",
	KeysDir: "/home/rschmicker/src/github.com/Open-Android/webserver/keys/",
}

func Index(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	templ, err := template.ParseFiles(c.TemplateDir + "index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templ.Execute(w, nil)
    //t := template.New("index")
    //err := t.ExecuteTemplate(w, c.TemplateDir + "index.html", nil)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", Index)
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
    cert, _ := filepath.Abs(c.KeysDir + "tls.crt")
    key, _ := filepath.Abs(c.KeysDir + "tls.key")
    log.Fatal(srv.ListenAndServeTLS(cert, key))
}