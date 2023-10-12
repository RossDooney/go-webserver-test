package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	r := chi.NewRouter()
	cfg := &apiConfig{
		fileserverHits: 0,
	}

	handle := http.FileServer((http.Dir(filepathRoot)))
	r.Handle("/app", http.StripPrefix("/app", cfg.middlewareMetricsInc(handle)))
	r.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(handle)))
	r.Get("/api/healthz", handlerReadiness)
	r.HandleFunc("/api/reset", cfg.handlerReset)
	r.Get("/api/metrics", cfg.handlerMetrics)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	body := "Hits: " + strconv.Itoa(cfg.fileserverHits)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	fmt.Fprint(w, body)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		fmt.Printf("%v \n", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}
