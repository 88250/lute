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
)

func main() {
	luteEngine := lute.New() // 默认已经启用 GFM 支持以及中文优化

	html, err := luteEngine.MarkdownStr("demo", "**Lute**")
	if nil != err {
		panic(err)
	}
	fmt.Println(html)
	// <p><strong>Lute</strong></p>
}
