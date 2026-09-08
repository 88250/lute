package parse

import (
	"fmt"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
)

// TableCellRichOptions 固定富文本存储语法，不受当前编辑器的 Markdown 开关影响。
func TableCellRichOptions() *Options {
	ret := NewOptions()
	ret.DisableTableCellRich = true
	ret.ProtyleWYSIWYG = true
	ret.TextMark = true
	ret.HTMLTag2TextMark = true
	ret.BlockRef = true
	ret.FileAnnotationRef = true
	ret.KramdownBlockIAL = true
	ret.KramdownSpanIAL = true
	ret.Emoji = false
	ret.SuperBlock = true
	ret.CustomBlock = true
	ret.GitConflict = true
	ret.Callout = true
	ret.Tabs = true
	ret.ImgPathAllowSpace = true
	ret.Sup = true
	ret.Sub = true
	ret.Tag = true
	ret.Mark = true
	ret.GFMStrikethrough1 = false
	ret.FullWidthStrikethrough = true
	ret.InlineMathAllowDigitAfterOpenMarker = true
	ret.Footnotes = false
	ret.ToC = false
	ret.HeadingID = false
	ret.Setext = false
	ret.YamlFrontMatter = false
	ret.LinkRef = false
	ret.IndentCodeBlock = false
	ret.ParagraphBeginningSpace = true
	ret.ArbitraryTaskListItemMarker = true
	ret.EnsureListItemParagraph = true
	return ret
}

// ParseTableCellRich 解析受支持的片段；片段中的块身份不进入文档块树。
func ParseTableCellRich(rich *ast.TableCellRich) (tree *Tree, err error) {
	if nil == rich {
		return nil, fmt.Errorf("table cell rich text is missing")
	}
	if err = rich.Validate(); nil != err {
		return
	}
	content, restoreStyles := protectTableCellRichStyles(rich.Content)
	tree = Parse("", []byte(content), TableCellRichOptions())
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}
		if n.IsBlock() {
			switch n.Type {
			case ast.NodeDocument, ast.NodeParagraph, ast.NodeHeading, ast.NodeList, ast.NodeListItem,
				ast.NodeBlockquote, ast.NodeCodeBlock, ast.NodeMathBlock, ast.NodeKramdownBlockIAL:
			default:
				err = fmt.Errorf("unsupported table cell rich text block [%s]", n.Type)
				return ast.WalkStop
			}
		}
		if ast.NodeCodeBlock == n.Type {
			info := strings.Fields(strings.ToLower(string(n.CodeBlockInfo)))
			if len(info) > 0 {
				switch info[0] {
				case "abc", "echarts", "flowchart", "graphviz", "infographic", "mermaid", "mindmap", "plantuml":
					err = fmt.Errorf("unsupported table cell rich text code language [%s]", info[0])
					return ast.WalkStop
				}
			}
		}
		return ast.WalkContinue
	})
	if nil != err {
		return nil, err
	}
	TextMarks2Inlines(tree)
	NestedInlines2FlattedSpansHybrid(tree, false)
	restoreStyles(tree)
	return
}

// CloneTableCellInline 复制行内节点及其属性，避免投影和源树共享可变节点。
func CloneTableCellInline(node *ast.Node) *ast.Node {
	clone := *node
	clone.Parent, clone.Previous, clone.Next, clone.FirstChild, clone.LastChild = nil, nil, nil, nil, nil
	clone.Children = nil
	clone.Tokens = append([]byte(nil), node.Tokens...)
	clone.KramdownIAL = nil
	for _, ial := range node.KramdownIAL {
		clone.KramdownIAL = append(clone.KramdownIAL, append([]string(nil), ial...))
	}
	for child := node.FirstChild; nil != child; child = child.Next {
		clone.AppendChild(CloneTableCellInline(child))
	}
	return &clone
}

// ProjectTableCellRich 将片段转换成行内节点，保留图片、引用、样式以及可读的列表标记。
func ProjectTableCellRich(tree *Tree) *ast.Node {
	ret, _ := ProjectTableCellRichWithSources(tree)
	return ret
}

// ProjectTableCellRichWithSources 同时记录行内节点的来源，供资源和引用改写同步回富文本源。
func ProjectTableCellRichWithSources(tree *Tree) (*ast.Node, map[*ast.Node]*ast.Node) {
	ret := &ast.Node{Type: ast.NodeTableCell}
	sources := map[*ast.Node]*ast.Node{}
	var mapSource func(*ast.Node, *ast.Node)
	mapSource = func(copy, original *ast.Node) {
		sources[copy] = original
		for c, o := copy.FirstChild, original.FirstChild; c != nil && o != nil; c, o = c.Next, o.Next {
			mapSource(c, o)
		}
	}
	appendText := func(text string) {
		for index, line := range strings.Split(text, "\n") {
			if index > 0 {
				ret.AppendChild(&ast.Node{Type: ast.NodeBr, Tokens: []byte("<br />")})
			}
			ret.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(line)})
		}
	}
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}
		switch n.Type {
		case ast.NodeParagraph, ast.NodeHeading, ast.NodeCodeBlock, ast.NodeMathBlock:
		default:
			return ast.WalkContinue
		}
		if nil != ret.FirstChild {
			ret.AppendChild(&ast.Node{Type: ast.NodeBr, Tokens: []byte("<br />")})
		}
		if ast.NodeParagraph == n.Type && n.ParentIs(ast.NodeListItem) {
			first := n.Parent.FirstChild
			if ast.NodeTaskListItemMarker == first.Type {
				first = first.Next
			}
			if first == n {
				marker := "- "
				if n.Parent.ListData != nil && n.Parent.ListData.Typ == 1 {
					marker = fmt.Sprintf("%d. ", n.Parent.ListData.Num)
				} else if task := n.Parent.ChildByType(ast.NodeTaskListItemMarker); nil != task {
					marker = "- [ ] "
					if task.TaskListItemChecked {
						marker = "- [x] "
					}
				}
				appendText(marker)
			}
		}
		if ast.NodeCodeBlock == n.Type {
			if content := n.ChildByType(ast.NodeCodeBlockCode); nil != content {
				code := &ast.Node{Type: ast.NodeTextMark, TextMarkType: "code",
					TextMarkTextContent: html.EscapeHTMLStr(strings.TrimSuffix(content.TokensStr(), "\n"))}
				ret.AppendChild(code)
				sources[code] = content
			}
		} else if ast.NodeMathBlock == n.Type {
			if content := n.ChildByType(ast.NodeMathBlockContent); nil != content {
				math := &ast.Node{Type: ast.NodeTextMark, TextMarkType: "inline-math",
					TextMarkInlineMathContent: html.EscapeHTMLStr(strings.TrimSpace(content.TokensStr()))}
				ret.AppendChild(math)
				sources[math] = content
			}
		} else {
			for child := n.FirstChild; nil != child; child = child.Next {
				if ast.NodeHeadingC8hMarker != child.Type && ast.NodeTaskListItemMarker != child.Type {
					copy := CloneTableCellInline(child)
					ret.AppendChild(copy)
					mapSource(copy, child)
				}
			}
		}
		return ast.WalkSkipChildren
	})
	return ret, sources
}

// ApplyTableCellRichProjection 仅在源成功解析后刷新派生的行内内容。
func ApplyTableCellRichProjection(cell *ast.Node) error {
	tree, err := ParseTableCellRich(cell.TableCellRich)
	if nil != err {
		return err
	}
	projection := ProjectTableCellRich(tree)
	for nil != cell.FirstChild {
		cell.FirstChild.Unlink()
	}
	for nil != projection.FirstChild {
		cell.AppendChild(projection.FirstChild)
	}
	return nil
}

func (tree *Tree) finalParseTableCellRich() {
	if tree.Context.ParseOption.DisableTableCellRich {
		return
	}
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeTableCell != n.Type {
			return ast.WalkContinue
		}
		for _, attribute := range n.KramdownIAL {
			if attribute[0] != ast.TableCellRichIAL {
				continue
			}
			encoded := n.IALAttr(ast.TableCellRichIAL)
			n.TableCellRich = ast.DecodeTableCellRich(encoded)
			n.RemoveIALAttr(ast.TableCellRichIAL)
			for c := n.FirstChild; nil != c; {
				next := c.Next
				if ast.NodeKramdownSpanIAL == c.Type {
					c.Unlink()
				}
				c = next
			}
			_ = ApplyTableCellRichProjection(n)
			break
		}
		return ast.WalkSkipChildren
	})
}
