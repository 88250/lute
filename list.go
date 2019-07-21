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

type ListType int

type List struct {
	*BaseNode

	Bullet bool
	Start  int
	Delim  string
	Tight  bool

	Marker            string
	StartIndentSpaces int
	IndentSpaces      int
}

func (n *List) Close() {
	if n.close {
		return
	}

	tight := true
	for child := n.FirstChild(); nil != child; child = child.Next() {
		if NodeListItem == child.Type() && !child.(*ListItem).Tight {
			tight = false
			break
		}
	}
	n.Tight = tight

	for child := n.FirstChild(); nil != child; child = child.Next() {
		if NodeListItem == child.Type() {
			child.(*ListItem).Tight = tight
		}

		child.Close()
	}

	n.close = true
}

func (t *Tree) parseList(line items) (ret Node) {
	n := t.parseListItem(line)
	if nil == n {
		return
	}

	li := n.(*ListItem)
	ret = &List{
		&BaseNode{typ: NodeList},
		li.Bullet,
		li.Start,
		li.Delim,
		li.Tight,
		li.Marker,
		li.StartIndentSpaces,
		li.IndentSpaces,
	}

	return
}
