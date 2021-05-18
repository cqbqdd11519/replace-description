package main

import (
	"github.com/tmax-cloud/replace-description/pkg/converter"
	"log"
	"os"
)

func main() {
	config := &converter.Config{}
	config.ParseFlags()

	// Validate
	if !config.Validate() {
		log.Println("flags are not valid")
		os.Exit(1)
	}

	// Open Schema
	c, err := converter.New(config)
	handleError(err)

	// Convert
	handleError(c.Convert())

	// Print
	handleError(c.Write())
}

func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
