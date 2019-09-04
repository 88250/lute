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

// +build !js

package lute

import (
	"bytes"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	chromalexers "github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// languagesNoHighlight 中定义的语言不要进行代码语法高亮。这些代码块会在前端进行渲染，比如各种图表。
var languagesNoHighlight = []string{"mermaid", "echarts", "abc"}

// renderCodeBlockHTML 进行代码块 HTML 渲染，实现语法高亮。
func (r *Renderer) renderCodeBlockHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		if nil != node.codeBlockInfo {
			infoWords := bytes.Fields(node.codeBlockInfo)
			language := infoWords[0]
			r.writeString("<pre><code class=\"language-")
			r.write(language)
			r.writeString("\">")
			rendered := false
			if r.option.CodeSyntaxHighlight && !noHighlight(fromItems(language)) {
				codeBlock := fromItems(tokens)
				var lexer chroma.Lexer
				if nil != language {
					lexer = chromalexers.Get(string(language))
				} else {
					lexer = chromalexers.Analyse(codeBlock)
				}
				if nil == lexer {
					lexer = chromalexers.Fallback
				}
				iterator, err := lexer.Tokenise(nil, codeBlock)
				if nil == err {
					formatter := chromahtml.New(chromahtml.PreventSurroundingPre(), chromahtml.WithClasses(), chromahtml.ClassPrefix("highlight-"))
					var b bytes.Buffer
					if err = formatter.Format(&b, styles.GitHub, iterator); nil == err {
						r.write(b.Bytes())
						rendered = true
						// 生成 CSS 临时调试用：
						//formatter.WriteCSS(os.Stdout, styles.GitHub)
						//os.Stdout.WriteString("\n")
					}
				}
			}

			if !rendered {
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		} else {
			rendered := false
			if r.option.CodeSyntaxHighlight {
				language := "fallback"
				codeBlock := fromItems(tokens)
				var lexer = chromalexers.Analyse(codeBlock)
				if nil == lexer {
					lexer = chromalexers.Fallback
				}
				language = lexer.Config().Name
				r.writeString("<pre><code class=\"language-" + language + "\">")

				iterator, err := lexer.Tokenise(nil, codeBlock)
				if nil == err {
					formatter := chromahtml.New(chromahtml.PreventSurroundingPre(), chromahtml.WithClasses(), chromahtml.ClassPrefix("highlight-"))
					var b bytes.Buffer
					if err = formatter.Format(&b, styles.GitHub, iterator); nil == err {
						r.write(b.Bytes())
						rendered = true
					}
				}

				if !rendered {
					tokens = escapeHTML(tokens)
					r.write(tokens)
				}
			} else {
				r.writeString("<pre><code>")
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		}
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	r.newline()
	return WalkContinue, nil
}

func noHighlight(language string) bool {
	for _, langNoHighlight := range languagesNoHighlight {
		if language == langNoHighlight {
			return true
		}
	}
	return false
}
