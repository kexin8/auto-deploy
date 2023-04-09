package sftp

import (
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHConfig struct {
	Address  string `json:"address"`  // IP address or hostname
	Username string `json:"username"` // Username
	Password string `json:"password"` // Password
	Timeout  int    `json:"timeout"`  // Timeout in seconds
}

// NewSshClient initializes the configuration
func (c SSHConfig) NewSshClient() (client *ssh.Client, err error) {
	client, err = ssh.Dial("tcp", c.Address, &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		Timeout:         time.Duration(c.Timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	return
}
