package steg

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"learning_steganography/src/steg/secrets"
	"log"
	"os"
)

const MaxInputFileSize = 64000000

type encodeProgramOptions struct {
	inputFilePath  *string
	outputFilePath *string
	secretFilePath *string
}

func (o encodeProgramOptions) Validate() error {
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

	if o.secretFilePath == nil {
		return errors.New("secret file must not be nil")
	}
	if *o.secretFilePath == "" {
		return errors.New("secret file must not be blank")
	}
	return nil
}

func EncodeCli() {
	encodeFlagSet := flag.NewFlagSet("encode", flag.ExitOnError)
	options := &encodeProgramOptions{}
	options.inputFilePath = encodeFlagSet.String("input", "", "the input file to read")
	options.outputFilePath = encodeFlagSet.String("output", "", "the output file to write")
	options.secretFilePath = encodeFlagSet.String("secret", "", "the secret file to hide")

	err := encodeFlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	err = options.Validate()
	if err != nil {
		log.Fatal(err)
	}
	Encode(options)
}

func Encode(options *encodeProgramOptions) {
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
	// Check secret file
	secretStats, err := os.Stat(*options.secretFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if !secretStats.Mode().IsRegular() {
		log.Fatal(errors.New("secret must be a regular file"))
	}
	log.Printf("secret file size:\t%d bytes", secretStats.Size())
	// Open secret file
	secretFile, err := os.Open(*options.secretFilePath)
	if secretFile != nil {
		defer secretFile.Close()
	}
	// Encode!
	encodePng(inputFile, outputFile, secretFile)
}

func encodePng(in io.Reader, out io.Writer, secret io.Reader) {
	inImg, err := png.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	outImg := encodeCommon(inImg, secret)

	err = png.Encode(out, outImg)
	if err != nil {
		log.Fatal(err)
	}
}

func encodeCommon(inImg image.Image, secret io.Reader) image.Image {
	var rIn, gIn, bIn, aIn uint32
	var rOut, gOut, bOut, aOut uint8
	// Buffer size is 3 so we never run out of data while writing a pixel
	secretReader := secrets.NewReader(secret, 3)

	rect := inImg.Bounds()
	outImg := image.NewNRGBA(rect)
	hasData := true
	// We can store 6 bits per pixel, and there are 8 bits per byte
	log.Printf("max encodable size:\t%d bytes", (rect.Dx()*rect.Dy()*6)/8)
	log.Print("Encoding...")
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			// Get the current pixel
			rIn, gIn, bIn, aIn = inImg.At(x, y).RGBA()
			// Convert each component from a 32-bit to 8-bit
			rOut = uint8(rIn >> 8)
			gOut = uint8(gIn >> 8)
			bOut = uint8(bIn >> 8)
			aOut = uint8(aIn >> 8)
			// Encode each component
			rOut, hasData = encodeColorComponent(rOut, secretReader, hasData)
			gOut, hasData = encodeColorComponent(gOut, secretReader, hasData)
			bOut, hasData = encodeColorComponent(bOut, secretReader, hasData)
			outImg.Set(x, y, color.NRGBA{R: rOut, G: gOut, B: bOut, A: aOut})
		}
	}
	log.Print("Done!")
	if hasData {
		log.Print("WARNING: Secret file was too big to encode in image")
	}
	return outImg
}

func encodeColorComponent(component uint8, reader *secrets.Reader, hasData bool) (uint8, bool) {
	// We can skip reading if we're already out of data
	if hasData {
		secretCrumb, err := reader.ReadCrumb()
		if err != nil {
			hasData = false
		}
		return (component & secrets.NegatedCrumbBitMask) | secretCrumb, hasData
	}
	return component, false
}
