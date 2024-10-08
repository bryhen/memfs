package main

import (
	"log"
	"net/http"

	"example/internal/handler"

	"github.com/bryhen/memfs"
)

func main() {
	fs, failed, err := memfs.New("./app/dist")
	if err != nil {
		panic(err)
	}

	// Prints: The following routes failed to load: invalid url.txt
	log.Println("The following routes failed to load: ", failed)

	// The default handler performs:
	// w.WriteHeader(http.StatusNotFound)
	// w.Write(http.StatusText(http.StatusNotFound))
	fs.NotFoundHandler(handler.NotFoundHandler)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /endpoint", handler.EndpointHandler)
	mux.Handle("GET /", fs)

	s := http.Server{
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Println("Server shutting down: ", err)
	}
}
