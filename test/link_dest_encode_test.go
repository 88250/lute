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
	"bytes"
	"testing"

	"github.com/88250/lute/html"
)

func TestLinkDestEncode(t *testing.T) {
	dest1 := []byte("http://foo.bar/测试")
	encoded := html.EncodeDestination(dest1)
	dest2 := html.DecodeDestination(encoded)
	if !bytes.Equal(dest1, dest2) {
		t.Fatalf("Link dest encode failed")
	}
}

func TestLinkDestDecode(t *testing.T) {
	dest1 := []byte("C:\\Users\\DL882\\Documents\\SiYuan\\data\\assets\\%E6%B5%8B%E8%AF%95-20210621160821-zxz1j35.pdf")
	decoded := html.DecodeDestination(dest1)
	if !bytes.Equal([]byte("C:\\Users\\DL882\\Documents\\SiYuan\\data\\assets\\测试-20210621160821-zxz1j35.pdf"), decoded) {
		t.Fatalf("Link dest encode failed")
	}
}
