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
	ExportICS             bool
	IncludedCourses       []string
	ExcludedCourses       []string
	IncludedTitleKeywords map[string]string
}

type Creds struct {
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

func NewCreds(m map[string]any) (*Creds, error) {
	creds := Creds{}
	for k, v := range m {
		err := SetField(&creds, k, v)
		if err != nil {
			return nil, err
		}
	}
	return &creds, nil
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
		inputOptsStdin(opts)
	}

	return opts, nil
}

func inputOptsStdin(opts *Options) error {
	fmt.Print("Domain :")
	fmt.Scanln(&opts.BaseDomain)

	return nil
}

func GetCreds() (*Creds, error) {
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

	creds, err := NewCreds(content["creds"])
	if err != nil {
		return nil, err
	}

	if creds.Username == "" || creds.Password == "" {
		inputCredsStdin(creds)
	}

	return creds, nil
}

func inputCredsStdin(creds *Creds) error {
	fmt.Print("Username: ")
	fmt.Scanln(&creds.Username)

	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	creds.Password = string(bytePassword)

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
