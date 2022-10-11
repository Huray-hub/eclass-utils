package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func readYaml(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func decodeYaml(yamlFile []byte) (map[string]string, error) {
	creds := make(map[string]string, 3)
	err := yaml.Unmarshal(yamlFile, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func GetConfiguration() (map[string]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(home + "/.config/eclass-deadlines-py/config.yaml")
	if err != nil {
		return nil, err
	}

	return decodeYaml(yamlFile)
}
