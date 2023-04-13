package test

import (
	"github.com/kexin8/auto-deploy/deploy"
	lsftp "github.com/kexin8/auto-deploy/sftp"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"testing"
)

func TestSftp(t *testing.T) {
	config := lsftp.SSHConfig{
		Address:  "",
		Username: "",
		Password: "",
		Timeout:  5,
	}

	sshcli, err := config.NewSshClient()
	if err != nil {
		t.Error(err)
	}

	sftpcli, _ := sftp.NewClient(sshcli,
		sftp.UseConcurrentReads(true),
		sftp.UseConcurrentWrites(true),
	)

	//upload file 'text.txt' to '/home/lifucang'
	targetfile, _ := sftpcli.Create("/home/text.jtl")
	defer targetfile.Close()

	srcfile := "D:\\download\\upgrade2500.jtl"
	file, _ := os.Open(srcfile)
	defer file.Close()

	//bufio.NewWriter(file).ReadFrom(targetfile)

	//targetfile.ReadFrom(file)

	//get file size
	fileinfo, _ := file.Stat()
	// create progress bar
	bar := progressbar.DefaultBytes(
		fileinfo.Size(), //total size, -1 means unknown
		"uploading",
	)
	// 上传文件并显示进度条
	//_, err = io.Copy(targetfile, io.TeeReader(file, bar))

	targetfile.ReadFrom(io.TeeReader(file, bar))

}

func TestSSh(t *testing.T) {
	config := lsftp.SSHConfig{
		Address:  "",
		Username: "",
		Password: "",
		Timeout:  5,
	}

	client, err := config.NewSshClient()
	if err != nil {
		t.Error(err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Error(err)
	}

	output, err := session.Output("cat hello.txt")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(output))
}

func TestPipeline(t *testing.T) {
	c, err := deploy.ReadConfig("../dyconfig.json")
	if err != nil {
		t.Error(err)
	}

	sshConf := lsftp.SSHConfig{
		Address:       c.Address,
		Username:      c.Username,
		Password:      c.Password,
		PublicKey:     c.PublicKey,
		PublicKeyPath: c.PublicKeyPath,
		Timeout:       c.Timeout,
	}
	gossh, err := sshConf.NewSshClient()
	if err != nil {
		t.Error(err)
	}

	session, _ := gossh.NewSession()
	defer session.Close()

	// 执行管道命令
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stderr
	//session.Stdin = os.Stdin

	//_, err = session.StdoutPipe()
	//if err != nil {
	//	t.Error(err)
	//}

	//err = session.Run("ls -l |grep deploy")
	//if err != nil {
	//	t.Error(err)
	//}

	output, err := session.CombinedOutput("ps -ef |grep tb_tools |grep -v grep |awk '{print $2}' |xargs kill -9")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(output))

}
