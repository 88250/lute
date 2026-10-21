# JavaScript 构建

`javascript/lute.min.js` 由 GopherJS 从 Go 源码生成，供思源编辑器前端使用。该文件随源码一起提交，修改 Go 源码后需要重新生成。

## 为什么不能直接用 `go build`

GopherJS 只支持某一个特定的 Go 次版本（当前 `v1.21.0` 对应 Go 1.21.x），而本模块的 `go.mod` 声明了更新的 Go 版本，部分依赖也要求更新的版本：

| 项 | 要求 |
| --- | --- |
| GopherJS `v1.21.0` | Go 1.21.x |
| `go.mod` | Go 1.26 |
| `github.com/alecthomas/chroma/v2 v2.27.0` | Go >= 1.25 |

三者在模块模式下无法同时满足，因此构建必须以 **GOPATH 模式** 运行：忽略 `go.mod`，用与 GopherJS 匹配的 Go 工具链编译。

另外，`github.com/gopherjs/gopherjs/js` 包必须与 GopherJS 二进制版本一致，否则会缺少 `js.Uint64High` 等运行时符号。

## 推荐做法：`javascript/build.ps1`

脚本会依次完成：下载匹配的 Go 工具链、用该工具链安装 GopherJS、准备临时 GOPATH（依赖用目录联接，不复制）、执行构建、按仓库约定处理产物。

```powershell
# 在仓库根目录
pwsh -File .\javascript\build.ps1
# 网络受限时使用国内代理
pwsh -File .\javascript\build.ps1 -GOPROXY https://goproxy.cn,direct -GOSUMDB sum.golang.google.cn
```

可用参数：`-GopherJSVersion`、`-GoVersion`、`-GOPROXY`、`-GOSUMDB`、`-KeepWorkDir`（保留临时目录便于排查）。

注意事项：

- 首次运行需要网络下载工具链与 GopherJS。
- 不能用 `GOSUMDB=off`，否则工具链模块校验失败；网络受限时应指定可用的校验数据库。
- 产物会被处理为纯 ASCII、CRLF 行尾，并去掉末尾的 `sourceMappingURL` 注释，与仓库中已有的产物保持一致。

## 手工步骤

其他平台可以按同样的思路手工执行，关键在于「匹配的 Go 工具链 + GOPATH 模式 + 对应版本的 js 包」：

```sh
# 1. 取得匹配的 Go 工具链 GOROOT
GOTOOLCHAIN=go1.21.13 go env GOROOT

# 2. 用该工具链安装 GopherJS
GOTOOLCHAIN=go1.21.13 GOBIN="$PWD/.build/bin" go install github.com/gopherjs/gopherjs@v1.21.0

# 3. 准备 GOPATH：本仓库与全部依赖链接到 $GOPATH/src/<module path>
#    依赖目录可用 go list -m -f '{{.Path}} {{.Dir}}' all 获取
#    排除上述列表中的 github.com/gopherjs/gopherjs，改用与编译器一致的 v1.21.0
#    将以下命令输出的目录链接到 $GOPATH/src/github.com/gopherjs/gopherjs
go list -m -f '{{.Dir}}' github.com/gopherjs/gopherjs@v1.21.0

# 4. 在仓库根目录构建 JavaScript 入口包
GOTOOLCHAIN=local GOPATH=... GOROOT=... GO111MODULE=off GOOS=js GOARCH=ecmascript \
  .build/bin/gopherjs build --tags javascript -o javascript/lute.min.js -m github.com/88250/lute/javascript
```

## 验证

生成的产物可以直接在 Node.js 中加载：

```sh
node -e "require('./javascript/lute.min.js'); const l = Lute.New(); l.SetProtyleWYSIWYG(true); console.log(l.BlockDOM2StdMd(l.Md2BlockDOM('- a\n  - b\n')))"
```
