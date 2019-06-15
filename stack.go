// Lute - A structural markdown engine.
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

type stack struct {
	items
	count int
}

func (s *stack) popMatch(token *item) (tokens []*item) {
	for i := s.count - 1; 0 <= i; i-- {
		t := s.items[i]
		if token.typ == t.typ && token.val == t.val {
			s.count = i
			tokens = append(s.items[i:], token)
			s.items = s.items[:i]

			return
		}
	}

	return nil
}

func (s *stack) push(token *item) {
	s.items = append(s.items[:s.count], token)
	s.count++
}

func (s *stack) pop() *item {
	if s.count == 0 {
		return &tokenEOF
	}

	s.count--

	return s.items[s.count]
}

func (s *stack) popAll() items {
	ret := s.items
	s.count = 0
	s.items = items{}

	return ret
}

func (s *stack) peek() *item {
	if s.count == 0 {
		return &tokenEOF
	}

	return s.items[s.count-1]
}

func (s *stack) isEmpty() bool {
	return 0 == len(s.items)
}
