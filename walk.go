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

// WalkStatus represents a current status of the Walk function.
type WalkStatus int

const (
	// WalkStop indicates no more walking needed.
	WalkStop = iota
	// WalkSkipChildren indicates that Walk wont walk on children of current node.
	WalkSkipChildren
	// WalkContinue indicates that Walk can continue to walk.
	WalkContinue
)

// Walker is a function that will be called when Walk find a new node.
// entering is set true before walks children, false after walked children.
// If Walker returns error, Walk function immediately stop walking.
type Walker func(n Node, entering bool) (WalkStatus, error)

// Walk walks a AST tree by the depth first search algorighm.
func Walk(n Node, walker Walker) error {
	status, err := walker(n, true)
	if err != nil || status == WalkStop {
		return err
	}
	if status != WalkSkipChildren {
		for c := n.FirstChild(); c != nil; c = c.Next() {
			if err := Walk(c, walker); err != nil {
				return err
			}
		}
	}
	status, err = walker(n, false)
	if err != nil || status == WalkStop {
		return err
	}
	return nil
}
