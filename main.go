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
	MaxLinks int
	MaxImgProcs int
	ShowErrors bool
}

func init() {
	flag.StringVar(&Config.InputFile, "i", "", "input csv")
	flag.StringVar(&Config.OutputFile, "o", "", "output csv")
	flag.IntVar(&Config.MaxImgProcs, "maximgprocs", 1000, "number of images to be processed concurrently (default 1000)")
	flag.IntVar(&Config.MaxLinks, "maxlinks", 400, "Number of links to read from the input stream at a time")
	flag.BoolVar(&Config.ShowErrors, "e", false, "Whether to include error messages in the output")
}

func main() {
	// command line parsing
	flag.Parse()
	if len(flag.Args()) < 1 && len(Config.InputFile) == 0 {
		fmt.Println("Usage: imgproc [options] links...")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	// get the list of links, either from provided CSV or command line args
	links := initLinks(Config.InputFile)
	defer links.Close()

	// init output formatter based on config
	outputter := initOutputter(Config.OutputFile)
	defer outputter.Close()

	// continuously read links from the source
	// we buffer linksChan to limit the amount of memory loaded from the file at once
	linksChan := make(chan string, Config.MaxLinks)
	go processLinks(links, linksChan)

	// continuously check the linksChan for new links, then process the images
	resChan := make(chan string)
	go processImages(linksChan, resChan)
	
	// results of the image processor come back on reschan for display in the main process	
	for res := range resChan {
		outputter.OutputSingle(res)
	}		
}

// initLinks initializes the links data source from CSV or CLI 
func initLinks(input string) Links {
	var links Links

	if len(input) > 0 {
		file, err := os.Open(input + ".csv") // HACK for now because flag ignores everything after the .
		if err != nil {
			log.Fatalf("Problem reading input CSV: %s", err.Error())
			os.Exit(-1)
		}

		links = NewCsvLinks(file)
	} else {
		links = NewArrayLinks(flag.Args())
	}

	return links
}

func initOutputter(output string) Outputter {
	var outputter Outputter
	if len(output) != 0 {
		outputter = &CsvOutputter{}
		
		err := outputter.Open(output + ".csv") // HACK for some reason flag doesn't pick up anything after the .
		if err != nil {
			log.Fatalf("Problem opening CSV file for output: %s", err.Error())
			return nil
		}
	} else {
		outputter = &CmdLineOutputter{}
	}

	return outputter
}

// processLinks streams the image links from the data source to a channel
func processLinks(links Links, out chan<- string) {
	for {
		link, err := links.GetNextLink()

		// stop processing on error or EOF
		if len(link) != 0 && err == nil {
			out <- link
		} else {
		 	close(out)
		 	return
		}	
	}
}

// processImages receives links from the channel and then loads the images in individual goroutines,
// it will return the results of each calculation on resChan and close the channel when all image processes complete
// and the links data source has been exhausted
func processImages(linksChan <-chan string, resChan chan<- string) {
	numProcess := 0
	procChan := make(chan string)
	for link := range linksChan {			
		go handleImageProcess(NewTopColorsProcessor(3), link, procChan)			
		numProcess++
	}

	// need to coordinate with processImage's child goroutines so we know when to close
	for p := range procChan {
		if p != "ERR" {
			resChan <- p
		}

		numProcess--
		if numProcess <= 0 {
			close(procChan)
			break
		}
	}
	close(resChan)
}

// handleImageProcess loads and processes an individual image using the provided ImageProcessor
// and outputs the results to a channel
// TODO could abstract image loading into its own concurrent operation
func handleImageProcess(p ImageProcessor, imgLink string, out chan<- string) {
	l := CreateLoader(imgLink)
	
	img, err := l.Load()
	if err != nil {
		if Config.ShowErrors {
			out <- err.Error()
		} else {
			out <- "ERR"
		}
		return
	}

	res := p.ProcessImage(img)

	out <- fmt.Sprintf("%s;%s", imgLink, res)	
}