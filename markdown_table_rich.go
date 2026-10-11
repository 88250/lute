package lute

import (
	"io"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

// NormalizeMarkdownTableCells 仅用于外部 Markdown 导入，将明确的块结构还原为已有的单元格富文本源。
// 已有载荷不重新解释；不支持或未闭合的结构保留原始行内内容。
func NormalizeMarkdownTableCells(tree *parse.Tree) {
	if tree == nil || tree.Context.ParseOption.DisableTableCellRich {
		return
	}
	defer func() { tree.Context.TableCellSources = nil }()
	ast.Walk(tree.Root, func(cell *ast.Node, entering bool) ast.WalkStatus {
		if !entering || cell.Type != ast.NodeTableCell || cell.TableCellRich != nil {
			return ast.WalkContinue
		}
		source := strings.TrimSpace(tree.Context.TableCellSources[cell])
		var content string
		if isMarkdownTableCellBlockHTML(source) {
			content = markdownTableCellHTMLContent(source)
		} else if isMarkdownTableCellTask(source) {
			content = markdownTableCellBreaks(source)
		} else {
			return ast.WalkSkipChildren
		}
		rich := &ast.TableCellRich{Spec: ast.TableCellRichSpec, Format: ast.TableCellRichFormat, Content: content}
		if _, err := parse.ParseTableCellRich(rich); err != nil {
			return ast.WalkSkipChildren
		}
		cell.TableCellRich = rich
		_ = parse.ApplyTableCellRichProjection(cell)
		if table := cell.Parent.Parent; table != nil {
			if table.Type == ast.NodeTableHead {
				table = table.Parent
			}
			if table != nil && table.Type == ast.NodeTable {
				table.SetIALAttr("custom-sy-table-rich", "1")
			}
		}
		return ast.WalkSkipChildren
	})
}

func markdownTableCellHTMLContent(source string) string {
	engine := New()
	engine.ParseOptions = parse.TableCellRichOptions()
	engine.SetKramdownIAL(true)
	engine.SetAutoSpace(false)
	engine.SetExportNormalizeTaskListMarker(false)
	fragment := engine.HTML2Tree(render.Sanitize(source))
	return string(render.NewFormatRenderer(fragment, engine.RenderOptions, engine.ParseOptions).Render())
}

// preserveHTMLTableCellStructure 在通用 HTML 预处理展开单元格内容前保存受支持的块结构。
func preserveHTMLTableCellStructure(n *html.Node) {
	if n.DataAtom == atom.Td || n.DataAtom == atom.Th {
		for _, attr := range n.Attr {
			if attr.Key == ast.TableCellRichAttribute {
				return
			}
		}
		var source strings.Builder
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			_ = html.Render(&source, child)
		}
		if fragment := strings.TrimSpace(source.String()); isMarkdownTableCellBlockHTML(fragment) {
			rich := &ast.TableCellRich{Spec: ast.TableCellRichSpec, Format: ast.TableCellRichFormat,
				Content: markdownTableCellHTMLContent(fragment)}
			if _, err := parse.ParseTableCellRich(rich); err == nil {
				n.Attr = append(n.Attr, &html.Attribute{Key: ast.TableCellRichAttribute, Val: rich.Encode()})
			}
		}
		return
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		preserveHTMLTableCellStructure(child)
	}
}

// isMarkdownTableCellBlockHTML 要求显式闭合和受支持的标签，避免 HTML 容错解析吞掉原文。
func isMarkdownTableCellBlockHTML(source string) bool {
	if !strings.HasPrefix(source, "<") {
		return false
	}
	var stack []atom.Atom
	hasBlock := false
	z := html.NewTokenizer(strings.NewReader(source))
	for {
		switch typ := z.Next(); typ {
		case html.ErrorToken:
			return z.Err() == io.EOF && len(stack) == 0 && hasBlock
		case html.TextToken:
		case html.StartTagToken, html.SelfClosingTagToken:
			tag := z.Token()
			switch tag.DataAtom {
			case atom.P, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Ul, atom.Ol, atom.Li, atom.Blockquote, atom.Pre:
				hasBlock = true
			case atom.Span, atom.Strong, atom.B, atom.Em, atom.I, atom.S, atom.Del, atom.U, atom.Mark, atom.Kbd, atom.Sub, atom.Sup, atom.Code, atom.A, atom.Br, atom.Img:
			case atom.Div:
				math := false
				for _, attr := range tag.Attr {
					if attr.Key == "class" && attr.Val == "language-math" {
						math = true
					}
				}
				if !math {
					return false
				}
				hasBlock = true
			case atom.Input:
				checkbox := false
				for _, attr := range tag.Attr {
					if attr.Key == "type" && strings.EqualFold(attr.Val, "checkbox") {
						checkbox = true
					}
				}
				if !checkbox {
					return false
				}
			default:
				return false
			}
			for _, attr := range tag.Attr {
				if strings.HasPrefix(attr.Key, "on") || attr.Key == ast.TableCellRichAttribute {
					return false
				}
			}
			if tag.DataAtom != atom.Br && tag.DataAtom != atom.Img && tag.DataAtom != atom.Input {
				if typ == html.SelfClosingTagToken {
					return false
				}
				stack = append(stack, tag.DataAtom)
			}
		case html.EndTagToken:
			tag := z.Token()
			if len(stack) == 0 || stack[len(stack)-1] != tag.DataAtom {
				return false
			}
			stack = stack[:len(stack)-1]
		default:
			return false
		}
	}
}

func isMarkdownTableCellTask(source string) bool {
	if len(source) < 5 || !strings.ContainsRune("-*+", rune(source[0])) {
		return false
	}
	source = source[1:]
	if len(source) == 0 || source[0] != ' ' && source[0] != '\t' {
		return false
	}
	source = strings.TrimLeft(source, " \t")
	return len(source) >= 3 && source[0] == '[' && source[2] == ']' &&
		parse.TableCellRichOptions().IsValidTaskListItemMarker(source[1]) &&
		(len(source) == 3 || source[3] == ' ' || source[3] == '\t')
}

// markdownTableCellBreaks 只展开代码之外的显式换行，保留反斜杠转义和代码中的 HTML。
func markdownTableCellBreaks(source string) string {
	var out strings.Builder
	code := 0
	for i := 0; i < len(source); {
		if source[i] == '\\' && i+1 < len(source) && code == 0 {
			out.WriteString(source[i : i+2])
			i += 2
			continue
		}
		if source[i] == '`' {
			end := i + 1
			for end < len(source) && source[end] == '`' {
				end++
			}
			if code == 0 {
				code = end - i
			} else if code == end-i {
				code = 0
			}
			out.WriteString(source[i:end])
			i = end
			continue
		}
		if code == 0 && source[i] == '<' {
			z := html.NewTokenizer(strings.NewReader(source[i:]))
			typ := z.Next()
			if (typ == html.StartTagToken || typ == html.SelfClosingTagToken) && z.Token().DataAtom == atom.Br {
				out.WriteByte('\n')
				i += len(z.Raw())
				continue
			}
		}
		out.WriteByte(source[i])
		i++
	}
	return out.String()
}
