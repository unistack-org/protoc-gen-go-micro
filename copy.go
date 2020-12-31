package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func writeFile(file *object.File, dir string, flag int, mode os.FileMode) error {
	path := filepath.Join(dir, file.Name)

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return err
	}

	w, err := os.OpenFile(path, flag, mode)
	if err != nil {
		return err
	}

	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("Err: failed to close file: %v", err)
		}
	}()

	r, err := file.Reader()
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Err: failed to close file: %v", err)
		}
	}()

	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	return nil
}
