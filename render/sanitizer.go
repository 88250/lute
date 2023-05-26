// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
)

// 没有实现可扩展的策略，仅过滤不安全的标签和属性。
// 鸣谢 https://github.com/microcosm-cc/bluemonday

var setOfElementsToSkipContent = map[string]interface{}{
	"frame":    nil,
	"frameset": nil,
	//"iframe":   nil,
	"noembed":  nil,
	"noframes": nil,
	"noscript": nil,
	"nostyle":  nil,
	"object":   nil,
	"script":   nil,
	"style":    nil,
	"title":    nil,
}

func Sanitize(str string) string {
	return string(sanitize([]byte(str)))
}

func sanitize(tokens []byte) []byte {
	var (
		buff                     bytes.Buffer
		skipElementContent       bool
		skippingElementsCount    int64
		mostRecentlyStartedToken string
	)

	caretLeftSpace := bytes.Contains(tokens, []byte(" "+editor.Caret))
	tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, []byte(editor.CaretReplacement))

	tokenizer := html.NewTokenizer(bytes.NewReader(tokens))
	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				ret := buff.Bytes()
				if caretLeftSpace {
					ret = bytes.ReplaceAll(ret, []byte("\""+editor.CaretReplacement), []byte("\" "+editor.CaretReplacement))
				} else {
					ret = bytes.ReplaceAll(ret, []byte("\" "+editor.CaretReplacement), []byte("\""+editor.CaretReplacement))
				}
				ret = bytes.ReplaceAll(ret, []byte(editor.CaretReplacement), editor.CaretTokens)
				return ret
			}

			return util.StrToBytes(err.Error())
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:
		case html.CommentToken:
		case html.StartTagToken:
			mostRecentlyStartedToken = token.Data

			if _, ok := setOfElementsToSkipContent[token.Data]; ok {
				skipElementContent = true
				skippingElementsCount++
				buff.WriteString(" ")
				break
			}

			if len(token.Attr) != 0 {
				token.Attr = sanitizeAttrs(token.Attr)
			}

			if !skipElementContent {
				// do not escape multiple query parameters
				if linkable(token.Data) {
					writeLinkableBuf(&buff, &token)
				} else {
					buff.WriteString(token.String())
				}
			}
		case html.EndTagToken:
			if mostRecentlyStartedToken == token.Data {
				mostRecentlyStartedToken = ""
			}

			if _, ok := setOfElementsToSkipContent[token.Data]; ok {
				skippingElementsCount--
				if skippingElementsCount == 0 {
					skipElementContent = false
				}
				buff.WriteString(" ")
				break
			}

			if !skipElementContent {
				buff.WriteString(token.String())
			}
		case html.SelfClosingTagToken:
			if len(token.Attr) != 0 {
				token.Attr = sanitizeAttrs(token.Attr)
			}

			if !skipElementContent {
				// do not escape multiple query parameters
				if linkable(token.Data) {
					writeLinkableBuf(&buff, &token)
				} else {
					buff.WriteString(token.String())
				}
			}
		case html.TextToken:
			if !skipElementContent {
				switch mostRecentlyStartedToken {
				case "script":
					// not encouraged, but if a policy allows JavaScript we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
				case "style":
					// not encouraged, but if a policy allows CSS styles we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
				default:
					// HTML escape the text
					buff.WriteString(token.String())
				}
			}
		}
	}
}

func linkable(elementName string) bool {
	switch elementName {
	case "a", "area", "blockquote", "img", "link", "script":
		return true
	default:
		return false
	}
}

func writeLinkableBuf(buff *bytes.Buffer, token *html.Token) {
	// do not escape multiple query parameters
	tokenBuff := bytes.NewBufferString("")
	tokenBuff.WriteString("<")
	tokenBuff.WriteString(token.Data)
	for _, attr := range token.Attr {
		if attr.Key == editor.CaretReplacement {
			tokenBuff.WriteString(" " + editor.CaretReplacement)
			continue
		}
		tokenBuff.WriteByte(' ')
		tokenBuff.WriteString(attr.Key)
		tokenBuff.WriteString(`="`)
		switch attr.Key {
		case "href", "src":
			tokenBuff.WriteString(html.EscapeString(attr.Val))
		default:
			// re-apply
			tokenBuff.WriteString(html.EscapeString(attr.Val))
		}
		tokenBuff.WriteByte('"')
	}
	if token.Type == html.SelfClosingTagToken {
		tokenBuff.WriteString(" /")
	}
	tokenBuff.WriteString(">")
	buff.WriteString(tokenBuff.String())
}

func sanitizeAttrs(attrs []*html.Attribute) (ret []*html.Attribute) {
	for _, attr := range attrs {
		if !allowAttr(attr.Key) {
			continue
		}
		if "src" == attr.Key || "srcdoc" == attr.Key || "srcset" == attr.Key || "href" == attr.Key {
			val := strings.ToLower(strings.TrimSpace(attr.Val))
			val = removeSpace(val)
			if strings.HasPrefix(val, "data:image/svg+xml") || strings.HasPrefix(val, "data:text/html") || strings.HasPrefix(val, "javascript") {
				continue
			}

			if newVal := html.UnescapeAttrVal(string(sanitize([]byte(val)))); val != newVal {
				continue
			}
		}

		ret = append(ret, attr)
	}
	return
}

func removeSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) || ' ' == r {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func allowAttr(attrName string) bool {
	for name := range eventAttrs {
		if attrName == name {
			return false
		}
	}
	return true
}

// HTML 事件属性。https://www.w3schools.com/tags/ref_eventattributes.asp
var eventAttrs = map[string]interface{}{
	// Window
	"onafterprint":   nil,
	"onbeforeprint":  nil,
	"onbeforeunload": nil,
	"onerror":        nil,
	"onhashchange":   nil,
	"onload":         nil,
	"onmessage":      nil,
	"onoffline":      nil,
	"ononline":       nil,
	"onpagehide":     nil,
	"onpageshow":     nil,
	"onpopstate":     nil,
	"onresize":       nil,
	"onstorage":      nil,
	"onunload":       nil,

	// Form
	"onblur":        nil,
	"onchange":      nil,
	"oncontextmenu": nil,
	"onfocus":       nil,
	"oninput":       nil,
	"oninvalid":     nil,
	"onreset":       nil,
	"onsearch":      nil,
	"onselect":      nil,
	"onsubmit":      nil,

	// Keyboard
	"onkeydown":  nil,
	"onkeypress": nil,
	"onkeyup":    nil,

	// Mouse
	"onclick":      nil,
	"ondblclick":   nil,
	"onmousedown":  nil,
	"onmousemove":  nil,
	"onmouseout":   nil,
	"onmouseover":  nil,
	"onmouseleave": nil,
	"onmouseenter": nil,
	"onmouseup":    nil,
	"onmousewheel": nil,
	"onwheel":      nil,

	// Drag
	"ondrag":      nil,
	"ondragend":   nil,
	"ondragenter": nil,
	"ondragleave": nil,
	"ondragover":  nil,
	"ondragstart": nil,
	"ondrop":      nil,
	"onscroll":    nil,

	// Clipboard
	"oncopy":  nil,
	"oncut":   nil,
	"onpaste": nil,

	// Media
	"onabort":          nil,
	"oncanplay":        nil,
	"oncanplaythrough": nil,
	"oncuechange":      nil,
	"ondurationchange": nil,
	"onemptied":        nil,
	"onended":          nil,
	"onloadeddata":     nil,
	"onloadedmetadata": nil,
	"onloadstart":      nil,
	"onpause":          nil,
	"onplay":           nil,
	"onplaying":        nil,
	"onprogress":       nil,
	"onratechange":     nil,
	"onseeked":         nil,
	"onseeking":        nil,
	"onstalled":        nil,
	"onsuspend":        nil,
	"ontimeupdate":     nil,
	"onvolumechange":   nil,
	"onwaiting":        nil,

	// Misc
	"ontoggle": nil,

	// SVG
	"onbegin":  nil,
	"onend":    nil,
	"onrepeat": nil,

	// meta
	"http-equiv": nil,

	// input
	"formaction": nil,
}
