package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func readYaml(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	return file, nil
}

func decodeYaml(yamlFile []byte) map[string]string {
	creds := make(map[string]string, 3)
	err := yaml.Unmarshal(yamlFile, &creds)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return creds
}

func GetConfiguration() (map[string]string, error) {
	yamlFile, err := readYaml("/home/pskiadas/.config/eclass-deadlines-py/config.yaml")
	if err != nil {
		return nil, err
	}
	return decodeYaml(yamlFile), nil
}
