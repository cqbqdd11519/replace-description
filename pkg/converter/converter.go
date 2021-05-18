package converter

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"github.com/tmax-cloud/replace-description/pkg/apis"
	"os"
	"regexp"
	"strings"
)

const (
	colNumKey = iota
	colNumEng
	colNumKor
)

var jsonPathReg = regexp.MustCompile(`(%?.*(\.yaml|\.json))(.*)`)

type Converter interface {
	Convert() error
	Write() error
}

type converter struct {
	cfg *Config

	sheet *xlsx.Sheet
	crd   *apis.CustomResourceDefinition
}

func New(cfg *Config) (*converter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	converter := &converter{cfg: cfg}

	// Open CRD file
	var err error
	converter.crd, err = Unmarshal(cfg.inputYAML)
	if err != nil {
		return nil, err
	}

	// Open schema file
	wb, err := xlsx.OpenFile(cfg.schemaXLSX)
	if err != nil {
		return nil, err
	}

	var exist bool
	converter.sheet, exist = wb.Sheet[cfg.schemaSheet]
	if !exist {
		return nil, fmt.Errorf("no sheet %s is found", cfg.schemaSheet)
	}

	return converter, nil
}

func (c *converter) Convert() error {
	maxRow := c.sheet.MaxRow
	for i := 0; i < maxRow; i++ {
		key, err := c.sheet.Cell(i, colNumKey)
		if err != nil {
			return err
		}
		m := jsonPathReg.FindAllStringSubmatch(key.Value, -1)
		if len(m) != 1 || len(m[0]) != 4 {
			return fmt.Errorf("key %s is not in form of %s", key.Value, jsonPathReg.String())
		}
		jsonPath := m[0][3]

		eng, err := c.sheet.Cell(i, colNumEng)
		if err != nil {
			return err
		}
		kor, err := c.sheet.Cell(i, colNumKor)
		if err != nil {
			return err
		}
		_ = eng.Value

		effPath := strings.TrimPrefix(jsonPath, ".spec.versions.schema.openAPIV3Schema")

		setValue(c.crd.Spec.Versions[0].Schema.OpenAPIV3Schema, strings.Split(effPath, "."), kor.Value)
	}
	return nil
}

func setValue(parent *apis.JSONSchemaProps, path []string, val string) {
	if len(path) <= 1 {
		parent.Description = val
		return
	}

	curPath := path[1]
	if curPath == "properties" {
		if len(path) == 3 {
			parent.Properties[path[2]].Description = val
		} else {
			setValue(parent.Properties[path[2]], path[2:], val)
		}
	} else if curPath == "items" {
		if len(path) == 2 {
			parent.Items.Schema.Description = val
		} else {
			setValue(parent.Items.Schema, path[1:], val)
		}
	}
}

func (c *converter) Write() error {
	// Print crd
	if c.cfg.inPlace {
		if err := Marshal(c.crd, c.cfg.inputYAML); err != nil {
			return err
		}
	} else {
		_, err := os.Stat(c.cfg.outputYAML)
		if err == nil {
			return fmt.Errorf("file %s already exists", c.cfg.outputYAML)
		}
		if err := Marshal(c.crd, c.cfg.outputYAML); err != nil {
			return err
		}
	}
	return nil
}
