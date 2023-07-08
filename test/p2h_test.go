// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var blockDOM2HTMLTest = []parseTest{

	{"1", "<div data-node-id=\"20230708221830-dbk6i2j\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20230708221825-6cebbvv\" data-type=\"NodeListItem\" class=\"li protyle-wysiwyg--select\" updated=\"20230708221901\" custom-t=\"{&amp;quot;11&amp;quot;}\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20230708221825-zfrqic3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20230708221825\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<ul id=\"20230708221830-dbk6i2j\" updated=\"20230708221830\">\n<li id=\"20230708221825-6cebbvv\" updated=\"20230708221901\" custom-t=\"&#123;&quot;11&quot;&#125;\">foo</li>\n</ul>\n"},
	{"0", "foo <span data-type=\"code\">​bar</span>​ baz", "<p id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\">foo <span data-type=\"code\">bar</span>\u200b baz</p>\n"},
}

func TestBlockDOM2HTML(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.ParseOptions.GitConflict = true
	luteEngine.ParseOptions.LinkRef = false
	luteEngine.SetAutoSpace(true)
	luteEngine.SetParagraphBeginningSpace(true)
	luteEngine.SetFileAnnotationRef(true)
	luteEngine.SetImgPathAllowSpace(true)
	luteEngine.SetSanitize(true)
	luteEngine.SetTextMark(true)
	luteEngine.SetTag(true)
	luteEngine.SetTextMark(true)
	luteEngine.SetHTMLTag2TextMark(true)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetMark(true)

	ast.Testing = true
	for _, test := range blockDOM2HTMLTest {
		md := luteEngine.BlockDOM2HTML(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
