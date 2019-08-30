// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/b3log/lute"
)

func markdown(this js.Value, args []js.Value) interface{} {
	markdownText := args[0].String()
	fmt.Printf("markdown text [%s]", markdownText)
	luteEngine := lute.New()
	html, _ := luteEngine.MarkdownStr("", markdownText)
	return html
}

// 运行 build 脚本后可生成 lute.wasm，可通过 brotli -o lute.wasm.br lute.wasm 进行压缩
//
// 通过如下方式启动一个 HTTP Server 然后即可访问查看最终效果
// 1. 安装 goexec：go get -u github.com/shurcooL/goexec
// 2. 启动 HTTP 服务：goexec "http.ListenAndServe(`:8080`, gzipped.FileServer(http.Dir(`.`)))"

func main() {
	js.Global().Set("markdown", js.FuncOf(markdown))
	select {}
}
