// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

type stack struct {
	items []*item
	count int
}

func (s *stack) popMatch(token *item) (tokens []*item) {
	for i := s.count - 1; 0 <= i; i-- {
		t := s.items[i]
		if token.typ == t.typ && token.val == t.val {
			s.count = i

			return s.items[i:]
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

func (s *stack) peek() *item {
	if s.count == 0 {
		return &tokenEOF
	}

	return s.items[s.count-1]
}

func (s *stack) isEmpty() bool {
	return 0 == len(s.items)
}
