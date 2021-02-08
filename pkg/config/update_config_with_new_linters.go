package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// UpdateConfigFileWithNewLinters adds new linters to the "linters" config in the file at the provided path
func UpdateConfigFileWithNewLinters(configFilePath string, newLinters []string) error {
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return errors.Wrap(err, "could not read config file")
	}

	var doc struct {
		Linters struct {
			Enable yaml.Node
			yaml.Node
		}
		yaml.Node
	}
	if err := yaml.Unmarshal(configData, &doc); err != nil {
		return errors.Wrapf(err, "failed to unmarshal config file %q", configFilePath)
	}

	// make a guess as to the indent size
	indentSpaces := 2
	for _, n := range doc.Content {
		indentSpaces = n.Column - 1
	}

	// create the "linters" mapping node if it doesn't exist
	linters := &doc.Linters
	if linters.IsZero() {
		linters.Node = yaml.Node{Kind: yaml.MappingNode}
		doc.Content = append(doc.Content, &doc.Linters.Node)
	}

	// create the "enable" sequence node if it doesn't exist
	enableNode := &doc.Linters.Enable
	if enableNode.IsZero() {
		linters.Content = append(linters.Content, enableNode)
		*enableNode = yaml.Node{Kind: yaml.SequenceNode}
	}

	for _, l := range newLinters {
		enableLinterNode := &yaml.Node{}
		enableLinterNode.SetString(l)
		doc.Linters.Enable.Content = append(doc.Linters.Enable.Content, enableLinterNode)
	}

	configFile, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %q for writing", configFilePath)
	}

	encoder := yaml.NewEncoder(configFile)
	encoder.SetIndent(indentSpaces)
	err = encoder.Encode(doc)
	if err == nil {
		err = encoder.Close()
	}
	if err != nil {
		err = configFile.Close()
	}
	return errors.Wrapf(err, "failed to update config file %q", configFilePath)
}
