// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"io/ioutil"
	"testing"
)

func BenchmarkLute(b *testing.B) {
	spec := "commonmark-0.29-spec"
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	tree, err := Parse("spec text", string(bytes))
	if nil != err {
		b.Fatalf("parse [%s] failed: %s", tree.name, err.Error())
	}

	renderer := NewHTMLRenderer()
	html, err := tree.Render(renderer)
	if nil != err {
		b.Fatalf("unexpected: %s", err)
	}

	ioutil.WriteFile(spec+".html", []byte(html), 0644)

	for i := 0; i < b.N; i++ {
		tree, err := Parse("spec text", string(bytes))
		if nil != err {
			b.Fatalf("parse [%s] failed: %s", tree.name, err.Error())
		}

		renderer := NewHTMLRenderer()
		if _, err := tree.Render(renderer); nil != err {
			b.Fatalf("parse [%s] failed: %s", tree.name, err.Error())
		}

	}
}
