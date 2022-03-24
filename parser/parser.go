package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Parse walks the directories tree and builds the final document.
func Parse(uri string) ([]Fragment, error) {
	if ok, _ := isDirectory(uri); ok {
		return parseDir(uri)
	}

	if ok, _ := isZipFile(uri); ok {
		return parseZip(uri)
	}

	return nil, fmt.Errorf("%s is not a zip archive, nor a directory, nor an http url", uri)
}

// parseDir parse a ZIP archive to create YAML blocks
func parseZip(path string) ([]Fragment, error) {
	ok, err := isZipFile(filepath.Clean(path))
	if !ok {
		return nil, fmt.Errorf("'%s' is not a zip archive", path)
	}
	if err != nil {
		return nil, err
	}

	fsys, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer fsys.Close()

	return parseFS(fsys)
}

// parseDir parse a directory to create YAML blocks
func parseDir(dir string) ([]Fragment, error) {
	ok, err := isDirectory(filepath.Clean(dir))
	if !ok {
		return nil, fmt.Errorf("'%s' is not a directory", dir)
	}
	if err != nil {
		return nil, err
	}

	return parseFS(os.DirFS(dir))
}

// parseFS walks the fileSystem parsing all the foders and YAML fragments
func parseFS(fsys fs.FS) ([]Fragment, error) {
	// walk the current directory
	blocks := make([]Fragment, 0)
	err := fs.WalkDir(fsys, ".", treeWalker(fsys, &blocks))
	if err != nil {
		return nil, err
	}

	return blocks, err
}

// treeWalker walks the FS parsing and assembling YAML fragments
func treeWalker(fsys fs.FS, blocks *[]Fragment) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == "." {
			return nil
		}

		level := depth(path)

		// a directory is a key
		if d.IsDir() {
			*blocks = append(*blocks, Fragment{
				path:    path,
				depth:   level,
				isKey:   true,
				content: filepath.Base(path),
			})

			return nil
		}

		file, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		src, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		*blocks = append(*blocks, Fragment{
			path:    path,
			depth:   level,
			isKey:   false,
			content: strings.TrimRightFunc(string(src), unicode.IsSpace),
		})

		return nil
	}
}

// depth compute the path levels
func depth(path string) int {
	parts := strings.Split(path, string(os.PathSeparator))
	return len(parts) - 1
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

// isZipFile checks if the file is a ZIP archive
func isZipFile(path string) (bool, error) {
	sig := []byte("PK\x03\x04")

	fp, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fp.Close()

	var header [32]byte
	_, err = io.ReadFull(fp, header[:])
	if err != nil {
		return false, err
	}

	if bytes.HasPrefix(header[:], sig) {
		return true, nil
	}

	return false, nil
}
