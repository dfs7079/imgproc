package main

import (
	//"bytes"
	//"encoding/base64"
	"fmt"
	"image"
	"io"
	"log"
	//"os"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
)

// abstract processor interface. ProcessImage() returns []byte to give flexibility with the output
type ImageProcessor interface {
	ProcessImage(i image.Image) []byte
}

// TopColorsProcessor calculates the top numColors most prevelant RGB colors in the image (ignores alpha channel)
type TopColorsProcessor struct {
	numColors int
	colorMap map[uint32]int
	outColors []uint32
}

func NewTopColorsProcessor(numColors int) *TopColorsProcessor {
	return &TopColorsProcessor{
		numColors,
		make(map[uint32]int),
		make([]uint32, numColors),
	}
}

func (t *TopColorsProcessor) insertColor(c uint32, index int) {
	if len(t.outColors) != t.numColors {
		t.outColors = make([]uint32, t.numColors)
	}
	
	// inserts the color to the list of top colors
	t.outColors = append(t.outColors, 0)
	copy(t.outColors[index+1:], t.outColors[index:])
	t.outColors[index] = c

	// trim the color list to the number of colors we're counting
	t.outColors = t.outColors[:t.numColors]
}

func (t *TopColorsProcessor) ProcessImage(img image.Image) []byte {
	// get image boundaries
	b := img.Bounds()

	// combines each color from the Image into a single uint32, using this as a hash key
	// for keeping count of the maximum color
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// truncate each channel to 8 bits
			r = r >> 8
			g = g >> 8
			b = b >> 8

			// calculate a single color value for the key, again ignoring alpha channel
			var index uint32 = (r << 24) + (g << 16) + (b << 8)

			// increment the count for this color
			val, ok := t.colorMap[index] 
			if !ok {
				t.colorMap[index] = 1
			} else {
				t.colorMap[index] = val + 1
			}
		}
	}

	// check each color and insert it into a list of maxes	
	for k, v := range t.colorMap {
		// loop thru the existing maxes list and insert where necessary
		for i := 0; i < t.numColors; i++ {
			if v > t.colorMap[t.outColors[i]] {
				t.insertColor(k, i)
				break // we only insert nax once per color to preserve the list
			}	
		}
	}

	// return a string represenation of the max colors
	out := fmt.Sprintf("0x%x", t.outColors[0])
	for i := 1; i < t.numColors; i++ {
		out = fmt.Sprintf("%s;0x%x", out, t.outColors[i])
	} 

	return []byte(out)
}


// DecodeImage takes a raw byte array and attempts to parse it into a golang Image type
func DecodeImage(r io.Reader) image.Image {
	image, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	return image
}