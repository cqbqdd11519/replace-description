package converter

import (
	"github.com/ghodss/yaml"
	"github.com/tmax-cloud/replace-description/pkg/apis"
	"io/ioutil"
)

func Unmarshal(filePath string) (*apis.CustomResourceDefinition, error) {
	obj := &apis.CustomResourceDefinition{}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func Marshal(crd *apis.CustomResourceDefinition, outputFilePath string) error {
	b, err := yaml.Marshal(crd)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(outputFilePath, b, 0755); err != nil {
		return err
	}
	return nil
}
