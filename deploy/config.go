package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"

	lsftp "github.com/kexin8/auto-deploy/sftp"
)

type Config struct {
	Address       string `json:"address"`                 // IP address or hostname
	Username      string `json:"username"`                // Username
	Password      string `json:"password"`                // Password
	PublicKey     string `json:"publicKey,omitempty"`     // PublicKey
	PublicKeyPath string `json:"publicKeyPath,omitempty"` // PublicKeyPath
	Timeout       int    `json:"timeout,omitempty"`       // Timeout in seconds,default 10s

	//Path of the file to be uploaded
	SrcFile string `json:"srcFile"`

	//deploy directory
	TargetDir string `json:"targetDir"`
	//Command to be executed before uploading a file
	//e.g. "rm -rf {dir}/*"
	PreCmd []string `json:"preCmd,omitempty"`
	//Command to be executed after uploading a file
	//e.g. "unzip -o {dir}/{file} -d {dir}"
	PostCmd []string `json:"postCmd,omitempty"`

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func (c *Config) Init() error {
	if err := c.validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	sshConf := lsftp.SSHConfig{
		Address:       c.Address,
		Username:      c.Username,
		Password:      c.Password,
		PublicKey:     c.PublicKey,
		PublicKeyPath: c.PublicKeyPath,
		Timeout:       c.Timeout,
	}
	sshClient, err := sshConf.NewSshClient()
	if err != nil {
		return err
	}

	sftpClient, err := lsftp.NewSftpClient(sshClient)
	if err != nil {
		return err
	}

	c.sshClient = sshClient
	c.sftpClient = sftpClient

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

	//src file is exist
	if _, err := os.Stat(c.SrcFile); err != nil {
		if os.IsNotExist(err) {
			return errors.New(c.SrcFile + " is not exist")
		}

		return err
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
