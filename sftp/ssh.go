package sftp

import (
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type SSHConfig struct {
	Address       string `json:"address"`                 // IP address or hostname
	Username      string `json:"username"`                // Username
	Password      string `json:"password"`                // Password
	PublicKey     string `json:"publicKey,omitempty"`     // PublicKey
	PublicKeyPath string `json:"publicKeyPath,omitempty"` // PublicKeyPath
	Timeout       int    `json:"timeout"`                 // Timeout in seconds
}

// NewSshClient initializes the configuration
func (c SSHConfig) NewSshClient() (client *ssh.Client, err error) {

	auth := make([]ssh.AuthMethod, 0)

	auth = append(auth, ssh.Password(c.Password))

	signers := make([]ssh.Signer, 0)
	if c.PublicKey != "" {
		key, err := ssh.ParsePrivateKey([]byte(c.PublicKey))
		if err != nil {
			return nil, err
		}
		signers = append(signers, key)
	}

	if c.PublicKeyPath != "" {
		filebytes, err := os.ReadFile(c.PublicKeyPath)
		if err != nil {
			return nil, err
		}

		key, err := ssh.ParsePrivateKey(filebytes)
		if err != nil {
			return nil, err
		}
		signers = append(signers, key)
	}

	if signers != nil && len(signers) > 0 {
		auth = append(auth, ssh.PublicKeys(signers...))
	}

	client, err = ssh.Dial("tcp", c.Address, &ssh.ClientConfig{
		User:            c.Username,
		Auth:            auth,
		Timeout:         time.Duration(c.Timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	return
}
