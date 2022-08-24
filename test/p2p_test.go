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

	{"7", "<div data-node-id=\"20220525090743-ueavg67\" data-node-index=\"1\" data-type=\"NodeHTMLBlock\" class=\"render-node\" updated=\"20220525090743\" data-subtype=\"block\"><div class=\"protyle-icons\"><span class=\"protyle-icon protyle-icon--first protyle-action__edit\"><svg><use xlink:href=\"#iconEdit\"></use></svg></span><span class=\"protyle-icon protyle-action__menu protyle-icon--last\"><svg><use xlink:href=\"#iconMore\"></use></svg></span></div><div><protyle-html data-content=\"&lt;testnode&gt;\n          &lt;name Value=&quot;1&quot; /&gt;\n          &lt;value Value=&quot;1&quot; /&gt;\n          &lt;description Value=&quot;1&quot; /&gt;\n          &lt;note Value=&quot;0&quot; /&gt;\n        &lt;/testnode&gt;\"></protyle-html><span style=\"position: absolute\">​</span></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "&lt;testnode&gt;\n          &lt;name Value=&quot;1&quot; /&gt;\n          &lt;value Value=&quot;1&quot; /&gt;\n          &lt;description Value=&quot;1&quot; /&gt;\n          &lt;note Value=&quot;0&quot; /&gt;\n        &lt;/testnode&gt;"},
	{"6", "<testnode>\n          <name Value=\"1\" />\n          <value Value=\"1\" />\n          <description Value=\"1\" />\n          <note Value=\"0\" />\n        </testnode>", "testnode          description\n          note"},
	{"5", "<div data-subtype=\"t\" data-node-id=\"20220325100509-t05nmj6\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\" updated=\"20220325100509\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20220325100509-skwrk4f\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20220325100509-r89zejn\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220325100509\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20220325100509-hhbowb2\" data-type=\"NodeListItem\" class=\"li\" updated=\"20220325100509\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20220325100509-2vhyu9m\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220325100509\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "foobar"},
	{"4", "<div data-node-id=\"20220303220959-l05fsv1\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\" updated=\"20220303220959\" data-eof=\"true\"><div class=\"protyle-action protyle-icons\"><span class=\"protyle-action__language\" contenteditable=\"false\"></span><span class=\"protyle-action__copy b3-tooltips b3-tooltips__nw\" aria-label=\"复制\"><svg><use xlink:href=\"#iconCopy\"></use></svg></span></div><div contenteditable=\"true\" spellcheck=\"false\" class=\"hljs protyle-linenumber\" data-render=\"true\" style=\"white-space: pre-wrap; word-break: break-all; font-variant-ligatures: none;\">var a = 1\n</div><span contenteditable=\"false\" class=\"protyle-linenumber__rows\"><span style=\"height:22px;\"></span></span><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", ""},
	{"3", "a<br>b<br>c", "a<br />b<br />c"},
	{"2", "<div data-node-id=\"20210726013400-wvzbzmq\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"block-ref\" data-subtype=\"d\" data-id=\"20210726013400-u53umzr\">foo</span></div><div class=\"protyle-attr\"></div></div>", "<span data-type=\"block-ref\" data-subtype=\"d\" data-id=\"20210726013400-u53umzr\">foo</span>"},
	{"1", "<div data-subtype=\"h2\" data-node-id=\"20210716100144-gx162jy\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"u\" data-node-id=\"20210716100144-28v7vkz\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210716100144-d1bczk7\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210716100144-eum0d82\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"u\" data-node-id=\"20210716100144-dzc9rre\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210716100144-ubzfkjs\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210716100144-9e1sop0\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">baz</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210716100144-hwxvsqf\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bazz</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "foobarbazbazz"},
	{"0", "<div data-subtype=\"h3\" data-node-id=\"20210716095835-fkiy2dh\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h3\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"a\" data-href=\"https://github.com/siyuan-note/siyuan/issues/bar\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<span data-type=\"a\" data-href=\"https://github.com/siyuan-note/siyuan/issues/bar\">foo</span>"},
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
	luteEngine.SetFileAnnotationRef(true)

	for _, test := range blockDOM2InlineBlockDOM {
		result := luteEngine.BlockDOM2InlineBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

var blockDOM2StdMd = []parseTest{

	{"7", "<span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"&amp;lt;foo&amp;gt;\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\"><span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.5782em;vertical-align:-0.0391em;\"></span><span class=\"mrel\">&lt;</span><span class=\"mspace\" style=\"margin-right:0.2778em;\"></span></span><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span><span class=\"mspace\" style=\"margin-right:0.2778em;\"></span><span class=\"mrel\">&gt;</span></span></span></span></span>", "$<foo>$\n"},
	{"6", "  foo", "  foo\n"},
	{"5", "&nbsp;&nbsp;foo", "  foo\n"},
	{"4", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><code>foo</code></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "* [ ] `foo`\n"},
	{"3", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><strong>foo</strong></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "* [ ] **foo**\n"},
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
	luteEngine.SetFileAnnotationRef(true)

	for _, test := range blockDOM2StdMd {
		result := luteEngine.BlockDOM2StdMd(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

var blockDOM2Content = []parseTest{

	{"6", "foo&lt;&quot;&nbsp;<span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\"></span>&nbsp;<strong style=\"color: var(--b3-font-color8);\">bar</strong>&nbsp;&lt;baz&gt;", "foo<\" foo bar <baz>"},
	{"5", "<div data-subtype=\"h1\" data-node-id=\"20220620223803-e5c7fez\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20220620231839\"><div contenteditable=\"true\" spellcheck=\"false\">foo&lt;\" <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\"><span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span></span></span></span></span> <strong style=\"color: var(--b3-font-color8);\">bar</strong> &lt;baz&gt;<wbr></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "foo<\" foo bar <baz>‸"},
	{"4", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><code>foo</code></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "foo"},
	{"3", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><strong>foo</strong></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "foo"},
	{"2", "<span data-type=\"tag\">foo</span> bar <em>foo</em> bar", "foo bar foo bar"},
	{"1", "foo <code>bar</code> baz", "foo bar baz"},
	{"0", "foo<u>bar</u>baz~~abc~~xyz", "foobarbaz~~abc~~xyz"},
}

func TestBlockDOM2Content(t *testing.T) {
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
	luteEngine.SetFileAnnotationRef(true)

	for _, test := range blockDOM2Content {
		result := luteEngine.BlockDOM2Content(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
