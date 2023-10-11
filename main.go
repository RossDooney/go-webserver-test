package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	handle := http.FileServer((http.Dir(filepathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handle)))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/reset", cfg.handlerReset)
	mux.HandleFunc("/metrics", cfg.handlerMetrics)

	corsMux := middlewareCors(mux)

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
