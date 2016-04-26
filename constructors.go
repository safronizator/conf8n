package conf8n

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	JSON = "json"
	YAML = "yaml"
)

// Creates Config instance from YAML-encoded data
func NewConfigFromYaml(data []byte) (*Config, error) {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	return NewConfig(m), nil
}

// Creates Config instance from JSON-encoded data
func NewConfigFromJson(data []byte) (*Config, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return NewConfig(m), nil
}

// Creates Config instance from data in file.
// Data encoding will be defined from file extension (".json" & ".yaml" supported for the moment)
func NewConfigFromFile(filename string) (*Config, error) {
	var f io.Reader
	var err error
	if f, err = os.Open(filename); err != nil {
		return nil, err
	}
	ext := strings.TrimLeft(strings.ToLower(filepath.Ext(filename)), ".")
	return NewConfigFromReader(f, ext)
}

// Creates Config instance with data from io.Reader. Specifying of incoming data format is required
func NewConfigFromReader(r io.Reader, format string) (*Config, error) {
	var data []byte
	var err error
	if data, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}
	switch format {
	case JSON:
		return NewConfigFromJson(data)
	case YAML:
		return NewConfigFromYaml(data)
	default:
		return nil, fmt.Errorf("Unknown config format: '%s'", format)
	}
}
