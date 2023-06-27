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
VERSION=`curl -s https://api.github.com/repos/kexin8/auto-deploy/releases/latest | grep tag_name | cut -d '"' -f 4`
URL="https://github.com/kexin8/auto-deploy/releases/download/$VERSION"
Proxy=$1

if [ -n "$Proxy" ]; then
    URL="$Proxy/$URL"
fi

DEPLOY_DIR=""

# 获取当前系统
os=`uname -s`
if [ $os == "Darwin" ]; then
    os="darwin"
    DEPLOY_DIR="$HOME/Applications/deploy"
elif [ $os == "Linux" ]; then
    os="linux"
    DEPLOY_DIR="/usr/bin/deploy"
else
    check "不支持的系统 $os"
fi

# 获取当前系统架构 amd64 or arm64
arch=`uname -m`
if [ $arch == "x86_64" ]; then
    arch="amd64"
elif [ $arch == "arm64" ]; then
    arch="arm64"
else
    check "不支持的系统架构 $arch"
fi

info "download $URL/deploy-$os-$arch.tgz to $DEPLOY_DIR"

tarFileTmpDir=$DEPLOY_DIR/tmp
mkdir -p $tarFileTmpDir
check "create tmp dir failed $tarFileTmpDir"

# 解压至临时目录，避免覆盖
curl $URL/deploy-$os-$arch.tgz | tar -zxf - -C $tarFileTmpDir
check "download deploy-$os-$arch.tgz failed"

# 复制文件到目标目录
cp $tarFileTmpDir/* $DEPLOY_DIR

# 删除临时目录
rm -rf $tarFileTmpDir

# 获取当前系统的环境变量Path，判断是否已经存在，不存在则添加
path=`echo $PATH | grep $DEPLOY_DIR`
if [ -z "$path" ]; then
    info "export PATH=$DEPLOY_DIR:$PATH" >> ~/.bash_profile
    # tips
    warn "please run 'source ~/.bash_profile'"
    warn "if you are using zsh, please run 'echo export PATH=$DEPLOY_DIR:'\$PATH' >> ~/.zshrc && source ~/.zshrc'"
fi
