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

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var mathTests = []parseTest{
	{"8", "lu$a^2 + b^2 = \\color{red}c^2$1te", "<p>lu$a^2 + b^2 = \\color{red}c^2$1te</p>\n"},
	{"7", "lu$1a^2 + b^2 = \\color{red}c^2$te", "<p>lu$1a^2 + b^2 = \\color{red}c^2$te</p>\n"},
	{"6", "lu$a^2 + b^2 = \\color{red}c^2$te$a^2$m", "<p>lu<span class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</span>te<span class=\"vditor-math\">a^2</span>m</p>\n"},
	{"5", "lu$a^2 + b^2 = \\color{red}c^2$te", "<p>lu<span class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</span>te</p>\n"},
	{"4", "lu$$a^2 + b^2 = \\color{red}c^2$$te", "<p>lu\n<div class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</div>\nte</p>\n"},
	{"3", "$$\na^2 + b^2 = \\color{red}c^2\n$$", "<div class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</div>\n"},
	{"2", "| $a^2 + b^2 = \\color{red}c^2$ | bar |\n| --- | --- |\n| baz | bim |\n", "<table>\n<thead>\n<tr>\n<th><span class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</span></th>\n<th>bar</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>baz</td>\n<td>bim</td>\n</tr>\n</tbody>\n</table>\n"},
	{"1", "$a^2 + b^2 = \\color{red}c^2$", "<p><span class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</span></p>\n"},
	{"0", "$$a^2 + b^2 = \\color{red}c^2$$", "<div class=\"vditor-math\">a^2 + b^2 = \\color{red}c^2</div>\n"},
}

func TestMath(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range mathTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
