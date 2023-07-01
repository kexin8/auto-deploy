package deploy

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Addr     string `json:"addr"`              // IP address or hostname, e.g. "127.0.0.1:22,127.0.0.2:22,..."
	User     string `json:"username"`          // Username
	Pass     string `json:"password"`          // Password
	PubKey   string `json:"pubKey,omitempty"`  // PublicKey
	Timeout  int    `json:"timeout,omitempty"` // Timeout, default: 5s
	SrcFiles string `json:"srcFiles"`          // Local files to be uploaded
	WorkDir  string `json:"workDir"`           // Remote working directory
	// ChangeWorkDir 命令执行之前是否将shell工作目录切换到工作目录下
	// Default: true
	ChangeWorkDir bool `json:"changeWorkDir,omitempty"`
	// PreCmds A command to run before src files is upload
	PreCmds []string `json:"preCmds"`
	// PostCmds A command to run after src files is upload
	PostCmds []string `json:"postCmds"`
}

// NewConfig returns a new Config
func NewConfig(p string) (*Config, error) {
	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	c := Config{
		ChangeWorkDir: true, // Default: true
		Timeout:       5,    // Default: 5s
	}
	if err := json.Unmarshal(file, &c); err != nil {
		return nil, err
	}

	// Validate the Config
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate validates the Config
func (c *Config) Validate() error {
	if c.Addr == "" {
		return errors.New("address can't be empty")
	}
	if c.User == "" {
		return errors.New("username can't be empty")
	}
	if c.Pass == "" && c.PubKey == "" {
		return errors.New("password and publicKey can't be empty at the same time")
	}
	if c.SrcFiles == "" {
		return errors.New("srcFiles can't be empty")
	}
	if c.WorkDir == "" {
		return errors.New("workDir can't be empty")
	}

	// check if the srcFiles exists
	for _, filepath := range strings.Split(c.SrcFiles, ",") {
		if _, err := os.Stat(filepath); err != nil {
			if os.IsNotExist(err) {
				return errors.New(filepath + " not exists")
			}
			return err
		}
	}

	return nil
}

// ExampleConfig Config
func ExampleConfig() *Config {
	return &Config{
		Addr:     "host1:port1,host2:port2,...",
		User:     "username",
		Pass:     "password",
		SrcFiles: "file1,file2,...",
		WorkDir:  "/path/to/remote/dir",
		PreCmds:  []string{"cmd1", "cmd2", "..."},
		PostCmds: []string{"cmd1", "cmd2", "..."},
	}
}

// ExampleAllConfig Config
func ExampleAllConfig() *Config {
	return &Config{
		Addr:          "host1:port1,host2:port2,...",
		User:          "username",
		Pass:          "password",
		PubKey:        "ssh public key",
		Timeout:       5,
		SrcFiles:      "file1,file2,...",
		WorkDir:       "/path/to/remote/dir",
		ChangeWorkDir: true,
		PreCmds:       []string{"cmd1", "cmd2", "..."},
		PostCmds:      []string{"cmd1", "cmd2", "..."},
	}
}
