package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromEnvDefaults(t *testing.T) {
	t.Setenv("ADDR", "")
	t.Setenv("CERT_FILE", "")
	t.Setenv("KEY_FILE", "")
	t.Setenv("PUBLIC_DIR", "")

	cfg := LoadConfigFromEnv()

	if cfg.Addr != defaultAddr {
		t.Fatalf("unexpected default ADDR: %s", cfg.Addr)
	}
	if cfg.CertFile != defaultCertFile {
		t.Fatalf("unexpected default CERT_FILE: %s", cfg.CertFile)
	}
	if cfg.KeyFile != defaultKeyFile {
		t.Fatalf("unexpected default KEY_FILE: %s", cfg.KeyFile)
	}
	if cfg.PublicDir != defaultPublic {
		t.Fatalf("unexpected default PUBLIC_DIR: %s", cfg.PublicDir)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	publicDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(publicDir, "index.html"), []byte("ok"), 0o600); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	handler := NewHandler(publicDir)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "ok" {
		t.Fatalf("expected body 'ok', got %q", rr.Body.String())
	}
}

func TestStaticResponseHasSecurityHeaders(t *testing.T) {
	publicDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(publicDir, "index.html"), []byte("<h1>Hello</h1>"), 0o600); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	handler := NewHandler(publicDir)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Cross-Origin-Opener-Policy"); got != "same-origin" {
		t.Fatalf("unexpected COOP header: %q", got)
	}
	if got := rr.Header().Get("Cross-Origin-Embedder-Policy"); got != "require-corp" {
		t.Fatalf("unexpected COEP header: %q", got)
	}
}

func TestStaticAssetsGetCacheControl(t *testing.T) {
	publicDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(publicDir, "app.js"), []byte("console.log('x')"), 0o600); err != nil {
		t.Fatalf("write app.js: %v", err)
	}

	handler := NewHandler(publicDir)
	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Cache-Control"); got != "public, max-age=3600" {
		t.Fatalf("unexpected Cache-Control header: %q", got)
	}
}

func TestGzipStaticTextResponse(t *testing.T) {
	publicDir := t.TempDir()
	content := strings.Repeat("hello world;", 100)
	if err := os.WriteFile(filepath.Join(publicDir, "app.js"), []byte(content), 0o600); err != nil {
		t.Fatalf("write app.js: %v", err)
	}

	handler := NewHandler(publicDir)
	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip encoding, got %q", got)
	}

	gzr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatalf("create gzip reader: %v", err)
	}
	defer gzr.Close()

	body, err := io.ReadAll(gzr)
	if err != nil {
		t.Fatalf("read gzip body: %v", err)
	}
	if string(body) != content {
		t.Fatalf("unexpected decompressed body")
	}
}
