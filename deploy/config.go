package deploy

import (
	"encoding/json"
	"errors"
	"github.com/kexin8/auto-deploy/log"
	"os"
	"strings"
)

type Config struct {
	Addr      string `json:"addr"`                // IP address or hostname, e.g. "127.0.0.1:22,127.0.0.2:22,..."
	User      string `json:"username"`            // Username
	Pass      string `json:"password"`            // Password
	PublicKey string `json:"publicKey,omitempty"` // PublicKey
	Timeout   int    `json:"timeout,omitempty"`   // Timeout, default: 5s
	SrcFile   string `json:"srcFile"`             // Local files to be uploaded
	WorkDir   string `json:"workDir"`             // Remote working directory
	// ChangeWorkDir 命令执行之前是否将shell工作目录切换到工作目录下
	// Default: true
	ChangeWorkDir bool `json:"changeWorkDir,omitempty"`
	// PreCmd A command to run before src files is upload
	PreCmd []string `json:"preCmd"`
	// PostCmd A command to run after src files is upload
	PostCmd []string `json:"postCmd"`

	////////////////////// Deprecated Fields //////////////////////
	Address       string `json:"address,omitempty"`       // Address is deprecated, use Addr instead
	PublicKeyPath string `json:"publicKeyPath,omitempty"` // PublicKeyPath is deprecated, use PublicKey instead
	TargetDir     string `json:"targetDir,omitempty"`     // TargetDir is deprecated, use WorkDir instead
}

// Deprecated is used to convert the deprecated fields to the new fields
// 兼容v1.2.0之前的配置文件
func (c *Config) Deprecated() {
	if c.Address != "" {
		log.Warn("dyconfig: `address` is deprecated, please use `addr`")
		c.Addr = c.Address
	}
	if c.PublicKeyPath != "" {
		log.Warn("dyconfig: `publicKeyPath` is deprecated, please use `publicKey`")
		c.PublicKey = c.PublicKeyPath
	}

	if c.TargetDir != "" {
		log.Warn("dyconfig: `targetDir` is deprecated, please use `workDir`")
		c.WorkDir = c.TargetDir
	}
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

	c.Deprecated() // 兼容v1.2.0之前的配置文件
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
	if c.Pass == "" && c.PublicKey == "" {
		return errors.New("password and publicKey can't be empty at the same time")
	}
	if c.SrcFile == "" {
		return errors.New("srcFiles can't be empty")
	}
	if c.WorkDir == "" {
		return errors.New("workDir can't be empty")
	}

	// check if the srcFiles exists
	for _, filepath := range strings.Split(c.SrcFile, ",") {
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
		Addr:    "host1:port1,host2:port2,...",
		User:    "username",
		Pass:    "password",
		SrcFile: "file1,file2,...",
		WorkDir: "/path/to/remote/dir",
		PreCmd:  []string{"cmd1", "cmd2", "..."},
		PostCmd: []string{"cmd1", "cmd2", "..."},
	}
}

// ExampleAllConfig Config
func ExampleAllConfig() *Config {
	return &Config{
		Addr:          "host1:port1,host2:port2,...",
		User:          "username",
		Pass:          "password",
		PublicKey:     "ssh public key",
		Timeout:       5,
		SrcFile:       "file1,file2,...",
		WorkDir:       "/path/to/remote/dir",
		ChangeWorkDir: true,
		PreCmd:        []string{"cmd1", "cmd2", "..."},
		PostCmd:       []string{"cmd1", "cmd2", "..."},
	}
}
