#!/usr/bin/env node

"use strict"

const request = require('request'),
    path = require('path'),
    process = require('process'),
    fs = require('fs'),
    tar = require('tar'),
    unzipper = require('unzipper'),
    exec = require('child_process').exec;

// Platform from Node's `os.platform()` to Golang's `GOOS`
const PLATFORMS = {
    "darwin": "darwin",
    "linux": "linux",
    "win32": "windows",
    "freebsd": "freebsd"
}
// Arch from Node's `os.arch()` to Golang's `GOARCH`
const ARCHS = {
    "ia32": "386",
    "x64": "amd64",
    "arm": "arm"
}

/**
 * 判断当前用户是否在中国, 通过ip-api.com查询
 * @returns {Promise<boolean>}
 */
function isChinaUser() {
    return new Promise((resolve, reject) => {
        request(`http://ip-api.com/json/?fields=countryCode`, {timeout: 5000}, (error, response, body) => {
            if (error) {
                reject(error);
            } else {
                const data = JSON.parse(body);
                resolve(data.countryCode === 'CN');
            }
        });
    })
}

/**
 * 获取binary文件安装路径
 * @returns {Promise<string>} binary文件安装路径
 */
async function getInstallPath() {
    // 命令行`npm config get prefix`获取
    return await new Promise((resolve, reject) => {
        exec('npm config get prefix', (error, stdout, stderr) => {
            if (error) {
                reject(error);
            } else {
                resolve(stdout.trim());
            }
        });
    })
}

/**
 * 校验binary配置
 * @param config {{name: string, version: string, url: string}}
 * @returns {string}
 */
function validateConfig(config) {

    if (!config || typeof config !== "object") {
        return "'binary' property in package.json must be an object";
    }

    if (!config.name || typeof config.name !== "string") {
        return "'name' property in binary config must be a string";
    }

    if (!config.url || typeof config.url !== "string") {
        return "'url' property in binary config must be a string";
    }
}

/**
 * 获取package.json中的binary配置
 * @returns {Promise<{name: string, version: string, url: string}>}
 */
async function getBinaryOpts() {
    const packageJsonPath = path.join(".", "package.json");
    if (!fs.existsSync(packageJsonPath)) {
        throw new Error("Unable to find package.json. Please run this script at root of the package you want to be installed")
    }
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath));

    const errorMsg = validateConfig(packageJson.binary);
    if (errorMsg) {
        throw new Error(errorMsg)
    }

    const name = process.platform === 'win32' ? packageJson.binary.name + '.exe' : packageJson.binary.name
    const version = packageJson.binary.version ? packageJson.binary.version : packageJson.version
    const {url, proxy} = packageJson.binary
    const ext = process.platform === 'win32' ? 'zip' : 'tar.gz'


    let downloadUrl = url;
    downloadUrl = downloadUrl.replace(/{{os}}/g, PLATFORMS[process.platform]);
    downloadUrl = downloadUrl.replace(/{{arch}}/g, ARCHS[process.arch]);
    downloadUrl = downloadUrl.replace(/{{version}}/g, version);
    downloadUrl = downloadUrl.replace(/{{ext}}/g, ext);

    const isChina = await isChinaUser()
    if (isChina) {
        downloadUrl = proxy + "/" + downloadUrl
    }
    return {name: name, version: version, url: downloadUrl}
}

/**
 * 下载&&解压zip文件
 * @param name 文件名
 * @param dir 解压目录
 * @param url 下载地址
 * @returns {Promise<void>}
 */
async function downloadAndUnzip(name, dir, url) {
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, {recursive: true})
    }

    const directory = await unzipper.Open.url(request, url);
    const file = directory.files
        .find(file => file.path.endsWith(name))

    if (!file) {
        throw new Error("Binary file not found in zip");
    }

    return new Promise((resolve, reject) => {
        file.stream()
            .pipe(fs.createWriteStream(path.join(dir, name)))
            .on('error', reject)
            .on('finish', resolve)
    });
}

/**
 * 下载&&解压tar.gz文件
 * @param name 文件名
 * @param dir 解压目录
 * @param url 下载地址
 * @returns {Promise<void>}
 */
async function downloadAndUntar(name, dir, url) {
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, {recursive: true})
    }
    return await new Promise((resolve, reject) => {
        request(url)
            .pipe(tar.x({
                    strip: 1, // 去除第一层目录
                    // 筛选指定文件
                    filter: (path, stat) => {
                        return path.endsWith(name)
                    },
                    cwd: dir,// 解压目录
                    sync: true // 同步解压
                })
            )
            .on('end', resolve)
            .on('error', reject)
    })
}

/**
 * 安装binary文件
 * @returns {Promise<void>}
 */
async function install() {
    const opts = await getBinaryOpts()

    const dir = 'bin'

    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, {recursive: true})
    }

    console.log(`Downloading ${opts.name} from ${opts.url} ...`)
    if (process.platform === 'win32') {
        await downloadAndUnzip(opts.name, dir, opts.url)
    } else {
        await downloadAndUntar(opts.name, dir, opts.url)
    }

    if (!fs.existsSync(path.join(dir, opts.name))) {
        throw new Error("Binary file not found at path: " + opts.name);
    }

    console.log(`Downloaded done`)

    const installDir = await getInstallPath()
    console.log(`Installing ${opts.name} to ${installDir} ...`)

    if (!fs.existsSync(installDir)) {
        fs.mkdirSync(installDir, {recursive: true})
    }

    fs.renameSync(path.join(dir, opts.name), path.join(installDir, opts.name))
    // fs.symlinkSync(path.join(process.cwd(), dir, opts.name), path.join(installDir, opts.name), 'file')

    console.log(`Installed done`)
    console.log('Please run `deploy -h` or `deploy --help` to get help')
}

/**
 * 卸载binary文件
 * @returns {Promise<void>}
 * @deprecated  <a href="https://github.com/npm/cli/issues/3042">npm v7+ 之后不再支持`uninstall` 生命周期</a>
 */
async function uninstall() {
    const opts = await getBinaryOpts();

    const installPath = await getInstallPath()

    console.log(`Uninstalling ${opts.name} from ${installPath} ...`)

    if (!fs.existsSync(path.join(installPath, opts.name))) {
        // throw new Error("Binary is not install")
        console.warn("Binary is not install")
        return
    }

    fs.unlinkSync(path.join(installPath, opts.name))
}

const cmds = {
    'install': install,
    'uninstall': uninstall
}

// 获取命令行参数，忽略前两个参数
const args = process.argv.slice(2)
if (!args || args.length === 0 || !cmds[args[0]]) {
    console.error('Please enter `install` or `uninstall`')
    process.exit(1)
}

cmds[args[0]]().then(() => {
    process.exit(0)
}).catch((e) => {
    console.error(args[0], 'failed', e)
    process.exit(1)
})



