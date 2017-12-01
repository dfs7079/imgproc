package main

import (
	"flag"
	"fmt"
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

func main() {
	l := NewFileLoader("foo.png")

	p := NewTopColorsProcessor(3)

	img := DecodeImage(l.GetReader())

	out := p.ProcessImage(img)

	fmt.Printf("foo.png:%s\n", out)
}