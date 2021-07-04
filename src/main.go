package main

import (
	"learning_steganography/src/steg"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelpText()
	}

	switch os.Args[1] {
	case "encode":
		steg.EncodeCli()
	case "decode":
		steg.DecodeCli()
	default:
		printHelpText()
	}
}

func printHelpText() {
	log.Print("Pass a valid option into the -program flag to use this tool. Valid options are as follows\n\tencode\n\t\tEncode hidden data in another file.\n\tdecode\n\t\tDecode a hidden file.\n\tlist\n\t\tShows the output you are seeing now.")
	os.Exit(1)
}
