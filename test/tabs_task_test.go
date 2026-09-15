package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
)

func TestTabsTaskRoundTrip(t *testing.T) {
	for _, marker := range []string{" ", "X", "x", "/", "-", "?", "!", "\"", "&", "<"} {
		l := tabsEngine()
		dom := l.Md2BlockDOM("::: tabs\n@tab **Task**\nBody\n:::\n", false)
		dom = strings.Replace(dom, `class="tab-item"`, `class="tab-item" tabs-task="`+html.EscapeAttrVal(marker)+`"`, 1)
		dom = strings.Replace(dom, `class="tabs"`, `class="tabs" tabs-task="true"`, 1)
		for round := 0; round < 3; round++ {
			tree := l.BlockDOM2Tree(dom)
			groups, items := tabNodes(tree.Root)
			if len(groups) != 1 || groups[0].IALAttr("tabs-task") != "true" {
				t.Fatalf("group task lost in AST: %s", dom)
			}
			if len(items) != 1 || items[0].IALAttr("tabs-task") != marker {
				t.Fatalf("marker %q round %d lost in AST: %s", marker, round, dom)
			}
			exported := l.BlockDOM2HTML(dom)
			htmlGroups, htmlItems := tabNodes(l.BlockDOM2Tree(l.HTML2BlockDOM(exported)).Root)
			if len(htmlGroups) != 1 || htmlGroups[0].IALAttr("tabs-task") != "true" {
				t.Fatalf("group task lost in HTML: %s", exported)
			}
			if len(htmlItems) != 1 || htmlItems[0].IALAttr("tabs-task") != marker {
				t.Fatalf("marker %q lost in HTML: %s", marker, exported)
			}
			md := l.BlockDOM2StdMd(dom)
			groups, items = tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
			if len(groups) != 1 || groups[0].IALAttr("tabs-task") != "true" {
				t.Fatalf("group task lost in Markdown: %s", md)
			}
			if len(items) != 1 || items[0].IALAttr("tabs-task") != marker {
				t.Fatalf("marker %q lost in Markdown: %s", marker, md)
			}
			dom = l.Md2BlockDOM(md, false)
		}
	}
}

func TestCustomTaskStatusSpinRoundTrip(t *testing.T) {
	l := tabsEngine()
	l.SetDataTask(true)
	l.SetArbitraryTaskListItemMarker(true)
	l.SetExportNormalizeTaskListMarker(true)
	for _, marker := range []string{" ", "X", "x", "/", "-", "?", "!", "\"", "&", "<"} {
		dom := l.Md2BlockDOM("* [ ] Task\n", false)
		dom = strings.Replace(dom, `data-task=" "`, `data-task="`+html.EscapeAttrVal(marker)+`"`, 1)
		for round := 0; round < 3; round++ {
			dom = l.SpinBlockDOM(dom)
			expected := marker
			if "x" == marker {
				expected = "X"
			}
			if !strings.Contains(dom, `data-task="`+html.EscapeAttrVal(expected)+`"`) {
				t.Fatalf("marker %q lost after editing: %s", marker, dom)
			}
		}
	}
}
