package server

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultAddr     = "0.0.0.0:8080"
	defaultCertFile = "./certs/srvr.crt"
	defaultKeyFile  = "./certs/priv.key"
	defaultPublic   = "./public"
)

type Config struct {
	Addr      string
	CertFile  string
	KeyFile   string
	PublicDir string
}

func LoadConfigFromEnv() Config {
	return Config{
		Addr:      getEnvOrDefault("ADDR", defaultAddr),
		CertFile:  getEnvOrDefault("CERT_FILE", defaultCertFile),
		KeyFile:   getEnvOrDefault("KEY_FILE", defaultKeyFile),
		PublicDir: getEnvOrDefault("PUBLIC_DIR", defaultPublic),
	}
}

func (c Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("ADDR cannot be empty")
	}

	if err := fileMustExist(c.CertFile); err != nil {
		return fmt.Errorf("CERT_FILE invalid: %w", err)
	}

	if err := fileMustExist(c.KeyFile); err != nil {
		return fmt.Errorf("KEY_FILE invalid: %w", err)
	}

	info, err := os.Stat(c.PublicDir)
	if err != nil {
		return fmt.Errorf("PUBLIC_DIR invalid: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("PUBLIC_DIR must be a directory: %s", c.PublicDir)
	}

	return nil
}

func BuildTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
}

func BuildHTTPServer(cfg Config) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           NewHandler(cfg.PublicDir),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig:         BuildTLSConfig(),
	}
}

func NewHandler(publicDir string) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(publicDir))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	staticChain := securityHeaders(cacheHeaders(gzipMiddleware(fileServer)))
	mux.Handle("/", staticChain)

	return mux
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		next.ServeHTTP(w, r)
	})
}

func cacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := strings.ToLower(filepath.Ext(r.URL.Path))
		switch ext {
		case ".html":
			w.Header().Set("Cache-Control", "no-cache")
		case ".js", ".css", ".wasm", ".json", ".pck", ".data", ".png", ".jpg", ".jpeg", ".webp", ".svg":
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		ext := strings.ToLower(filepath.Ext(r.URL.Path))
		switch ext {
		case ".html", ".js", ".css", ".json", ".txt", ".svg":
		default:
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Del("Content-Length")

		gzw := gzip.NewWriter(w)
		defer gzw.Close()

		gzrw := &gzipResponseWriter{ResponseWriter: w, writer: gzw}
		next.ServeHTTP(gzrw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.writer.Write(p)
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func fileMustExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("expected file but got directory: %s", path)
	}
	return nil
}
