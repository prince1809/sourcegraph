package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// updateFileIfDifferent will automatically update the file if the contents are
// different. If it does an update ok is true.
func updateFileDifferent(path string, content []byte) (bool, error) {
	current, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		// If the file doesn't exist we write a new file.
		return false, err
	}

	if bytes.Equal(current, content) {
		return false, nil
	}

	// We write to a tempfile first to do the atomic update (via rename)
	f, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return false, err
	}
	// We always remove the tempFile. In the happy case it won't exist.
	defer os.Remove(f.Name())

	if n, err := f.Write(content); err != nil {
		f.Close()
		return false, err
	} else if n != len(content) {
		return false, io.ErrShortWrite
	}

	// fsync to ensure the disk contents are written. This is important, since
	// we are not guaranteed that os.Rename is recorded to disk after f's
	// contents.
	if err := f.Sync(); err != nil {
		f.Close()
		return false, err
	}
	if err := f.Close(); err != nil {
		return false, err
	}
	return true, os.Rename(f.Name(), path)
}
