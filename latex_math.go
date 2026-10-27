package lute

import (
	"strings"

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/parse"
)

type latexMathDOMText struct {
	node       *html.Node
	start, end int
}

// Normalize only visible text runs, before HTML marker escaping. Formatting around a formula
// survives, and a formula split across inline formatting nodes is still recognized once.
func normalizeLatexMathDOM(root *html.Node, inline bool) {
	if root == nil || skipLatexMathDOM(root) {
		return
	}
	var refs []latexMathDOMText
	var text strings.Builder
	type run struct {
		refs []latexMathDOMText
		text string
	}
	var runs []run
	boundary := false
	flush := func() {
		if len(refs) > 0 {
			runs = append(runs, run{refs, text.String()})
		}
		refs = nil
		text.Reset()
	}
	var collect func(*html.Node)
	collect = func(n *html.Node) {
		if n.Type == html.TextNode || n.DataAtom == atom.Br {
			value := n.Data
			if n.DataAtom == atom.Br {
				value = "\n"
			}
			start := text.Len()
			text.WriteString(value)
			refs = append(refs, latexMathDOMText{n, start, text.Len()})
			return
		}
		if skipLatexMathDOM(n) {
			boundary = true
			flush()
			return
		}
		if isLatexMathInlineDOM(n) {
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				collect(child)
			}
			return
		}
		flush()
		boundary = true
		normalizeLatexMathDOM(n, inline)
	}
	// Snapshot children: replacements can insert siblings while runs are flushed.
	var children []*html.Node
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, child)
	}
	for _, child := range children {
		collect(child)
	}
	flush()
	for _, run := range runs {
		replaceLatexMathDOMRun(root, run.refs, run.text, inline, !boundary && len(runs) == 1)
	}
}

func skipLatexMathDOM(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Pre, atom.Code, atom.Kbd, atom.Script, atom.Style, atom.Textarea, atom.Math, atom.Svg, atom.A:
		return true
	}
	typ := strings.Fields(strings.ToLower(domAttr(n, "data-type")))
	for _, field := range typ {
		if field == "inline-math" || field == "nodemathblock" || field == "code" {
			return true
		}
	}
	class := strings.ToLower(domAttr(n, "class"))
	return strings.Contains(class, "katex") || strings.Contains(class, "mathjax") || domAttr(n, "data-tex") != ""
}

func domAttr(n *html.Node, name string) string {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func isLatexMathInlineDOM(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Span, atom.Strong, atom.B, atom.Em, atom.I, atom.U, atom.S, atom.Del, atom.Sub, atom.Sup, atom.Font, atom.Small:
		return true
	}
	return false
}

func replaceLatexMathDOMRun(root *html.Node, refs []latexMathDOMText, text string, inline, display bool) {
	type match struct {
		start, end int
		content    string
		block      bool
	}
	var matches []match
	tokens := []byte(text)
	for i := 0; i+1 < len(tokens); i++ {
		if tokens[i] == '$' {
			marker := "$"
			if tokens[i+1] == '$' {
				marker = "$$"
			}
			for j := i + len(marker); j < len(tokens); j++ {
				if tokens[j] == '\\' {
					j++
					continue
				}
				if strings.HasPrefix(text[j:], marker) {
					i = j + len(marker) - 1
					break
				}
			}
			continue
		}
		if tokens[i] != '\\' {
			continue
		}
		if tokens[i+1] == '\\' {
			i++
			continue
		}
		block := tokens[i+1] == '['
		if block && !display {
			continue
		}
		if !block && (!inline || tokens[i+1] != '(') {
			continue
		}
		content, consumed := parse.MatchLatexMath(tokens[i:], block)
		if consumed == 0 {
			continue
		}
		// Display formulas must occupy the whole text run in their block container.
		if block && (strings.TrimSpace(text[:i]) != "" || strings.TrimSpace(text[i+consumed:]) != "") {
			continue
		}
		matches = append(matches, match{i, i + consumed, string(content), block})
		i += consumed - 1
	}
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		math := &html.Node{Type: html.ElementNode, DataAtom: atom.Span, Data: "span", Attr: []*html.Attribute{
			{Key: "data-type", Val: "inline-math"}, {Key: "data-content", Val: m.content},
		}}
		if m.block {
			math.DataAtom, math.Data = atom.Div, "div"
			math.Attr = []*html.Attribute{{Key: "class", Val: "language-math"}}
			math.AppendChild(&html.Node{Type: html.TextNode, Data: strings.TrimSpace(m.content)})
			// The formula occupies this entire block. Lift it out of inline wrappers so
			// HTML-to-AST conversion cannot produce a block nested inside bold/italic text.
			for root.FirstChild != nil {
				root.RemoveChild(root.FirstChild)
			}
			root.AppendChild(math)
			return
		}
		inserted := false
		for _, ref := range refs {
			if ref.end <= m.start || ref.start >= m.end {
				continue
			}
			n := ref.node
			if n.DataAtom == atom.Br {
				n.Parent.RemoveChild(n)
				continue
			}
			start, end := max(0, m.start-ref.start), min(len(n.Data), m.end-ref.start)
			prefix, suffix := n.Data[:start], n.Data[end:]
			n.Data = prefix
			if !inserted {
				n.InsertAfter(math)
				inserted = true
			}
			if suffix != "" {
				if ref.end >= m.end {
					if ref.start <= m.start {
						math.InsertAfter(&html.Node{Type: html.TextNode, Data: suffix})
					} else {
						n.Data = suffix
					}
				}
			}
		}
	}
}
