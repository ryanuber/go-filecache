package filecache

import (
	"os"
	"time"
)

// Updater is an interface for a function which will handle updating
// the file when it is expired.
type Updater func(string) error

// Cache is a wrapper around os.File providing simple file caching.
type Cache struct {
	// The path to the file to cache
	Path string

	// The maximum elapsed time since the last file update.
	MaxAge time.Duration

	// The function used to handle updating the cached file.
	UpdateFunc Updater
}

// New is a shortcut function for making a new file cache
func New(path string, maxAge time.Duration, updater Updater) *Cache {
	cache := &Cache{
		Path:       path,
		MaxAge:     maxAge,
		UpdateFunc: updater,
	}
	return cache
}

// Expired is a predicate which determines if the file should be updated
func (f *Cache) Expired() bool {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return true
	}

	expireTime := fi.ModTime().Add(f.MaxAge)
	return time.Now().After(expireTime)
}

// Update calls the file update function (if present) on the cached file
func (f *Cache) Update() error {
	if f.UpdateFunc == nil {
		return nil
	}

	return f.UpdateFunc(f.Path)
}

// Get is used to retrieve an os.File handle. If the cache file has expired,
// this method will update it before opening it and returning the handle.
func (f *Cache) Get() (*os.File, error) {
	if f.Expired() {
		if err := f.Update(); err != nil {
			return nil, err
		}
	}

	return os.Open(f.Path)
}
