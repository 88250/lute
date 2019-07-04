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
	"testing"
)

func TestDelimiterStack(t *testing.T) {
	t1 := &Text{NodeType: NodeText, Value: "*"}
	e1 := &delimiter{node: t1, typ: "*", num: 1}
	t2 := &Text{NodeType: NodeText, Value: "lute"}
	e2 := &delimiter{node: t2, typ: ""}
	t3 := &Text{NodeType: NodeText, Value: "*"}
	e3 := &delimiter{node: t3, typ: "", num: 1}

	s := &delimiterStack{}
	s.push(e1)
	s.push(e2)

	if "lute" != s.pop().node.(*Text).Value {
		t.Error("unexpected stack item")
	}

	s.push(e3)

	if "*" != s.peek().node.(*Text).Value {
		t.Error("unexpected stack item")
	}
	if "*" != s.pop().node.(*Text).Value {
		t.Error("unexpected stack item")
	}

	s.push(e1)
	s.push(e2)

	elements := s.popMatch(e1)
	if t1.Value != elements[0].node.(*Text).Value || t2.Value != elements[1].node.(*Text).Value || t3.Value != elements[2].node.(*Text).Value {
		t.Error("unexpected stack item")
	}
}
