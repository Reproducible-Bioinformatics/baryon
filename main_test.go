package main

import (
	"os"
	"path"
	"testing"
)

func Test_getFile(t *testing.T) {
	a, err := getFile("")
	if err != nil {
		t.Fatal("Got error", err)
	}
	if a != os.Stdin {
		t.Fatal("Got wrong file:", a.Name())
	}

	dir := t.TempDir()

	filename := path.Join(dir, "arandomfile")
	f, err := os.Create(filename)
	_, _ = f.WriteString("test")
	f.Close()

	a, err = getFile(filename)
	if err != nil {
		t.Fatal("Got error", err)
	}
	if a.Name() != filename {
		t.Fatal("Got wrong file:", a.Name())
	}
}
