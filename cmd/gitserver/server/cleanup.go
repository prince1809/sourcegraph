package server

import (
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SetupAndClearTmp sets up the tempdir for ReposDir as well as clearing it
// out. It returns the temporary directory location.
func (s *Server) SetupAndClearTmp() (string, error) {
	// Additionally we create directories with the prefix .tmp-old which are
	// asynchronously removed. We do not remove in place since it may be a
	// slow operation to block on. Our tmp dir will be ${s.ReposDir}/.tmp
	dir := filepath.Join(s.ReposDir, tempDirName) // .tmp
	oldPrefix := tempDirName + "-old"
	if _, err := os.Stat(dir); err == nil {
		// Rename the current tmp file so we can asynchronously remove it. Use
		// a consistent pattern so if we get interrupted, so we can clean it
		// another name.
		oldTmp, err := ioutil.TempDir(s.ReposDir, oldPrefix)
		if err != nil {
			return "", err
		}
		// oldTmp dir exists, so we need to use a child of oldTmp as the
		// rename target.
		if err := os.Rename(dir, filepath.Join(oldTmp, tempDirName)); err != nil {
			return "", err
		}
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// Asynchronously remove old temporary directories
	files, err := ioutil.ReadDir(s.ReposDir)
	if err != nil {
		log15.Error("failed to do tmp cleanup", "error", err)
	} else {
		for _, f := range files {
			// Remove older .tmp directories as well as our older tmp-
			// directories we would place into ReposDir.
			if !strings.HasPrefix(f.Name(), oldPrefix) && !strings.HasPrefix(f.Name(), "tmp-") {
				continue
			}
			go func(path string) {
				if err := os.RemoveAll(path); err != nil {
					log15.Error("cleanup: failed to remove old temporary directory", "path", "error", err)
				}
			}(filepath.Join(s.ReposDir, f.Name()))
		}
	}
	return dir, nil
}
