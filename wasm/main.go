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
	_ = luteEngine
	return html
}

// 导出 TinyGo 函数。
//go:export md
func md(markdownText string) string {
	fmt.Println("TinyGo！")
	return "TinyGo！"
}

func main() {
	js.Global().Set("markdown", js.FuncOf(markdown))
	select {}
}
