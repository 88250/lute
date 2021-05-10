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

import "strings"

func isFileExt(pos, length int, runes *[]rune) bool {
	max := pos + maxCommonFileTypeLen
	if max > length {
		max = length
	}

	ext := string((*runes)[pos:max])
	for j := 0; j < commonFileTypesLen; j++ {
		if strings.HasPrefix(ext, commonFileTypes[j]) {
			return true
		}
	}
	return false
}

var commonFileTypesLen = len(commonFileTypes)
var maxCommonFileTypeLen = 10 // textbundle

// commonFileTypes 列出了常见的文件后缀，主要用于判断是否需要将英文句号.转换为中文句号。
var commonFileTypes = []string{
	// 图片

	"jpg",
	"png",
	"gif",
	"webp",
	"cr2",
	"tif",
	"bmp",
	"heif",
	"jxr",
	"psd",
	"ico",
	"dwg",

	// 视频

	"mp4",
	"m4v",
	"mkv",
	"webm",
	"mov",
	"avi",
	"wmv",
	"mpg",
	"flv",
	"3gp",

	// 音频

	"mid",
	"mp3",
	"m4a",
	"ogg",
	"flac",
	"wav",
	"amr",
	"aac",

	// 压缩包

	"epub",
	"zip",
	"tar",
	"rar",
	"gz",
	"bz2",
	"7z",
	"xz",
	"pdf",
	"exe",
	"swf",
	"rtf",
	"iso",
	"eot",
	"ps",
	"sqli",
	"nes",
	"crx",
	"cab",
	"deb",
	"ar",
	"Z",
	"lz",
	"rpm",
	"elf",
	"dcm",

	// 文件

	"doc",
	"docx",
	"xls",
	"xlsx",
	"ppt",
	"pptx",
	"md",
	"txt",

	// 字体

	"woff",
	"woff2",
	"ttf",
	"otf",

	// 应用程序

	"wasm",
	"exe",

	// 编程语言

	"html",
	"js",
	"css",
	"go",
	"java",

	// 其他

	"textbundle",
}
