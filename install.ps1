<#
.SYNOPSIS
    Deploy installer.
.DESCRIPTION
    The installer of Deploy. For details please check the website and wiki.
.PARAMETER DeployDir
    Specifies Deploy root path.
    If not specified, Deploy will be installed to '$env:USERPROFILE\deploy'.
.PARAMETER NoProxy
    Bypass system proxy during the installation.
.PARAMETER Proxy
    Specifies proxy to use during the installation.
.PARAMETER ProxyCredential
    Specifies credential for the given proxy.
.PARAMETER ProxyUseDefaultCredentials
    Use the credentials of the current user for the proxy server that is specified by the -Proxy parameter.
.PARAMETER RunAsAdmin
    Force to run the installer as administrator.
.LINK
    https://github.com/kexin8/auto-deploy
#>
param(
    [String] $DeployDir, # 安装目录
    [Switch] $NoProxy, # 不使用代理
    [Uri] $Proxy, # 代理地址
    [System.Management.Automation.PSCredential] $ProxyCredential, # 代理认证
    [Switch] $ProxyUseDefaultCredentials, # 代理使用默认认证
    [Switch] $RunAsAdmin # 是否以管理员身份运行
)

# Disable StrictMode in this script
Set-StrictMode -Off

function Write-InstallInfo
{
    param(
        [Parameter(Mandatory = $True, Position = 0)]
        [String] $String,
        [Parameter(Mandatory = $False, Position = 1)]
        [System.ConsoleColor] $ForegroundColor = $host.UI.RawUI.ForegroundColor
    )

    $backup = $host.UI.RawUI.ForegroundColor

    if ($ForegroundColor -ne $host.UI.RawUI.ForegroundColor)
    {
        $host.UI.RawUI.ForegroundColor = $ForegroundColor
    }

    Write-Output "$String"

    $host.UI.RawUI.ForegroundColor = $backup
}

function Deny-Install
{
    param(
        [String] $message,
        [Int] $errorCode = 1
    )

    Write-InstallInfo -String $message -ForegroundColor DarkRed
    Write-InstallInfo "Abort."

    # Don't abort if invoked with iex that would close the PS session
    if ($IS_EXECUTED_FROM_IEX)
    {
        break
    }
    else
    {
        exit $errorCode
    }
}

function Test-ValidateParameter
{

    if ($null -eq $Proxy -and ($null -ne $ProxyCredential -or $ProxyUseDefaultCredentials))
    {
        Deny-Install "Provide a valid proxy URI for the -Proxy parameter when using the -ProxyCredential or -ProxyUseDefaultCredentials."
    }

    if ($ProxyUseDefaultCredentials -and $null -ne $ProxyCredential)
    {
        Deny-Install "ProxyUseDefaultCredentials is conflict with ProxyCredential. Don't use the -ProxyCredential and -ProxyUseDefaultCredentials together."
    }
}

function Test-IsAdministrator
{
    return ([Security.Principal.WindowsPrincipal]`
            [Security.Principal.WindowsIdentity]::GetCurrent()`
    ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator) -and $env:USERNAME -ne 'WDAGUtilityAccount'
}

function Optimize-SecurityProtocol
{
    # .NET Framework 4.7+ has a default security protocol called 'SystemDefault',
    # which allows the operating system to choose the best protocol to use.
    # If SecurityProtocolType contains 'SystemDefault' (means .NET4.7+ detected)
    # and the value of SecurityProtocol is 'SystemDefault', just do nothing on SecurityProtocol,
    # 'SystemDefault' will use TLS 1.2 if the webrequest requires.
    $isNewerNetFramework = ([System.Enum]::GetNames([System.Net.SecurityProtocolType]) -contains 'SystemDefault')
    $isSystemDefault = ([System.Net.ServicePointManager]::SecurityProtocol.Equals([System.Net.SecurityProtocolType]::SystemDefault))

    # If not, change it to support TLS 1.2
    if (!($isNewerNetFramework -and $isSystemDefault))
    {
        # Set to TLS 1.2 (3072), then TLS 1.1 (768), and TLS 1.0 (192). Ssl3 has been superseded,
        # https://docs.microsoft.com/en-us/dotnet/api/system.net.securityprotocoltype?view=netframework-4.5
        [System.Net.ServicePointManager]::SecurityProtocol = 3072 -bor 768 -bor 192
        Write-Verbose "SecurityProtocol has been updated to support TLS 1.2"
    }
}

function Expand-ZipArchive {
    param(
        [String] $path,
        [String] $to
    )

    if (!(Test-Path $path)) {
        Deny-Install "Unzip failed: can't find $path to unzip."
    }

    # Check if the zip file is locked, by antivirus software for example
    $retries = 0
    while ($retries -le 10) {
        if ($retries -eq 10) {
            Deny-Install "Unzip failed: can't unzip because a process is locking the file."
        }
        if (Test-isFileLocked $path) {
            Write-InstallInfo "Waiting for $path to be unlocked by another process... ($retries/10)"
            $retries++
            Start-Sleep -Seconds 2
        } else {
            break
        }
    }

    # Workaround to suspend Expand-Archive verbose output,
    # upstream issue: https://github.com/PowerShell/Microsoft.PowerShell.Archive/issues/98
    $oldVerbosePreference = $VerbosePreference
    $global:VerbosePreference = 'SilentlyContinue'

    # Disable progress bar to gain performance
    $oldProgressPreference = $ProgressPreference
    $global:ProgressPreference = 'SilentlyContinue'

    # PowerShell 5+: use Expand-Archive to extract zip files
    Microsoft.PowerShell.Archive\Expand-Archive -Path $path -DestinationPath $to -Force
    $global:VerbosePreference = $oldVerbosePreference
    $global:ProgressPreference = $oldProgressPreference
}

function Get-Env
{
    param(
        [String] $name,
        [Switch] $global
    )

    $RegisterKey = if ($global)
    {
        Get-Item -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager'
    }
    else
    {
        Get-Item -Path 'HKCU:'
    }

    $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment')
    $RegistryValueOption = [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames
    $EnvRegisterKey.GetValue($name, $null, $RegistryValueOption)
}

function Write-Env
{
    param(
        [String] $name,
        [String] $val,
        [Switch] $global
    )

    $RegisterKey = if ($global)
    {
        Get-Item -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager'
    }
    else
    {
        Get-Item -Path 'HKCU:'
    }

    $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment', $true)
    if ($val -eq $null)
    {
        $EnvRegisterKey.DeleteValue($name)
    }
    else
    {
        $RegistryValueKind = if ( $val.Contains('%'))
        {
            [Microsoft.Win32.RegistryValueKind]::ExpandString
        }
        elseif ($EnvRegisterKey.GetValue($name))
        {
            $EnvRegisterKey.GetValueKind($name)
        }
        else
        {
            [Microsoft.Win32.RegistryValueKind]::String
        }
        $EnvRegisterKey.SetValue($name, $val, $RegistryValueKind)
    }
}

function Add-DeployDirToPath
{
    # Get $env:PATH of current user
    $userEnvPath = Get-Env 'PATH'

    if ($userEnvPath -notmatch [Regex]::Escape($DEPLOY_DIR))
    {
        $h = (Get-PSProvider 'FileSystem').Home
        if (!$h.EndsWith('\'))
        {
            $h += '\'
        }

        if (!($h -eq '\'))
        {
            $friendlyPath = "$DEPLOY_DIR" -Replace ([Regex]::Escape($h)), "~\"
            Write-InstallInfo "Adding $friendlyPath to your path."
        }
        else
        {
            Write-InstallInfo "Adding $DEPLOY_DIR to your path."
        }

        # For future sessions
        Write-Env 'PATH' "$userEnvPath;$DEPLOY_DIR"
        # For current session
        $env:PATH = "$env:PATH;$DEPLOY_DIR"
    }
}

function Get-Downloader
{
    $downloadSession = New-Object System.Net.WebClient

    # Set proxy to null if NoProxy is specificed
    if ($NoProxy)
    {
        $downloadSession.Proxy = $null
    }
    elseif ($Proxy)
    {
        # Prepend protocol if not provided
        if (!$Proxy.IsAbsoluteUri)
        {
            $Proxy = New-Object System.Uri("http://" + $Proxy.OriginalString)
        }

        $Proxy = New-Object System.Net.WebProxy($Proxy)

        if ($null -ne $ProxyCredential)
        {
            $Proxy.Credentials = $ProxyCredential.GetNetworkCredential()
        }
        elseif ($ProxyUseDefaultCredentials)
        {
            $Proxy.UseDefaultCredentials = $true
        }

        $downloadSession.Proxy = $Proxy
    }

    return $downloadSession
}

function Install-Deploy
{
    Write-InstallInfo "Initializing..."
    # Validate install parameters
    Test-ValidateParameter
    # Enable TLS 1.2
    Optimize-SecurityProtocol

    Write-InstallInfo "Installing Deploy..."
    Write-Verbose "$BaseUrl"

    $downloader = Get-Downloader

    $deployZipFile = "$DEPLOY_DIR\deploy-windows-amd64.zip"
    # 创建目录
    if (!(Test-Path $DEPLOY_DIR))
    {
        New-Item -Type Directory $DEPLOY_DIR | Out-Null
    }

    # 输出deployTarFile
    Write-InstallInfo "Downloading Deploy from $url to $deployZipFile ..."
    $downloader.downloadFile($URL, $deployZipFile)

    Write-InstallInfo "Extracting ..."

    $deployUnzipTmpDir = "$DEPLOY_DIR\_tmp"
    if (!(Test-Path $deployUnzipTmpDir))
    {
        New-Item -Type Directory $deployUnzipTmpDir | Out-Null
    }
    Write-Verbose "Extracting $deployZipFile to $deployUnzipTmpDir"

    # 解压，如果解压失败，删除临时文件
    #    tar -xzf $deployTarFile -C $deployUnTarFileTmpDir | Out-Null
    Expand-ZipArchive $deployZipFile $deployUnzipTmpDir
    if ($LastExitCode -ne 0)
    {
        Deny-Install "Failed to extract $deployZipFile to $deployUnzipTmpDir"
        Remove-Item $deployZipFile -Force
        Remove-Item $deployUnzipTmpDir -Force -Recurse
        return
    }

    Copy-Item "$deployUnzipTmpDir\deploy_windows_x86_64\*" $DEPLOY_DIR -Recurse -Force

    # 删除临时文件
    Remove-Item $deployZipFile -Force
    Remove-Item $deployUnzipTmpDir -Force -Recurse

    # 设置环境变量
    Add-DeployDirToPath
}

$VERSION = Invoke-WebRequest -Uri "https://api.github.com/repos/kexin8/auto-deploy/releases/latest" -UseBasicParsing | ConvertFrom-Json | Select-Object -ExpandProperty tag_name

$URL = "https://github.com/kexin8/auto-deploy/releases/download/$VERSION/deploy_windows_x86_64.zip"

$DEPLOY_DIR = "$env:LOCALAPPDATA\deploy"

Install-Deploy