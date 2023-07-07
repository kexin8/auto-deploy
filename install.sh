#!/bin/bash
################################################################
# @Author: kexin8
# @Date: 2023-06-27
# @Description: install auto-deploy tool
################################################################

function warn() {
  echo -e "\033[33m$1\033[0m"
}

function info() {
  echo $1
}

function error() {
  echo -e "\033[31m$1\033[0m"
}

# 判断命令是否执行成功
function check() {
  if [ $? -ne 0 ]; then
    error $1
    exit 1
  fi
}

# 获取版本号
VERSION=$(curl -s https://api.github.com/repos/kexin8/auto-deploy/releases/latest | grep tag_name | cut -d '"' -f 4)
URL="https://github.com/kexin8/auto-deploy/releases/download/$VERSION"
Proxy=$1

if [ -n "$Proxy" ]; then
  URL="$Proxy/$URL"
fi

DEPLOY_DIR=""

# 获取当前系统
os=$(uname -s)
if [ $os == "Darwin" ]; then
  os="darwin"
  DEPLOY_DIR="$HOME/Applications/deploy"
elif [ $os == "Linux" ]; then
  os="linux"
  DEPLOY_DIR="/usr/bin/deploy"
else
  error "不支持的系统 $os"
  exit 1
fi

# 获取当前系统架构 x86_64 or arm64
arch=$(uname -m)
if [ $arch == "x86_64" ]; then
  arch="amd64"
elif [ $arch == "arm64" ]; then
  arch="arm64"
else
  error "不支持的架构 $arch"
  exit 1
fi

DownloadUrl="$URL/deploy_"$os"_"$arch".tar.gz"
info "download $DownloadUrl to $DEPLOY_DIR"

tarFileTmpDir=$DEPLOY_DIR/tmp

# 创建临时目录，并判断是否有权限
if [ -w $DEPLOY_DIR ]; then
  mkdir -p $tarFileTmpDir
else
  warn "permission denied, auto change to sudo"
  sudo mkdir -p $tarFileTmpDir
fi

# 解压至临时目录，避免覆盖
info "download $DownloadUrl"
curl $DownloadUrl | tar -zxf - -C $tarFileTmpDir
check "download $DownloadUrl failed"

# 复制文件到目标目录
cp $tarFileTmpDir/deploy_"$os"_"$arch"/* $DEPLOY_DIR

# 删除临时目录
rm -rf $tarFileTmpDir

# 获取当前系统的环境变量Path，判断是否已经存在，不存在则添加
path=$(echo $PATH | grep $DEPLOY_DIR)
if [ -z "$path" ]; then
  info "export PATH=$DEPLOY_DIR:$PATH" >>~/.bash_profile
  # tips
  warn "please run 'source ~/.bash_profile'"
  warn "if you are using zsh, please run 'echo export PATH=$DEPLOY_DIR:'\$PATH' >> ~/.zshrc && source ~/.zshrc'"
fi
