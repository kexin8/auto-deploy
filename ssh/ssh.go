package sftp

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type Config struct {
	Addr    string        // IP address or hostname, e.g. "127.0.0.1:22"
	User    string        // Username
	Pass    string        // Password
	PubKey  string        // PublicKey
	Timeout time.Duration // Timeout
}

func (c Config) Validate() error {
	if c.Addr == "" {
		return errors.New("address cannot be empty")
	}

	if c.User == "" {
		return errors.New("username cannot be empty")
	}

	if c.Pass == "" && c.PubKey == "" {
		return errors.New("password and public key cannot be empty")
	}
	return nil
}

func NewClient(c *Config) (client *ssh.Client, err error) {

	if err := c.Validate(); err != nil {
		return nil, err
	}

	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}

	var auth = make([]ssh.AuthMethod, 0)
	if c.Pass != "" {
		auth = append(auth, ssh.Password(c.Pass))
	}
	if c.PubKey != "" {
		var publicKey []byte
		if _, err := os.Stat(c.PubKey); err == nil {
			// is a public key file
			publicKey, err = os.ReadFile(c.PubKey)
			if err != nil {
				return nil, err
			}
		} else {
			// is a public key string
			publicKey = []byte(c.PubKey)
		}

		key, err := ssh.ParsePrivateKey(publicKey)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(key))
	}

	config := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            c.User,
		Auth:            auth,
		Timeout:         c.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err = ssh.Dial("tcp", c.Addr, config)
	return
}
