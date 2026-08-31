// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of the License at
//
//     http://license.coscl.org.cn/MulanPSL2
//
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestProtyleTagPunctuationBoundary(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetTextMark(true)
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetTag(true)

	inputs := []string{
		"Per the exit condition, #83 shifts state. Escalated via the close (`20260819002715-enfbv43`, Routine Close), and whether this warrants breaking the standing \"#83 has no dedicated card\" convention.",
		"#83 shifts state. See (#113)",
		"#83 shifts state. See [#113]",
	}
	for _, input := range inputs {
		dom := luteEngine.Md2BlockDOM(input, false)
		if strings.Contains(dom, `data-type="tag`) {
			t.Fatalf("tag should not span a punctuation-prefixed hash reference\ninput\n\t%q\ngot\n\t%q", input, dom)
		}
		if markdown := strings.TrimSpace(luteEngine.BlockDOM2StdMd(dom)); input != markdown {
			t.Fatalf("markdown should round-trip unchanged\nexpected\n\t%q\ngot\n\t%q", input, markdown)
		}
	}
}

func TestProtyleTagPunctuationBoundaryKeepsValidSyntax(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetTextMark(true)
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetTag(true)

	for _, input := range []string{"#foo#", "(#foo#)", `"#foo#"`, "#foo bar#"} {
		dom := luteEngine.Md2BlockDOM(input, false)
		if !strings.Contains(dom, `data-type="tag"`) {
			t.Fatalf("valid tag syntax should remain supported\ninput\n\t%q\ngot\n\t%q", input, dom)
		}
	}

	adjacentTagsDOM := luteEngine.Md2BlockDOM("#foo##bar#", false)
	if count := strings.Count(adjacentTagsDOM, `data-type="tag"`); 2 != count {
		t.Fatalf("adjacent tags should remain separate\nexpected tag count\n\t2\ngot\n\t%d\nDOM\n\t%q", count, adjacentTagsDOM)
	}

	strongDOM := luteEngine.Md2BlockDOM("**foo.**1", false)
	if !strings.Contains(strongDOM, `data-type="strong"`) {
		t.Fatalf("Protyle punctuation handling should continue to support strong syntax\ngot\n\t%q", strongDOM)
	}
}
