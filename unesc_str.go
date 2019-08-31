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

// +build !wasm

package lute

import "html"

func unescapeString(tokens items) (ret items) {
	if nil == tokens {
		return
	}

	tokens = toItems(html.UnescapeString(fromItems(tokens))) // FIXME: 此处应该用内部的实体转义方式
	length := len(tokens)
	ret = make(items, 0, length)
	for i := 0; i < length; i++ {
		if tokens.isBackslashEscapePunct(i) {
			ret = ret[:len(ret)-1]
		}
		ret = append(ret, tokens[i])
	}
	return
}
