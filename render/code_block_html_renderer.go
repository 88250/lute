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

package render

import (
	"bytes"
	"go/format"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	chromalexers "github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// languagesNoHighlight 中定义的语言不要进行代码语法高亮。这些代码块会在前端进行渲染，比如各种图表。
var languagesNoHighlight = []string{"mermaid", "echarts", "abc", "graphviz"}

func (r *HtmlRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if !node.IsFencedCodeBlock {
		// 缩进代码块处理
		r.Newline()
		rendered := false
		tokens := node.FirstChild.Tokens
		if r.Option.CodeSyntaxHighlight {
			rendered = highlightChroma(tokens, "", r)
			if !rendered {
				tokens = util.EscapeHTML(tokens)
				r.Write(tokens)
			}
		} else {
			r.WriteString("<pre><code>")
			tokens = util.EscapeHTML(tokens)
			r.Write(tokens)
		}
		r.WriteString("</code></pre>")
		r.Newline()
		return ast.WalkStop
	}
	r.Newline()
	return ast.WalkContinue
}

// renderCodeBlockCode 进行代码块 HTML 渲染，实现语法高亮。
func (r *HtmlRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if 0 < len(node.Previous.CodeBlockInfo) {
			infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
			language := util.BytesToStr(infoWords[0])
			rendered := false
			if isGo(language) {
				// Go 代码块自动格式化 https://github.com/b3log/lute/issues/37
				if buf, err := format.Source(tokens); nil == err {
					tokens = buf
				}
			}

			if "mindmap" == language {
				json := r.renderMindmap(tokens)
				r.WriteString("<pre><code data-code=\"")
				r.Write(util.EscapeHTML(json))
				r.WriteString("\" class=\"language-mindmap\">")
				r.Write(util.EscapeHTML(tokens))
				rendered = true
			} else {
				if r.Option.CodeSyntaxHighlight && !noHighlight(language) {
					rendered = highlightChroma(tokens, language, r)
				}
			}

			if !rendered {
				r.WriteString("<pre><code class=\"language-")
				r.WriteString(language)
				r.WriteString("\">")
				tokens = util.EscapeHTML(tokens)
				r.Write(tokens)
			}
		} else {
			rendered := false
			if r.Option.CodeSyntaxHighlight {
				rendered = highlightChroma(tokens, "", r)
				if !rendered {
					tokens = util.EscapeHTML(tokens)
					r.Write(tokens)
				}
			} else {
				if r.Option.CodeSyntaxHighlightDetectLang {
					language := detectLanguage(tokens)
					if "" != language {
						r.WriteString("<pre><code class=\"language-" + language)
					} else {
						r.WriteString("<pre><code>")
					}
				} else {
					r.WriteString("<pre><code>")
				}
				tokens = util.EscapeHTML(tokens)
				r.Write(tokens)
			}
		}
		return ast.WalkSkipChildren
	}
	r.WriteString("</code></pre>")
	return ast.WalkStop
}

func highlightChroma(tokens []byte, language string, r *HtmlRenderer) (rendered bool) {
	codeBlock := util.BytesToStr(tokens)
	var lexer chroma.Lexer
	if "" != language {
		lexer = chromalexers.Get(language)
	} else {
		lexer = chromalexers.Analyse(codeBlock)
	}
	if nil == lexer {
		lexer = chromalexers.Fallback
	} else {
		language = lexer.Config().Aliases[0]
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, codeBlock)
	if nil == err {
		chromahtmlOpts := []chromahtml.Option{
			chromahtml.PreventSurroundingPre(true),
			chromahtml.ClassPrefix("highlight-"),
		}
		if !r.Option.CodeSyntaxHighlightInlineStyle {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithClasses(true))
		}
		if r.Option.CodeSyntaxHighlightLineNum {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithLineNumbers(true))
		}
		formatter := chromahtml.New(chromahtmlOpts...)
		style := styles.Get(r.Option.CodeSyntaxHighlightStyleName)
		var b bytes.Buffer
		if err = formatter.Format(&b, style, iterator); nil == err {
			if !r.Option.CodeSyntaxHighlightInlineStyle {
				r.WriteString("<pre>")
			} else {
				r.WriteString("<pre style=\"" + chromahtml.StyleEntryToCSS(style.Get(chroma.Background)) + "\">")
			}
			if "" != language {
				r.WriteString("<code class=\"language-" + language)
			} else {
				r.WriteString("<code class=\"")
			}
			if !r.Option.CodeSyntaxHighlightInlineStyle {
				if "" != language {
					r.WriteByte(lex.ItemSpace)
				}
				r.WriteString("highlight-chroma")
			}
			r.WriteString("\">")
			r.Write(b.Bytes())
			rendered = true
		}
	}
	return
}

func noHighlight(language string) bool {
	for _, langNoHighlight := range languagesNoHighlight {
		if language == langNoHighlight {
			return true
		}
	}
	return false
}

func isGo(language string) bool {
	return strings.EqualFold(language, "go") || strings.EqualFold(language, "golang")
}

// github.com/src-d/enry/v2 不怎么准确
//
//var candidateLangs = []string{
//	"bash", "csharp", "cpp", "css", "go", "html", "xml", "java", "js", "json", "kotlin", "less", "lua", "makefile", "markdown",
//	"nginx", "objc", "php", "properties", "python", "ruby", "rust", "scss", "sql", "shell", "toml", "ts", "yaml", "swift",
//	"dart", "gradle", "julia", "matlab",
//}
//
//func detectLanguage(code []byte) (language string) {
//	language, _ = enry.GetLanguageByClassifier(code, candidateLangs)
//	return
//}

func detectLanguage(code []byte) string {
	lexer := chromalexers.Analyse(util.BytesToStr(code))
	if nil == lexer {
		return ""
	}
	return lexer.Config().Aliases[0]
}
