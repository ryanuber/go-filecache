go-filecache
============

This is a very small file cache library in Golang. It allows you to read a file
from disk with a conditional update based on last modified time from
`os.Stat()`.

Updates are handled by passing in a function to perform the updates.

# Example

The following example simply dumps a timestamp into a file. We will set the max
age to 5 seconds, and run the example in a bash loop to demonstrate the cache
timeout:

```
$ while sleep 1; do go run cachetest.go; done
2014-08-09 13:13:17.435246233 -0700 PDT
2014-08-09 13:13:17.435246233 -0700 PDT
2014-08-09 13:13:17.435246233 -0700 PDT
2014-08-09 13:13:17.435246233 -0700 PDT
2014-08-09 13:13:22.508095087 -0700 PDT
2014-08-09 13:13:22.508095087 -0700 PDT
2014-08-09 13:13:22.508095087 -0700 PDT
2014-08-09 13:13:22.508095087 -0700 PDT
2014-08-09 13:13:27.584300053 -0700 PDT
```

And the code:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ryanuber/go-filecache"
)

func main() {
	// Here we define a function which handles updating our file. In this case,
	// we just dump a timestamp into the file.
	updater := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		stamp := time.Now().String()
		_, err = f.Write([]byte(stamp))
		return err
	}

	// Create a new file cache, passing in the path to the file, the maximum
	// age, and the updater function.
	fc := filecache.New("testcache", 5*time.Second, updater)

	// Retrieve a file handle from the cache. The updater function is invoked
	// during this call if the max age is exceeded.
	fh, err := fc.Get()
	if err != nil {
		return
	}

	content, err := ioutil.ReadAll(fh)
	if err != nil {
		return
	}

	fmt.Println(string(content))
}
```
