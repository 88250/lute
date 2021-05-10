// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"

	"github.com/88250/lute/util"
)

func (r *BaseRenderer) LinkPath(dest []byte) []byte {
	dest = r.RelativePath(dest)
	dest = r.PrefixPath(dest)
	return dest
}

func (r *BaseRenderer) PrefixPath(dest []byte) []byte {
	if "" == r.Options.LinkPrefix {
		return dest
	}

	linkPrefix := util.StrToBytes(r.Options.LinkPrefix)
	ret := append(linkPrefix, dest...)
	return ret
}

func (r *BaseRenderer) RelativePath(dest []byte) []byte {
	if "" == r.Options.LinkBase {
		return dest
	}

	// 强制将 %5C 即反斜杠 \ 转换为斜杠 / 以兼容 Windows 平台上使用的路径
	dest = bytes.ReplaceAll(dest, []byte("%5C"), []byte("\\"))
	if !r.isRelativePath(dest) {
		return dest
	}

	linkBase := util.StrToBytes(r.Options.LinkBase)
	if !bytes.HasSuffix(linkBase, []byte("/")) {
		linkBase = append(linkBase, []byte("/")...)
	}
	ret := append(linkBase, dest...)
	if bytes.Equal(linkBase, ret) {
		return []byte("")
	}
	return ret
}

func (r *BaseRenderer) isRelativePath(dest []byte) bool {
	if 1 > len(dest) {
		return true
	}

	if '/' == dest[0] {
		return false
	}
	return !bytes.Contains(dest, []byte(":/")) && !bytes.Contains(dest, []byte(":\\"))
}
