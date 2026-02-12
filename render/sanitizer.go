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
	"strings"

	"github.com/88250/lute/util"
	"github.com/microcosm-cc/bluemonday"
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

func SanitizeSrc(src string) string {
	img := strings.ReplaceAll(src, "\"", "__@QUOTE@__")
	img = strings.ReplaceAll(img, " ", "__@SPACE@__")
	img = strings.ReplaceAll(img, "#", "__@HASH@__")
	img = "<img src=\"" + img + "\"></img>"

	sanitizer := newSanitizer()
	img = sanitizer.Sanitize(img)
	img = string(util.TagSrcStr((img)))
	img = strings.ReplaceAll(img, "__@QUOTE@__", "\"")
	img = strings.ReplaceAll(img, "__@SPACE@__", " ")
	img = strings.ReplaceAll(img, "__@HASH@__", "#")
	return img
}

func Sanitize(str string) string {
	return newSanitizer().Sanitize(str)
}

func sanitize(tokens []byte) []byte {
	return newSanitizer().SanitizeBytes(tokens)
}

func newSanitizer() *bluemonday.Policy {
	ret := bluemonday.NewPolicy()
	ret.AllowStandardAttributes()
	ret.AllowDataAttributes()
	ret.AllowStandardURLs()
	ret.AllowImages()
	ret.AllowLists()
	ret.AllowStyling()
	ret.AllowTables()
	ret.AllowAttrs("align").OnElements("p", "div")
	ret.AllowAttrs("src", "scrolling", "border", "frameborder", "framespacing", "allowfullscreen", "data-subtype", "updated").OnElements("iframe")
	ret.AllowAttrs("content").OnElements("meta")
	ret.AllowAttrs("type", "allowscriptaccess").OnElements("embed")
	ret.AllowAttrs("loading").OnElements("img")
	ret.AllowAttrs("controls", "autoplay", "loop", "muted", "src").OnElements("video", "audio")
	ret.AllowElements("details", "summary", "video", "source", "audio", "embed")
	ret.AllowAttrs("open").OnElements("details")
	return ret
}
