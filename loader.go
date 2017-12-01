package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
)

type Loader interface {
	GetReader() io.Reader
}

// FileLoader handles opening local files
type FileLoader struct {
	filename string
}

func NewFileLoader(filename string) *FileLoader {
	return &FileLoader{
		filename,
	}
}

func (f *FileLoader) GetReader() io.Reader {
	if f == nil || len(f.filename) == 0 {
		log.Fatal("Uninstantiated FileaLoader or missing filename")
		return nil
	}
	
	reader, err := os.Open(f.filename)
	if err != nil {
		log.Fatalf("Problem opening file %s: %s", f.filename, err.Error())
		return nil
	}

	return reader
}

// HttpLoader retrieves data from a remote URL
type HttpLoader struct {
	url string
}

func (h *HttpLoader) GetReader() io.Reader {
	if h == nil || len(h.url) == 0 {
		log.Fatal("Uninstantiated HttpLoader or missing URL")
		return nil
	}

	res, err := http.Get(h.url)
	if err != nil {
		log.Fatalf("Problem accessing url %s: %s", h.url, err.Error())
		return nil
	}

	return base64.NewDecoder(base64.StdEncoding, res.Body) 
}