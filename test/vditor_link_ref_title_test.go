package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

func TestVditorLinkRefTitleRoundTrip(t *testing.T) {
	for _, mode := range []string{"ir", "wysiwyg", "sv"} {
		for _, fixture := range []struct {
			name     string
			markdown string
		}{
			{"used", "[text][ref]\n\n[ref]: /target \"Reference title\"\n"},
			{"unused", "[ref]: /target \"Unused title\"\n"},
			{"image", "![alt][ref]\n\n[ref]: /image.png \"Image title\"\n"},
			{"mixed", "[ref]\n\n[ref]: /used \"Used title\"\n[unused]: /other \"Unused title\"\n[empty]: /empty \"\"\n[untitled]: /plain\n"},
			{"quotes", "[ref]: /target 'A \"quoted\" title'\n"},
			{"escaped", "[ref]: /target \"Escaped \\\"quote\\\" and \\\\ slash\"\n"},
			{"entities", "[ref]: /target \"A &amp; B &quot;quoted&quot;\"\n"},
			{"html", "[ref]: /target \"<img src=x onerror=alert(1)> & text\"\n"},
			{"multiline", "[ref]: /target (First line\nsecond line)\n"},
			{"duplicate", "[ref]: /first \"First title\"\n[ref]: /second \"Second title\"\n"},
		} {
			t.Run(mode+"/"+fixture.name, func(t *testing.T) {
				engine := lute.New()
				var toDOM, toMarkdown, spin func(string) string
				switch mode {
				case "ir":
					engine.SetVditorIR(true)
					toDOM, toMarkdown, spin = engine.Md2VditorIRDOM, engine.VditorIRDOM2Md, engine.SpinVditorIRDOM
				case "wysiwyg":
					engine.SetVditorWYSIWYG(true)
					toDOM, toMarkdown, spin = engine.Md2VditorDOM, engine.VditorDOM2Md, engine.SpinVditorDOM
				case "sv":
					engine.SetVditorSV(true)
					toDOM = engine.Md2VditorSVDOM
					toMarkdown = func(dom string) string { return util.DomText(util.ParseHTML(dom)) }
					spin = func(dom string) string { return engine.SpinVditorSVDOM(toMarkdown(dom)) }
				}
				want := referenceDefinitionHTML(fixture.markdown)
				if len(want) == 0 {
					t.Fatal("fixture has no reference definition")
				}
				dom := toDOM(fixture.markdown)
				for round := 0; round < 3; round++ {
					if fixture.name == "html" && strings.Contains(dom, "<img") {
						t.Fatalf("reference title became HTML: %s", dom)
					}
					markdown := toMarkdown(dom)
					if got := referenceDefinitionHTML(markdown); !reflect.DeepEqual(got, want) {
						t.Fatalf("round %d definitions changed\nwant: %#v\ngot: %#v\nMarkdown: %q\nDOM: %s", round, want, got, markdown, dom)
					}
					dom = spin(dom)
				}
			})
		}
	}
}

// Render each definition through a reference link to compare its decoded target
// and title, including unused definitions and different valid title delimiters.
func referenceDefinitionHTML(markdown string) []vditorReferenceDefinition {
	engine := lute.New()
	tree := parse.Parse("", []byte(markdown), engine.ParseOptions)
	var definitions []vditorReferenceDefinition
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && node.Type == ast.NodeLinkRefDef {
			label := string(node.Tokens)
			definition := vditorReferenceDefinition{
				label:       label,
				destination: string(node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens),
				html:        engine.MarkdownStr("", "[reference]["+label+"]\n\n"+markdown),
			}
			if title := node.FirstChild.ChildByType(ast.NodeLinkTitle); nil != title {
				definition.hasTitle = true
				definition.title = string(html.UnescapeBytes(title.Tokens))
			}
			definitions = append(definitions, definition)
		}
		return ast.WalkContinue
	})
	return definitions
}

type vditorReferenceDefinition struct {
	label, destination, title, html string
	hasTitle                        bool
}

func TestVditorInlineLinkAndImageTitles(t *testing.T) {
	markdown := "[link](/link \"Link title\") ![image](/image.png \"Image title\")\n"
	engine := lute.New()
	want := engine.MarkdownStr("", markdown)
	for mode, got := range map[string]string{
		"ir":      engine.VditorIRDOM2Md(engine.Md2VditorIRDOM(markdown)),
		"wysiwyg": engine.VditorDOM2Md(engine.Md2VditorDOM(markdown)),
		"sv":      util.DomText(util.ParseHTML(engine.Md2VditorSVDOM(markdown))),
	} {
		if html := engine.MarkdownStr("", got); html != want {
			t.Fatalf("%s changed inline link/image titles\nwant: %q\ngot: %q", mode, want, html)
		}
	}
}
