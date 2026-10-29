package test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestDocumentTitleAllInlineSources(t *testing.T) {
	engine := newTableRichLute()
	paragraph := &ast.Node{Type: ast.NodeParagraph, ID: "20261025000000-abcdefg"}
	paragraph.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte("literal ** / ／ & 中文")})
	for _, mark := range []string{"text", "strong", "em", "u", "s", "mark", "sup", "sub", "kbd", "code", "tag",
		"a", "block-ref", "inline-math", "inline-memo", "file-annotation-ref", "strong em mark"} {
		node := &ast.Node{Type: ast.NodeTextMark, TextMarkType: mark, TextMarkTextContent: "visible " + mark}
		switch mark {
		case "a":
			node.TextMarkAHref, node.TextMarkATitle = "https://example.com/a?b=1&c=2", "link & title"
		case "block-ref":
			node.TextMarkBlockRefID, node.TextMarkBlockRefSubtype = "20261025000001-abcdefg", "d"
		case "inline-math":
			node.TextMarkInlineMathContent, node.TextMarkTextContent = `x^2 + \frac{a}{b}`, ""
		case "inline-memo":
			node.TextMarkInlineMemoContent = "hidden note\nwith \"quotes\" and <markup>"
		case "file-annotation-ref":
			node.TextMarkFileAnnotationRefID = "assets/paper.pdf/20261025000002-abcdefg"
		case "strong em mark":
			node.TextMarkFlashcardOcclusionID = "20261025000003-abcdefg"
		}
		paragraph.AppendChild(node)
	}
	styled := &ast.Node{Type: ast.NodeTextMark, TextMarkType: "strong", TextMarkTextContent: "styled"}
	styled.SetIALAttr("style", "color: red; font-family: 'A & B';")
	paragraph.AppendChild(styled)
	paragraph.AppendChild(&ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: []byte(`{: style="color: red; font-family: 'A & B';"}`)})
	imageTree := parse.Parse("", []byte(`![alt](assets/image.png "image title"){: style="width: 123px;"}`), engine.ParseOptions)
	image := imageTree.Root.FirstChild.FirstChild
	paragraph.AppendChild(image)
	paragraph.AppendChild(&ast.Node{Type: ast.NodeSoftBreak})
	paragraph.AppendChild(&ast.Node{Type: ast.NodeHardBreak})
	paragraph.AppendChild(&ast.Node{Type: ast.NodeBr, Tokens: []byte("<br />")})
	backslash := &ast.Node{Type: ast.NodeBackslash}
	backslash.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: []byte("*")})
	paragraph.AppendChild(backslash)
	paragraph.AppendChild(&ast.Node{Type: ast.NodeHTMLEntity, Tokens: []byte("©"), HtmlEntityTokens: []byte("&copy;")})
	emoji := &ast.Node{Type: ast.NodeEmoji}
	emojiUnicode := &ast.Node{Type: ast.NodeEmojiUnicode, Tokens: []byte("😄")}
	emojiUnicode.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: []byte(":smile:")})
	emoji.AppendChild(emojiUnicode)
	paragraph.AppendChild(emoji)

	source, err := render.NewDocumentTitle(paragraph)
	if nil != err {
		t.Fatal(err)
	}
	if paragraph.ID == "" || styled.IALAttr("style") != "color: red; font-family: 'A & B';" || backslash.Parent != paragraph {
		t.Fatal("source construction modified the input nodes")
	}
	for iteration := 0; iteration < 3; iteration++ {
		restored, parseErr := parse.ParseDocumentTitle(source)
		if nil != parseErr {
			t.Fatal(parseErr)
		}
		if restored.ID != "" || restored.Parent != nil || restored.ChildByType(ast.NodeBackslash) == nil {
			t.Fatal("title acquired a block identity or lost an escape")
		}
		next, encodeErr := render.NewDocumentTitle(restored)
		if nil != encodeErr || next.Content != source.Content {
			t.Fatalf("title source round trip changed on iteration %d: %v\n%s\n%s", iteration, encodeErr, source.Content, next.Content)
		}
		source = next
	}
	for _, expected := range []string{"TextMarkInlineMemoContent", "TextMarkInlineMathContent", "TextMarkFileAnnotationRefID",
		"TextMarkFlashcardOcclusionID", "image title", "123px", "NodeBackslashContent", "HtmlEntityTokens"} {
		if !strings.Contains(source.Content, expected) {
			t.Errorf("source lost %s", expected)
		}
	}
	copyA, _ := parse.ParseDocumentTitle(source)
	copyB, _ := parse.ParseDocumentTitle(source)
	copyA.FirstChild.Tokens[0] = 'X'
	if string(copyA.FirstChild.Tokens) == string(copyB.FirstChild.Tokens) || !strings.Contains(source.Content, "literal") {
		t.Fatal("temporary title projections share mutable content")
	}
}

func TestDocumentTitleEmptyImageAlt(t *testing.T) {
	engine := newTableRichLute()
	for _, markdown := range []string{"# ![](assets/test.png)", "# ![alt]()", "# ![]()"} {
		tree := parse.Parse("", []byte(markdown), engine.ParseOptions)
		source, err := render.NewDocumentTitle(tree.Root.FirstChild)
		if err != nil {
			t.Fatalf("%q: %v", markdown, err)
		}
		if _, err = parse.ParseDocumentTitle(source); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDocumentTitleHeadingAndNativeDOM(t *testing.T) {
	engine := newTableRichLute()
	for _, markdown := range []string{
		`## **bold** *em* ~~strike~~ ==mark== ^sup^ ~sub~ [link](https://example.com "title")`,
		"## `code` #tag# $a^2$ ((20261025000000-abcdefg 'dynamic'))",
		`## ![alt](assets/image.png "title"){: style="width: 50%;"} \*literal\*`,
	} {
		t.Run(markdown, func(t *testing.T) {
			tree := parse.Parse("", []byte(markdown), engine.ParseOptions)
			heading := tree.Root.FirstChild
			source, err := render.NewDocumentTitle(heading)
			if nil != err {
				t.Fatal(err)
			}
			if strings.Contains(source.Content, "NodeHeadingC8hMarker") || heading.Type != ast.NodeHeading {
				t.Fatal("heading marker entered source or input heading was changed")
			}
			paragraph, _ := parse.ParseDocumentTitle(source)
			paragraph.ID = ast.NewNodeID()
			paragraph.SetIALAttr("id", paragraph.ID)
			projection := &parse.Tree{Root: &ast.Node{Type: ast.NodeDocument}}
			projection.Root.AppendChild(paragraph)
			dom := engine.Tree2BlockDOM(projection, engine.RenderOptions, engine.ParseOptions)
			restored := engine.BlockDOM2Tree(dom).Root.FirstChild
			if _, err = render.NewDocumentTitle(restored); nil != err {
				t.Fatalf("native editor DOM cannot produce title source: %v\n%s\n%s", err, dom,
					render.NewJSONRenderer(engine.BlockDOM2Tree(dom), engine.RenderOptions, engine.ParseOptions).Render())
			}
		})
	}
}

func TestDocumentTitleJSONPersistenceAndLegacy(t *testing.T) {
	engine := newTableRichLute()
	tree := parse.Parse("", []byte("body"), engine.ParseOptions)
	legacy := string(render.NewJSONRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
	if strings.Contains(legacy, "DocumentTitle") {
		t.Fatal("ordinary document acquired a title source")
	}
	title := &ast.Node{Type: ast.NodeParagraph}
	title.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(strings.Repeat("long title ", 100))})
	source, err := render.NewDocumentTitle(title)
	if nil != err {
		t.Fatal(err)
	}
	tree.Root.DocumentTitle = source
	data := render.NewJSONRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
	var restored ast.Node
	if json.Unmarshal(data, &restored) != nil || !reflect.DeepEqual(source, restored.DocumentTitle) || len(restored.Children) != 1 {
		t.Fatalf("root title source or body changed: %s", data)
	}
	empty, err := render.NewDocumentTitle(&ast.Node{Type: ast.NodeParagraph})
	if nil != err || empty.Content != "[]" {
		t.Fatalf("empty title is not an explicit empty array: %#v, %v", empty, err)
	}
}

func TestDocumentTitleRejectsInvalidSources(t *testing.T) {
	for _, content := range []string{
		`null`, `{}`, `[null]`, `[] {}`,
		`[{"Type":"NodeText","Type":"NodeParagraph","Data":"private"}]`,
		`[{"Type":"NodeText","type":"NodeParagraph"}]`,
		`[{"Type":"NodeText","Future":"private"}]`,
		`[{"Type":"NodeText","ID":"20261025000000-abcdefg"}]`,
		`[{"Type":"NodeText","Properties":{"id":""}}]`,
		`[{"Type":"NodeText","Properties":{"style":"a","style":"b"}}]`,
		`[{"Type":"NodeText","Data":"\u0000"}]`,
		`[{"Type":"NodeText","Data":null}]`, `[{"Type":"NodeText","Properties":{"style":null}}]`,
		`[{"Type":"NodeText","Data":"\ud800"}]`, `[{"Type":"NodeText","Data":"\udc00"}]`,
		`[{"Type":"NodeText","Data":"\ud800\u0020"}]`,
		`[{"Type":"FutureNode"}]`, `[{"Type":"NodeTableCell"}]`, `[{"Type":"NodeParagraph"}]`,
		`[{"Type":"NodeInlineHTML","Data":"<span>"}]`,
		`[{"Type":"NodeLinkDest","Data":"private"}]`,
		`[{"Type":"NodeText","Children":[{"Type":"NodeText","Data":"nested"}]}]`,
		`[{"Type":"NodeImage","Children":[{"Type":"NodeLinkDest","Data":"assets/a.png"}]}]`,
		`[{"Type":"NodeBackslash"}]`,
		`[{"Type":"NodeTextMark","TextMarkType":"future","TextMarkTextContent":"private"}]`,
		`[{"Type":"NodeTextMark","TextMarkType":"block-ref","TextMarkBlockRefID":"private"}]`,
		"[" + strings.Repeat(`{"Type":"NodeText"},`, ast.DocumentTitleMaxNodes) + `{"Type":"NodeText"}]`,
		`[{"Type":"NodeText","Data":"` + strings.Repeat("x", ast.DocumentTitleMaxBytes) + `"}]`,
	} {
		source := &ast.DocumentTitle{Spec: 1, Format: "inline-ast", Content: content}
		if _, err := parse.ParseDocumentTitle(source); nil == err || strings.Contains(err.Error(), "private") {
			t.Fatalf("invalid source was accepted or exposed: %.150s", content)
		}
		if source.Content != content {
			t.Fatal("rejected source was modified")
		}
	}
	for _, data := range []string{
		`null`, `{}`, `{"spec":1,"format":"inline-ast"}`, `{"spec":1,"format":"inline-ast","content":null}`,
		`{"spec":1,"format":"inline-ast","content":"[]","extra":true}`,
		`{"spec":2,"spec":1,"format":"inline-ast","content":"[]"}`,
		`{"spec":1,"format":"future","content":"[]"}`, `{"spec":2,"format":"inline-ast","content":"[]"}`,
	} {
		if _, err := ast.ParseDocumentTitleJSON([]byte(data)); nil == err {
			t.Errorf("invalid envelope accepted: %s", data)
		}
	}
	if _, err := ast.ParseDocumentTitleJSON([]byte(`{"spec":1,"format":"inline-ast","content":"[]"}`)); nil != err {
		t.Fatal(err)
	}
	for _, content := range []string{`[{"Type":"NodeText","Data":"\ud83d\ude04"}]`, `[{"Type":"NodeText","Data":"\\ud800"}]`} {
		if _, err := parse.ParseDocumentTitle(&ast.DocumentTitle{Spec: 1, Format: "inline-ast", Content: content}); err != nil {
			t.Fatalf("valid Unicode or literal escape rejected: %s", content)
		}
	}
}
