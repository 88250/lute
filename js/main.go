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

package main

import (
	"fmt"
	"github.com/b3log/lute"
	"github.com/gopherjs/gopherjs/js"
)

func markdown(markdownText string) string {
	luteEngine := lute.New()
	html, err := luteEngine.MarkdownStr("", markdownText)
	if nil != err {
		fmt.Println(err)
	}
	return html
}

func format(markdownText string) string {
	luteEngine := lute.New()
	formatted, err := luteEngine.FormatStr("", markdownText)
	if nil != err {
		fmt.Println(err)
	}
	return formatted
}

func getEmojis(imgStaticPath string) map[string]string {
	return lute.New().GetEmojis(imgStaticPath)
}

func main() {
	js.Global.Set("lute", make(map[string]interface{}))
	lute := js.Global.Get("lute")
	lute.Set("markdown", markdown)
	lute.Set("getEmojis", getEmojis)
	lute.Set("format", format)
}
