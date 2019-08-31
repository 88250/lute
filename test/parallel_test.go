// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package test

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/b3log/lute"
)

func TestParallel(t *testing.T) {
	data0, err := ioutil.ReadFile("../test/commonmark-spec.md")
	if nil != err {
		t.Fatalf("read test text failed: " + err.Error())
	}

	data1, err := ioutil.ReadFile("../test/case1.md")
	if nil != err {
		t.Fatalf("read test text failed: " + err.Error())
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			luteEngine := lute.New()
			_, err := luteEngine.Markdown("", data0)
			if nil != err {
				t.Fatalf(err.Error())
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			luteEngine := lute.New()
			_, err := luteEngine.Markdown("", data1)
			if nil != err {
				t.Fatalf(err.Error())
			}
		}()
	}
	wg.Wait()
}
