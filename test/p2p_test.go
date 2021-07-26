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
	"testing"

	"github.com/88250/lute"
)

var blockDOM2InlineBlockDOM = []parseTest{

	{"2", "<div data-node-id=\"20210726013400-wvzbzmq\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"block-ref\" data-id=\"20210726013400-u53umzr\" data-anchor=\"\">foo</span></div><div class=\"protyle-attr\"></div></div>", "<span data-type=\"block-ref\" data-id=\"20210726013400-u53umzr\" data-anchor=\"\"></span>"},
	{"1", "<div data-subtype=\"h2\" data-node-id=\"20210716100144-gx162jy\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"u\" data-node-id=\"20210716100144-28v7vkz\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210716100144-d1bczk7\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210716100144-eum0d82\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"u\" data-node-id=\"20210716100144-dzc9rre\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210716100144-ubzfkjs\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210716100144-9e1sop0\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">baz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210716100144-hwxvsqf\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bazz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "foobarbazbazz"},
	{"0", "<div data-subtype=\"h3\" data-node-id=\"20210716095835-fkiy2dh\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h3\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"a\" data-href=\"https://github.com/siyuan-note/siyuan/issues/bar\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<span data-type=\"a\" data-href=\"https://github.com/siyuan-note/siyuan/issues/bar\">foo</span>"},
}

func TestBlockDOM2InlineBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.SetTag(true)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)
	luteEngine.SetIndentCodeBlock(false)
	luteEngine.SetEmojiSite("http://127.0.0.1:6806/stage/protyle/images/emoji")
	luteEngine.SetAutoSpace(true)
	luteEngine.SetParagraphBeginningSpace(true)

	for _, test := range blockDOM2InlineBlockDOM {
		result := luteEngine.BlockDOM2InlineBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}


var blockDOM2StdMd = []parseTest{

	{"2", "<span data-type=\"tag\">foo</span> bar <em>foo</em> bar", "#foo# bar *foo* bar\n"},
	{"1", "foo <code>bar</code> baz", "foo `bar` baz\n"},
	{"0", "foo<u>bar</u>baz~~abc~~xyz", "foo<u>bar</u>baz~~abc~~xyz\n"},
}

func TestBlockDOM2StdMd(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.SetTag(true)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)
	luteEngine.SetIndentCodeBlock(false)
	luteEngine.SetEmojiSite("http://127.0.0.1:6806/stage/protyle/images/emoji")
	luteEngine.SetAutoSpace(true)
	luteEngine.SetParagraphBeginningSpace(true)

	for _, test := range blockDOM2StdMd {
		result := luteEngine.BlockDOM2StdMd(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
