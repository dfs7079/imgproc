package main

import (
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"regexp"
)

type Loader interface {
	Load() (image.Image, error)
}

// regex used to determine if a link is an HTTP url
const httpRegex string = "^http[s]{0,1}://.*\\.[^/]{2,3}/.*"

// loadingErr is a helper to format a specific error message
func loadingErr(s1, s2 string) error {
	return errors.New(fmt.Sprintf("Problem loading %s: %s", s1, s2))
}

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

func (f *FileLoader) Load() (image.Image, error) {
	if f == nil || len(f.filename) == 0 {
		return nil, errors.New("Uninstantiated FileLoader or missing filename")
	}
	
	reader, err := os.Open(f.filename)
	if err != nil {
		return nil, loadingErr(f.filename, err.Error())
	}

	img, err := decodeImage(reader)
	if err != nil {
		return nil, loadingErr(f.filename, err.Error())
	}

	return img, nil
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

func (h *HttpLoader) Load() (image.Image, error) {
	if h == nil || len(h.url) == 0 {		
		return nil, errors.New("Uninstantiated HttpLoader or missing URL")
	}

	res, err := http.Get(h.url)
	if err != nil {
		return nil, loadingErr(h.url, err.Error())
	}

	img, err := decodeImage(res.Body)
	if err != nil {
		return nil, loadingErr(h.url, err.Error())
	}

	return img, nil
}

// decodeImage reads raw bytes from an io.Reader and attempts to parse them into a golang Image type
func decodeImage(r io.Reader) (image.Image, error) {
	image, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return image, nil
}