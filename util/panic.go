// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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

import (
	"errors"
	"runtime/debug"
)

// RecoverPanic recovers a panic.
func RecoverPanic(err *error) {
	if e := recover(); nil != e {
		stack := debug.Stack()
		errMsg := ""
		switch x := e.(type) {
		case error:
			errMsg = x.Error()
		case string:
			errMsg = x
		default:
			errMsg = "unknown panic"
		}
		if nil != err {
			*err = errors.New("PANIC RECOVERED: " + errMsg + "\n\t" + string(stack) + "\n")
		}
	}
}
