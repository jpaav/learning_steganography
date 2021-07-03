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

type EncodeProgramOptions struct {
	InputFilePath  *string
	OutputFilePath *string
	SecretFilePath *string
}

func (o EncodeProgramOptions) Validate() error {
	if o.InputFilePath == nil {
		return errors.New("input file must not be nil")
	}
	if *o.InputFilePath == "" {
		return errors.New("input file must not be blank")
	}

	if o.OutputFilePath == nil {
		return errors.New("output file must not be nil")
	}
	if *o.OutputFilePath == "" {
		return errors.New("output file must not be blank")
	}

	if o.SecretFilePath == nil {
		return errors.New("secret file must not be nil")
	}
	if *o.SecretFilePath == "" {
		return errors.New("secret file must not be blank")
	}
	return nil
}

func EncodeCli() {
	encodeFlagSet := flag.NewFlagSet("encode", flag.ExitOnError)
	options := &EncodeProgramOptions{}
	options.InputFilePath = encodeFlagSet.String("input", "", "the input file to read")
	options.OutputFilePath = encodeFlagSet.String("output", "", "the output file to write")
	options.SecretFilePath = encodeFlagSet.String("secret", "", "the secret file to hide")

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

func Encode(options *EncodeProgramOptions) {
	// Check input file
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
	// Open input file
	inputFile, err := os.Open(*options.InputFilePath)
	if inputFile != nil {
		defer inputFile.Close()
	}
	// Open output file
	outputFile, err := os.OpenFile(*options.OutputFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if outputFile != nil {
		defer outputFile.Close()
	}
	// Check secret file
	secretStats, err := os.Stat(*options.SecretFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if !secretStats.Mode().IsRegular() {
		log.Fatal(errors.New("secret must be a regular file"))
	}
	log.Printf("secret file size:\t%d bytes", secretStats.Size())
	// Open secret file
	secretFile, err := os.Open(*options.SecretFilePath)
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
	var err error
	var secretCrumb byte
	// Buffer size is 3 so we never run out of data while writing a pixel
	secretReader := secrets.NewReader(secret, 3)

	rect := inImg.Bounds()
	outImg := image.NewNRGBA(rect)
	hasData := true
	index := 0
	pixelCount := rect.Dx() * rect.Dy()
	// We can store 6 bits per pixel, and there are 8 bits per byte
	log.Printf("max encodable size:\t%d bytes", (pixelCount*6)/8)
	log.Print("Encoding...")
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			// We can skip reading if we're already out of data
			if hasData {
				secretCrumb, err = secretReader.ReadCrumb()
				if err != nil {
					hasData = false
				}
			}

			// Get the current pixel
			rIn, gIn, bIn, aIn = inImg.At(x, y).RGBA()
			// Convert each component from a 32-bit to 8-bit
			rOut = uint8(rIn >> 8)
			gOut = uint8(gIn >> 8)
			bOut = uint8(bIn >> 8)
			aOut = uint8(aIn >> 8)
			if hasData {
				rOut = (rOut & secrets.NegatedCrumbBitMask) | secretCrumb
			}
			outImg.Set(x, y, color.NRGBA{R: rOut, G: gOut, B: bOut, A: aOut})
			index += 1
		}
	}
	log.Print("Done!")
	if hasData {
		log.Print("WARNING: Secret file was too big to encode in image")
	}
	return outImg
}
