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
	"regexp"
	"strings"

	"github.com/88250/lute/html"
)

var (
	// tableCaptionMarkupPattern 匹配表格标题的 markup 形态：IAL 中保存的是转义后的 <caption> 标记。
	tableCaptionMarkupPattern = regexp.MustCompile(`(?is)^\s*<caption\b([^>]*)>(.*)</caption>\s*$`)
	// tableCaptionAttrPattern 提取 title 标记上的属性，未识别的属性一律丢弃。
	tableCaptionAttrPattern = regexp.MustCompile(`(?is)([a-z][a-z0-9-]*)\s*=\s*"([^"]*)"`)
	// tableCaptionBottomStylePattern 仅接受 caption-side: bottom，可带 !important。
	tableCaptionBottomStylePattern = regexp.MustCompile(`(?is)^\s*caption-side\s*:\s*bottom\s*(?:!\s*important\s*)?;?\s*$`)
	// tableCaptionTagPattern 匹配标题文本中残留的标签，用于剥离未转义的 HTML。
	tableCaptionTagPattern = regexp.MustCompile(`(?is)<[^>]*>`)
)

// renderTableCaption 渲染表格标题的 <caption> 元素。
// IAL 中保存的是转义后的 <caption> 标记，先按既有格式解码一次，再只保留标题文本和
// caption-side: bottom 样式，其余标签、属性一律丢弃，避免标题被当作 HTML 执行。
// editable 为 true 时输出 contenteditable="false"，用于编辑器 BlockDOM；预览和导出不输出。
func renderTableCaption(caption string, editable bool) string {
	caption = html.UnescapeHTMLStr(caption)
	inner, bottom, contentEditable := sanitizeTableCaption(caption)

	var attrs strings.Builder
	if editable && contentEditable {
		attrs.WriteString(` contenteditable="false"`)
	}
	if bottom {
		attrs.WriteString(` style="caption-side: bottom;"`)
	}
	return "<caption" + attrs.String() + ">" + inner + "</caption>"
}

// sanitizeTableCaption 从解码后的标题中提取安全文本、位置样式和 contenteditable 标记。
// 标题不是 <caption> 标记时按纯文本处理；文本中残留的标签会被剥离，剩余的 < 会被转义。
func sanitizeTableCaption(caption string) (inner string, bottom, contentEditable bool) {
	if matches := tableCaptionMarkupPattern.FindStringSubmatch(caption); nil != matches {
		inner = matches[2]
		for _, attr := range tableCaptionAttrPattern.FindAllStringSubmatch(matches[1], -1) {
			switch strings.ToLower(attr[1]) {
			case "style":
				if tableCaptionBottomStylePattern.MatchString(attr[2]) {
					bottom = true
				}
			case "contenteditable":
				if "false" == strings.ToLower(strings.TrimSpace(attr[2])) {
					contentEditable = true
				}
			}
		}
	} else {
		inner = caption
	}

	inner = tableCaptionTagPattern.ReplaceAllString(inner, "")
	inner = strings.ReplaceAll(inner, "<", "&lt;")
	return
}
