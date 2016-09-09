package infrastructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration holds the values necessary to configure the application
type Configuration struct {
	Port         string   `yaml:"port"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	RedirectURI  string   `yaml:"redirectURI"`
	Scopes       []string `yaml:"scopes,flow"`
}

// GetConfiguration reads the file with the configuration and returns an struct
// with the fields
func GetConfiguration(path string) (*Configuration, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil

}
