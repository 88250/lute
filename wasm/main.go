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
	"syscall/js"

	"github.com/b3log/lute"
)

func markdown(this js.Value, args []js.Value) interface{} {
	markdownText := args[0].String()
	luteEngine := lute.New()
	html, _ := luteEngine.MarkdownStr("", markdownText)
	_ = luteEngine
	return html
}

// 导出 TinyGo 函数。
//go:export md
func md(markdownText string) string {
	return "Lute"
}

func main() {
	js.Global().Set("lute", make(map[string]interface{}))
	lute := js.Global().Get("lute")
	lute.Set("markdown", js.FuncOf(markdown))

	select {}
}
