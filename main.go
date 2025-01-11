package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(cfg.fileserverHits)
		cfg.fileserverHits.Add(1)
		fmt.Print(cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareResetMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)   
	})
}


func main() {
	apiCfg := &apiConfig{}
	multiPlex := http.NewServeMux()
	normDir := http.Dir("./")
	fileServer := http.FileServer(normDir)
	strippedFileServer := http.StripPrefix("/app", fileServer)
	multiPlex.Handle("/app/", apiCfg.middlewareMetricsInc(strippedFileServer))
	multiPlex.Handle("/reset", apiCfg.middlewareResetMetrics(fileServer))
	
	multiPlex.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		hits := apiCfg.fileserverHits.Load()
		fmt.Print(hits)
		io.WriteString(w, fmt.Sprintf("Hits: %d\n", hits))
	})


	multiPlex.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	})
	//fmt.Printf("%T",multiPlex)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        multiPlex,
	}
	log.Fatal(s.ListenAndServe())
}