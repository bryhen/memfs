package memfs

import (
	"net/http"
)

type FileServer struct {
	m        map[string][]byte
	notFound http.Handler
}

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if bs, ok := fs.m[r.URL.Path]; ok {
		w.Write(bs)
	} else {
		fs.notFound.ServeHTTP(w, r)
	}
}

func (fs *FileServer) SetNotFoundHandler(h http.Handler) {
	fs.notFound = h
}
