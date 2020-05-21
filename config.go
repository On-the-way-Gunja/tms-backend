package main

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfig(path string, validator func(i interface{}) error) (*ConfigFormat, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := ConfigFormat{}
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if validator != nil {
		if err := validator(config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}
