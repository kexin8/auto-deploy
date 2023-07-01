package test

import (
	gssh "github.com/kexin8/auto-deploy/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"testing"
)

const (
	address   = ""
	username  = ""
	publicKey = ""
)

func TestPubKeyLogin(t *testing.T) {
	//t.Log(os.Getenv("SSH_AUTH_SOCK"))

	key, err := os.ReadFile(publicKey)
	if err != nil {
		t.Error(err)
	}

	signers, err := ssh.ParsePrivateKey(key)
	if err != nil {
		t.Error(err)
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(signers)}
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		t.Error(err)
	}

	defer session.Close()

	output, err := session.Output("ls -l")
	if err != nil {
		t.Log(string(output))
		t.Error(err)
	}
	t.Log(string(output))
}

func TestPipeline(t *testing.T) {

	config := gssh.Config{
		Addr: "127.0.0.1:2222",
		User: "zhengzi",
		Pass: "zxy.1314",
	}

	client, err := gssh.NewClient(&config)
	if err != nil {
		t.Error(err)
	}

	session, err := gssh.NewSession(client)
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	var commands = []string{
		"pwd",
		"cd /home",
		"ls -l",
		"pwd",
		"exit",
	}

	for _, command := range commands {
		t.Log("$ ", command)

		output, err := gssh.CommandOutput(command)
		if err != nil {
			t.Error(err)
		}

		if len(output) > 0 {
			t.Log(output)
		}
	}

	if err := session.Wait(); err != nil {
		t.Log(err)
	}

}

func TestTerminal(t *testing.T) {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		t.Error(err)
	}

	t.Log(width, height)
}
