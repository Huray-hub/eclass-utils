package internal

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Options struct {
	BaseDomain            string
	PlainText             bool
	IgnoreExpired         bool
	ExportICS             bool
	IncludedCourses       []string
	ExcludedCourses       []string
	IncludedTitleKeywords map[string]string
}

type Credentials struct {
	Username string
	Password string
}

func NewOptions(m map[string]any) (*Options, error) {
	opts := Options{}
	var err error
	for k, v := range m {
		err = SetField(&opts, k, v)
		if err != nil {
			return nil, err
		}
	}
	return &opts, nil
}

func NewCredentials(m map[string]any) (*Credentials, error) {
	credentials := Credentials{}
	for k, v := range m {
		err := SetField(&credentials, k, v)
		if err != nil {
			return nil, err
		}
	}
	return &credentials, nil
}

func SetField(obj any, name string, value any) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return fmt.Errorf("provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func readYaml(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func decodeYaml(yamlFile []byte) (map[string]map[string]any, error) {
	var content map[string]map[string]any

	err := yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetOptions() (*Options, error) {
	cfgPath, err := configPath()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(cfgPath)
	if err != nil {
		return nil, err
	}

	content, err := decodeYaml(yamlFile)
	if err != nil {
		return nil, err
	}

	opts, err := NewOptions(content["options"])
	if err != nil {
		return nil, err
	}

	if opts.BaseDomain == "" {
		err := inputOptsStdin(opts)
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func inputOptsStdin(opts *Options) error {
	fmt.Print("Domain :")
	_, err := fmt.Scanln(&opts.BaseDomain)
	if err != nil {
		return err
	}

	return nil
}

func GetCredentials() (*Credentials, error) {
	cfgPath, err := configPath()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(cfgPath)
	if err != nil {
		return nil, err
	}

	content, err := decodeYaml(yamlFile)
	if err != nil {
		return nil, err
	}

	credentials, err := NewCredentials(content["credentials"])
	if err != nil {
		return nil, err
	}

	if credentials.Username == "" || credentials.Password == "" {
		err := inputCredentialsStdin(credentials)
		if err != nil {
			return nil, err
		}
	}

	return credentials, nil
}

func inputCredentialsStdin(credentials *Credentials) error {
	fmt.Print("Username: ")
	_, err := fmt.Scanln(&credentials.Username)
	if err != nil {
		return err
	}

	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	credentials.Password = string(bytePassword)

	return nil
}

func configPath() (string, error) {
	homeConfig, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	cfgPath := homeConfig + "/eclass-utils/config.yaml"
	if _, err = os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		cfgPath = "../config/default-config.yaml"
	}

	return cfgPath, nil
}
