package internal

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Options struct {
	BaseDomain      string
	PlainText       bool
	ExportICS       bool
	IncludedCourses []string
	ExcludedCourses []string
}

type Creds struct {
	Username string
	Password string
}

func NewOptions(m map[string]any) (*Options, error) {
	opts := Options{}
	for k, v := range m {
		err := SetField(&opts, k, v)
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
		return errors.New("provided value type didn't match obj field type")
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
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(home + "/eclass-deadlines-py/config.yaml")
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
	return opts, nil
}

func GetCreds() (*Creds, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	yamlFile, err := readYaml(home + "/eclass-deadlines-py/config.yaml")
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
	return creds, nil
}
