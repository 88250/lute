// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

type cmdTest struct {
	cmd string
	*parseTest
}

var vditorIRBlockDOMListCommandTests = []*cmdTest{

	//{"tab0", &parseTest{"5", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130232132-m4xi7u5\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130232133-use7p8k\"><p data-block=\"0\" data-node-id=\"20210130232159-7lp8ceo\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130232139-y2ocdmk\"><p data-block=\"0\" data-node-id=\"20210130232159-s0bef3c\" data-type=\"p\"><wbr>​bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130232139-a2kpuve\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130232150-bhw66nk\"><p data-block=\"0\" data-node-id=\"20210130232150-cfdf9nz\" data-type=\"p\">baz</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\">\u200b<wbr>bar</p></li></ul>"}},
	//{"tab0", &parseTest{"4", "<ol data-marker=\"1.\" data-block=\"0\" data-node-id=\"20210130223346-movsux7\" data-type=\"ol\"><li data-marker=\"1.\" data-node-id=\"20210130223345-c1dbsod\"><p data-block=\"0\" data-node-id=\"20210130223345-3pob39p\" data-type=\"p\">foo</p></li><li data-marker=\"2.\" data-node-id=\"20210130223346-1o48sdg\"><p data-block=\"0\" data-node-id=\"20210130223347-zmiu4ai\" data-type=\"p\"><wbr></p></li></ol>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\">\u200b<wbr>bar</p></li></ul>"}},
	//{"stab", &parseTest{"3", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130222701-uxus28e\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\"><wbr>bar</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\">\u200b<wbr>bar</p></li></ul>"}},
	//{"enter", &parseTest{"2", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222701-uxus28e\" data-type=\"p\">bar<wbr></p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130222707-0a8o2t4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222707-e5fngvv\"><p data-block=\"0\" data-node-id=\"20210130222723-j3vxsk9\" data-type=\"p\">baz</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222701-uxus28e\" data-type=\"p\">bar</p></li><li data-marker=\"*\" data-node-id=\"20060102150405-1a2b3c4\"><p data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">\u200b<wbr></p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130222707-0a8o2t4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222707-e5fngvv\"><p data-block=\"0\" data-node-id=\"20210130222723-j3vxsk9\" data-type=\"p\">baz</p></li></ul></li></ul>"}},
	{"tab0", &parseTest{"1", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130222504-3o12agi\" data-type=\"p\">​<wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222616-95qmu4z\"><p data-block=\"0\" data-node-id=\"20210130222616-2ixnd7w\" data-type=\"p\">baz</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><wbr>bar</p></li></ul><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222616-95qmu4z\"><p data-block=\"0\" data-node-id=\"20210130222616-2ixnd7w\" data-type=\"p\">baz</p></li></ul></li></ul>"}},
	{"tab0", &parseTest{"0", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"p\"><wbr>bar</p></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><wbr>bar</p></li></ul></li></ul>"}},
}

func TestVditorIRBlockDOMListCommand(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorIR(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.SetAutoSpace(false)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)

	ast.Testing = true
	for _, test := range vditorIRBlockDOMListCommandTests {
		html := luteEngine.VditorIRBlockDOMListCommand(test.from, test.cmd)

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
	ast.Testing = false
}
