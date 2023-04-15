package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Address       string `json:"address"`                 // IP address or hostname,multiple address use comma split
	Username      string `json:"username"`                // Username
	Password      string `json:"password"`                // Password
	PublicKey     string `json:"publicKey,omitempty"`     // PublicKey
	PublicKeyPath string `json:"publicKeyPath,omitempty"` // PublicKeyPath
	Timeout       int    `json:"timeout,omitempty"`       // Timeout in seconds,default 10s

	//Path of the file to be uploaded,multiple file use comma split e.g. "a.txt,b.txt"
	SrcFile string `json:"srcFile"`

	//The path of the file to be uploaded after the command is executed
	TargetDir string `json:"targetDir"`
	//Command to be executed before uploading a file
	//e.g. "rm -rf {dir}/*"
	PreCmd []string `json:"preCmd,omitempty"`
	//Command to be executed after uploading a file
	//e.g. "unzip -o {dir}/{file} -d {dir}"
	PostCmd []string `json:"postCmd,omitempty"`
}

func (c *Config) Init() error {
	if err := c.validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	return nil
}

// validate validates the configuration.
func (c *Config) validate() error {
	if c.Address == "" {
		return errors.New("address is required")
	}

	if c.Username == "" {
		return errors.New("username is required")
	}

	if c.SrcFile == "" {
		return errors.New("srcFile is required")
	}

	//src files is exist
	for _, path := range strings.Split(c.SrcFile, ",") {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file %s is not exist", path)
			}
			return err
		}
	}

	if c.TargetDir == "" {
		return errors.New("targetDir is required")
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// default timeout is 10 seconds
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	return &config, nil
}
