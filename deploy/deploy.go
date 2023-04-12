package deploy

import (
	"fmt"
	"github.com/fatih/color"
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
	fmt.Println("Pre command ...")
	if err := execCommands(c.sshClient, c.PreCmd...); err != nil {
		return err
	}

	//3.上传文件
	fmt.Println("Upload file to remote ...")

	//创建目标目录
	if err := sftpcli.MkdirAll(c.TargetDir); err != nil {
		return err
	}

	paths := strings.Split(c.SrcFile, ",")
	for i, path := range paths {
		err2 := c.upload(path, i, len(paths))
		if err2 != nil {
			return err2
		}
	}

	//4.执行后置命令
	fmt.Println("Post command ...")
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
	targetFile, err := c.sftpClient.Create(filepath.Join(c.TargetDir, filename))
	if err != nil {
		return err
	}
	defer targetFile.Close()

	//进度条
	bar := progressbar.NewOptions64(info.Size(),
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		//progressbar.OptionSetWidth(25),
		progressbar.OptionFullWidth(),
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
		//执行命令
		if err := execCommand(client, command); err != nil {
			// [command] command					FAILED
			fmt.Printf("%s `%s`					%s\r\n",
				color.YellowString("[command]"), command, color.RedString("FAILED"))
			return err
		}
		fmt.Printf("%s `%s`					%s\r\n",
			color.YellowString("[command]"), command, color.GreenString("OK"))
	}
	return nil
}

func execCommand(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	if err != nil {
		return err
	}

	if len(output) > 0 {
		fmt.Println(string(output))
	}

	return nil
}
