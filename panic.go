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

// +build !javascript

package lute

import (
	"errors"
	"runtime/debug"
)

// Recover recovers a panic.
func recoverPanic(err *error) {
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
