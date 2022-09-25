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

//
//import (
//	"os"
//	"sync"
//	"testing"
//
//	"github.com/88250/lute"
//)
//
//func TestParallel(t *testing.T) {
//	data0, err := os.ReadFile("../test/commonmark-spec.md")
//	if nil != err {
//		t.Fatalf("read test text failed: " + err.Error())
//	}
//
//	data1, err := os.ReadFile("../test/case1.md")
//	if nil != err {
//		t.Fatalf("read test text failed: " + err.Error())
//	}
//
//	wg := sync.WaitGroup{}
//	for i := 0; i < 50; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			luteEngine := lute.New()
//			luteEngine.Markdown("", data0)
//		}()
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			luteEngine := lute.New()
//			luteEngine.Markdown("", data1)
//		}()
//	}
//	wg.Wait()
//}
