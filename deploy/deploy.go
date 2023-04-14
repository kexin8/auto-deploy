package deploy

import (
	"fmt"
	"github.com/kexin8/auto-deploy/log"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UploadFile uploads a file to the remote server
func (c *Config) UploadFile() (err error) {
	sftpcli := c.sftpClient
	//上传文件至远程服务器指定目录
	//2.执行前置命令
	log.Info("Pre command ...")
	if err := execCommands(c.sshClient, c.PreCmd...); err != nil {
		return err
	}

	//3.上传文件
	log.Info("Upload file to remote ...")

	//创建目标目录
	if err := sftpcli.MkdirAll(c.TargetDir); err != nil {
		return err
	}

	paths := strings.Split(c.SrcFile, ",")
	for i, path := range paths {
		if err := c.upload(path, i, len(paths)); err != nil {
			return err
		}
	}

	//4.执行后置命令
	log.Info("Post command ...")
	if err := execCommands(c.sshClient, c.PostCmd...); err != nil {
		return err
	}

	return
}

func (c *Config) upload(path string, number, total int) error {
	//获取源文件信息
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	filename := filepath.Base(path)

	//创建目标文件
	targetFile, err := c.sftpClient.Create(c.TargetDir + "/" + filename)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	//进度条
	bar := progressbar.NewOptions64(info.Size(),
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(25),
		//progressbar.OptionFullWidth(),
		//progressbar.OptionSetDescription("[cyan][1/1][reset] "+filename+" "),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] %s ", number+1, total, filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	//上传文件
	if _, err := io.Copy(targetFile, io.TeeReader(file, bar)); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func execCommands(client *ssh.Client, cmd ...string) error {
	for _, command := range cmd {
		log.InfoShell(command)
		//执行命令
		if err := execCommand(client, command); err != nil {
			return err
		}
	}
	return nil
}

func execCommand(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 如何解决无法使用管道的问题
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Run(cmd); err != nil {
		return err
	}

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Error(string(output))
		return err
	}

	if len(output) > 0 {
		log.Info(string(output))
	}

	return nil
}
