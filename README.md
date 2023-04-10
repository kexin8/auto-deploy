## 安装

### 脚本安装（推荐）

#### Window

```shell
# Optional: Needed to run a remote script the first time
# 可选：第一次运行远程脚本时需要
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm https://github.com/kexin8/auto-deploy/releases/download/install/install.ps1 | iex
# if you can't access github, you can use proxy
# 如果访问github很慢，可以使用镜像代理
irm https://ghproxy.com/https://github.com/kexin8/auto-deploy/releases/download/install/install.ps1 | iex
```

#### Linux & Mac
```shell
curl https://github.com/kexin8/auto-deploy/releases/download/install/install.sh | sh
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
