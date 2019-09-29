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

// +build js

package lute

// item 描述了词法分析的一个 token。
type item struct {
	term byte // 源码字节
	ln   int  // 源码行号
	col  int  // 源码列号
}

// items 定义了 token 数组。
type items []item

// nilItem 返回一个空值 token。
func nilItem() item {
	return item{term: 0}
}

// isNilItem 判断 item 是否为空值。
func isNilItem(item item) bool {
	return 0 == item.term
}

// newItem 构造一个 token。
func newItem(term byte, ln, col int) item {
	return item{term: term, ln: ln, col: col}
}

// term 返回 item 的词素。
func term(item item) byte {
	return item.term
}

func setTerm(tokens *items, i int, term byte) {
	(*tokens)[i].term = term
}

// strToItems 将 str 转为 items。
func strToItems(str string) (ret items) {
	ret = make(items, 0, len(str))
	length := len(str)
	for i := 0; i < length; i++ {
		ret = append(ret, item{term: str[i]})
	}
	return
}

// itemsToStr 将 items 转为 string。
func itemsToStr(items items) string {
	return string(itemsToBytes(items))
}

// itemsToBytes 将 items 转为 []byte。
func itemsToBytes(items items) (ret []byte) {
	length := len(items)
	for i := 0; i < length; i++ {
		ret = append(ret, term(items[i]))
	}
	return
}

// bytesToItems 将 bytes 转为 items。
func bytesToItems(bytes []byte) (ret items) {
	ret = make(items, 0, len(bytes))
	length := len(bytes)
	for i := 0; i < length; i++ {
		ret = append(ret, item{term: bytes[i]})
	}
	return
}
