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

func TestGetLinkDest(t *testing.T) {
	luteEngine := lute.New()

	dest := luteEngine.GetLinkDest("www.bing.com/search?q=\"你好\"")
	if "http://www.bing.com/search?q=%22%E4%BD%A0%E5%A5%BD%22" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("www.bing.com/search?q=@Hello")
	if "http://www.bing.com/search?q=@Hello" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("www.bing.com/search?q=\"Hello\"")
	if "http://www.bing.com/search?q=%22Hello%22" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("https://foo.com/bar:baz")
	if "https://foo.com/bar:baz" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("[foo](bar.com)")
	if "bar.com" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("[foo](siyuan://blocks/20220817180757-c57m8qi)")
	if "siyuan://blocks/20220817180757-c57m8qi" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("[foo](https://abc.to/)")
	if "https://abc.to/" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("[foo](https://abc.pm/)")
	if "https://abc.pm/" != dest {
		t.Fatalf("get link dest failed")
	}

	dest = luteEngine.GetLinkDest("file://D:\\Admin\\Downloads\\隔壁叔叔过年还在玩思源.jpg")
	if "file://D:\\Admin\\Downloads\\隔壁叔叔过年还在玩思源.jpg" != dest {
		t.Fatalf("get link dest failed")
	}
}

func TestIsValidLinkDest(t *testing.T) {
	luteEngine := lute.New()

	if luteEngine.IsValidLinkDest("[foo](bar.com)") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("siyuan://blocks/20220817180757-c57m8qi") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://abc.to/") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://abc.pm/") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://abc.dev") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://www.notion.so/AUR-9b010a17ca2d4996a898a801426b0585") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("http://127.0.0.1:6806") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://ld246.com") {
		t.Fatalf("check link dest failed")
	}

	if luteEngine.IsValidLinkDest("https://ld246") {
		t.Fatalf("check link dest failed")
	}

	if luteEngine.IsValidLinkDest("ld246.com") {
		t.Fatalf("check link dest failed")
	}

	if !luteEngine.IsValidLinkDest("https://www.electronjs.org/docs/api/shell") {
		t.Fatalf("check link dest failed")
	}
}
