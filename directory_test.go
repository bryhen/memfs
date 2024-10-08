package memfs

import (
	"log"
	"testing"
)

func TestReadDir(t *testing.T) {
	ps, bs, invlds, err := readDir("example/app/dist")
	log.Println("PATHS:")
	for _, p := range ps {
		log.Printf("\t%s", p)
	}
	log.Print("\n\n")

	log.Printf("Files Read: %d", len(bs))
	log.Print("\n\n")

	log.Println("Invalids:")
	for _, p := range invlds {
		log.Printf("\t%s", p)
	}
	log.Print("\n\n")

	log.Printf("ERR: %+v", err)

}
