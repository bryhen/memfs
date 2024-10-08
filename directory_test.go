package memfs

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	testDir       = "./test"
	testFileStr   = "test"
	testFileBsStr = "4.00 Bytes"

	benchDir     = "bench"
	benchFileStr = "0123456789"
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

func TestFileSize(t *testing.T) {
	bs, sk, err := DirSize(testDir)
	if err != nil {
		t.Error(err)
		return
	}

	if len(sk) != 0 {
		t.Error(fmt.Errorf("found invalid url"))
		return
	}

	if bs != int64(len(testFileStr)) {
		t.Error(fmt.Errorf("wrong byte size %d", bs))
		return
	}

	t.Log("File size is correct")
}

func TestFileSizeStr(t *testing.T) {
	bs, sk, err := DirSizeStr(testDir)
	if err != nil {
		t.Error(err)
		return
	}

	if len(sk) != 0 {
		t.Error(fmt.Errorf("found invalid url"))
		return
	}

	if bs != testFileBsStr {
		t.Error(fmt.Errorf("wrong byte size %s", bs))
		return
	}

	t.Logf("File size is correct at: %s", bs)
}

func BenchmarkXxx(b *testing.B) {
	n := setupBenchTemp()
	defer cleanupBenchTemp(n)
	time.Sleep(time.Second * 5)
}

func setupBenchTemp() string {
	name, err := os.MkdirTemp("./", benchDir)
	if err != nil {
		panic(err)
	}

	return name
}

func cleanupBenchTemp(name string) {
	os.RemoveAll(name)
}
