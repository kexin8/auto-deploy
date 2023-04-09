package sftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSftpClient(sshClient *ssh.Client) (*sftp.Client, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return client, nil
}
