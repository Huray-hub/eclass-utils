package config

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Credentials Credentials `yaml:"credentials"`
	Options     Options     `yaml:"options"`
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password,omitempty"`
}

type Options struct {
	BaseDomain      string   `yaml:"baseDomain,omitempty"`
	PlainText       bool     `yaml:"plainText"`
	IgnoreExpired   bool     `yaml:"ignoreExpired"`
	ExportICS       bool     `yaml:"exportICS"`
	ExcludedCourses []string `yaml:"excludedCourses,omitempty"`

	ExcludedAssignmentsByKeyword map[string][]string `yaml:"excludedAssignmentsByKeyword,omitempty"`
	// ExcludedAssignmentsByKeyword []struct {
	//        []
	// 	CourseID string   `yaml:"courseID"`
	// 	Keywords []string `yaml:"keywords"`
	// } `yaml:"excludedAssignmentsByKeyword,omitempty"`
}

func readYaml(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func decodeYaml(yamlFile []byte) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Import() (*Options, *Credentials, error) {
	configPath, err := path()
	if err != nil {
		return nil, nil, err
	}

	yamlFile, err := readYaml(configPath)
	if err != nil {
		return nil, nil, err
	}

	config, err := decodeYaml(yamlFile)
	if err != nil {
		return nil, nil, err
	}

	opts, err := extractOptions(config)
	if err != nil {
		return nil, nil, err
	}

	creds, err := extractCredentials(config)
	if err != nil {
		return nil, nil, err
	}

	return opts, creds, nil
}

func extractOptions(config *Config) (*Options, error) {
	opts := &config.Options
	if opts.BaseDomain == "" {
		err := inputStdin(&opts.BaseDomain, "Domain of the university :")
		// err := inputOptsStdin(opts)
		if err != nil {
			return nil, err
		}
	}
	return opts, nil
}

func extractCredentials(config *Config) (*Credentials, error) {
	creds := &config.Credentials
	if creds.Username == "" {
		err := inputStdin(&creds.Username, "Username :")
		if err != nil {
			return nil, err
		}
	}
	if creds.Password == "" {
		err := inputPasswordStdin(&creds.Password)
		if err != nil {
			return nil, err
		}
	}

	return creds, nil
}

func inputStdin(value *string, message string) error {
	fmt.Print(message + ": ")
	_, err := fmt.Scanln(value)
	if err != nil {
		return err
	}

	return nil
}

func inputPasswordStdin(password *string) error {
	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	*password = string(bytePassword)
	return nil
}

func path() (string, error) {
	homeConfig, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := homeConfig + "/eclass-utils/config.yaml"
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = "../config/default-config.yaml"
	}

	return path, nil
}
