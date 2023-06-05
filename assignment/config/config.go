package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Credentials auth.Credentials `yaml:"credentials"`
	Options     Options          `yaml:"options"`
}

type Options struct {
	PlainText           bool                `yaml:"plainText"`
	IncludeExpired      bool                `yaml:"includeExpired"`
	ExportICS           bool                `yaml:"exportICS"`
	ExcludedAssignments map[string][]string `yaml:"excludedAssignments"`
	course.Options      `yaml:",inline"`
}

// Import function will read options and credentials from the
// given path of a yaml file. If the config file is missing, it will
// be created with default values.
func Import(configPath string) (*Config, error) {
	err := ensurePath(configPath)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Import function will read options and credentials from the
// given path of a yaml file. If the config file is missing, it will
// be created with default values.
func ImportDefault() (*Config, error) {
	configPath, err := defaultPath()
	if err != nil {
		return nil, err
	}

	return Import(configPath)
}

// Ensure function will check for required configuration values
// that are missing. If they do, they will be requested from Stdin.
func Ensure(cfg *Config) error {
	updateOpts, err := ensureOptions(&cfg.Options)
	if err != nil {
		return err
	}

	updateCreds, err := ensureCredentials(&cfg.Credentials)
	if err != nil {
		return err
	}

	if updateOpts || updateCreds {
		path, err := defaultPath()
		if err != nil {
			return err
		}

		err = Export(path, *cfg, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// Export function creates a config.yaml file at the specified path using the config
// struct provided.
func Export(configPath string, config Config, parents bool) error {
	if parents {
		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return err
		}

	}

	yamlFile, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, yamlFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ensureOptions(opts *Options) (bool, error) {
	updateDomain := false
	for opts.BaseDomain == "" || !isValidDomain(opts.BaseDomain) {
		err := inputStdin(&opts.BaseDomain, "Domain")
		if err != nil {
			return false, err
		}
		updateDomain = true
	}

	return updateDomain, nil
}

func isValidDomain(baseDomain string) bool {
	if !strings.Contains(baseDomain, ".gr") || !strings.Contains(baseDomain, "eclass") {
		fmt.Println("Invalid domain. Try eclass.<yourcollege>.gr")
		return false
	}
	return true
}

func ensureCredentials(creds *auth.Credentials) (bool, error) {
	updateUsername, err := ensureUsername(creds)
	if err != nil {
		return false, err
	}

	updatePassword, err := ensurePassword(creds)
	if err != nil {
		return false, err
	}

	if updateUsername || updatePassword {
		var decision string
		err := inputStdin(&decision, "Store credentials in config file? y/N")
		if err != nil {
			return false, err
		}

		if decision == "yes" || decision == "y" {
			return true, err
		}
	}

	return false, nil
}

func ensureUsername(creds *auth.Credentials) (bool, error) {
	if creds.UsernameEmpty() {
		err := inputStdin(&creds.Username, "Username")
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func ensurePassword(creds *auth.Credentials) (bool, error) {
	if creds.PasswordEmpty() {
		err := inputPasswordStdin(&creds.Password)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func inputStdin(value *string, message string) error {
	fmt.Print(message + ": ")
	_, err := fmt.Scanln(value)
	if err != nil && err.Error() != "unexpected newline" {
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
	fmt.Println()
	return nil
}

func defaultPath() (string, error) {
	homeConfig, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeConfig, "eclass-utils", "config.yaml"), nil
}

func ensurePath(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return Export(path, newDefault(), true)
	}

	return nil
}

func newDefault() Config {
	return Config{
		Credentials: auth.Credentials{
			Username: "",
			Password: "",
		},
		Options: Options{
			PlainText:           false,
			IncludeExpired:      false,
			ExportICS:           false,
			ExcludedAssignments: map[string][]string{},
			Options: course.Options{
				ExcludedCourses: map[string]struct{}{},
			},
		},
	}
}
