
## 用法

```shell
# 初始化配置
deploy init [-a]
# 修改配置文件
vim dyconfig.json
# 执行
deploy
```

## 安装

### 脚本安装（推荐）

#### Window

##### PowerShell

```shell
# Optional: Needed to run a remote script the first time
# 可选：第一次运行远程脚本时需要
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm https://github.com/kexin8/auto-deploy/releases/download/install/install.ps1 | iex
# if you can't access github, you can use proxy
# 如果访问github很慢，可以使用镜像代理
irm https://github.com/kexin8/auto-deploy/releases/download/install/install.ps1 -Proxy '<host>:<ip>' | iex
```

*国内访问*
```shell
# 可选：第一次运行远程脚本时需要
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm https://ghproxy.com/https://github.com/kexin8/auto-deploy/releases/download/install/install_ZH-CN.ps1 | iex
```

#### Linux & Mac
```shell
curl -fsSL https://github.com/kexin8/auto-deploy/releases/download/install/install.sh | sh
```

*国内访问*
```shell
curl -fsSL https://ghproxy.com/https://github.com/kexin8/auto-deploy/releases/download/install/install.sh | sh https://ghproxy.com
```

### 手动下载

#### Window

##### PowerShell（推荐）
```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-windows-amd64.tgz

# 解压
tar -xvzf deploy-windows-amd64.tgz -C /your/path
# 设置环境变量
[environment]::SetEnvironmentvariable("PATH", "$([environment]::GetEnvironmentvariable("Path", "User"));/your/path", "User")
```

##### Cmd

```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-windows-amd64.tgz

# 解压
tar -xvzf deploy-windows-amd64.tgz -C /your/path
# 手动添加环境变量
```



#### Linux

```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-linux-amd64.tgz
# 解压
tar -zxvf deploy-linux-x64.tgz -C /your/path/
# 设置环境变量（追加）
export PATH=$PATH:/your/path
```

### Mac
```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-darwin-amd64.tgz
# 解压
tar -zxvf deploy-darwin-amd64.tgz -C /your/path/
# 设置环境变量
export PATH=$PATH:/your/path
```

## 命令说明

#### init

> 初始化配置

`deploy init`会再当前路径下创建一个名为`dyconfig.json`的配置文件

```shell
# init dyconfig.json
> deploy init
# Or select Display all configurations
> deploy init -a
```

dyconfig.json

```json
{
	"address": "localhost:22",
	"username": "your_username",
	"password": "your_password",
	"publicKey": "your_pubkey",
	"publicKeyPath": "your_pubkey_path",
	"timeout": 10,
	"srcFile": "your need to deploy file",
	"targetDir": "remote target dir",
	"preCmd": [
		"Before uploading a file"
	],
	"postCmd": [
		"After uploading the file"
	]
}
```

**配置说明**

| 名称          | 必填 | 说明                                 | 样例         |
| ------------- | ---- | ------------------------------------ | ------------ |
| address       | Y    | 服务器地址<br />`<host>:<ip>`        | localhost:22 |
| username      | Y    | 用户名                               |              |
| password      | N    | 密码                                 |              |
| publicKey     | N    | 公钥                                 |              |
| publicKeyPath | N    | 公钥文件路径                         |              |
| timeout       | N    | 超时时间，单位s，默认10s             |              |
| srcFile       | Y    | 需要上传的文件路径，多个使用逗号分隔 |              |
| targetDir     | Y    | 服务器目标路径                       |              |
| preCmd        | N    | 上传文件前需要执行的命令             |              |
| postCmd       | N    | 上传后需要执行的命令                 |              |

#### version

> 查看版本

```shell
> deploy version
>
deploy version v0.2.4 windows/amd64
```



#### help

> 帮助

```shell
> deploy --help
# or
> deploy -h
NAME:
   deploy - this is a simple cli app that automates deploy

USAGE:
   deploy [\path\to\config.json]

DESCRIPTION:
   This is a simple cli app that automates deploy.
   e.g. This is a common way to perform deploy, according to dyconfig.json in the current path
     deploy
   This is manually specifying the configuration file
     deploy \path\to\config.json

COMMANDS:
   init
   version, v  Show version
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

