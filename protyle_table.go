package lute

import (
	"strconv"

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/util"
)

// normalizeProtyleTableSpans 为合并单元格覆盖的位置补齐占位，保留已有单元格及其属性。
func normalizeProtyleTableSpans(table *html.Node) {
	for section := table.FirstChild; nil != section; section = section.NextSibling {
		if atom.Thead != section.DataAtom && atom.Tbody != section.DataAtom && atom.Tfoot != section.DataAtom {
			continue
		}
		var rows []*html.Node
		for row := section.FirstChild; nil != row; row = row.NextSibling {
			if atom.Tr == row.DataAtom {
				rows = append(rows, row)
			}
		}
		// 跨度只在当前行组内生效，数值记录每列被覆盖的结束行。
		occupied := map[int]int{}
		for rowIndex, row := range rows {
			column := 0
			for cell := row.FirstChild; nil != cell; {
				next := cell.NextSibling
				if atom.Td != cell.DataAtom && atom.Th != cell.DataAtom {
					cell = next
					continue
				}
				if !hasDOMClass(cell, "fn__none") {
					for occupied[column] > rowIndex {
						cell.InsertBefore(newProtyleTablePlaceholder(section.DataAtom))
						column++
					}
					colspan := protyleTableSpan(cell, "colspan", 1000)
					rowspan := protyleTableSpan(cell, "rowspan", len(rows)-rowIndex)
					for offset := 0; offset < colspan; offset++ {
						occupied[column+offset] = rowIndex + rowspan
					}
				}
				column++
				cell = next
			}
			for occupied[column] > rowIndex {
				row.AppendChild(newProtyleTablePlaceholder(section.DataAtom))
				column++
			}
		}
	}
}

func protyleTableSpan(cell *html.Node, name string, limit int) int {
	value, err := strconv.Atoi(util.DomAttrValue(cell, name))
	if nil != err || value < 1 {
		return 1
	}
	if value > limit {
		return limit
	}
	return value
}

func newProtyleTablePlaceholder(section atom.Atom) *html.Node {
	cell := newPlaceholderTD()
	if atom.Thead == section {
		cell.DataAtom = atom.Th
		cell.Data = atom.Th.String()
	}
	return cell
}
