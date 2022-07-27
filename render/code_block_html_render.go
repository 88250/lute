// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

//go:build !javascript
// +build !javascript

package render

import (
	"bytes"
	"go/format"
	"strings"

	"github.com/88250/lute/html"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	chromalexers "github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func (r *HtmlRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()

	if !node.IsFencedCodeBlock {
		if entering {
			// 缩进代码块处理
			rendered := false
			tokens := node.FirstChild.Tokens
			if r.Options.CodeSyntaxHighlight {
				rendered = highlightChroma(node, tokens, "", r)
				if !rendered {
					tokens = html.EscapeHTML(tokens)
					r.Write(tokens)
				}
			} else {
				var attrs [][]string
				r.handleKramdownBlockIAL(node)
				attrs = append(attrs, node.KramdownIAL...)
				r.Tag("pre", attrs, false)
				r.WriteString("<code>")
				tokens = html.EscapeHTML(tokens)
				r.Write(tokens)
			}
			r.WriteString("</code></pre>")
			return ast.WalkSkipChildren
		} else {
			return ast.WalkContinue
		}
	}
	return ast.WalkContinue
}

// renderCodeBlockCode 进行代码块 HTML 渲染，实现语法高亮。
func (r *HtmlRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	var language string
	if 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = util.BytesToStr(infoWords[0])
	}
	preDiv := NoHighlight(language)
	if entering {
		var attrs [][]string
		r.handleKramdownBlockIAL(node.Parent)
		attrs = append(attrs, node.Parent.KramdownIAL...)

		tokens := node.Tokens
		if 0 < len(node.Previous.CodeBlockInfo) {
			rendered := false
			if isGo(language) {
				// Go 代码块自动格式化 https://github.com/b3log/lute/issues/37
				if buf, err := format.Source(tokens); nil == err {
					tokens = buf
				}
			}

			if "mindmap" == language {
				json := EChartsMindmap(tokens)
				r.WriteString("<div data-code=\"")
				r.Write(json)
				r.WriteString("\" class=\"language-mindmap\">")
				r.Write(html.EscapeHTML(tokens))
				rendered = true
			} else {
				if r.Options.CodeSyntaxHighlight && !preDiv {
					rendered = highlightChroma(node.Parent, tokens, language, r)
				}
			}

			if !rendered {
				if preDiv {
					r.WriteString("<div class=\"language-")
				} else {
					r.Tag("pre", attrs, false)
					r.WriteString("<code class=\"language-")
				}
				r.WriteString(language)
				r.WriteString("\">")
				tokens = html.EscapeHTML(tokens)
				r.Write(tokens)
			}
		} else {
			rendered := false
			if r.Options.CodeSyntaxHighlight {
				rendered = highlightChroma(node.Parent, tokens, "", r)
				if !rendered {
					tokens = html.EscapeHTML(tokens)
					r.Write(tokens)
				}
			} else {
				r.Tag("pre", attrs, false)
				if r.Options.CodeSyntaxHighlightDetectLang {
					language := detectLanguage(tokens)
					if "" != language {
						r.WriteString("<code class=\"language-" + language)
					} else {
						r.WriteString("<code>")
					}
				} else {
					r.WriteString("<code>")
				}
				tokens = html.EscapeHTML(tokens)
				r.Write(tokens)
			}
		}
	} else {
		if preDiv {
			r.WriteString("</div>")
		} else {
			r.WriteString("</code></pre>")
		}
	}
	return ast.WalkContinue
}

func highlightChroma(codeNode *ast.Node, tokens []byte, language string, r *HtmlRenderer) (rendered bool) {
	var attrs [][]string
	r.handleKramdownBlockIAL(codeNode)
	attrs = append(attrs, codeNode.KramdownIAL...)

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
		if !r.Options.CodeSyntaxHighlightInlineStyle {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithClasses(true))
		}
		if r.Options.CodeSyntaxHighlightLineNum {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithLineNumbers(true))
		}
		formatter := chromahtml.New(chromahtmlOpts...)
		style := styles.Get(r.Options.CodeSyntaxHighlightStyleName)
		var b bytes.Buffer
		if err = formatter.Format(&b, style, iterator); nil == err {
			if !r.Options.CodeSyntaxHighlightInlineStyle {
				r.Tag("pre", attrs, false)
			} else {
				attrs = append(attrs, []string{"style", chromahtml.StyleEntryToCSS(style.Get(chroma.Background))})
				r.Tag("pre", attrs, false)
			}
			if "" != language {
				r.WriteString("<code class=\"language-" + language)
			} else {
				r.WriteString("<code class=\"")
			}
			if !r.Options.CodeSyntaxHighlightInlineStyle {
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
