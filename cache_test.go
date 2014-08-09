package filecache

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	data1 := []byte{0}
	data2 := []byte{1}

	updater := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(data2)
		return err
	}

	f1, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f1.Name())

	f1.Write(data1)

	fc := New(f1.Name())
	fc.MaxAge = 1 * time.Second
	fc.UpdateFunc = updater

	f2, err := fc.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	content, err := ioutil.ReadAll(f2)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(content, data1) {
		t.Fatalf("bad: %s", content)
	}

	// Wait for cache to expire...
	time.Sleep(1 * time.Second)

	f3, err := fc.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	content, err = ioutil.ReadAll(f3)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(content, data2) {
		t.Fatalf("bad: %#v", content)
	}
}
