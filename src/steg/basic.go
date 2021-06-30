package steg

import (
	"errors"
	"flag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

const MaxInputFileSize = 64000000

type BasicProgramOptions struct {
	InputFilePath  *string
	OutputFilePath *string
}

func (o BasicProgramOptions) Validate() error {
	if o.InputFilePath == nil {
		return errors.New("input file must not be nil")
	}
	if *o.InputFilePath == "" {
		return errors.New("input file must not be blank")
	}
	return nil
}

func BasicCli() {
	basicProgram := flag.NewFlagSet("basic", flag.ExitOnError)
	options := &BasicProgramOptions{}
	options.InputFilePath = basicProgram.String("input", "", "the input file to read")
	options.OutputFilePath = basicProgram.String("output", "", "the output file to write")

	err := basicProgram.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	err = options.Validate()
	if err != nil {
		log.Fatal(err)
	}
	Basic(options)
}

func Basic(options *BasicProgramOptions) {
	inputStats, err := os.Stat(*options.InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if !inputStats.Mode().IsRegular() {
		log.Fatal(errors.New("input must be a regular file"))
	}
	if inputStats.Size() > MaxInputFileSize {
		log.Fatal(errors.New("input must not exceed 64 MB"))
	}

	inputFile, err := os.Open(*options.InputFilePath)
	if inputFile != nil {
		defer inputFile.Close()
	}

	img, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	basicImg(img)
}

func basicImg(img image.Image) {
	rect := img.Bounds()
	var r, g, b, a uint32
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			r, g, b, a = img.At(x, y).RGBA()
			log.Printf("r: %d, b: %d, g: %d, a: %d", r, g, b, a)
		}
	}
}
