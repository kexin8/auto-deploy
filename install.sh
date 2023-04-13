#!/bin/bash

# 获取版本号
VERSION=`curl -s https://api.github.com/repos/kexin8/auto-deploy/releases/latest | grep tag_name | cut -d '"' -f 4`
URL="https://github.com/kexin8/auto-deploy/releases/download/$VERSION"

DEPLOY_DIR="/usr/local/bin/deploy"

# 获取当前系统
os=`uname -s`
if [ $os == "Darwin" ]; then
    os="darwin"
elif [ $os == "Linux" ]; then
    os="linux"
else
    echo "不支持的系统 $os"
    exit 1
fi

# 获取当前系统架构 amd64 or arm64
arch=`uname -m`
if [ $arch == "x86_64" ]; then
    arch="amd64"
elif [ $arch == "arm64" ]; then
    arch="arm64"
else
    echo "不支持的系统架构 $arch"
    exit 1
fi

echo "download $URL/deploy-$os-$arch.tgz to $DEPLOY_DIR"

tarFileTmpDir=$DEPLOY_DIR/tmp
mkdir -p $tarFileTmpDir

# 解压至临时目录，避免覆盖
curl $URL/deploy-$os-$arch.tgz | tar -zxvf - -C $tarFileTmpDir

# 复制文件到目标目录
cp $tarFileTmpDir/* $DEPLOY_DIR

# 删除临时目录
rm -rf $tarFileTmpDir

# 获取当前系统的环境变量Path，判断是否已经存在，不存在则添加
path=`echo $PATH | grep "$DEPLOY_DIR"`
if [ -z "$path" ]; then
    echo "export PATH=$DEPLOY_DIR:\$PATH" >> ~/.bash_profile
    # shellcheck disable=SC1090
    source ~/.bash_profile
fi
