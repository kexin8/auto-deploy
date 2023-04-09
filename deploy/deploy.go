package deploy

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"path/filepath"
)

// UploadFile uploads a file to the remote server
func (c *DeployConfig) UploadFile() (err error) {
	sftpcli := c.sftpClient
	//上传文件至远程服务器指定目录
	//2.执行前置命令
	fmt.Println("Pre command ...")
	if err := execList(c.sshClient, c.PreCmd); err != nil {
		return err
	}

	//3.上传文件
	fmt.Println("Upload file to remote ...")
	filename := filepath.Base(c.SrcFile)
	if err := sftpcli.MkdirAll(c.TargetDir); err != nil {
		return err
	}
	targetfile, err := sftpcli.Create(path.Join(c.TargetDir, filename))
	if err != nil {
		return err
	}

	srcfile, err := os.Open(c.SrcFile)
	if err != nil {
		return err
	}

	srcfileinfo, err := srcfile.Stat()
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(srcfileinfo.Size(),
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		//progressbar.OptionSetWidth(25),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan][1/1][reset] "+filename+" "),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	_, err = targetfile.ReadFrom(io.TeeReader(srcfile, bar))
	if err != nil {
		return err
	}

	fmt.Println()
	//4.执行后置命令
	fmt.Println("Post command ...")
	if err := execList(c.sshClient, c.PostCmd); err != nil {
		return err
	}

	return
}

func execList(client *ssh.Client, cmd []string) error {
	for _, command := range cmd {
		//执行命令
		if err := exec(client, command); err != nil {
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

func exec(client *ssh.Client, cmd string) error {
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
