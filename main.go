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
	url           = "https://github.com/kexin8/auto-deploy"
)

var (
	Version = "v0.0.0"
	Os      = "windows"
	Arch    = "amd64"
)

func main() {

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("deploy version %s %s/%s\r\n", ctx.App.Version, Os, Arch)
	}

	app := &cli.App{
		Name: "deploy",
		Description: `This is a simple cli app that automates deploy.
e.g. This is a common way to perform deploy, according to dyconfig.json in the current path
	deploy
This is manually specifying the configuration file
	deploy \path\to\config.json`,
		Usage:     "this is a simple cli app that automates deploy",
		UsageText: `deploy [\path\to\config.json]`,
		Version:   Version,
		Action: func(ctx *cli.Context) error {

			// 进行版本检查，如果有新版本则提示更新
			//latestVersion, err := deploy.GetLatestVersion()
			//if err != nil {
			//	// 不影响正常使用
			//	log.ErrorF("check latest version failed: %s", err)
			//}

			//if latestVersion != "" && latestVersion != Version {
			//	log.Info(color.YellowString("latest version %s is available, please update to the latest version: %s", latestVersion, url))
			//}

			// 仅提醒用户访问网址，不影响正常使用
			log.Info(color.YellowString("please visit %s for more information", url))

			profile := ctx.Args().First()
			if profile == "" {
				//检查当前目录是否存在配置文件 pdconfig.json
				_, err := os.Stat(defConfigName)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("dyconfig.json does not exist, please use 'deploy init' to initialize")
					}
					return err
				}
				profile = defConfigName
			}

			//读取配置文件
			config, err := deploy.ReadConfig(profile)
			if err != nil {
				return err
			}

			if err := config.Init(); err != nil {
				return err
			}
			if err := config.UploadFiles(); err != nil {
				return err
			}
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "a",
						DefaultText: "false",
						Usage:       "All configurations",
					},
				},
				Action: func(ctx *cli.Context) (err error) {

					isAllConfig := ctx.Bool("a")

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

					config := deploy.Config{
						Address:   "localhost:22",
						Username:  "your_username",
						Password:  "your_password",
						SrcFile:   "your need to deploy file",
						TargetDir: "remote target dir",
					}

					if isAllConfig {
						config = deploy.Config{
							Address:       "localhost:22",
							Username:      "your_username",
							Password:      "your_password",
							PublicKey:     "your_pubkey",
							PublicKeyPath: "your_pubkey_path",
							Timeout:       10,
							SrcFile:       "your need to deploy file",
							TargetDir:     "remote target dir",
							PreCmd:        []string{"Before uploading a file"},
							PostCmd:       []string{"After uploading the file"},
						}
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
			{
				Name:  "upgrade",
				Usage: "upgrade deploy",
				Action: func(ctx *cli.Context) error {
					latestVersion, err := deploy.GetLatestVersion()
					if err != nil {
						return err
					}

					if latestVersion == Version {
						log.Info("deploy is already the latest version")
						return nil
					}

					log.InfoF("latest  version: %s", latestVersion)
					log.InfoF("current version: %s", Version)
					log.Info(color.YellowString("latest version %s is available, please update to the latest version: %s", latestVersion, url))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
