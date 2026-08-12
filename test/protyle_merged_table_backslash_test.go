// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

const mergedTableWithBackslashBlockDOM = `<div data-type="NodeTable" colgroup="|">` +
	`<div contenteditable="true" spellcheck="false"><table><colgroup><col><col></colgroup>` +
	`<thead><tr><th>1</th><th>reels</th></tr></thead><tbody>` +
	`<tr><td rowspan="2">4</td><td><span data-type="backslash">*</span>raj-s</td></tr>` +
	`<tr><td class="fn__none"></td><td><span data-type="backslash">*</span>[r]aj</td></tr>` +
	`</tbody></table><div class="protyle-action__table"><div class="table__resize"></div>` +
	`<div class="table__select"></div></div></div>` +
	`<div class="protyle-attr" contenteditable="false">&#8203;</div></div>`

func TestBlockDOM2StdMdMergedTablePreservesBackslashContent(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)

	markdown := luteEngine.BlockDOM2StdMd(mergedTableWithBackslashBlockDOM)
	if !strings.Contains(markdown, `<td>\*raj-s</td>`) || !strings.Contains(markdown, `<td>\*[r]aj</td>`) {
		t.Fatalf("expected escaped characters in merged table, got %q", markdown)
	}
}

func TestBlockDOM2RichHTMLMergedTablePreservesBackslashContent(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)

	html := luteEngine.BlockDOM2RichHTML(mergedTableWithBackslashBlockDOM)
	if !strings.Contains(html, `<td>*raj-s</td>`) || !strings.Contains(html, `<td>*[r]aj</td>`) {
		t.Fatalf("expected escaped character content in merged table, got %q", html)
	}
}
