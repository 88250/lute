// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of the License at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestFlashcardOcclusionIDRoundTrip(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHTMLTag2TextMark(true)
	luteEngine.SetProtyleWYSIWYG(true)
	blockDOM := `<div data-node-id="20260812000000-abcdefg" data-type="NodeParagraph" class="p"><div contenteditable="true"><span data-type="mark" data-occlusion-id="occlusion-a">alpha</span><span data-type="mark" data-occlusion-id="occlusion-b">beta</span></div><div class="protyle-attr" contenteditable="false"></div></div>`
	tree := luteEngine.BlockDOM2Tree(blockDOM)
	first := tree.Root.FirstChild.FirstChild
	if first.TextMarkFlashcardOcclusionID != "occlusion-a" || first.Next.TextMarkFlashcardOcclusionID != "occlusion-b" {
		t.Fatalf("flashcard occlusion IDs were not parsed: first=%q second=%q",
			first.TextMarkFlashcardOcclusionID, first.Next.TextMarkFlashcardOcclusionID)
	}
	rendered := luteEngine.RenderNodeBlockDOM(tree.Root.FirstChild)
	if !strings.Contains(rendered, `data-occlusion-id="occlusion-a"`) ||
		!strings.Contains(rendered, `data-occlusion-id="occlusion-b"`) {
		t.Fatalf("flashcard occlusion IDs were not rendered: %s", rendered)
	}
	markdown := luteEngine.BlockDOM2Md(blockDOM)
	if !strings.Contains(markdown, `data-occlusion-id="occlusion-a"`) ||
		!strings.Contains(markdown, `data-occlusion-id="occlusion-b"`) {
		t.Fatalf("flashcard occlusion IDs were not persisted to Markdown: %s", markdown)
	}
	roundTrip := luteEngine.Md2BlockDOM(markdown, false)
	if !strings.Contains(roundTrip, `data-occlusion-id="occlusion-a"`) ||
		!strings.Contains(roundTrip, `data-occlusion-id="occlusion-b"`) {
		t.Fatalf("flashcard occlusion IDs did not survive Markdown round trip: %s", roundTrip)
	}
}
