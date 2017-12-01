package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Loader interface {
	GetReader() io.Reader
}

const httpRegex string = "^http[s]{0,1}://.*\\.[^/]{2,3}/.*"

// CreateLoader is a factory that instantiates the correct form of Loader based on the link format
func CreateLoader(link string) Loader {
	if ok, _ := regexp.MatchString(httpRegex, link); ok {
		return NewHttpLoader(link)
	}
	// everything else just tries to open a file for now
	return NewFileLoader(link)
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
		log.Fatalf("Problem opening file: %s", err.Error())
		return nil
	}

	return reader
}

// HttpLoader retrieves data from a remote URL
type HttpLoader struct {
	url string
}

func NewHttpLoader(url string) *HttpLoader {
	return &HttpLoader{
		url,
	}
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

	return res.Body 
}