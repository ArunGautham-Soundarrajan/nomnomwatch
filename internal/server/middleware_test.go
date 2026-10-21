package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogRequests(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello")) })
	mux.HandleFunc("GET /fail", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {})
	h := logRequests(mux)

	tests := []struct {
		path      string
		htmx      bool
		wantLog   bool
		wantLevel string
		wantCode  int
	}{
		{path: "/ok?year=2026", htmx: true, wantLog: true, wantLevel: "INFO", wantCode: 200},
		{path: "/fail", wantLog: true, wantLevel: "ERROR", wantCode: 500},
		{path: "/missing", wantLog: true, wantLevel: "WARN", wantCode: 404},
		{path: "/health", wantLog: false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			buf.Reset()
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.htmx {
				req.Header.Set("HX-Request", "true")
			}
			h.ServeHTTP(httptest.NewRecorder(), req)

			if !tt.wantLog {
				if buf.Len() != 0 {
					t.Fatalf("expected no log at info level, got %s", buf.String())
				}
				return
			}
			var entry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
				t.Fatalf("log line is not JSON: %v: %s", err, buf.String())
			}
			if entry["level"] != tt.wantLevel || entry["status"] != float64(tt.wantCode) || entry["htmx"] != tt.htmx {
				t.Errorf("got level=%v status=%v htmx=%v, want %s %d %v", entry["level"], entry["status"], entry["htmx"], tt.wantLevel, tt.wantCode, tt.htmx)
			}
		})
	}
}
