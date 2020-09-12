// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"github.com/88250/lute/util"
)

var openCurlyBraceColon = util.StrToBytes("{: ")
var emptyIAL = util.StrToBytes("{:}")

func (t *Tree) parseKramdownIAL() (ret [][]string) {
	content := t.Context.currentLine[t.Context.nextNonspace:]
	if curlyBracesStart := bytes.Index(content, []byte("{:")); 0 <= curlyBracesStart {
		content = content[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(content, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		content = content[:len(content)-2]
		for {
			valid, remains, attr, name, val := t.parseTagAttr(content)
			if !valid {
				break
			}

			content = remains
			if 1 > len(attr) {
				break
			}

			ret = append(ret, []string{util.BytesToStr(name), util.BytesToStr(val)})
		}
	}
	return
}
