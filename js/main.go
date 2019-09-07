// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package main

import (
	"github.com/b3log/lute"
	"github.com/gopherjs/gopherjs/js"
)

func markdown(markdownText string, options map[string]interface{}) string {
	gfm := true
	softBreak2HardBreak := true
	autoSpace := true
	fixTermTypo := true
	emoji := true
	if 0 < len(options) {
		if v, ok := options["gfm"]; ok {
			gfm = v.(bool)
		}
		if v, ok := options["softBreak2HardBreak"]; ok {
			softBreak2HardBreak = v.(bool)
		}
		if v, ok := options["autoSpace"]; ok {
			autoSpace = v.(bool)
		}
		if v, ok := options["fixTermTypo"]; ok {
			fixTermTypo = v.(bool)
		}
		if v, ok := options["emoji"]; ok {
			emoji = v.(bool)
		}
	}

	luteEngine := lute.New(
		lute.GFM(gfm),
		lute.SoftBreak2HardBreak(softBreak2HardBreak),
		lute.CodeSyntaxHighlight(false, false, false, "github"),
		lute.AutoSpace(autoSpace),
		lute.FixTermTypo(fixTermTypo),
		lute.Emoji(emoji),
	)

	html, err := luteEngine.MarkdownStr("", markdownText)
	if nil != err {
		return err.Error()
	}
	return html
}

func format(markdownText string) string {
	luteEngine := lute.New()
	formatted, err := luteEngine.FormatStr("", markdownText)
	if nil != err {
		return err.Error()
	}
	return formatted
}

func New(options map[string]interface{}) *js.Object {
	luteEngine := lute.New()
	if 0 < len(options) {
		if v, ok := options["gfm"]; ok {
			lute.GFM(v.(bool))(luteEngine)
		}
		if v, ok := options["softBreak2HardBreak"]; ok {
			luteEngine.SoftBreak2HardBreak = v.(bool)
		}
		if v, ok := options["autoSpace"]; ok {
			luteEngine.AutoSpace = v.(bool)
		}
		if v, ok := options["fixTermTypo"]; ok {
			luteEngine.FixTermTypo = v.(bool)
		}
		if v, ok := options["emoji"]; ok {
			luteEngine.Emoji = v.(bool)
		}
		if v, ok := options["emojis"]; ok {
			luteEngine.PutEmojis(v.(map[string]string))
		}
		if v, ok := options["emojiSite"]; ok {
			luteEngine.EmojiSite = v.(string)
		}
	}
	return js.MakeWrapper(luteEngine)
}

func main() {
	js.Global.Set("lute", make(map[string]interface{}))
	l := js.Global.Get("lute")
	l.Set("new", New)
}
