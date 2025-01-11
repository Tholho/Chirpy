package main

import (
	//"fmt"
	"log"
	"net/http"
)

func main() {
	multiPlex := http.NewServeMux()
	normDir := http.Dir("./")
	fileServer := http.FileServer(normDir)
	multiPlex.Handle("/", fileServer)
	//fmt.Printf("%T",multiPlex)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        multiPlex,
	}
	log.Fatal(s.ListenAndServe())
}