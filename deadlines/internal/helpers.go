package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	username   string
	password   string
	baseDomain string
}

func readYaml(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func decodeYaml(yamlFile []byte) (map[string]string, error) {
	cfg := make(map[string]map[string]string, 3)
	err := yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func GetConfiguration() (map[string]string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(home + "/eclass-deadlines-py/config.yaml")
	if err != nil {
		return nil, err
	}

	return decodeYaml(yamlFile)
}
