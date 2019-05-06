package vfsutil

import "os"

// ArchiveCacheDir is the location on disk that archives are cached. It is
// configurable so that in production we can point it into CACHE_DIR
var ArchiveCacheDir = "/tmp/vfsutil-archive-cache"

// Evicter implements Evict
type Evicter interface {
	// Evict evicts an item from a cache.
	Evict()
}

type cachedFile struct {
	// File is an open FD to the fetched data
	File *os.File

	// path is the disk path for File
	path string
}

// Evict will remove the file from the cache. It does not close File. It also
// does not protect against other open readers or concurrent fetches.
func (f *cachedFile) Evict() {
	panic("implement me")
}

