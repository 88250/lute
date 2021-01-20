// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
)

func GitConflictContinue(gitConflictBlock *ast.Node, context *Context) int {
	if context.isGitConflictClose() {
		context.finalize(gitConflictBlock)
		return 2
	}
	return 0
}

func (context *Context) gitConflictFinalize(gitConflictBlock *ast.Node) {
	tokens := gitConflictBlock.Tokens
	lines := bytes.Split(tokens, []byte("=======\n"))
	parts := lines[0]
	localParts := bytes.Split(parts, []byte("\n"))
	openMarkerTokens := localParts[0]
	local := bytes.Join(localParts[1:], []byte("\n"))
	local = bytes.TrimSpace(local)
	remote := bytes.TrimSpace(lines[1])
	closeMarkerTokens := bytes.TrimSpace(context.currentLine)
	gitConflictBlock.Tokens = nil
	gitConflictBlock.AppendChild(&ast.Node{Type: ast.NodeGitConflictOpenMarker, Tokens: openMarkerTokens})
	gitConflictBlock.AppendChild(&ast.Node{Type: ast.NodeGitConflictLocalContent, Tokens: local})
	gitConflictBlock.AppendChild(&ast.Node{Type: ast.NodeGitConflictSepMarker})
	gitConflictBlock.AppendChild(&ast.Node{Type: ast.NodeGitConflictRemoteContent, Tokens: remote})
	gitConflictBlock.AppendChild(&ast.Node{Type: ast.NodeGitConflictCloseMarker, Tokens: closeMarkerTokens})
}

func (t *Tree) parseGitConflict() (ok bool) {
	return bytes.HasPrefix(t.Context.currentLine, []byte("<<<<<<<"))
}

func (context *Context) isGitConflictClose() bool {
	return bytes.HasPrefix(context.currentLine, []byte(">>>>>>>"))
}
