package config

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Credentials Credentials `yaml:"credentials"`
	Options     Options     `yaml:"options"`
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Options struct {
	BaseDomain          string              `yaml:"baseDomain"`
	PlainText           bool                `yaml:"plainText"`
	IgnoreExpired       bool                `yaml:"ignoreExpired"`
	ExportICS           bool                `yaml:"exportICS"`
	ExcludedCourses     map[string]struct{} `yaml:"excludedCourses"`
	ExcludedAssignments map[string][]string `yaml:"excludedAssignments"`
}

// Import function will read options and credentials from the
// config.yaml file. If the config file is missing, it will
// be created with default values.
func Import() (*Options, *Credentials, error) {
	configPath, err := path()
	if err != nil {
		return nil, nil, err
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, err
	}

	config, err := decodeYaml(yamlFile)
	if err != nil {
		return nil, nil, err
	}

	return &config.Options, &config.Credentials, nil
}

func decodeYaml(yamlFile []byte) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Ensure function will check for required configuration values
// that are missing. If they do, they will be requested from Stdin.
func Ensure(opts *Options, creds *Credentials) error {
	updateOpts, err := ensureOptions(opts)
	if err != nil {
		return err
	}

	updateCreds, err := ensureCredentials(creds)
	if err != nil {
		return err
	}

	if updateOpts || updateCreds {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		path := filepath.Join(configDir, "eclass-utils", "config.yaml")

		cfg := &Config{Options: *opts}
		if updateCreds {
			cfg.Credentials = *creds
		} else {
			cfg.Credentials = *newDefaultCredentials()
		}

		createConfig(path, cfg)
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
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head("https://" + baseDomain)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		fmt.Println("Invalid domain")
		return false
	}
	return true
}

func ensureCredentials(creds *Credentials) (bool, error) {
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

func ensureUsername(creds *Credentials) (bool, error) {
	if creds.Username == "" {
		err := inputStdin(&creds.Username, "Username")
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func ensurePassword(creds *Credentials) (bool, error) {
	if creds.Password == "" {
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

func path() (string, error) {
	homeConfig, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeConfig, "eclass-utils", "config.yaml")
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = createConfig(path, newDefault())
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func createConfig(configPath string, config *Config) error {
	os.MkdirAll(filepath.Dir(configPath), 0644)

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

func newDefault() *Config {
	return &Config{
		Credentials: *newDefaultCredentials(),
		Options: Options{
			BaseDomain:          "",
			PlainText:           false,
			IgnoreExpired:       true,
			ExportICS:           false,
			ExcludedCourses:     map[string]struct{}{},
			ExcludedAssignments: map[string][]string{},
		},
	}
}

func newDefaultCredentials() *Credentials {
	return &Credentials{
		Username: "",
		Password: "",
	}
}
