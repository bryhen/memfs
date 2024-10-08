package memfs

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	testDir       = "./test"
	testFileStr   = "test"
	testFileBsStr = "4.00 Bytes"

	benchDir          = "bench"
	benchFileN        = 100
	benchFileStrIters = 10_000
	benchFileStr      = "0123456789"
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

func BenchmarkFS(b *testing.B) {
	name, bs := setupBenchTemp()
	defer cleanupBenchTemp(name)
	if b.N == 1 {
		b.N = 100
	}

	httpFS := http.FileServer(http.Dir(name))
	fileServer, _, err := New(name)
	if err != nil {
		panic(err)
	}

	size, _, err := DirSize(name)
	if err != nil {
		b.Error(err)
		return
	} else if size != int64(bs) {
		b.Error(fmt.Errorf("wrong bytes have %d want %d", size, bs))
	} else {
		b.Logf("Wrote: %d bytes", bs)
	}

	st := time.Now()
	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchFileStrIters; j++ {
			r.URL.Path = "/" + fmt.Sprint(j)
			w := httptest.NewRecorder()
			httpFS.ServeHTTP(w, r)
		}
	}

	b.Logf("httpFS: %d microseconds", time.Since(st).Microseconds())

	st = time.Now()
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchFileStrIters; j++ {
			r.URL.Path = "/" + fmt.Sprint(j)
			w := httptest.NewRecorder()
			fileServer.ServeHTTP(w, r)
		}
	}

	b.Logf("FileServer: %d microseconds", time.Since(st).Microseconds())
	b.Logf("Iters this round: %d\n\n", b.N)

}

func setupBenchTemp() (string, int) {
	name, err := os.MkdirTemp("./", benchDir)
	if err != nil {
		panic(err)
	}

	var tot int

	for i := 0; i < benchFileN; i++ {
		f, err := os.Create(name + "/" + fmt.Sprint(i))
		if err != nil {
			return name, 0
		}

		for j := 0; j < benchFileStrIters; j++ {
			n, err := f.WriteString(benchFileStr)
			if err != nil {
				return name, 0
			}

			tot += n
		}

		if err := f.Close(); err != nil {
			panic(err)
		}
	}

	return name, tot
}

func cleanupBenchTemp(name string) {
	os.RemoveAll(name)
}
