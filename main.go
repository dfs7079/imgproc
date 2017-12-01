package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var Config struct {
	InputFile string
	OutputFile string
	NumGoRoutines int
}

func init() {
	flag.StringVar(&Config.InputFile, "i", "", "input csv")
	flag.StringVar(&Config.OutputFile, "o", "", "output csv")
	flag.IntVar(&Config.NumGoRoutines, "c", 4, "number of concurrent goroutines for image loading/processing (default 4)")
}

//https://1.bp.blogspot.com/-h9n7ExEQBBc/WdgemuGqb1I/AAAAAAAAlGk/kasPu4PNzMcbLdQMNpjvuA3k_TZDy33jwCLcBGAs/s1600/DSC_2126-20171005_1438%2B-%2BNearing%2BPeak%2BColor%2BIn%2BParts%2Bof%2BHigh%2BKnob%2BMassif%2B-%2BWayne%2BBrowning%2BPhotograph%2BPNG.png

func main() {
	// command line parsing
	flag.Parse()
	if len(flag.Args()) < 1 && len(Config.InputFile) == 0 {
		fmt.Println("Usage: imgproc [options] links...")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	// get the list of links, either from provided CSV or command line args
	var links Links
	if len(Config.InputFile) > 0 {
		file, err := os.Open(Config.InputFile + ".csv")
		if err != nil {
			log.Fatalf("Problem reading input CSV: %s", err.Error())
			os.Exit(-1)
		}
		defer file.Close()

		links = NewCsvLinks(file)
	} else {
		links = NewArrayLinks(flag.Args())
	}
	
	// loop through the links, wait for memory to become available
	// then spin up a new goroutine to handle the image
	var err error
	for nextLink, err := links.GetNextLink(); len(nextLink) > 0 && err == nil; nextLink, err = links.GetNextLink() {
		l := CreateLoader(nextLink)

		p := NewTopColorsProcessor(3)

		img := DecodeImage(l.GetReader())

		out := p.ProcessImage(img)

		fmt.Printf("%s:%s\n", nextLink, out)	
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}