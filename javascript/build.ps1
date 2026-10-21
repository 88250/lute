# Lute - A structured markdown engine.
# Copyright (c) 2019-present, b3log.org
#
# Lute is licensed under the Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#     http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v2 for more details.
#
# Build javascript/lute.min.js.
#
# GopherJS only supports one specific Go minor version (1.21 at the time of writing), while this
# module declares a newer Go version and depends on modules that require it. GopherJS therefore
# cannot build in module mode: the module graph is ignored and the build runs in GOPATH mode with a
# Go toolchain matching GopherJS. This script prepares the toolchain, the GopherJS binary and a
# temporary GOPATH, then writes the artifact in the format the repository expects (ASCII, without
# the sourceMappingURL comment).

param(
    [string]$GopherJSVersion = "v1.21.0",
    [string]$GoVersion = "1.21.13",
    [string]$GOPROXY = "",
    [string]$GOSUMDB = "",
    [switch]$KeepWorkDir
)

$ErrorActionPreference = "Stop"

# 让 GOPATH 模式下也能解析本模块与全部依赖，用目录联接避免复制依赖
function Add-ModuleLink {
    param(
        [string]$SrcRoot,
        [string]$Path,
        [string]$Dir
    )

    $target = Join-Path $SrcRoot ($Path -replace '/', '\')
    if (Test-Path $target) {
        return
    }
    $parent = Split-Path -Parent $target
    New-Item -ItemType Directory -Path $parent -Force | Out-Null
    New-Item -ItemType Junction -Path $target -Target $Dir | Out-Null
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$jsDir = Join-Path $repoRoot "javascript"
$goMod = Join-Path $repoRoot "go.mod"

if (-not (Test-Path $goMod)) {
    throw "go.mod not found, run this script from the repository"
}

$modulePath = (Select-String -Path $goMod -Pattern '^module\s+(\S+)').Matches[0].Groups[1].Value
$toolchain = "go$GoVersion"

$saved = @{
    GOTOOLCHAIN = $env:GOTOOLCHAIN
    GOPROXY     = $env:GOPROXY
    GOSUMDB     = $env:GOSUMDB
    GOBIN       = $env:GOBIN
    GOPATH      = $env:GOPATH
    GOROOT      = $env:GOROOT
    GO111MODULE = $env:GO111MODULE
    GOOS        = $env:GOOS
    GOARCH      = $env:GOARCH
    GOFLAGS     = $env:GOFLAGS
    PATH        = $env:PATH
}

$workDir = Join-Path $env:TEMP ("lute-js-build-" + [System.Guid]::NewGuid().ToString("N").Substring(0, 8))
New-Item -ItemType Directory -Path $workDir -Force | Out-Null

try {
    if ("" -ne $GOPROXY) {
        $env:GOPROXY = $GOPROXY
    }
    if ("" -ne $GOSUMDB) {
        $env:GOSUMDB = $GOSUMDB
    }

    Write-Host "==> Resolving Go $GoVersion toolchain"
    # 前序环境可能遗留交叉编译或 GOPATH 设置，模块模式下必须清空
    $env:GO111MODULE = "on"
    Remove-Item Env:GOROOT -ErrorAction SilentlyContinue
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    $env:GOTOOLCHAIN = $toolchain
    $goroot = (& go env GOROOT).Trim()
    if (-not (Test-Path (Join-Path $goroot "src\runtime"))) {
        throw "Go $GoVersion toolchain not found at $goroot"
    }

    Write-Host "==> Installing GopherJS $GopherJSVersion with Go $GoVersion"
    $goBin = Join-Path $workDir "gobin"
    New-Item -ItemType Directory -Path $goBin -Force | Out-Null
    $env:GOBIN = $goBin
    & go install "github.com/gopherjs/gopherjs@$GopherJSVersion"
    if (0 -ne $LASTEXITCODE) {
        throw "go install gopherjs failed with exit code $LASTEXITCODE"
    }
    $gopherjs = Join-Path $goBin "gopherjs.exe"
    if (-not (Test-Path $gopherjs)) {
        $gopherjs = Join-Path $goBin "gopherjs"
    }

    Write-Host "==> Preparing GOPATH"
    $env:GOTOOLCHAIN = "auto"
    # 只读模块模式，避免构建过程改写 go.sum
    $env:GOFLAGS = "-mod=readonly"
    & go mod download
    if (0 -ne $LASTEXITCODE) {
        throw "go mod download failed with exit code $LASTEXITCODE"
    }

    $gopath = Join-Path $workDir "gopath"
    $srcRoot = Join-Path $gopath "src"
    New-Item -ItemType Directory -Path $srcRoot -Force | Out-Null

    # 复制本模块源码到 GOPATH，避免构建过程写入工作区
    $selfDir = Join-Path $srcRoot ($modulePath -replace '/', '\')
    New-Item -ItemType Directory -Path $selfDir -Force | Out-Null
    & robocopy $repoRoot $selfDir /E /XD .git node_modules /NFL /NDL /NJH /NJS /NP | Out-Null
    if (7 -lt $LASTEXITCODE) {
        throw "copying sources failed with exit code $LASTEXITCODE"
    }
    $global:LASTEXITCODE = 0
    # js 包必须与 GopherJS 二进制版本一致，否则会缺少新版运行时符号
    $gopherjsModule = "github.com/gopherjs/gopherjs"
    $gopherjsDir = (& go list -m -f "{{.Dir}}" "$gopherjsModule@$GopherJSVersion").Trim()
    if (0 -ne $LASTEXITCODE -or "" -eq $gopherjsDir) {
        throw "cannot resolve $gopherjsModule@$GopherJSVersion"
    }
    Add-ModuleLink -SrcRoot $srcRoot -Path $gopherjsModule -Dir $gopherjsDir

    $modules = & go list -m -f "{{.Path}}|{{.Dir}}" all
    if (0 -ne $LASTEXITCODE) {
        throw "go list -m failed with exit code $LASTEXITCODE"
    }
    foreach ($line in $modules) {
        if ($line -notmatch '^(.+)\|(.+)$') {
            continue
        }
        $path = $Matches[1]
        $dir = $Matches[2]
        if ("" -eq $dir -or $path -eq $modulePath -or $path -eq $gopherjsModule) {
            continue
        }
        Add-ModuleLink -SrcRoot $srcRoot -Path $path -Dir $dir
    }

    Write-Host "==> Building lute.min.js"
    $env:GOPATH = $gopath
    $env:GOROOT = $goroot
    $env:GO111MODULE = "off"
    $env:GOOS = "js"
    $env:GOARCH = "ecmascript"
    $env:GOTOOLCHAIN = "local"
    $env:PATH = (Join-Path $goroot "bin") + [System.IO.Path]::PathSeparator + $saved.PATH

    $buildDir = Join-Path (Join-Path $srcRoot ($modulePath -replace '/', '\')) "javascript"
    # 清掉旧产物，避免构建失败时误用上次的结果
    Remove-Item (Join-Path $buildDir "lute.min.js") -Force -ErrorAction SilentlyContinue
    Remove-Item (Join-Path $buildDir "lute.min.js.map") -Force -ErrorAction SilentlyContinue
    Push-Location $buildDir
    try {
        & $gopherjs build --tags javascript -o lute.min.js -m
        if (0 -ne $LASTEXITCODE) {
            throw "gopherjs build failed with exit code $LASTEXITCODE"
        }
    } finally {
        Pop-Location
    }

    $artifact = Join-Path $buildDir "lute.min.js"
    if (-not (Test-Path $artifact)) {
        throw "gopherjs did not produce $artifact"
    }

    Write-Host "==> Post-processing artifact"
    $lines = [System.IO.File]::ReadAllLines($artifact)
    $kept = $lines | Where-Object { -not $_.StartsWith("//# sourceMappingURL=") }
    # 仓库中的产物为纯 ASCII 且使用 CRLF，保证不同平台产出一致
    foreach ($line in $kept) {
        foreach ($ch in $line.ToCharArray()) {
            if ([int]$ch -gt 127) {
                throw "non-ASCII character found in artifact, aborting to avoid breaking it"
            }
        }
    }
    $content = ($kept -join "`r`n") + "`r`n"
    [System.IO.File]::WriteAllText((Join-Path $jsDir "lute.min.js"), $content, [System.Text.Encoding]::ASCII)
    Copy-Item (Join-Path $buildDir "lute.min.js.map") (Join-Path $jsDir "lute.min.js.map") -Force

    Write-Host "==> Done: javascript/lute.min.js"
} finally {
    foreach ($key in $saved.Keys) {
        if ($null -eq $saved[$key]) {
            Remove-Item "Env:$key" -ErrorAction SilentlyContinue
        } else {
            Set-Item "Env:$key" $saved[$key]
        }
    }
    if (-not $KeepWorkDir) {
        Remove-Item $workDir -Recurse -Force -ErrorAction SilentlyContinue
    } else {
        Write-Host "work dir kept: $workDir"
    }
}
