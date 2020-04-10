// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
	"io"
	"net/url"
	"strings"
)

// 没有实现可扩展的策略，仅过滤不安全的标签和属性。
// 鸣谢 https://github.com/microcosm-cc/bluemonday

var setOfElementsToSkipContent = map[string]interface{}{
	"frame":    nil,
	"frameset": nil,
	"iframe":   nil,
	"noembed":  nil,
	"noframes": nil,
	"noscript": nil,
	"nostyle":  nil,
	"object":   nil,
	"script":   nil,
	"style":    nil,
	"title":    nil,
}

var allowedAttrs = map[string]interface{}{
	"id":     nil,
	"title":  nil,
	"alt":    nil,
	"href":   nil,
	"src":    nil,
	"class":  nil,
	"value":  nil,
	"align":  nil,
	"height": nil,
	"width":  nil,
}

func sanitize(tokens []byte) []byte {
	var (
		buff                     bytes.Buffer
		skipElementContent       bool
		skippingElementsCount    int64
		skipClosingTag           bool
		closingTagToSkipStack    []string
		mostRecentlyStartedToken string
	)

	tokenizer := html.NewTokenizer(bytes.NewReader(tokens))
	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				return buff.Bytes()
			}

			return util.StrToBytes(err.Error())
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:
		case html.CommentToken:
		case html.StartTagToken:
			mostRecentlyStartedToken = strings.ToLower(token.Data)

			if _, ok := setOfElementsToSkipContent[token.Data]; ok {
				skipElementContent = true
				skippingElementsCount++
				buff.WriteString(" ")
				break
			}

			if len(token.Attr) != 0 {
				token.Attr = sanitizeAttrs(token.Attr)
			}

			if len(token.Attr) == 0 {
				skipClosingTag = true
				closingTagToSkipStack = append(closingTagToSkipStack, token.Data)
				buff.WriteString(" ")
				break
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
			if mostRecentlyStartedToken == strings.ToLower(token.Data) {
				mostRecentlyStartedToken = ""
			}

			if skipClosingTag && closingTagToSkipStack[len(closingTagToSkipStack)-1] == token.Data {
				closingTagToSkipStack = closingTagToSkipStack[:len(closingTagToSkipStack)-1]
				if len(closingTagToSkipStack) == 0 {
					skipClosingTag = false
				}
				buff.WriteString(" ")
				break
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
		tokenBuff.WriteByte(' ')
		tokenBuff.WriteString(attr.Key)
		tokenBuff.WriteString(`="`)
		switch attr.Key {
		case "href", "src":
			u, err := sanitizedUrl(attr.Val)
			if err == nil {
				tokenBuff.WriteString(u)
			} else {
				// fallthrough
				tokenBuff.WriteString(html.EscapeString(attr.Val))
			}
		default:
			// re-apply
			tokenBuff.WriteString(html.EscapeString(attr.Val))
		}
		tokenBuff.WriteByte('"')
	}
	if token.Type == html.SelfClosingTagToken {
		tokenBuff.WriteString("/")
	}
	tokenBuff.WriteString(">")
	buff.WriteString(tokenBuff.String())
}

func sanitizedUrl(val string) (string, error) {
	u, err := url.Parse(val)
	if err != nil {
		return "", err
	}
	// sanitize the url query params
	sanitizedQueryValues := make(url.Values, 0)
	queryValues := u.Query()
	for k, vals := range queryValues {
		sk := html.EscapeString(k)
		for _, v := range vals {
			sv := escapeUrlComponent(v)
			sanitizedQueryValues.Set(sk, sv)
		}
	}
	u.RawQuery = sanitizedQueryValues.Encode()
	// u.String() will also sanitize host/scheme/user/pass
	return u.String(), nil
}

const escapedURLChars = "'<>\"\r"

func escapeUrlComponent(val string) string {
	w := bytes.NewBufferString("")
	i := strings.IndexAny(val, escapedURLChars)
	for i != -1 {
		if _, err := w.WriteString(val[:i]); err != nil {
			return w.String()
		}
		var esc string
		switch val[i] {
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		val = val[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return w.String()
		}
		i = strings.IndexAny(val, escapedURLChars)
	}
	w.WriteString(val)
	return w.String()
}

func sanitizeAttrs(attrs []html.Attribute) (ret []html.Attribute) {
	for _, attr := range attrs {
		if _, ok := allowedAttrs[attr.Key]; ok {
			ret = append(ret, attr)
		}
	}
	return
}
