## 安装

### 手动下载

#### Window

##### PowerShell
```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-windows-amd64.tgz

# 解压
tar -xvzf deploy-windows-amd64.tgz -C /your/path
# 设置环境变量
$env:Path="$env:Path;/your/path"
```

#### Linux
```shell
# 下载
wget https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-linux-amd64.tgz
# 解压
tar -zxvf deploy-linux-x64.tgz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### Mac
```shell
# 下载
curl https://github.com/kexin8/auto-deploy/releases/download/{latest-version}/deploy-darwin-amd64.tgz
# 解压
tar -zxvf deploy-darwin-amd64.tgz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### 脚本安装（推荐）

#### Window

暂无

#### Linux & Mac
```shell
curl https://github.com/kexin8/auto-deploy/releases/download/{version}/install.sh | sh
```
