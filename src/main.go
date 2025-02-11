package main

import (
        "fmt"
        "net/http"
        "log"
        "time"
        "crypto/tls"
)

func main() {
        // use environment variables ?
        certFile := "./certs/srvr.crt"
        privKey := "./certs/priv.key"

        mux := http.NewServeMux()
        fsHandler := http.FileServer(http.Dir("./public"))

        mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
                w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
                w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
                fsHandler.ServeHTTP(w, req)
        })

        srvr := &http.Server {
                Addr: "0.0.0.0:8080",
                Handler: mux,
                ReadTimeout: time.Second * 10,
                WriteTimeout: time.Second *10,
                TLSConfig: &tls.Config {
                        InsecureSkipVerify: true,
                },
        }

        fmt.Printf("\n---\nServer is running at http://%s\n", srvr.Addr)

        err := srvr.ListenAndServeTLS(certFile, privKey); if err != nil {
                log.Fatal(err)
        }
}
