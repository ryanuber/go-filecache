package filecache

import (
	"os"
	"time"
)

// FileCacheUpdater is an interface for a function which will handle updating
// the file when it is expired.
type FileCacheUpdater func(string) error

// FileCache is a wrapper around os.File providing simple file caching.
type FileCache struct {
	// The path to the file to cache
	Path string

	// The maximum elapsed time since the last file update.
	MaxAge time.Duration

	// The function used to handle updating the cached file.
	UpdateFunc FileCacheUpdater
}

// New returns a new FileCache with reasonable defaults.
func New(path string) *FileCache {
	cache := &FileCache{
		Path:   path,
		MaxAge: 24 * time.Second,
	}
	return cache
}

// Expired is a predicate which determines if the file should be updated
func (f *FileCache) Expired() bool {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return true
	}

	expireTime := fi.ModTime().Add(f.MaxAge)
	return time.Now().After(expireTime)
}

// Update calls the file update function (if present) on the cached file
func (f *FileCache) Update() error {
	if f.UpdateFunc == nil {
		return nil
	}

	return f.UpdateFunc(f.Path)
}

// Get is used to retrieve an os.File handle. If the cache file has expired,
// this method will update it before opening it and returning the handle.
func (f *FileCache) Get() (*os.File, error) {
	if f.Expired() {
		if err := f.Update(); err != nil {
			return nil, err
		}
	}

	return os.Open(f.Path)
}
