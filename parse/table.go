// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

func (context *Context) parseTable(paragraph *ast.Node) (retParagraph, retTable *ast.Node) {
	var tokens []byte
	length := len(paragraph.Tokens)
	lineCnt := 0
	for i := 0; i < length; i++ {
		if context.ParseOption.ProtyleWYSIWYG {
			lines := lex.Split(paragraph.Tokens, lex.ItemNewline)
			delimRowIndex := context.findTableDelimRow(lines)
			if 1 > delimRowIndex {
				return
			}

			aligns := context.parseTableDelimRow(lex.TrimWhitespace(lines[delimRowIndex]))
			if nil == aligns {
				return
			}

			if 2 == length && 1 == len(aligns) && 0 == aligns[0] && !bytes.Contains(tokens, []byte("|")) {
				// 对于 Protyle 来说，这里应该是可以不必判断的，但为了保险，还是保留该判断逻辑
				// 具体细节可参考下方 GFM Table 解析的注释
				return
			}

			var headRows []*ast.Node
			for j := 0; j < delimRowIndex; j++ {
				headRow := context.parseTableRow(lex.TrimWhitespace(lines[j]), aligns, true)
				if nil == headRow {
					return
				}
				headRows = append(headRows, headRow)
				for th := headRow.FirstChild; nil != th; th = th.Next {
					ialStart := bytes.Index(th.Tokens, []byte("{:"))
					if 0 != ialStart {
						continue
					}

					subTokens := th.Tokens[ialStart:]
					if pos, ial := context.parseKramdownSpanIAL(subTokens); 0 < len(ial) {
						ialTokens := subTokens[:pos+1]
						if bytes.Contains(ialTokens, []byte("span")) || bytes.Contains(ialTokens, []byte("fn__none")) || // 合并单元格
							bytes.Contains(ialTokens, []byte("width:")) /* width: 是为了兼容遗留数据 */ {
							th.KramdownIAL = ial
							th.Tokens = th.Tokens[len(ialTokens):]
							spanIAL := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
							th.PrependChild(spanIAL)
						}
					}
				}
			}

			retTable = &ast.Node{Type: ast.NodeTable, TableAligns: aligns}
			retTable.TableAligns = aligns
			retTable.AppendChild(context.newTableHead(headRows))

			for j := delimRowIndex + 1; j < len(lines); j++ {
				line := lex.TrimWhitespace(lines[j])
				tableRow := context.parseTableRow(line, aligns, false)
				if nil == tableRow {
					return
				}
				if context.ParseOption.KramdownSpanIAL {
					for td := tableRow.FirstChild; nil != td; td = td.Next {
						ialStart := bytes.Index(td.Tokens, []byte("{:"))
						if 0 != ialStart {
							continue
						}

						subTokens := td.Tokens[ialStart:]
						if pos, ial := context.parseKramdownSpanIAL(subTokens); 0 < len(ial) {
							ialTokens := subTokens[:pos+1]
							if bytes.Contains(ialTokens, []byte("span")) || bytes.Contains(ialTokens, []byte("fn__none")) || // 合并单元格
								bytes.Contains(ialTokens, []byte("width:")) /* width: 是为了兼容遗留数据 */ {
								td.KramdownIAL = ial
								td.Tokens = td.Tokens[len(ialTokens):]
								spanIAL := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
								td.PrependChild(spanIAL)
							}
						}
					}
				}
				retTable.AppendChild(tableRow)
			}
			return
		} else {
			if lex.ItemNewline == paragraph.Tokens[i] || 0 == i {
				if 0 == i {
					tokens = paragraph.Tokens[i:]
				} else {
					tokens = paragraph.Tokens[i+1:]
				}
				if table := context.parseTable0(tokens); nil != table {
					if 0 < lineCnt {
						retParagraph = &ast.Node{Type: ast.NodeParagraph, Tokens: paragraph.Tokens[0:i]}
					}
					retTable = table
					retTable.Tokens = tokens
					break
				}
			}
		}
		lineCnt++
	}
	return
}

func (context *Context) parseTable0(tokens []byte) (ret *ast.Node) {
	lines := lex.Split(tokens, lex.ItemNewline)
	length := len(lines)
	if 2 > length {
		return
	}

	delimRow := lex.TrimWhitespace(lines[1])
	if 2 > len(delimRow) {
		// 换行+冒号会被识别为表格 https://github.com/88250/lute/issues/198
		return
	}

	aligns := context.parseTableDelimRow(delimRow)
	if nil == aligns {
		return
	}

	if 2 == length && 1 == len(aligns) && 0 == aligns[0] && !bytes.Contains(tokens, []byte("|")) {
		// 如果只有两行并且对齐方式是默认对齐且没有 | 时（foo\n---）就和 Setext 标题规则冲突了
		// 但在块级解析时显然已经尝试进行解析 Setext 标题，还能走到这里说明 Setetxt 标题解析失败，
		// 所以这里也不能当作表进行解析了，返回普通段落
		return
	}

	headRow := context.parseTableRow(lex.TrimWhitespace(lines[0]), aligns, true)
	if nil == headRow {
		return
	}

	if context.ParseOption.KramdownSpanIAL {
		for th := headRow.FirstChild; nil != th; th = th.Next {
			ialStart := bytes.LastIndex(th.Tokens, []byte("{:"))
			if 0 > ialStart {
				continue
			}
			subTokens := th.Tokens[ialStart:]
			if pos, ial := context.parseKramdownSpanIAL(subTokens); 0 < len(ial) {
				th.KramdownIAL = ial
				ialTokens := subTokens[:pos+1]
				th.Tokens = th.Tokens[:len(th.Tokens)-len(ialTokens)]
				spanIAL := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
				th.InsertAfter(spanIAL)
				th = th.Next
			}
		}
	}

	ret = &ast.Node{Type: ast.NodeTable, TableAligns: aligns}
	ret.TableAligns = aligns
	ret.AppendChild(context.newTableHead([]*ast.Node{headRow}))
	for i := 2; i < length; i++ {
		line := lex.TrimWhitespace(lines[i])
		tableRow := context.parseTableRow(line, aligns, false)
		if nil == tableRow {
			return
		}
		if context.ParseOption.KramdownSpanIAL {
			for th := tableRow.FirstChild; nil != th; th = th.Next {
				ialStart := bytes.LastIndex(th.Tokens, []byte("{:"))
				if 0 > ialStart {
					continue
				}
				subTokens := th.Tokens[ialStart:]
				if pos, ial := context.parseKramdownSpanIAL(subTokens); 0 < len(ial) {
					th.KramdownIAL = ial
					ialTokens := subTokens[:pos+1]
					th.Tokens = th.Tokens[:len(th.Tokens)-len(ialTokens)]
					spanIAL := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
					th.InsertAfter(spanIAL)
					th = th.Next
				}
			}
		}
		ret.AppendChild(tableRow)
	}
	return
}

func (context *Context) newTableHead(headRows []*ast.Node) *ast.Node {
	ret := &ast.Node{Type: ast.NodeTableHead}
	for _, headRow := range headRows {
		tr := &ast.Node{Type: ast.NodeTableRow}
		ret.AppendChild(tr)
		for c := headRow.FirstChild; nil != c; {
			next := c.Next
			tr.AppendChild(c)
			c = next
		}
	}
	return ret
}

func inInline(tokens []byte, i int, mathOrCodeMarker byte) bool {
	if i+1 >= len(tokens) || i < 1 {
		return false
	}

	start := bytes.IndexByte(tokens[:i], mathOrCodeMarker)
	startClosed := 0 == bytes.Count(tokens[:i], []byte{mathOrCodeMarker})%2
	if startClosed {
		return false
	}
	end := bytes.IndexByte(tokens[i+1:], mathOrCodeMarker)
	return -1 < start && -1 < end
}

func (context *Context) parseTableRow(line []byte, aligns []int, isHead bool) (ret *ast.Node) {
	ret = &ast.Node{Type: ast.NodeTableRow, TableAligns: aligns}

	if idx := bytes.Index(line, []byte("\\|")); 0 < idx {
		if inInline(line, idx, lex.ItemDollar) || inInline(line, idx, lex.ItemBacktick) {
			line = bytes.ReplaceAll(line, []byte("\\|"), []byte("\\&#124;"))
		}
	}

	cols := lex.SplitWithoutBackslashEscape(line, lex.ItemPipe)
	if 1 > len(cols) {
		return nil
	}
	if lex.IsBlank(cols[0]) {
		cols = cols[1:]
	}
	if len(cols) > 0 && lex.IsBlank(cols[len(cols)-1]) {
		cols = cols[:len(cols)-1]
	}

	colsLen := len(cols)
	alignsLen := len(aligns)
	if isHead && colsLen > alignsLen { // 分隔符行定义了表的列数，如果表头列数还大于这个列数，则说明不满足表格式
		return nil
	}

	var i int
	var col []byte
	for ; i < colsLen && i < alignsLen; i++ {
		col = lex.TrimWhitespace(cols[i])
		col = bytes.ReplaceAll(col, []byte("&#124;"), []byte("|"))
		cell := &ast.Node{Type: ast.NodeTableCell, TableCellAlign: aligns[i]}
		cell.Tokens = col
		ret.AppendChild(cell)
	}

	// 可能需要补全剩余的列
	for ; i < alignsLen; i++ {
		cell := &ast.Node{Type: ast.NodeTableCell, TableCellAlign: aligns[i]}
		ret.AppendChild(cell)
	}
	return
}

func (context *Context) findTableDelimRow(lines [][]byte) (index int) {
	length := len(lines)
	if 2 > length {
		return -1
	}

	for i, line := range lines {
		if nil != context.parseTableDelimRow(line) {
			index = i
			return
		}
	}
	return -1
}

func (context *Context) parseTableDelimRow(line []byte) (aligns []int) {
	length := len(line)
	if 1 > length {
		return nil
	}

	var token byte
	var i int
	for ; i < length; i++ {
		token = line[i]
		if lex.ItemPipe != token && lex.ItemHyphen != token && lex.ItemColon != token && lex.ItemSpace != token {
			return nil
		}
	}

	if idx := bytes.Index(line, []byte("\\|")); 0 < idx {
		if inInline(line, idx, lex.ItemDollar) || inInline(line, idx, lex.ItemBacktick) {
			line = bytes.ReplaceAll(line, []byte("\\|"), []byte("\\&#124;"))
		}
	}

	cols := lex.SplitWithoutBackslashEscape(line, lex.ItemPipe)
	if lex.IsBlank(cols[0]) {
		cols = cols[1:]
	}
	if len(cols) > 0 && lex.IsBlank(cols[len(cols)-1]) {
		cols = cols[:len(cols)-1]
	}

	var alignments []int
	for _, col := range cols {
		col = lex.TrimWhitespace(col)
		col = bytes.ReplaceAll(col, []byte("&#124;"), []byte("|"))
		if 1 > length || nil == col {
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

func (context *Context) tableDelimAlign(col []byte) int {
	length := len(col)
	if 1 > length {
		return -1
	}

	var left, right bool
	first := col[0]
	left = lex.ItemColon == first
	last := col[length-1]
	right = lex.ItemColon == last

	i := 1
	var token byte
	for ; i < length-1; i++ {
		token = col[i]
		if lex.ItemHyphen != token {
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
