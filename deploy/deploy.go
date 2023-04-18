package deploy

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/kexin8/auto-deploy/log"
	lsftp "github.com/kexin8/auto-deploy/sftp"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	line     = "------------------------------------------------------------------------"
	longLine = "--------------------------< %s >--------------------------"
	sortLine = "--- %s ---"
)

// UploadFiles uploads files to the remote server
func (c *Config) UploadFiles() (err error) {
	// 记录每个address操作消耗的时间
	var (
		exectimes = make(map[string]time.Duration)
		sumtime   time.Duration
	)
	for _, address := range strings.Split(c.Address, ",") {
		start := time.Now()
		if err := c.UploadFile(address); err != nil {
			return err
		}
		exectimes[address] = time.Since(start)
		sumtime += exectimes[address]
	}
	// Summary of the deployment
	log.Info(line)
	log.Info("Summary:")
	log.Info("")

	for _, address := range strings.Split(c.Address, ",") {
		//总长度固定52
		//计算需要补充的.的个数
		dotNum := 52 - len(address)
		dotStr := ""
		for i := 0; i < dotNum; i++ {
			dotStr += "."
		}
		log.InfoF("%s %s %s [  %f s]", address, dotStr, color.GreenString("SUCCESS"), exectimes[address].Seconds())
	}
	log.Info(line)
	log.Info(color.GreenString("DEPLOY SUCCESS"))
	log.Info(line)
	log.InfoF("Total time: %f s", sumtime.Seconds())
	log.InfoF("Finished at: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Info(line)
	return
}

func (c *Config) UploadFile(address string) error {
	log.InfoF(longLine, address)
	log.InfoF("Deploying to %s ...", address)
	//log.InfoF(longLine, address)
	// 创建远程服务器连接，包含ssh和sftp
	gssh, gsftp, err := c.newRemoteClient(address)
	if err != nil {
		return err
	}
	defer gssh.Close()
	defer gsftp.Close()

	log.SuccessF("Connected to %s", address)

	// 前置命令
	log.InfoF(sortLine, "Pre command")
	if err := execCommands(gssh, c.PreCmd...); err != nil {
		return err
	}

	// 上传文件
	log.InfoF(sortLine, "Upload file to remote")
	paths := strings.Split(c.SrcFile, ",")
	for i, p := range paths {
		if err := c.upload(p, err, gsftp, i, paths); err != nil {
			return err
		}
	}

	// 后置命令
	log.InfoF(sortLine, "Post command")
	if err := execCommands(gssh, c.PostCmd...); err != nil {
		return err
	}
	return nil
}

func (c *Config) upload(p string, err error, gsftp *sftp.Client, i int, paths []string) error {
	var (
		srcFile     *os.File
		srcFileInfo os.FileInfo

		filename = filepath.Base(p)
		filesize int64
	)

	if srcFileInfo, err = os.Stat(p); err != nil {
		return err
	}
	filesize = srcFileInfo.Size()

	if srcFile, err = os.Open(p); err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标目录
	if err := gsftp.MkdirAll(c.TargetDir); err != nil {
		return err
	}
	// 创建目标文件
	targetFile, err := gsftp.Create(c.TargetDir + "/" + filename)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	//进度条
	bar := progressbar.NewOptions64(filesize,
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		//progressbar.OptionSetWidth(),
		progressbar.OptionFullWidth(),
		//progressbar.OptionSetDescription("[cyan][1/1][reset] "+filename+" "),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] %s ", i+1, len(paths), filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	//上传文件
	if _, err := io.Copy(targetFile, io.TeeReader(srcFile, bar)); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func (c *Config) newRemoteClient(address string) (*ssh.Client, *sftp.Client, error) {
	//创建ssh客户端
	sshConfig := lsftp.SSHConfig{
		Address:       address,
		Username:      c.Username,
		Password:      c.Password,
		PublicKey:     c.PublicKey,
		PublicKeyPath: c.PublicKeyPath,
		Timeout:       c.Timeout,
	}

	gssh, err := sshConfig.NewSshClient()
	if err != nil {
		return nil, nil, err
	}

	//创建sftp客户端
	gsftp, err := sftp.NewClient(gssh)
	if err != nil {
		return nil, nil, err
	}

	return gssh, gsftp, nil
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
