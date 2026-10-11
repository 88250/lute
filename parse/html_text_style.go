// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/util"
)

// HTMLTextStyle 只保存可映射到文本标记的外部样式，不携带页面布局或资源引用。
type HTMLTextStyle struct {
	Color, Background, Family string
	Size                      float64
	SizeUnit                  string
}

var htmlTextColor = regexp.MustCompile(`(?i)^(#[0-9a-f]{3,4}|#[0-9a-f]{6}|#[0-9a-f]{8}|[a-z]+|(?:rgb|rgba|hsl|hsla)\([0-9.,%+ /-]+\))$`)
var htmlTextThemeColor = regexp.MustCompile(`^(?i:var)\(\s*--b3-font-(?:color|background)[1-9][0-9]*\s*\)$`)
var htmlTextSize = regexp.MustCompile(`^([0-9]*\.?[0-9]+)(px|pt|pc|in|cm|mm|em|rem|%)$`)

// isHTMLTextColor 保留内置主题色引用，使复制的文本样式随主题切换。
func isHTMLTextColor(value string) bool {
	return htmlTextColor.MatchString(value) || htmlTextThemeColor.MatchString(value)
}

// HTMLStyleDeclarations 在引号和括号之外拆分声明，保留同一属性的优先级。
func HTMLStyleDeclarations(style string) map[string]string {
	ret := map[string]string{}
	important := map[string]bool{}
	start, depth := 0, 0
	var quote rune
	escaped := false
	for i, c := range style + ";" {
		if escaped {
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			continue
		}
		if c == '(' {
			depth++
		}
		if c == ')' {
			depth--
		}
		if c != ';' || depth != 0 {
			continue
		}
		key, value, ok := strings.Cut(style[start:i], ":")
		start = i + 1
		if !ok {
			continue
		}
		key, value = strings.ToLower(strings.TrimSpace(key)), strings.TrimSpace(value)
		priority := false
		if idx := strings.LastIndex(value, "!"); idx >= 0 && strings.EqualFold(strings.TrimSpace(value[idx+1:]), "important") {
			priority, value = true, strings.TrimSpace(value[:idx])
		}
		// WPS 使用纯色 background 简写；转换为背景色后仍按声明顺序及优先级处理。
		if key == "background" && isHTMLTextColor(value) {
			key = "background-color"
		}
		if !important[key] || priority {
			ret[key], important[key] = value, priority
		}
	}
	return ret
}

func (s HTMLTextStyle) CSS() string {
	var ret strings.Builder
	for _, pair := range [][2]string{{"color", s.Color}, {"background-color", s.Background}, {"font-family", s.Family}} {
		if pair[1] != "" {
			ret.WriteString(pair[0] + ": " + pair[1] + ";")
		}
	}
	if s.SizeUnit != "" {
		ret.WriteString("font-size: " + strconv.FormatFloat(s.Size, 'f', 6, 64) + s.SizeUnit + ";")
	}
	return ret.String()
}

// ResolveHTMLTextStyle 将相对字号按父级计算；没有绝对基准时保留相对编辑器字号的 em 值。
func ResolveHTMLTextStyle(n *html.Node, parent HTMLTextStyle) HTMLTextStyle {
	s := parent
	rawStyle := util.DomAttrValue(n, "style")
	if rawStyle == "" && n.DataAtom != atom.Small && n.DataAtom != atom.Font {
		return s
	}
	decl := HTMLStyleDeclarations(rawStyle)
	if n.DataAtom == atom.Font {
		for _, pair := range [][2]string{{"color", "color"}, {"face", "font-family"}} {
			if _, ok := decl[pair[1]]; !ok {
				decl[pair[1]] = util.DomAttrValue(n, pair[0])
			}
		}
		if _, ok := decl["font-size"]; !ok {
			if size, err := strconv.Atoi(util.DomAttrValue(n, "size")); err == nil && size >= 1 && size <= 7 {
				decl["font-size"] = []string{"", "x-small", "small", "medium", "large", "x-large", "xx-large", "xxx-large"}[size]
			}
		}
	}
	for _, property := range []struct {
		key string
		dst *string
	}{{"color", &s.Color}, {"background-color", &s.Background}, {"font-family", &s.Family}} {
		key, dst := property.key, property.dst
		value, ok := decl[key]
		if !ok || value == "" {
			continue
		}
		switch strings.ToLower(value) {
		case "inherit", "unset":
			continue
		case "initial", "revert", "revert-layer":
			*dst = ""
			continue
		case "currentcolor":
			if key == "color" {
				continue
			}
			value = s.Color
		}
		if key == "font-family" {
			if safeHTMLFontFamily(value) {
				*dst = value
			}
		} else if isHTMLTextColor(value) {
			*dst = value
		}
	}
	value := strings.ToLower(decl["font-size"])
	if value == "" && n.DataAtom == atom.Small {
		value = "smaller"
	}
	if value != "" {
		base, unit := parent.Size, parent.SizeUnit
		if unit == "" {
			base, unit = 1, "em"
		}
		switch value {
		case "inherit", "unset":
		case "initial", "medium", "revert", "revert-layer":
			s.Size, s.SizeUnit = 1, "rem"
		case "smaller":
			s.Size, s.SizeUnit = base/1.2, unit
		case "larger":
			s.Size, s.SizeUnit = base*1.2, unit
		default:
			if ratio, ok := map[string]float64{"xx-small": .6, "x-small": .75, "small": .888889, "large": 1.2, "x-large": 1.5, "xx-large": 2, "xxx-large": 3}[value]; ok {
				s.Size, s.SizeUnit = ratio, "rem"
			} else if match := htmlTextSize.FindStringSubmatch(value); match != nil {
				num, _ := strconv.ParseFloat(match[1], 64)
				if num > 0 && num <= 10000 {
					switch match[2] {
					case "em":
						s.Size, s.SizeUnit = base*num, unit
					case "%":
						s.Size, s.SizeUnit = base*num/100, unit
					case "rem", "px":
						s.Size, s.SizeUnit = num, match[2]
					default:
						s.Size, s.SizeUnit = num*map[string]float64{"pt": 96.0 / 72, "pc": 16, "in": 96, "cm": 96 / 2.54, "mm": 96 / 25.4}[match[2]], "px"
					}
				}
			}
		}
	}
	return s
}

// ApplyHTMLTextSemantics 保留样式文本的 CSS 强调及祖先链接，避免平铺后丢失组合格式。
func ApplyHTMLTextSemantics(node *ast.Node, dom *html.Node) {
	weight, italic := "", ""
	for p := dom; p != nil; p = p.Parent {
		decl := HTMLStyleDeclarations(util.DomAttrValue(p, "style"))
		if weight == "" {
			weight = strings.ToLower(decl["font-weight"])
		}
		if italic == "" {
			italic = strings.ToLower(decl["font-style"])
		}
		decoration := decl["text-decoration"] + " " + decl["text-decoration-line"]
		for _, pair := range [][2]string{{"underline", "u"}, {"line-through", "s"}} {
			if strings.Contains(decoration, pair[0]) && !node.ContainTextMarkTypes(pair[1]) {
				node.TextMarkType += " " + pair[1]
			}
		}
	}
	numericWeight, _ := strconv.Atoi(weight)
	if (weight == "bold" || weight == "bolder" || numericWeight >= 600) && !node.ContainTextMarkTypes("strong") {
		node.TextMarkType += " strong"
	}
	if (italic == "italic" || strings.HasPrefix(italic, "oblique")) && !node.ContainTextMarkTypes("em") {
		node.TextMarkType += " em"
	}
}

func safeHTMLFontFamily(value string) bool {
	if len(value) > 1024 {
		return false
	}
	var quote rune
	for _, c := range value {
		if c == '\'' || c == '"' {
			if quote == c {
				quote = 0
			} else if quote == 0 {
				quote = c
			}
			continue
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) && c != '-' && c != '_' && c != ',' {
			return false
		}
	}
	return quote == 0
}

func HTMLNodeTextStyle(n *html.Node) HTMLTextStyle {
	var parents []*html.Node
	for p := n; p != nil; p = p.Parent {
		parents = append(parents, p)
	}
	s := HTMLTextStyle{}
	for i := len(parents) - 1; i >= 0; i-- {
		s = ResolveHTMLTextStyle(parents[i], s)
	}
	return s
}

// ApplyHTMLTextStyle 将样式附着到原生文本标记，保留链接及组合格式的数据。
func ApplyHTMLTextStyle(n *ast.Node, style string) {
	if style == "" {
		return
	}
	if !ensureHTMLTextMark(n) {
		return
	}
	n.SetIALAttr("style", style)
	if n.Next != nil && n.Next.Type == ast.NodeKramdownSpanIAL {
		n.Next.Tokens = IAL2Tokens(n.KramdownIAL)
	} else {
		n.InsertAfter(&ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: IAL2Tokens(n.KramdownIAL)})
	}
}

func ensureHTMLTextMark(n *ast.Node) bool {
	switch n.Type {
	case ast.NodeText, ast.NodeLinkText, ast.NodeHTMLEntity:
		if n.ParentIs(ast.NodeImage) {
			return false
		}
		n.TextMarkType = "text"
		n.TextMarkTextContent = string(html.EscapeHTML(n.Tokens))
		n.Tokens = nil
		n.Type = ast.NodeTextMark
	case ast.NodeTextMark:
	default:
		return false
	}
	for p := n.Parent; p != nil && !p.IsBlock(); p = p.Parent {
		if p.Type == ast.NodeLink {
			if !n.ContainTextMarkTypes("a") {
				n.TextMarkType = "a " + n.TextMarkType
			}
			if dest := p.ChildByType(ast.NodeLinkDest); dest != nil {
				n.TextMarkAHref = dest.TokensStr()
			}
			if title := p.ChildByType(ast.NodeLinkTitle); title != nil {
				n.TextMarkATitle = title.TokensStr()
			}
		}
	}
	return true
}
