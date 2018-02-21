package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"github.com/Open-Android/webserver/utils"
	"github.com/Open-Android/webserver/api"
)

func main() {
	keyFlag := flag.String("key", "", "Location to HTTPS key.")
	certFlag := flag.String("cert", "", "Location to HTTPS certificate.")
	flag.Parse()
	if (len(*keyFlag) == 0) || (len(*certFlag) == 0) {
		utils.PrintUsage()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.Query)
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
	log.Fatal(srv.ListenAndServeTLS(*certFlag, *keyFlag))
}
