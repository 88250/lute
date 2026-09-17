package lute

import (
	"strings"

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/util"
)

// normalizeCodeBlockBanner 在转换子节点前提取代码横幅中的语言，避免横幅生成普通段落。
func normalizeCodeBlockBanner(n *html.Node) {
	if atom.Div != n.DataAtom || !domClassContains(n, "md-code-block") {
		return
	}
	var banner, pre *html.Node
	for child := n.FirstChild; nil != child; child = child.NextSibling {
		if domClassContains(child, "md-code-block-banner") {
			if nil != banner {
				return
			}
			banner = child
		} else if atom.Pre == child.DataAtom {
			if nil != pre {
				return
			}
			pre = child
		}
	}
	if nil == banner || nil == pre {
		return
	}
	var text strings.Builder
	codeBannerText(banner, &text)
	language := strings.TrimSpace(text.String())
	if !isCodeBannerLanguage(language) {
		return
	}

	codes := util.DomChildrenByType(pre, atom.Code)
	if 1 < len(codes) {
		return
	}
	var code *html.Node
	if 1 == len(codes) {
		code = codes[0]
	} else {
		code = &html.Node{Type: html.ElementNode, DataAtom: atom.Code, Data: "code"}
		for nil != pre.FirstChild {
			child := pre.FirstChild
			child.Unlink()
			code.AppendChild(child)
		}
		pre.AppendChild(code)
	}
	// 已有语言由原有代码块解析逻辑处理，横幅只补充缺失的信息。
	preClass := util.DomAttrValue(pre, "class")
	hasPreLanguage := strings.Contains(preClass, "language-") || "" != util.DomAttrValue(pre, "data-language") ||
		("" != preClass && !strings.ContainsAny(preClass, " \t\r\n-_") && "fallback" != preClass && "chroma" != preClass)
	if "" == util.DomAttrValue(code, "class") && !hasPreLanguage {
		util.SetDomAttrValue(code, "class", "language-"+language)
	}
	banner.Unlink()
}

// codeBannerText 排除横幅中的交互控件和图标，保留语言标签文本。
func codeBannerText(n *html.Node, text *strings.Builder) {
	if atom.Button == n.DataAtom || atom.A == n.DataAtom || atom.Svg == n.DataAtom ||
		"button" == util.DomAttrValue(n, "role") || "true" == util.DomAttrValue(n, "aria-hidden") ||
		domClassContains(n, "ds-icon-button") {
		return
	}
	if html.TextNode == n.Type {
		text.WriteString(n.Data)
		text.WriteByte(' ')
	}
	for child := n.FirstChild; nil != child; child = child.NextSibling {
		codeBannerText(child, text)
	}
}

func isCodeBannerLanguage(language string) bool {
	if "" == language || len(language) > 64 {
		return false
	}
	for i, c := range language {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || i > 0 &&
			(c >= '0' && c <= '9' || strings.ContainsRune("#+.", c)) {
			continue
		}
		return false
	}
	switch strings.ToLower(language) {
	case "copy", "copied", "expand", "collapse", "download", "run":
		return false
	}
	return true
}
