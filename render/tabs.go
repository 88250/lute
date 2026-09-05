package render

import (
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

func (r *FormatRenderer) renderTabs(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		if ast.NodeTabs == node.Type {
			r.WriteString(":::tabs")
		} else {
			r.WriteString(":::tab")
			if "" != node.TabItemTitle {
				r.WriteString(" " + strings.ReplaceAll(strings.ReplaceAll(node.TabItemTitle, "\r", " "), "\n", " "))
			}
		}
	} else {
		r.WriteString(":::")
	}
	r.Newline()
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTabs(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering && ast.NodeTabItem == node.Type && "" != node.TabItemTitle {
		title := calloutInlineTree(node.TabItemTitle, r.ParseOptions)
		parse.TextMarks2Inlines(title)
		r.Write(NewProtyleExportMdRenderer(title, r.Options, r.ParseOptions).Render())
		r.WriteString("\n\n")
	}
	return ast.WalkContinue
}

func (r *ProtyleRenderer) renderTabs(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{}
		class := "tabs"
		if ast.NodeTabItem == node.Type {
			class = "tab-item"
		}
		r.blockNodeAttrs(node, &attrs, class)
		r.Tag("div", attrs, false)
		if ast.NodeTabs == node.Type {
			r.WriteString("<div class=\"tabs-header protyle-action\" contenteditable=\"false\"></div>")
		} else {
			r.WriteString("<div class=\"tab-item-info callout-info\" contenteditable=\"false\"><span class=\"tab-item-title callout-title\" contenteditable=\"true\" spellcheck=\"")
			r.WriteString(strconv.FormatBool(r.Options.Spellcheck))
			r.WriteString("\">")
			r.Write(NewProtyleRenderer(calloutInlineTree(node.TabItemTitle, r.ParseOptions), r.Options, r.ParseOptions).Render())
			r.WriteString("</span></div><div class=\"tab-item-content\">")
		}
	} else {
		if ast.NodeTabItem == node.Type {
			r.WriteString("</div>")
		}
		r.renderIAL(node)
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

// renderTabsHTML 先输出全部内容，只有交互成功初始化后才隐藏未选中的正文。
func (r *BaseRenderer) renderTabsHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		class := "tabs"
		if ast.NodeTabItem == node.Type {
			class = "tab-item"
		}
		attrs := [][]string{{"class", class}, {"data-type", node.Type.String()}}
		if "" != node.ID {
			attrs = append(attrs, []string{"data-node-id", node.ID}, []string{"id", node.ID})
		}
		for _, name := range []string{"tabs-active-id", "tabs-position"} {
			if value := node.IALAttr(name); "" != value {
				attrs = append(attrs, []string{name, value})
			}
		}
		r.Tag("div", attrs, false)
		if ast.NodeTabs == node.Type {
			r.WriteString("<div class=\"tabs-header protyle-action\"></div>")
		} else {
			r.WriteString("<div class=\"tab-item-info callout-info\"><span class=\"tab-item-title callout-title\">")
			r.Write(NewHtmlRenderer(calloutInlineTree(node.TabItemTitle, r.ParseOptions), r.Options, r.ParseOptions).Render())
			r.WriteString("</span></div><div class=\"tab-item-content\">")
		}
	} else {
		if ast.NodeTabItem == node.Type {
			r.WriteString("</div>")
		}
		r.WriteString("</div>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *BaseRenderer) renderTabsDocx(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering && ast.NodeTabItem == node.Type && "" != node.TabItemTitle {
		r.WriteString("<p>")
		r.Write(NewHtmlRenderer(calloutInlineTree(node.TabItemTitle, r.ParseOptions), r.Options, r.ParseOptions).Render())
		r.WriteString("</p>\n")
	}
	return ast.WalkContinue
}
