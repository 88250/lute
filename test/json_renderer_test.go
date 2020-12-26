// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"fmt"
	"testing"

	"github.com/88250/lute"
)

var JSONRendererTests = []parseTest{
	{"测试普通文本","普通文本测试","[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"普通文本测试\"}]}]"},
	{"测试行内代码", "`console.log(\"Hello World\")`", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"CodeSpan\",\"value\":\"console.log(\\\"Hello World\\\")\"}]}]\n"},
	{"测试代码块", "```js\nconsole.log(\"Hello World\")\n```\n", "[{\"type\":\"CodeBlock\",\"value\":\"console.log(\\\"Hello World\\\")\\n\",\"language\":\"js\"}]"},
	{"测试数学块", "$$\na + b = c\n$$\n", "[{\"type\":\"MathBlock\",\"value\":\"a + b = c\"}]"},
	{"测试行内数学公式", "$a + b = c$", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"InlineMath\",\"value\":\"a + b = c\"}]}]\n"},
	{"测试斜体", "*测试斜体*", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Emphasis\",\"children\":[{\"type\":\"Text\",\"value\":\"测试斜体\"}]}]}]\n"},
	{"测试加粗", "**测试粗体**", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Strong\",\"children\":[{\"type\":\"Text\",\"value\":\"测试粗体\"}]}]}]\n"},
	{"测试引用块", "> 测试引用块", "[{\"flag\":\"Blockquote\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"测试引用块\"}]}]}]\n"},
	{"测试标题", "# 一级标题", "[{\"type\":\"Heading\",\"value\":\"h1\",\"children\":[{\"type\":\"Text\",\"value\":\"一级标题\"}]}]\n"},
	{"测试无序列表", "- item1\n- item2\n- item3\n", "[{\"type\":\"List\",\"value\":\"ul\",\"children\":[{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item1\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item2\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item3\"}]}]}]}]\n"},
	{"测试有序列表", "1. item1\n2. item2\n3. item3\n", "[{\"type\":\"List\",\"value\":\"ol\",\"children\":[{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item1\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item2\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Text\",\"value\":\"item3\"}]}]}]}]\n"},
	//{"测试分割线", "---", ""},
	//{"测试软换行", "测试换行\\n", ""},
	{"测试HTML块", "<div>\nHTML块\n</div>\n", "[{\"type\":\"HTMLBlock\",\"value\":\"<div>\\nHTML块\\n</div>\"}]\n"},
	{"测试行内HTML", "<a>行内HTML</a>", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"InlineHTML\",\"value\":\"<a>\"},{\"type\":\"Text\",\"value\":\"行内HTML\"},{\"type\":\"InlineHTML\",\"value\":\"</a>\"}]}]\n"},
	{"测试链接", "[链接文本](链接)", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Link\",\"value\":\"%E9%93%BE%E6%8E%A5\",\"title\":\"链接文本\"}]}]\n"},
	{"测试图片", "![图片文本](图片链接)", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"Image\",\"value\":\"%E5%9B%BE%E7%89%87%E9%93%BE%E6%8E%A5\",\"title\":\"图片文本\"}]}]\n"},
	{"测试删除线", "~~删除线~~", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Strikethrough\",\"children\":[{\"type\":\"Text\",\"value\":\"删除线\"}]}]}]\n"},
	{"测试TaskList", "- [X] item1\n- [ ] item2\n- [X] item3\n", "[{\"type\":\"List\",\"value\":\"ul\",\"children\":[{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"TaskListItemMarker\",\"value\":\"true\"},{\"type\":\"Text\",\"value\":\" item1\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"TaskListItemMarker\",\"value\":\"false\"},{\"type\":\"Text\",\"value\":\" item2\"}]}]},{\"flag\":\"ListItem\",\"children\":[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"TaskListItemMarker\",\"value\":\"true\"},{\"type\":\"Text\",\"value\":\" item3\"}]}]}]}]\n"},
	{"测试表格", "| 表头 | 标题 |\n| --- | --- |\n| item1 | item2 |\n| item3 | item4 |\n", "[{\"flag\":\"Table\",\"children\":[{\"flag\":\"TableHead\",\"children\":[{\"flag\":\"TableRow\",\"children\":[{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"表头\"}]},{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"标题\"}]}]}]},{\"flag\":\"TableRow\",\"children\":[{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"item1\"}]},{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"item2\"}]}]},{\"flag\":\"TableRow\",\"children\":[{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"item3\"}]},{\"type\":\"TableCell\",\"value\":\"left\",\"children\":[{\"type\":\"Text\",\"value\":\"item4\"}]}]}]}]\n"},
	{"测试emoji", ":cn:", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"EmojiUnicode\",\"value\":\"🇨🇳\"}]}]\n"},
	{"测试HTML实体符号", "&copy;", "[{\"flag\":\"Paragraph\",\"children\":[{\"type\":\"HTMLEntity\",\"value\":\"©\"}]}]"},
	{"测试yaml", "---\nyaml测试\n---\n", "[{\"type\":\"YamlFrontMatter\",\"value\":\"yaml测试\"}]"},
	//{"测试块引用", "((id \"内容块引用\")", ""},
	{"测试高亮", "==高亮==", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Mark\",\"children\":[{\"type\":\"Text\",\"value\":\"高亮\"}]}]}]\n"},
	{"测试上标", "^上标^", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Sup\",\"children\":[{\"type\":\"Text\",\"value\":\"上标\"}]}]}]"},
	{"测试下标", "~下标~", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Sub\",\"children\":[{\"type\":\"Text\",\"value\":\"下标\"}]}]}]\n"},
	{"测试内容块查询嵌入", "!{{ SELECT * FROM blocks WHERE content LIKE '%待办%' }}", "[{\"type\":\"BlockQueryEmbed\",\"value\":\"SELECT * FROM blocks WHERE content LIKE \\'%待办%\\'\"},]\n"},
	{"测试内容块嵌入节点", "!((id \"text\"))", "[{\"type\":\"BlockEmbed\",\"value\":\"text\"}]\n"},
	{"测试标签", "#标签测试#", "[{\"flag\":\"Paragraph\",\"children\":[{\"flag\":\"Tag\",\"children\":[{\"type\":\"Text\",\"value\":\"标签测试\"}]}]}]\n"},
}

func TestJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.Sup = true
	luteEngine.Sub = true
	luteEngine.Mark = true
 	luteEngine.KramdownIAL = false
 	luteEngine.Footnotes = true
 	luteEngine.BlockRef = true
 	luteEngine.Tag = true
 	luteEngine.SoftBreak2HardBreak = true

	for _, test := range JSONRendererTests {
		jsonStr := luteEngine.RenderJSON(test.from)
		fmt.Println(jsonStr)
	//	if test.to != jsonStr {
	//		t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, jsonStr, test.from)
	//	}
	}
}
