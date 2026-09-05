package render

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/parse"
)

func (r *FormatRenderer) renderTabs(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if ast.NodeTabs == node.Type {
		r.WriteString(strings.Repeat(":", r.tabsFenceLength(node)))
		if entering {
			r.WriteString(" tabs")
		}
		r.Newline()
		return ast.WalkContinue
	}
	if !entering {
		return ast.WalkContinue
	}

	r.WriteString("@tab")
	if "" != node.ID && node.ParentIs(ast.NodeTabs) && node.ID == node.Parent.IALAttr("tabs-active-id") {
		r.WriteString(":active")
	}
	if "" != node.TabItemTitle {
		r.WriteString(" " + strings.ReplaceAll(strings.ReplaceAll(node.TabItemTitle, "\r", " "), "\n", " "))
	}
	r.Newline()
	if r.Options.KramdownBlockIAL {
		if 0 < len(node.KramdownIAL) {
			r.Write(parse.IAL2Tokens(node.KramdownIAL))
		} else if nil != node.Next && ast.NodeKramdownBlockIAL == node.Next.Type {
			r.Write(node.Next.Tokens)
		}
		r.Newline()
	}
	r.WriteByte('\n')
	return ast.WalkContinue
}

// tabsFenceLength 根据后代页签层数生成围栏，确保每个外层围栏都长于内层。
func (r *FormatRenderer) tabsFenceLength(node *ast.Node) int {
	if nil == r.tabsFenceLengths {
		r.tabsFenceLengths = map[*ast.Node]int{}
		depths := map[*ast.Node]int{}
		ast.Walk(r.Tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering {
				return ast.WalkContinue
			}
			depth := depths[n]
			if ast.NodeTabs == n.Type {
				depth++
				r.tabsFenceLengths[n] = depth + 2
			}
			if nil != n.Parent && depths[n.Parent] < depth {
				depths[n.Parent] = depth
			}
			return ast.WalkContinue
		})
	}
	if length := r.tabsFenceLengths[node]; 3 <= length {
		return length
	}
	return 3
}

// escapeTabsParagraphMarkers 保留段落中的字面标记，标准导出也保护平铺正文和标题。
func escapeTabsParagraphMarkers(node *ast.Node, tokens []byte, withinTabsOnly bool) []byte {
	if !node.ParentIs(ast.NodeParagraph) && (withinTabsOnly || !node.ParentIs(ast.NodeDocument)) {
		return tokens
	}
	if withinTabsOnly {
		var tabs *ast.Node
		for parent := node.Parent.Parent; nil != parent; parent = parent.Parent {
			if ast.NodeTabs == parent.Type {
				tabs = parent
				break
			}
		}
		if nil == tabs {
			return tokens
		}
	}

	lineStart := nil == node.Previous || ast.NodeSoftBreak == node.Previous.Type ||
		ast.NodeHardBreak == node.Previous.Type
	lines := bytes.Split(tokens, []byte{'\n'})
	changed := false
	for i, line := range lines {
		if (0 < i || lineStart) && isTabsParagraphMarker(line) {
			leading := len(line) - len(bytes.TrimLeft(line, " \t"))
			escaped := append([]byte{}, line[:leading]...)
			escaped = append(escaped, '\\')
			lines[i] = append(escaped, line[leading:]...)
			changed = true
		}
	}
	if !changed {
		return tokens
	}
	return bytes.Join(lines, []byte{'\n'})
}

func isTabsParagraphMarker(line []byte) bool {
	line = bytes.TrimSpace(line)
	line = bytes.TrimPrefix(line, editor.CaretTokens)
	for _, marker := range []string{"@tab", "@tab:active"} {
		if bytes.Equal(line, []byte(marker)) || bytes.HasPrefix(line, []byte(marker+" ")) ||
			bytes.HasPrefix(line, []byte(marker+"\t")) {
			return true
		}
	}
	count := 0
	for count < len(line) && ':' == line[count] {
		count++
	}
	if count < 3 {
		return false
	}
	if count == len(line) {
		return true
	}
	return (' ' == line[count] || '\t' == line[count]) && "tabs" == string(bytes.TrimSpace(line[count:]))
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
