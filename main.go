package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"featureflagservice/internal/api"
	"featureflagservice/internal/middleware"
	"featureflagservice/internal/store"
)

func main() {
	st := store.New()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /flags", api.CreateFlag(st))
	mux.HandleFunc("GET /flags", api.ListFlags(st))
	mux.HandleFunc("GET /flags/{key}", api.GetFlag(st))
	mux.HandleFunc("PUT /flags/{key}", api.UpdateFlag(st))
	mux.HandleFunc("DELETE /flags/{key}", api.DeleteFlag(st))
	mux.HandleFunc("GET /flags/{key}/evaluate", api.EvaluateFlag(st))
	mux.HandleFunc("GET /healthz", healthz)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           middleware.Logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("featureflagservice listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

// healthz serves the liveness/readiness probe: 200 {"status":"ok"}.
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
