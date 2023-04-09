#!/bin/bash

VERSION="{version}"
URL="https://github.com/kexin8/auto-deploy/releases/download/$VERSION"

# 获取当前系统
os=`uname -s`
if [ $os == "Darwin" ]; then
    os="darwin"
elif [ $os == "Linux" ]; then
    os="linux"
else
    echo "不支持的系统"
    exit 1
fi

# 解压
curl $URL/deploy-$os-amd64.tgz | tar -zxvf - -C /usr/local/bin/

# 设置环境变量
echo "export DEPLOY=/usr/local/bin" >> ~/.bash_profile
