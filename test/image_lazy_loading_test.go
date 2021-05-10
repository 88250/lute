// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"testing"

	"github.com/88250/lute"
)

var imageLazyLoadingTests = []parseTest{

	{"0", "![foo](bar.png)", "<p><img src=\"/images/img-loading.svg\" data-src=\"bar.png\" alt=\"foo\" /></p>\n"},
}

func TestImageLazyLoading(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderOptions.ImageLazyLoading = "/images/img-loading.svg"

	for _, test := range imageLazyLoadingTests {
		md := luteEngine.MarkdownStr("", test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
