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
	"strings"
)

func (r *BaseRenderer) EncodeLinkSpace(dest string) string {
	// Improve export of Markdown hyperlink spaces and markers https://github.com/siyuan-note/siyuan/issues/9792
	return strings.ReplaceAll(dest, " ", "%20")
}

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

	if !r.isRelativePath(dest) {
		return dest
	}

	if !bytes.HasPrefix(dest, []byte("assets/")) &&
		(strings.HasPrefix(r.Options.LinkBase, "https://assets.b3logfile.com/siyuan/") || strings.HasPrefix(r.Options.LinkBase, "https://assets.liuyun.io/siyuan/")) {
		return dest
	}

	dest = bytes.ReplaceAll(dest, []byte("%5C"), []byte("\\"))
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
	return !bytes.Contains(dest, []byte(":/")) && !bytes.Contains(dest, []byte(":\\")) && !bytes.Contains(dest, []byte(":%5C"))
}
