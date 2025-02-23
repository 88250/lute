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
	"github.com/88250/lute/ast"
)

var inlineMd2BlockDOM = []parseTest{

	{"4", "<sub>foo</sub>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"sub\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"3", "<sup>foo</sup>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"sup\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"2", "<kbd>foo</kbd>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">\u200b<span data-type=\"kbd\">\u200bfoo</span>\u200b</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"1", "<span data-type=\"strong\">foo</span>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"strong\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"0", "~**foo**~\u200b~bar~\n", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"sub strong\">foo</span>\u200b<span data-type=\"sub\">bar</span>\n</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestInlineMd2BlockDOM(t *testing.T) {
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
	luteEngine.SetTextMark(true)
	luteEngine.SetHTMLTag2TextMark(true)

	ast.Testing = true
	for _, test := range inlineMd2BlockDOM {
		result := luteEngine.InlineMd2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

var blockDOM2InlineBlockDOM = []parseTest{

	{"9", "&gt;2&gt;1", "&gt;2&gt;1"},
	{"8", "<div data-node-id=\"20231108094403-j76zf5u\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\" updated=\"20231108094403\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /><col /></colgroup><thead><tr><th align=\"right\">1</th><th align=\"right\">2</th></tr></thead><tbody><tr><td align=\"right\">3</td><td align=\"right\">4</td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "1 2 3 4"},
	{"7", "<div data-node-id=\"20220525090743-ueavg67\" data-node-index=\"1\" data-type=\"NodeHTMLBlock\" class=\"render-node\" updated=\"20220525090743\" data-subtype=\"block\"><div class=\"protyle-icons\"><span class=\"protyle-icon protyle-icon--first protyle-action__edit\"><svg><use xlink:href=\"#iconEdit\"></use></svg></span><span class=\"protyle-icon protyle-action__menu protyle-icon--last\"><svg><use xlink:href=\"#iconMore\"></use></svg></span></div><div><protyle-html data-content=\"&lt;testnode&gt;\n          &lt;name Value=&quot;1&quot; /&gt;\n          &lt;value Value=&quot;1&quot; /&gt;\n          &lt;description Value=&quot;1&quot; /&gt;\n          &lt;note Value=&quot;0&quot; /&gt;\n        &lt;/testnode&gt;\"></protyle-html><span style=\"position: absolute\">​</span></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "&lt;testnode&gt;\n          &lt;name Value=&quot;1&quot; /&gt;\n          &lt;value Value=&quot;1&quot; /&gt;\n          &lt;description Value=&quot;1&quot; /&gt;\n          &lt;note Value=&quot;0&quot; /&gt;\n        &lt;/testnode&gt;"},
	{"6", "<testnode>\n          <name Value=\"1\" />\n          <value Value=\"1\" />\n          <description Value=\"1\" />\n          <note Value=\"0\" />\n        </testnode>", "&lt;testnode&gt;          &lt;description&gt;\n          &lt;note&gt;"},
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

	{"17", "<span data-type=\"a\" data-href=\"bar baz\">foo</span>", "[foo](bar%20baz)\n"},
	{"16", "<div data-node-id=\"20231229113324-zgenpx4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p protyle-wysiwyg--select\" updated=\"20231229113524\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\\\\\nbar\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\"><span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span></span><span class=\"mspace newline\"></span><span class=\"base\"><span class=\"strut\" style=\"height:0.6944em;\"></span><span class=\"mord mathnormal\">ba</span><span class=\"mord mathnormal\" style=\"margin-right:0.02778em;\">r</span></span></span></span></span>​</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "$foo\\\\ bar$\n"},
	{"15", "foo <span data-type=\"em inline-math\" data-subtype=\"math\" data-content=\"bar\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\"><span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.6595em;\"></span><span class=\"mord mathnormal\">bar</span></span></span></span></span> baz", "foo *$bar$* baz\n"},
	{"14", "foo<span data-type=\"strong\">bar </span>bar", "foo**bar** bar\n"},
	{"13", "<div data-node-id=\"20231028232041-5dhtaps\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231028234123\"><div contenteditable=\"true\" spellcheck=\"false\">foo<span data-type=\"strong\">.bar</span></div><div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div>", "foo **.bar**\n"},
	{"12", "<div data-node-id=\"20231028232041-5dhtaps\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231028234123\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"strong\">foo.</span>bar</div><div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div>", "**foo.** bar\n"},
	{"11", "fo<span data-type=\"text\" id=\"\" style=\"color: var(--b3-font-color8);\">o </span><span data-type=\"inline-math text\" data-subtype=\"math\" data-content=\"2>1\" contenteditable=\"false\" class=\"render-node\" id=\"\" style=\"color: var(--b3-font-color8);\" data-render=\"true\"><span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.6835em;vertical-align:-0.0391em;\"></span><span class=\"mord\">2</span><span class=\"mspace\" style=\"margin-right:0.2778em;\"></span><span class=\"mrel\">&gt;</span><span class=\"mspace\" style=\"margin-right:0.2778em;\"></span></span><span class=\"base\"><span class=\"strut\" style=\"height:0.6444em;\"></span><span class=\"mord\">1</span></span></span></span>​</span><span data-type=\"text\" id=\"\" style=\"color: var(--b3-font-color8);\"> b</span>az", "fo<span data-type=\"text\" style=\"color: var(--b3-font-color8);\">o </span>$2>1$<span data-type=\"text\" style=\"color: var(--b3-font-color8);\"> b</span>az\n"},
	{"10", "<span data-type=\"kbd\">\u200bfoo</span>", "<kbd>foo</kbd>\n"},
	{"9", "<span data-type=\"sub strong\">foo</span><span data-type=\"sub\">bar</span>", "<sub>**foo**</sub>\u200b<sub>bar</sub>\n"},
	{"8", "<span data-type=\"strong em\">foo</span>", "***foo***\n"},
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

var blockDOM2Md = []parseTest{

	{"0", "<div data-node-id=\"20220922151247-vp1f2n4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220922153740\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"block-ref\" data-subtype=\"d\" data-id=\"20220922151244-p6ask52\">foo</span> bar </div><div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div>", "((20220922151244-p6ask52 'foo')) bar \n{: id=\"20220922151247-vp1f2n4\" updated=\"20220922153740\"}\n"},
}

func TestBlockDOM2Md(t *testing.T) {
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

	for _, test := range blockDOM2Md {
		result := luteEngine.BlockDOM2Md(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

var blockDOM2Content = []parseTest{

	{"11", "<div data-subtype=\"u\" data-node-id=\"20230427113011-68v0j4h\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20230427113011\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20230427113011-3xijm44\" data-type=\"NodeListItem\" class=\"li\" updated=\"20230427113011\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20230427113011-7o9ligu\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20230427113011\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20230427113011-zyyjwgn\" data-type=\"NodeListItem\" class=\"li\" updated=\"20230427113011\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20230427113011-as0k13v\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20230427113011\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "foo\nbar"},
	{"10", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">\u200b<span data-type=\"kbd\">\u200bfoo</span>\u200b</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "\u200bfoo\u200b"},
	{"9", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"strong\">foo</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "foo"},
	{"8", "<div data-node-id=\"20221012153945-e1aclg3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20221012153947\"><div contenteditable=\"true\" spellcheck=\"false\">&lt;span&gt;</div><div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div>", "<span>"},
	{"7", "<div data-node-id=\"20221011094818-dg6ktfw\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"em\">foo</span><span data-type=\"em inline-math\" data-subtype=\"math\" data-content=\"bar\" contenteditable=\"false\" class=\"render-node\"></span><span data-type=\"em\">baz</span></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "foobarbaz"},
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

var blockDOM2EscapeMarkerContent = []parseTest{

	{"1", "<div data-node-id=\"20230517211428-a5q3t4e\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20230517211428\"><div contenteditable=\"true\" spellcheck=\"false\">foo <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"bar^2\" contenteditable=\"false\" class=\"render-node\"></span> baz^3</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "foo bar\\^2 baz\\^3"},
	{"0", "<div data-node-id=\"20230517201616-1i9a98t\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20230517201616\"><div contenteditable=\"true\" spellcheck=\"false\">Since <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"0 \\leqq s-r&lt;d\" contenteditable=\"false\" class=\"render-node\"></span> we must have <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"s-r=0\" contenteditable=\"false\" class=\"render-node\"></span>. The cyclic subgroup generated by <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"a\" contenteditable=\"false\" class=\"render-node\"></span> has order <span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"d\" contenteditable=\"false\" class=\"render-node\"></span>.</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "Since 0 \\\\leqq s-r\\<d we must have s-r\\=0. The cyclic subgroup generated by a has order d."},
}

func TestBlockDOM2EscapeMarkerContent(t *testing.T) {
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

	for _, test := range blockDOM2EscapeMarkerContent {
		result := luteEngine.BlockDOM2EscapeMarkerContent(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
