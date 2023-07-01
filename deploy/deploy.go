package deploy

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/kexin8/auto-deploy/log"
	gssh "github.com/kexin8/auto-deploy/ssh"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
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

func Init(c *Config, addr string) (sshClient *ssh.Client, sftpClient *sftp.Client, err error) {
	config := &gssh.Config{
		Addr:    addr,
		User:    c.User,
		Pass:    c.Pass,
		PubKey:  c.PublicKey,
		Timeout: time.Duration(c.Timeout) * time.Second,
	}
	// 初始化sshClient和sftpClient
	sshClient, err = gssh.NewClient(config)
	if err != nil {
		return
	}

	sftpClient, err = sftp.NewClient(sshClient)

	return
}

// Deploy the specified files to the remote server
func Deploy(c *Config) (err error) {
	// 记录每个服务器发布&&部署消耗的时间
	var (
		deploytimes = make(map[string]time.Duration) // {"addr": time}
		sumtime     time.Duration
	)

	for _, addr := range strings.Split(c.Addr, ",") {
		startime := time.Now() // 记录开始时间

		if err := deploy(c, addr); err != nil {
			return err
		}

		deploytimes[addr] = time.Since(startime)
		sumtime += deploytimes[addr]
	}

	// 打印汇总信息
	printSummary(deploytimes, sumtime)
	return
}

func deploy(c *Config, addr string) (err error) {
	log.InfoF(longLine, addr)
	log.InfoF("Deploying to %s ...", addr)
	sshClient, sftpClient, err := Init(c, addr)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	log.SuccessF("Connected to %s", addr)

	// pre commands
	log.InfoF(sortLine, "Pre command")
	if err := runCommands(sshClient, c.PreCmd, c.WorkDir, c.ChangeWorkDir); err != nil {
		return err
	}

	// upload files
	log.InfoF(sortLine, "Upload file to remote")
	srcFilePaths := strings.Split(c.SrcFile, ",")
	for i, fp := range srcFilePaths {
		if err := upload(sftpClient, fp, c.WorkDir, i, len(srcFilePaths)); err != nil {
			return err
		}
	}

	// post commands
	log.InfoF(sortLine, "Post command")
	if err := runCommands(sshClient, c.PostCmd, c.WorkDir, c.ChangeWorkDir); err != nil {
		return err
	}
	return
}

// printSummary prints the summary of the deployment
func printSummary(deploytimes map[string]time.Duration, sumtime time.Duration) {
	log.Info(line)
	log.Info("Summary:")
	log.Info("")

	for addr, deploytime := range deploytimes {
		//总长度固定52
		//计算需要补充的.的个数
		dotNum := 52 - len(addr)
		dotStr := ""
		for i := 0; i < dotNum; i++ {
			dotStr += "."
		}
		log.InfoF("%s %s %s [  %f s]", addr, dotStr, color.GreenString("SUCCESS"), deploytime.Seconds())
	}
	log.Info(line)
	log.Info(color.GreenString("DEPLOY SUCCESS"))
	log.Info(line)
	log.InfoF("Total time: %f s", sumtime.Seconds())
	log.InfoF("Finished at: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Info(line)
}

// runCommands runs the commands on the remote server
func runCommands(client *ssh.Client, commands []string, workDir string, changeWorkDir bool) (err error) {

	if commands == nil || len(commands) == 0 {
		return
	}

	session, err := gssh.NewSession(client)
	if err != nil {
		return err
	}
	defer session.Close()

	_ = gssh.Run("mkdir -p " + workDir)

	if changeWorkDir {
		_ = gssh.Run("cd " + workDir)
	}

	var output string
	for _, cmd := range commands {
		log.InfoShell(cmd)
		output, _ = gssh.CommandOutput(cmd)
		log.Info(output)
	}

	_ = gssh.Run("exit")

	if err := session.Wait(); err != nil {
		return errors.New(output)
	}
	return
}

// upload uploads the specified file to the remote server
func upload(sftpClient *sftp.Client, srcfile, targetDir string, i, total int) (err error) {

	var (
		file     *os.File
		fileinfo os.FileInfo

		filename = filepath.Base(srcfile)
		filesize int64
	)

	if file, err = os.Open(srcfile); err != nil {
		return err
	}
	defer file.Close()

	if fileinfo, err = file.Stat(); err != nil {
		return err
	}
	filesize = fileinfo.Size()

	// 创建目标文件
	targetfile, err := sftpClient.Create(targetDir + "/" + filename)
	if err != nil {
		return err
	}
	defer targetfile.Close()

	//进度条
	width, err := getTermWidth()
	if err != nil {
		width = 80
	}

	width = int(float64(width) * 0.6)
	bar := progressbar.NewOptions64(filesize,
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(width),
		//progressbar.OptionFullWidth(),
		//progressbar.OptionSetDescription("[cyan][1/1][reset] "+filename+" "),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] %s ", i+1, total, filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	//上传文件
	if _, err := io.Copy(targetfile, io.TeeReader(file, bar)); err != nil {
		return err
	}
	fmt.Println()

	return
}

func getTermWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}
