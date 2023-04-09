## 安装

### 手动下载

#### Window
```shell
# 下载
curl https://xxx.com/deploy/releases/download/{version}/deployment.exe -d /your/path/
# 设置环境变量
setx DEPLOY /your/path
```

#### Linux
```shell
# 下载
curl https://xxx.com/releases/download/{version}/deployment-linux-x64.tar.gz
# 解压
tar -zxvf deployment-linux-x64.tar.gz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### Mac
```shell
# 下载
curl https://xxx.com/releases/download/{version}/deployment-darwin-x64.tar.gz
# 解压
tar -zxvf deployment-darwin-x64.tar.gz -C /your/path/
# 设置环境变量
export DEPLOY=/your/path
```

### 脚本安装（推荐）

#### Window

暂无

#### Linux & Mac
```shell
curl https://xxx.com/releases/download/{version}/deployment-install.sh | sh
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

## 开发指南
### 打包
#### Windows
```shell
# linux 下编译 windows10 可执行文件
GOOS=windows GOARCH=amd64 go build -o deploy.exe
```