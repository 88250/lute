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
	cmd            string
	param1, param2 string
	*parseTest
}

var vditorIRBlockDOMListCommandTests = []*cmdTest{

	{"tab0", "", "", &parseTest{"6", "<ol start=\"9\" data-marker=\"9.\" data-block=\"0\" data-node-id=\"20210208171000-kabmgdy\" data-type=\"ol\"><li data-marker=\"9.\" data-node-id=\"20210208171012-jm0ot60\"><p data-block=\"0\" data-node-id=\"20210208171012-o5e8nns\" data-type=\"p\">foo</p></li><li data-marker=\"10.\" data-node-id=\"20210208171012-vbrlc1j\"><p data-block=\"0\" data-node-id=\"20210208171012-8a8l0fe\" data-type=\"p\"><wbr></p></li></ol>", "<ol start=\"9\" data-marker=\"9.\" data-block=\"0\" data-node-id=\"20210208171000-kabmgdy\" data-type=\"ol\"><li data-marker=\"9.\" data-node-id=\"20210208171012-jm0ot60\"><p data-block=\"0\" data-node-id=\"20210208171012-o5e8nns\" data-type=\"p\">foo</p><ol data-marker=\"1.\" data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ol\"><li data-marker=\"1.\" data-node-id=\"20210208171012-vbrlc1j\"><p data-block=\"0\" data-node-id=\"20210208171012-8a8l0fe\" data-type=\"p\"><wbr></p></li></ol></li></ol>"}},
	{"tab2", "20210202180100-7bwdnda", "20210202180056-0v7vgiu", &parseTest{"5", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202170528-ojhmr62\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202160953-lv57abp\"><p data-block=\"0\" data-node-id=\"20210202175601-r14dwjn\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202180100-7bwdnda\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202180100-56r9236\"><p data-block=\"0\" data-node-id=\"20210202180100-khs50fn\" data-type=\"p\">foo2</p></li></ul></li><li data-marker=\"*\" data-node-id=\"20210202120408-9pvfz9m\"><p data-block=\"0\" data-node-id=\"20210202175604-woh50qg\" data-type=\"p\"><wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202180056-0v7vgiu\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202120410-lapp5km\"><p data-block=\"0\" data-node-id=\"20210202120410-37lsqmj\" data-type=\"p\">baz</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202170528-ojhmr62\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202160953-lv57abp\"><p data-block=\"0\" data-node-id=\"20210202175601-r14dwjn\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202180056-0v7vgiu\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202180100-56r9236\"><p data-block=\"0\" data-node-id=\"20210202180100-khs50fn\" data-type=\"p\">foo2</p></li><li data-marker=\"*\" data-node-id=\"20210202120408-9pvfz9m\"><p data-block=\"0\" data-node-id=\"20210202175604-woh50qg\" data-type=\"p\"><wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202120410-lapp5km\"><p data-block=\"0\" data-node-id=\"20210202120410-37lsqmj\" data-type=\"p\">baz</p></li></ul></li></ul></li></ul>"}},
	{"tab2", "undefined", "20210202120411-wvrov3m", &parseTest{"4", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202170528-ojhmr62\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202160953-lv57abp\"><p data-block=\"0\" data-node-id=\"20210202160953-h6zuk4x\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210202120408-9pvfz9m\"><p data-block=\"0\" data-node-id=\"20210202170202-zetod6k\" data-type=\"p\"><wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202120411-wvrov3m\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202120410-lapp5km\"><p data-block=\"0\" data-node-id=\"20210202120410-37lsqmj\" data-type=\"p\">baz</p></li></ul></li><li data-marker=\"*\" data-node-id=\"20210202170345-hgv4mh9\"><p data-block=\"0\" data-node-id=\"20210202170344-r8llpkk\" data-type=\"p\">bazz</p></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202170528-ojhmr62\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202160953-lv57abp\"><p data-block=\"0\" data-node-id=\"20210202160953-h6zuk4x\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210202120411-wvrov3m\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202120408-9pvfz9m\"><p data-block=\"0\" data-node-id=\"20210202170202-zetod6k\" data-type=\"p\"><wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210202120410-lapp5km\"><p data-block=\"0\" data-node-id=\"20210202120410-37lsqmj\" data-type=\"p\">baz</p></li></ul></li></ul></li><li data-marker=\"*\" data-node-id=\"20210202170345-hgv4mh9\"><p data-block=\"0\" data-node-id=\"20210202170344-r8llpkk\" data-type=\"p\">bazz</p></li></ul>"}},
	{"tab0", "", "", &parseTest{"3", "<ol data-marker=\"1.\" data-block=\"0\" data-node-id=\"20210130223346-movsux7\" data-type=\"ol\"><li data-marker=\"1.\" data-node-id=\"20210130223345-c1dbsod\"><p data-block=\"0\" data-node-id=\"20210130223345-3pob39p\" data-type=\"p\">foo</p></li><li data-marker=\"2.\" data-node-id=\"20210130223346-1o48sdg\"><p data-block=\"0\" data-node-id=\"20210130223347-zmiu4ai\" data-type=\"p\"><wbr></p></li></ol>", "<ol data-marker=\"1.\" data-block=\"0\" data-node-id=\"20210130223346-movsux7\" data-type=\"ol\"><li data-marker=\"1.\" data-node-id=\"20210130223345-c1dbsod\"><p data-block=\"0\" data-node-id=\"20210130223345-3pob39p\" data-type=\"p\">foo</p><ol data-marker=\"1.\" data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ol\"><li data-marker=\"1.\" data-node-id=\"20210130223346-1o48sdg\"><p data-block=\"0\" data-node-id=\"20210130223347-zmiu4ai\" data-type=\"p\"><wbr></p></li></ol></li></ol>"}},
	{"stab", "", "", &parseTest{"2", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130222701-uxus28e\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\"><wbr>bar</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130222621-8o6zx51\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130222701-11iujzf\"><p data-block=\"0\" data-node-id=\"20210130222822-80j0pqu\" data-type=\"p\"><wbr>bar</p></li></ul>"}},
	{"tab0", "", "", &parseTest{"1", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130222504-3o12agi\" data-type=\"p\">​<wbr>bar</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130222616-95qmu4z\"><p data-block=\"0\" data-node-id=\"20210130222616-2ixnd7w\" data-type=\"p\">baz</p></li></ul></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130222504-3o12agi\" data-type=\"p\"><wbr>bar</p></li><li data-marker=\"*\" data-node-id=\"20210130222616-95qmu4z\"><p data-block=\"0\" data-node-id=\"20210130222616-2ixnd7w\" data-type=\"p\">baz</p></li></ul></li></ul>"}},
	{"tab0", "", "", &parseTest{"0", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p></li><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"p\"><wbr>bar</p></li></ul>", "<ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20210130213415-0iexx8w\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213415-rbeto7n\"><p data-block=\"0\" data-node-id=\"20210130213415-3pdjf5w\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210130213421-4ym3unk\"><p data-block=\"0\" data-node-id=\"20210130213422-gu8cqov\" data-type=\"p\"><wbr>bar</p></li></ul></li></ul>"}},
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
		html := luteEngine.VditorIRBlockDOMListCommand(test.from, test.cmd, test.param1, test.param2)

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
	ast.Testing = false
}
