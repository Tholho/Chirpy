package main

import (
	//"fmt"
	"log"
	"net/http"
	"io"
)

func main() {
	multiPlex := http.NewServeMux()
	normDir := http.Dir("./")
	fileServer := http.FileServer(normDir)
	strippedFileServer := http.StripPrefix("/app", fileServer)
	multiPlex.Handle("/", strippedFileServer)
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