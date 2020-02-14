// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import "bytes"

func (context *Context) parseToC(paragraph *Node) *Node {
	lines := split(paragraph.tokens, itemNewline)
	if 1 != len(lines) {
		return nil
	}

	content := bytes.TrimSpace(lines[0])
	if !bytes.EqualFold(content, []byte("[toc]")) {
		return nil
	}

	return &Node{typ: NodeToC}
}
