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
	testURL      string = "http://x.com"
	defaultName  string = "index"

	kb float64 = 1024
	mb float64 = kb * kb
	gb float64 = mb * kb
)

// Used in tests to check the directory size for compatiblity with a server.
//
// Returns the directory size in bytes, ie 1123830123 or 1024.
// For a readable return, use memfs.DirSizeStr(dir).
func DirSize(dir string) (bytes int64, skippedFiles []string, err error) {
	dir = trimDir(dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		return bytes, skippedFiles, err
	}

	for _, f := range files {
		p := dir + forwardSlash + f.Name()

		if f.IsDir() {

			if bs, sk, nErr := DirSize(p); nErr != nil {
				return bytes, skippedFiles, nErr
			} else {
				bytes += bs
				skippedFiles = append(skippedFiles, sk...)
			}

		} else {

			_, err := url.ParseRequestURI(testURL + p)
			if err != nil {
				skippedFiles = append(skippedFiles, p)
				continue
			}

			f, err := os.Stat(p)
			if err != nil {
				return bytes, skippedFiles, err
			} else {
				bytes += f.Size()
			}
		}
	}
	return
}

// Used in tests to check the directory size for compatiblity with a server.
//
// Returns the directory size as a string, ie "1.12 Gigabytes" or "1.00 Kilobytes".
func DirSizeStr(dir string) (bytes string, skippedFiles []string, err error) {
	bs, sk, err := DirSize(dir)
	if err != nil {
		return bytes, skippedFiles, err
	}
	fSize := float64(bs)
	skippedFiles = sk
	if fSize > gb {
		return fmt.Sprintf("%.2f Gigabytes", fSize/gb), skippedFiles, nil
	} else if fSize > mb {
		return fmt.Sprintf("%.2f Megabytes", fSize/mb), skippedFiles, nil
	} else if fSize > kb {
		return fmt.Sprintf("%.2f Kilobytes", fSize/kb), skippedFiles, nil
	} else {
		return fmt.Sprintf("%.2f Bytes", fSize), skippedFiles, nil
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
		rootP := dir + forwardSlash
		p := rootP + f.Name()

		if f.IsDir() {

			if ps, fbs, invlds, err := readDir(p); err != nil {
				return nil, nil, nil, err
			} else {
				paths = append(paths, ps...)
				fileBs = append(fileBs, fbs...)
				invalid = append(invalid, invlds...)
			}

		} else {

			url, err := url.ParseRequestURI(testURL + p)
			if err != nil {
				invalid = append(invalid, invalidPath(p))
				continue
			}

			fBs, err := os.ReadFile(p)
			if err != nil {
				return nil, nil, nil, err
			} else {
				paths = append(paths, url.Path)
				fileBs = append(fileBs, fBs)
			}

			if strings.HasPrefix(f.Name(), defaultName) {
				paths = append(paths, rootP)
				fileBs = append(fileBs, fBs)
			}
		}
	}

	return paths, fileBs, invalid, nil
}
