package memfs

import (
	"net/http"
	"strings"
)

var notFound = []byte(http.StatusText(http.StatusNotFound))

// Creates a new FileServer that has loaded the directory.
func New(dir string) (*FileServer, []invalidPath, error) {
	if strings.HasPrefix(dir, "."+forwardSlash) {
		dir = dir[2:]
	}

	paths, files, invalid, err := readDir(dir)
	if err != nil {
		return nil, []invalidPath{}, err
	}

	fileServer := &FileServer{
		m:        map[string][]byte{},
		notFound: http.HandlerFunc(defaultNotFound),
	}

	for i, p := range paths {
		fileServer.m[p] = files[i]
	}

	return fileServer, invalid, nil
}

func defaultNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(notFound)
}
