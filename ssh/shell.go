package sftp

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
)

// Shell opens an interactive shell on the remote host.

var (
	start = false // session start flag

	in  chan string
	out chan string
)

// NewSession creates a new session for this client.
func NewSession(c *ssh.Client) (*ssh.Session, error) {
	start = false // prevent the session from starting

	session, err := c.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	w, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	r, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	MuxShell(w, r)
	if err := session.Start("/bin/sh"); err != nil {
		return nil, err
	}

	<-out //ignore the shell output

	start = true
	return session, nil
}

// MuxShell opens an interactive shell on the remote host.
func MuxShell(w io.Writer, r io.Reader) {

	in = make(chan string, 1)
	out = make(chan string, 1)

	var wg sync.WaitGroup

	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			_, _ = w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()
	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				close(in)
				close(out)
				return
			}
			t += n

			if buf[t-2] == '$' { //assuming the $PS1 == 'sh-4.3$ '
				// 过滤掉$
				out <- string(buf[:t-2])
				t = 0
				wg.Done()
			}
		}
	}()
}

// CommandOutput runs a command on the remote host and returns its output.
func CommandOutput(cmd string) (string, error) {
	if !start {
		return "", errors.New("shell : session not started")
	}

	in <- cmd
	return <-out, nil
}

// Run runs a command on the remote host.
func Run(cmd string) error {
	_, err := CommandOutput(cmd)
	return err
}
