## 安装

### 手动下载

#### Window
```shell
# 下载
curl https://xxx.com/deploy/releases/download/{version}/deploy.exe -d /your/path/
# 设置环境变量
setx DEPLOY /your/path
```

#### Linux
```shell
# 下载
curl https://xxx.com/releases/download/{version}/deploy-linux-x64.tar.gz
# 解压
tar -zxvf deploy-linux-x64.tar.gz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### Mac
```shell
# 下载
curl https://xxx.com/releases/download/{version}/deploy-darwin-x64.tar.gz
# 解压
tar -zxvf deploy-darwin-x64.tar.gz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### 脚本安装（推荐）

#### Window

暂无

#### Linux & Mac
```shell
curl https://xxx.com/releases/download/{version}/deploy-install.sh | sh
```

deploy-install.sh 脚本内容如下：
```shell
#!/bin/bash

URL="https://xxx.com/releases/download/{version}/deploy-{os}-x64.tar.gz"

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
curl $URL | tar -zxvf - -C /usr/local/bin/

# 设置环境变量
echo "export DEPLOY=/usr/local/bin" >> ~/.bash_profile
```