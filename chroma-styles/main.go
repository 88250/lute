// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
)

// 生成 Chroma 样式。
func main() {
	dir := "chroma-styles"
	prefix := "highlight-"
	formatter := chromahtml.New(chromahtml.WithClasses(true), chromahtml.ClassPrefix(prefix))
	var b bytes.Buffer
	names := styles.Names()
	for _, name := range names {
		formatter.WriteCSS(&b, styles.Get(name))
		os.WriteFile(filepath.Join(dir, name)+".css", b.Bytes(), 0644)
		b.Reset()
	}

	fmt.Println("[\"" + strings.Join(names, "\", \"") + "\"]")
}
