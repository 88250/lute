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
	"os"
	"testing"
	"time"

	"github.com/88250/lute"
	"github.com/88250/lute/parse"
)

func TestBlock(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorIR(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetAutoSpace(false)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)

	data, err := os.ReadFile("commonmark-spec-kramdown.md")
	if nil != err {
		t.Fatalf("read test data failed: %s", err)
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	tree := parse.Block("", data, luteEngine.ParseOptions)
	elapsed := time.Now().UnixNano()/int64(time.Millisecond) - now
	t.Logf("blocks [%d], doc blocks [%d], ellapsed [%d]ms", tree.BlockCount(), tree.DocBlockCount(), elapsed)

	now = time.Now().UnixNano() / int64(time.Millisecond)
	tree = parse.Parse("", data, luteEngine.ParseOptions)
	elapsed = time.Now().UnixNano()/int64(time.Millisecond) - now
	t.Logf("blocks [%d], doc blocks [%d], ellapsed [%d]ms", tree.BlockCount(), tree.DocBlockCount(), elapsed)
}
