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
	case "basic":
		steg.BasicCli()
	default:
		printHelpText()
	}
}

func printHelpText() {
	log.Print("Pass a valid option into the -program flag to use this tool. Valid options are as follows\n\tbasic\n\t\tThe most basic tool to hide data in a file.\n\tlist\n\t\tShows the output you are seeing now.")
	os.Exit(1)
}
