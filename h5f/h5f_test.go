package h5f

import (
	"os"
	"sync"
	"testing"
)
import "github.com/valoox/h5go/h5f"

// Tests the H5F package
func TestFCreate(t *testing.T) {
	const testfile = "./empty.h5"
	fid, err := h5f.Create(testfile, h5f.TRUNC,
		h5f.DefaultCreate, h5f.DefaultAccess)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testfile)
	if err := h5f.Close(fid); err != nil {
		t.Fatal(err)
	}
}

// Creates a file and opens it
func TestFOpen(t *testing.T) {
	wg := &sync.WaitGroup{}
	const testfile = "./shared.h5"
	// Creates the file initially
	fid, err := h5f.Create(testfile, h5f.TRUNC|h5f.DEBUG,
		h5f.DefaultCreate, h5f.DefaultAccess)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testfile)
	open := func() {
		defer wg.Done()
		if id, err := h5f.Open(testfile, h5f.RO|h5f.DEBUG,
			h5f.DefaultAccess); err != nil {
			t.Fatal(err)
		} else if err := h5f.Close(id); err != nil {
			t.Fatal(err)
		}
	}
	const n = 2
	for i := 0; i < n; i++ {
		wg.Add(1)
		go open()
	}
	wg.Wait()
	if err := fid.Close(); err != nil {
		t.Fatal(err)
	}
}
