package utils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go_blueprints/chapter_1/chat/trace"
)

// Configuration is a simple struct that contains a Tracer to allow the display
// of debug
type Configuration struct {
	Tracer trace.Tracer
}

// Provider is a structure containing an appId, a secret and a callback
// AppID and secret are given by your developer console
// callback is the link to which we go back once the user is identified
type Provider struct {
	Name     string `json:"name"`
	AppID    string `json:"id"`
	Secret   string `json:"secret"`
	Callback string `json:"callback"`
}

// NewConfiguration returns a Configuration object with a nilTracer by default
func NewConfiguration() *Configuration {
	return &Configuration{
		Tracer: trace.Off(),
	}
}

// readJSON take a slice of providers in parameter and returns a filepath.WalkFunc
// function that will be called for each file inside a specific folder
func (c *Configuration) readJSON(providers *[]Provider) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		c.Tracer.Trace("Decoding file: " + path)
		provider := Provider{}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&provider)
		if err != nil {
			return err
		}
		*providers = append(*providers, provider)
		return nil
	}
}

// GetProviderConfig searches for the configuration in /config/ and sets the params
// for the provider
// This is more a proof of concept that needs a little bit more refinement
// For the moment I only have one provider, however in a near future I could add
// a few more simply by adding a few things in the function readJSON
func (c *Configuration) GetProviderConfig() (*[]Provider, error) {
	providers := []Provider{}
	if err := filepath.Walk("./config", c.readJSON(&providers)); err != nil {
		return nil, err
	}
	return &providers, nil
}
