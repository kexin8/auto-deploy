package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/kexin8/auto-deploy/deploy"
	"github.com/kexin8/auto-deploy/log"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path"
	"path/filepath"
)

const (
	defConfigName = "dyconfig.json"
)

func main() {
	app := &cli.App{
		Name: "deploy",
		Description: `This is a simple cli app that automates deploy.
e.g. This is a common way to perform deploy, according to dyconfig.json in the current path
	deploy
This is manually specifying the configuration file
	deploy \path\to\config.json`,
		Usage:     "this is a simple cli app that automates deploy",
		UsageText: `deploy [\path\to\config.json]`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "enable debug mode",
				Hidden:      true,
				Value:       false,
				DefaultText: "false",
			},
		},
		Action: func(ctx *cli.Context) error {

			if ctx.Bool("debug") {
				log.NewLogger(true)
			}

			profile := ctx.Args().First()
			if profile == "" {
				//检查当前目录是否存在配置文件 pdconfig.json
				_, err := os.Stat(defConfigName)
				if err != nil {
					if os.IsNotExist(err) {
						return cli.Exit("dyconfig.json does not exist, please use 'deploy init' to initialize", -1)
					}
					return err
				}
				profile = defConfigName
			}

			//读取配置文件
			config, err := deploy.ReadConfig(profile)
			if err != nil {
				return cli.Exit(fmt.Sprintf("read config file %s failed : %s", profile, err), -1)
			}

			if err := config.Init(); err != nil {
				return cli.Exit(fmt.Sprintf("init failed : %s", err), -1)
			}
			if err := config.UploadFile(); err != nil {
				return cli.Exit(fmt.Sprintf("upload file failed : %s", err), -1)
			}

			if err == nil {
				//filename DEPLOY					SUCCESS
				green := color.New(color.FgGreen).SprintFunc()
				fmt.Printf("%s DEPLOY					%s", filepath.Base(config.SrcFile), green("SUCCESS"))
			}

			fmt.Println("END.")

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "init",
				Description: `Initialize a new deploy configuration file.
e.g. The usual way to config an app
		deploy init
The specified application directory has been initially configured
		deploy init \path\to\app
`,
				UsageText: `deploy init [\path\to\app]`,
				Action: func(ctx *cli.Context) (err error) {

					appath := ctx.Args().First()
					if appath == "" {
						//获取当前目录
						dir, err := os.Getwd()
						if err != nil {
							return err
						}
						appath = dir
					}

					if appath, err = filepath.Abs(appath); err != nil {
						return err
					}

					//fmt.Println("appath:" + appath)

					config := deploy.DeployConfig{
						Address:   "localhost:22",
						Username:  "your_username",
						Password:  "your_password",
						SrcFile:   "your need to deploy file",
						TargetDir: "remote target dir",
					}

					//写入配置文件
					confjson, err := json.MarshalIndent(config, "", "\t")
					if err != nil {
						return err
					}

					//写入文件
					dpyconfig, err := os.Create(path.Join(appath, defConfigName))
					if err != nil {
						return err
					}
					//goland:noinspection GoUnhandledErrorResult
					defer dpyconfig.Close()
					if _, err := io.WriteString(dpyconfig, string(confjson)); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
