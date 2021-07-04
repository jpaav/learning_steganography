package steg

import (
	"errors"
	"flag"
	"image"
	"image/png"
	"io"
	"learning_steganography/src/steg/secrets"
	"log"
	"os"
)

type decodeProgramOptions struct {
	inputFilePath  *string
	outputFilePath *string
}

func (o decodeProgramOptions) Validate() error {
	if o.inputFilePath == nil {
		return errors.New("input file must not be nil")
	}
	if *o.inputFilePath == "" {
		return errors.New("input file must not be blank")
	}

	if o.outputFilePath == nil {
		return errors.New("output file must not be nil")
	}
	if *o.outputFilePath == "" {
		return errors.New("output file must not be blank")
	}

	return nil
}

func DecodeCli() {
	decodeFlagSet := flag.NewFlagSet("decode", flag.ExitOnError)
	options := &decodeProgramOptions{}
	options.inputFilePath = decodeFlagSet.String("input", "", "the input file to read")
	options.outputFilePath = decodeFlagSet.String("output", "", "the output file to write")

	err := decodeFlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	err = options.Validate()
	if err != nil {
		log.Fatal(err)
	}
	Decode(options)
}

func Decode(options *decodeProgramOptions) {
	// Check input file
	inputStats, err := os.Stat(*options.inputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if !inputStats.Mode().IsRegular() {
		log.Fatal(errors.New("input must be a regular file"))
	}
	if inputStats.Size() > MaxInputFileSize {
		log.Fatal(errors.New("input must not exceed 64 MB"))
	}
	// Open input file
	inputFile, err := os.Open(*options.inputFilePath)
	if inputFile != nil {
		defer inputFile.Close()
	}
	// Open output file
	outputFile, err := os.OpenFile(*options.outputFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if outputFile != nil {
		defer outputFile.Close()
	}
	// Decode!
	decodePng(inputFile, outputFile)
}

func decodePng(in io.Reader, out io.Writer) {
	inImg, err := png.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	decodeCommon(inImg, out)
}

func decodeCommon(inImg image.Image, out io.Writer) {
	var rIn, gIn, bIn uint32
	var rOut, gOut, bOut uint8
	var err error

	secretWriter := secrets.NewWriter(out, 3)

	rect := inImg.Bounds()
	log.Print("Decoding...")
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			// Get the current pixel
			rIn, gIn, bIn, _ = inImg.At(x, y).RGBA()
			// Convert each component from a 32-bit to 8-bit and mask them
			rOut = uint8(rIn>>8) & secrets.CrumbBitMask
			gOut = uint8(gIn>>8) & secrets.CrumbBitMask
			bOut = uint8(bIn>>8) & secrets.CrumbBitMask
			// Store them in the byte buffer
			err = secretWriter.WriteCrumb(rOut)
			if err != nil {
				log.Fatal(err)
			}
			err = secretWriter.WriteCrumb(gOut)
			if err != nil {
				log.Fatal(err)
			}
			err = secretWriter.WriteCrumb(bOut)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = secretWriter.Sink()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}
