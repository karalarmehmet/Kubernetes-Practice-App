package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Message represents a simple message payload.
type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

// in-memory store (thread-safe)
var (
	messages   = make([]Message, 0, 16)
	mu         sync.RWMutex
	nextID     = 1
	readyAfter = time.Now().Add(2 * time.Second) // readiness gate for demo
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// K8s probes
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("live"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Before(readyAfter) {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		  <!doctype html>
		  <html><body style="font-family:sans-serif">
			<h1>Go API is running 🚀</h1>
			<ul>
			  <li><a href="/api/health">/api/health</a></li>
			  <li><a href="/api/version">/api/version</a></li>
			  <li>Mesajlar: <code>GET /api/messages</code>, <code>POST /api/messages</code></li>
			</ul>
		  </body></html>
		`))
	})

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"version": "1.0.0",
		})
	})

	mux.HandleFunc("/api/messages", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mu.RLock()
			defer mu.RUnlock()
			respondJSON(w, http.StatusOK, messages)
		case http.MethodPost:
			var body struct {
				Text string `json:"text"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			prefix := os.Getenv("MESSAGE_PREFIX")
			mu.Lock()
			msg := Message{ID: nextID, Text: prefix + body.Text, CreatedAt: time.Now().UTC()}
			nextID++
			messages = append(messages, msg)
			mu.Unlock()
			respondJSON(w, http.StatusCreated, msg)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	port := getenv("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("API listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Simple request logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// Very simple CORS middleware for the demo
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
