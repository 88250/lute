// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// +build !javascript

package util

import "unsafe"

// BytesToStr 快速转换 []byte 为 string。
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// StrToBytes 快速转换 string 为 []byte。
func StrToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// BytesShowLength 获取字节数组展示为UTF8字符串时的长度
func BytesShowLength(bytes []byte) int {
	length := 0
	for i := 0; i < len(bytes); i++ {
		//按位与 11000000 为 10000000 则表示为utf8字节首位
		if (bytes[i] & 0xc0) != 0x80 {
			if bytes[i] < 0x7f {
				length++
			} else {
				length += 2
			}
		}
	}
	return length
}
