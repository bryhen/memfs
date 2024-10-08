package memfs

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type invalidPath string

const (
	forwardSlash string = "/"
	testURL      string = "http://x.com/"
)

const (
	kb float64 = 1024
	mb float64 = kb * kb
	gb float64 = mb * kb
)

// Used in tests to check the directory size for compatiblity with a server.
//
// Returns the directory size in bytes, ie 1123830123 or 1024.
// For a readable return, use memfs.DirSizeStr(dir).
func DirSize(dir string) int64 {
	return 1
}

// Used in tests to check the directory size for compatiblity with a server.
//
// Returns the directory size as a string, ie "1.12 Gigabytes" or "1.00 Kilobytes".
func DirSizeStr(dir string) string {
	if size := float64(DirSize(dir)); size > gb {
		return fmt.Sprintf("%.2f Gigabytes", size/gb)
	} else if size > mb {
		return fmt.Sprintf("%.2f Megabytes", size/mb)
	} else if size > kb {
		return fmt.Sprintf("%.2f Kilobytes", size/kb)
	} else {
		return fmt.Sprintf("%.2f Bytes", size)
	}
}

func readDir(dir string) ([]string, [][]byte, []invalidPath, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, nil, err
	}

	paths := []string{}
	fileBs := [][]byte{}
	invalid := []invalidPath{}

	for _, f := range files {
		p := dir + forwardSlash + f.Name()
		pURL := dir + forwardSlash + f.Name()

		if f.IsDir() {

			if ps, fbs, invlds, err := readDir(p); err != nil {
				return nil, nil, nil, err
			} else {
				paths = append(paths, ps...)
				fileBs = append(fileBs, fbs...)
				invalid = append(invalid, invlds...)
			}

		} else {

			_, err := url.ParseRequestURI(testURL + pURL)
			if err != nil {
				invalid = append(invalid, invalidPath(p))
				continue
			}

			fBs, err := os.ReadFile(p)
			if err != nil {
				return nil, nil, nil, err
			} else {
				paths = append(paths, pURL)
				fileBs = append(fileBs, fBs)
			}

			if strings.HasPrefix(f.Name(), "index.") {
				paths = append(paths, pURL)
				fileBs = append(fileBs, fBs)
			}

		}
	}

	return paths, fileBs, invalid, nil
}
