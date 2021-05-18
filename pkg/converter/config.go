package converter

import (
	"flag"
	"log"
)

type Config struct {
	inputYAML  string
	outputYAML string

	schemaXLSX  string
	schemaSheet string

	inPlace bool
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.inputYAML, "input", "", "input .yaml file path")
	flag.StringVar(&c.outputYAML, "output", "", "output .yaml file path")

	flag.BoolVar(&c.inPlace, "i", false, "convert yaml file in-place")

	flag.StringVar(&c.schemaXLSX, "schemaFile", "", "schema .xlsx file")
	flag.StringVar(&c.schemaSheet, "schemaSheet", "", "schema sheet name")

	flag.Parse()
}

func (c *Config) Validate() bool {
	if c.inputYAML == "" {
		log.Println("input .yaml file is not specified")
		return false
	}

	if c.inPlace && c.outputYAML != "" {
		log.Println("in-place option and output .yaml file are mutually exclusive")
		return false
	}

	if !c.inPlace && c.outputYAML == "" {
		log.Println("output .yaml file is not specified")
		return false
	}

	if c.schemaXLSX == "" {
		log.Println("schema .xlsx file is not specified")
		return false
	}

	if c.schemaSheet == "" {
		log.Println("schema sheet name is not specified")
		return false
	}

	return true
}
