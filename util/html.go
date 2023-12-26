// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

import (
	"bytes"
	"strings"

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
)

func GetTextMarkTextDataWithoutEscapeSingleQuote(n *html.Node) (content string) {
	content = DomText(n)
	content = strings.ReplaceAll(content, editor.Zwsp, "")
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	content = html.EscapeHTMLStr(content)
	return
}

func GetTextMarkTextData(n *html.Node) (content string) {
	content = GetTextMarkTextDataWithoutEscapeSingleQuote(n)
	content = strings.ReplaceAll(content, "'", "&apos;")
	return
}

func GetTextMarkInlineMemoData(n *html.Node) (content string) {
	content = DomAttrValue(n, "data-inline-memo-content")
	content = strings.ReplaceAll(content, editor.Zwsp, "")
	content = strings.ReplaceAll(content, "\n", editor.IALValEscNewLine)
	content = html.UnescapeHTMLStr(content)
	return
}

func GetTextMarkAData(n *html.Node) (href, title string) {
	href = DomAttrValue(n, "data-href")
	href = html.EscapeHTMLStr(href)
	title = DomAttrValue(n, "data-title")
	title = html.EscapeHTMLStr(title)
	return
}

func GetTextMarkInlineMathData(n *html.Node) (content string) {
	content = DomAttrValue(n, "data-content")
	content = strings.ReplaceAll(content, "\n", editor.IALValEscNewLine)
	content = html.UnescapeHTMLStr(content)
	content = strings.ReplaceAll(content, editor.Zwsp, "")
	return
}

func GetTextMarkBlockRefData(n *html.Node) (id, subtype string) {
	id = DomAttrValue(n, "data-id")
	subtype = DomAttrValue(n, "data-subtype")
	if "" == subtype {
		subtype = "s"
	}
	return
}

func GetTextMarkFileAnnotationRefData(n *html.Node) (id string) {
	id = DomAttrValue(n, "data-id")
	return
}

func DomChildByType(n *html.Node, dataAtom atom.Atom) *html.Node {
	for c := n.FirstChild; nil != c; c = c.NextSibling {
		if c.DataAtom == dataAtom {
			return c
		}
	}
	return nil
}

func DomHTML(n *html.Node) []byte {
	if nil == n {
		return nil
	}
	buf := &bytes.Buffer{}
	html.Render(buf, n)
	return bytes.ReplaceAll(buf.Bytes(), []byte(editor.Zwsp), nil)
}

func DomText(n *html.Node) string {
	buf := &bytes.Buffer{}
	if html.TextNode == n.Type {
		buf.WriteString(n.Data)
		return buf.String()
	}
	for child := n.FirstChild; nil != child; child = child.NextSibling {
		domText0(child, buf)
	}
	return buf.String()
}

func domText0(n *html.Node, buffer *bytes.Buffer) {
	if nil == n {
		return
	}
	if dataRender := DomAttrValue(n, "data-render"); "1" == dataRender || "2" == dataRender {
		return
	}

	if "svg" == n.Namespace {
		return
	}

	isTempMark := false
	if 0 == n.DataAtom && html.ElementNode == n.Type {
		// 可能是自定义标签
		parent := n.Parent
		if nil == parent {
			return
		}
		if atom.Span != parent.DataAtom {
			return
		}

		if !IsTempMarkSpan(parent) {
			// Protyle 中的搜索高亮标记需要保留 https://github.com/siyuan-note/siyuan/issues/9821
			return
		}

		isTempMark = true
	}

	switch n.DataAtom {
	case 0:
		if isTempMark {
			buffer.WriteString("<" + n.Data + ">")
		} else {
			buffer.WriteString(n.Data)
		}
	case atom.Br:
		buffer.WriteString("\n")
	}

	for child := n.FirstChild; nil != child; child = child.NextSibling {
		domText0(child, buffer)
	}
}

func IsTempMarkSpan(n *html.Node) bool {
	dataType := DomAttrValue(n, "data-type")
	return "search-mark" == dataType || "virtual-block-ref" == dataType

}

func DomAttrValue(n *html.Node, attrName string) string {
	if nil == n {
		return ""
	}

	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func DomCustomAttrs(n *html.Node) (ret map[string]string) {
	ret = map[string]string{}
	for _, attr := range n.Attr {
		if strings.HasPrefix(attr.Key, "custom-") {
			ret[attr.Key] = attr.Val
		}
	}
	if 1 > len(ret) {
		return nil
	}
	return
}
