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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var kitymindJSONRendererTests = []parseTest{

	{"0", "# foo\n{: id=\"20210110005758-m303ovi\"}\n\nbar\n{: id=\"20210110115402-21ltd5v\"}\n\nbaz **bazz**\n{: id=\"20210110115405-17ng22v\"}\n\n* {: id=\"20210110131437-zrwhxvj\"}list\n  * {: id=\"20210110131439-mhegwqy\"}list2\n    * {: id=\"20210110131444-3zxfmg3\"}list3\n      {: id=\"20210110131527-f6tvieh\"}\n    * {: id=\"20210110131453-728ma0a\"}list33\n      {: id=\"20210110131527-dbykywu\"}\n\n      > test\n      > {: id=\"20210110131528-cb1cdfe\"}\n      >\n      {: id=\"20210110131527-k1r271n\"}\n    {: id=\"20210110131444-9bnlizr\"}\n  {: id=\"20210110131443-f1xhkyp\"}\n{: id=\"20210110131435-3en3eha\"}\n\n* {: id=\"20210110131551-duxknch\"}[ ] todo1\n  * {: id=\"20210110131554-lv5luco\"}[ ] **todo**{: style=\"color: rgb(190, 7, 18);\"}2\n  {: id=\"20210110131555-9n8shka\"}\n{: id=\"20210110131547-j84elev\"}\n\n\n{: id=\"20201228004131-bys3g5x\" type=\"doc\"}\n", "{\"root\":{\"data\":{\"text\":\"文档名 TODO\",\"id\":\"20201228004131-bys3g5x\",\"type\":\"NodeDocument\"},\"children\":[{\"data\":{\"text\":\"# foo\",\"id\":\"20210110005758-m303ovi\",\"type\":\"NodeHeading\"},\"children\":[{\"data\":{\"text\":\"bar\",\"id\":\"20210110115402-21ltd5v\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"baz **bazz**\",\"id\":\"20210110115405-17ng22v\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131435-3en3eha\",\"type\":\"NodeList\"},\"children\":[{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131437-zrwhxvj\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"list\",\"id\":\"\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131443-f1xhkyp\",\"type\":\"NodeList\"},\"children\":[{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131439-mhegwqy\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"list2\",\"id\":\"\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131444-9bnlizr\",\"type\":\"NodeList\"},\"children\":[{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131444-3zxfmg3\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"list3\",\"id\":\"20210110131527-f6tvieh\",\"type\":\"NodeParagraph\"},\"children\":[]}]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131453-728ma0a\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"list33\",\"id\":\"20210110131527-dbykywu\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"> test>\",\"id\":\"20210110131527-k1r271n\",\"type\":\"NodeBlockquote\"},\"children\":[{\"data\":{\"text\":\"test\",\"id\":\"20210110131528-cb1cdfe\",\"type\":\"NodeParagraph\"},\"children\":[]}]}]}]}]}]}]}]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131547-j84elev\",\"type\":\"NodeList\"},\"children\":[{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131551-duxknch\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"[ ] todo1\",\"id\":\"\",\"type\":\"NodeParagraph\"},\"children\":[]},{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131555-9n8shka\",\"type\":\"NodeList\"},\"children\":[{\"data\":{\"text\":\"* {: id=20210110...\",\"id\":\"20210110131554-lv5luco\",\"type\":\"NodeListItem\"},\"children\":[{\"data\":{\"text\":\"[ ] **todo**2\",\"id\":\"\",\"type\":\"NodeParagraph\"},\"children\":[]}]}]}]}]}]}]}}"},
}

func TestKityMinderJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.KramdownIAL = true

	ast.Testing = true
	for _, test := range kitymindJSONRendererTests {
		jsonStr := luteEngine.RenderKityMinderJSON(test.from)
		t.Log(jsonStr)
		if test.to != jsonStr {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, jsonStr, test.from)
		}
	}
}
