package main

import (
	"os"
	"errors"
)

type Links interface {
	GetNextLink() (string, error)
}

type ArrayLinks struct {
	links []string
	index int
}

func NewArrayLinks(array []string) *ArrayLinks {
	return &ArrayLinks{
		links: array,
		index: 0,
	}
}

func (a *ArrayLinks) GetNumLinks() int {
	if a == nil {
		return -1
	}
	return len(a.links)
}

func (a *ArrayLinks) GetNextLink() (string, error) {
	if a == nil {
		return "", errors.New("ArrayLinks not instantiated")
	}

	if a.index < a.GetNumLinks() {
		a.index++
		return a.links[a.index-1], nil
	}

	return "", nil
}

// represents links streamed from an ascii-encoded csv file
type CsvLinks struct{
	fh *os.File
}

func NewCsvLinks(file *os.File) *CsvLinks {
	return &CsvLinks{
		fh: file,
	}
}

func (c *CsvLinks) GetNextLink() (string, error) {
	if c.fh == nil {
		return "", errors.New("CsvLinks::GetNextLink invalid file handle")
	}

	var err error
	b := make([]byte, 1)
	link := ""
	for n, err := c.fh.Read(b); n != 0 && err == nil; n, err = c.fh.Read(b) {
		if b[0] != byte(',') {
			link = link + string(b)
			continue
		}
		break
	}

	return link, err
}