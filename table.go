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

// Table 描述了表节点结构。
type Table struct {
	*BaseNode
	Aligns []int // 从左到右每个表格节点的对齐方式，0：默认对齐，1：左对齐，2：居中对齐，3：右对齐
}

// TableHead 描述了表头节点结构。
type TableHead struct {
	*BaseNode
}

// TableRow 描述了表行节点结构。
type TableRow struct {
	*BaseNode
	Aligns []int
}

// TableCell 描述了表格节点结构。
type TableCell struct {
	*BaseNode
	Aligns int
}

func (context *Context) parseTable(lines []items) (ret *Table) {
	length := len(lines)
	if 2 > length {
		return
	}

	aligns := context.parseTableDelimRow(lines[1].trim())
	if nil == aligns {
		return
	}

	headRow := context.parseTableRow(lines[0].trim(), aligns, true)
	if nil == headRow {
		return
	}

	ret = &Table{&BaseNode{typ: NodeTable}, aligns}
	ret.Aligns = aligns
	ret.AppendChild(ret, context.newTableHead(headRow))
	for i := 2; i < length; i++ {
		tableRow := context.parseTableRow(lines[i].trim(), aligns, false)
		ret.AppendChild(ret, tableRow)
	}

	return
}

func (context *Context) newTableHead(headRow *TableRow) *TableHead {
	ret := &TableHead{&BaseNode{typ: NodeTableHead}}
	for c := headRow.FirstChild(); c != nil; {
		next := c.Next()
		ret.AppendChild(ret, c)
		c = next
	}
	return ret
}

func (context *Context) parseTableRow(line items, aligns []int, isHead bool) (ret *TableRow) {
	ret = &TableRow{&BaseNode{typ: NodeTableRow}, aligns}
	cols := line.splitWithoutBackslashEscape(itemPipe)
	if isBlank(cols[0]) {
		cols = cols[1:]
	}
	if len(cols) > 0 && isBlank(cols[len(cols)-1]) {
		cols = cols[:len(cols)-1]
	}

	for i, col := range cols {
		col = col.trim()
		cell := &TableCell{&BaseNode{typ: NodeTableCell}, aligns[i]}
		cell.tokens = col
		ret.AppendChild(ret, cell)
	}
	return
}

func (context *Context) parseTableDelimRow(line items) (aligns []int) {
	length := len(line)
	var token byte
	var i, pipes int
	for ; i < length; i++ {
		token = line[i]
		if itemPipe != token && itemHyphen != token && itemColon != token && itemSpace != token {
			return nil
		}
		if itemPipe == token {
			pipes++
		}
	}
	if 1 > pipes {
		return nil
	}

	cols := line.splitWithoutBackslashEscape(itemPipe)
	if isBlank(cols[0]) {
		cols = cols[1:]
	}
	if len(cols) > 0 && isBlank(cols[len(cols)-1]) {
		cols = cols[:len(cols)-1]
	}

	var alignments []int
	for _, col := range cols {
		col = col.trim()
		if 1 > length {
			return nil
		}

		align := context.tableDelimAlign(col)
		if -1 == align {
			return nil
		}
		alignments = append(alignments, align)
	}
	return alignments
}

func (context *Context) tableDelimAlign(col items) int {
	var left, right bool
	length := len(col)
	first := col[0]
	left = itemColon == first
	last := col[length-1]
	right = itemColon == last

	i := 1
	var token byte
	for ; i < length-1; i++ {
		token = col[i]
		if itemHyphen != token {
			return -1
		}
	}

	if left && right {
		return 2
	}
	if left {
		return 1
	}
	if right {
		return 3
	}

	return 0
}
