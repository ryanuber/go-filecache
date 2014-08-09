package filecache

import (
	"bytes"
	"errors"
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

	fc := New(f1.Name(), 1*time.Second, updater)

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

	// Wait for cache to expire
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

func TestCacheNoFile(t *testing.T) {
	fc := New("/no-such-file", 0, nil)
	if !fc.Expired() {
		t.Fatalf("expected non-existent file to be expired")
	}
}

func TestCacheNoUpdateFunc(t *testing.T) {
	data1 := []byte{0}

	f1, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f1.Name())

	f1.Write(data1)

	fc := New(f1.Name(), 1*time.Second, nil)

	// Wait for cache to expire
	time.Sleep(1 * time.Second)

	f2, err := fc.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	content, err := ioutil.ReadAll(f2)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(content, data1) {
		t.Fatalf("bad: %#v", content)
	}
}

func TestCacheError(t *testing.T) {
	updater := func(path string) error {
		return errors.New("test error")
	}

	f1, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f1.Name())

	fc := New(f1.Name(), 1*time.Second, updater)

	// Wait for cache to expire
	time.Sleep(1 * time.Second)

	if _, err := fc.Get(); err == nil {
		t.Fatalf("expected error from update func")
	}
}
