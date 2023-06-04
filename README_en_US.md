<p align = "center">
<img alt="Lute" src="https://b3log.org/images/brand/lute-128.png">
<br><br>
A structured Markdown engine that supports Go and JavaScript
<br><br>
<a title="Build Status" target="_blank" href="https://github.com/88250/lute/actions/workflows/gotest.yml"><img src="https://img.shields.io/github/actions/workflow/status/88250/lute/gotest.yml?style=flat-square"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/88250/lute"><img src="https://goreportcard.com/badge/github.com/88250/lute?style=flat-square"></a>
<a title="Coverage Status" target="_blank" href="https://coveralls.io/github/88250/lute"><img src="https://img.shields.io/coveralls/github/88250/lute.svg?style=flat-square&color=CC9933"></a>
<a title="Code Size" target="_blank" href="https://github.com/88250/lute"><img src="https://img.shields.io/github/languages/code-size/88250/lute.svg?style=flat-square"></a>
<a title="MulanPSL" target="_blank" href="https://github.com/88250/lute/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MulanPSL-orange.svg?style=flat-square"></a>
<br>
<a title="GitHub Commits" target="_blank" href="https://github.com/88250/lute/commits/master"><img src="https://img.shields.io/github/commit-activity/m/88250/lute.svg?style=flat-square"></a>
<a title="Last Commit" target="_blank" href="https://github.com/88250/lute/commits/master"><img src="https://img.shields.io/github/last-commit/88250/lute.svg?style=flat-square&color=FF9900"></a>
<a title="GitHub Pull Requests" target="_blank" href="https://github.com/88250/lute/pulls"><img src="https://img.shields.io/github/issues-pr-closed/88250/lute.svg?style=flat-square&color=FF9966"></a>
<a title="Hits" target="_blank" href="https://github.com/88250/hits"><img src="https://hits.b3log.org/88250/lute.svg"></a>
<br><br>
<a title="GitHub Watchers" target="_blank" href="https://github.com/88250/lute/watchers"><img src="https://img.shields.io/github/watchers/88250/lute.svg?label=Watchers&style=social"></a>  
<a title="GitHub Stars" target="_blank" href="https://github.com/88250/lute/stargazers"><img src="https://img.shields.io/github/stars/88250/lute.svg?label=Stars&style=social"></a>  
<a title="GitHub Forks" target="_blank" href="https://github.com/88250/lute/network/members"><img src="https://img.shields.io/github/forks/88250/lute.svg?label=Forks&style=social"></a>  
<a title="Author GitHub Followers" target="_blank" href="https://github.com/88250"><img src="https://img.shields.io/github/followers/88250.svg?label=Followers&style=social"></a>
</p>

<p align="center">
<a href="https://github.com/88250/lute/blob/master/README.md">中文</a>
</p>

## 💡 Introduction

[Lute](https://github.com/88250/lute) is a structured Markdown engine that fully implements the latest [GFM](https://github.github.com/gfm/) / [CommonMark](https://commonmark.org) standard, better support for Chinese context.

Welcome to [Lute Official Discussion Forum](https://ld246.com/tag/lute) to learn more.

## 📽️ Background

I have been using other Markdown engines before, and they are more or less "defective":

* Inconsistent support for standard specifications
* The processing of "strange" text is very time-consuming and even hangs
* Support for Chinese is not good enough

Lute's goal is to build a structured Markdown engine that implements GFM/CM specifications and provides better support for Chinese. The so-called "structured" refers to the construction of an abstract syntax tree from the input MD text, HTML output, text formatting, etc. through the operation tree.
The realization of the specification is to ensure that there is no ambiguity in Markdown rendering, so that the same Markdown text can be processed by the Markdown engine to achieve the same result, which is very important.

There are not many engines that implement specifications. I want to see if I can write one, which is one of Lute's motivations. There are many opinions on the Internet about how to implement a Markdown engine:

* Some people say that Markdown is suitable for regular analysis, because the grammar rules are too simple
* Some people say that Markdown can be handled by the compilation principle, but the rule is too difficult to maintain

I agree with the latter, because regular expressions is indeed too difficult to maintain and has low operating efficiency. The most important reason is that the core parsing algorithm of the Markdown engine that conforms to the GFM/CM specification cannot be written in regular, because the rules defined by the specification are too complicated.

Finally, another important motivation is that the B3log open source community needs its own Markdown engine:

* [Solo](https://github.com/88250/solo), [Pipe](https://github.com/88250/pipe), [Sym](https://github.com/88250/symphony ) Markdown rendering with uniform effects is required, and performance is very important
* [Vditor](https://github.com/Vanessa219/vditor) needs a structured engine as support to achieve the next generation of Markdown editor

## ✨  Features

* Implement the latest version of GFM/CM specifications
* Zero regular expressions, very fast
* Built-in code block syntax highlighting
* Better support for Chinese context
* Terminology spelling correction
* Markdown format
* Emoji analysis
* HTML to Markdown
* Custom rendering function
* Support JavaScript

## 🗃 Showcases

* [LianDi](https://ld246.com)
* [SiYuan](https://github.com/siyuan-note/siyuan)
* [Vditor](https://github.com/Vanessa219/vditor)
* [Sym](https://github.com/88250/symphony)
* [Solo](https://github.com/88250/solo)
* [Pipe](https://github.com/88250/pipe)

## 🇨🇳 Chinese context optimization

* Enhanced automatic link recognition
* Automatically insert spaces between Chinese and Western languages

## ♍ Format

The formatting function can format "untidy" Markdown text into a unified style. In scenarios that require public editing, a unified typography style makes it easier for everyone to collaborate.

<details>
<summary>Click here to expand the formatting example.</summary>
<br>
Markdown: 

````````markdown
# ATX 标题也有可能需要格式化的 ##
一个简短的段落。

Setext 说实话我不喜欢 Setext 标题
----
0. 有序列表可以从 0 开始
0. 应该自增序号的
1.   对齐对齐对齐

我们再来看看另一个有序列表。
1. 没空行的情况下序号要从 1 开始才能打断段落开始一个新列表
3. 虽然乱序不影响渲染
2. 但是随意写序号容易引起误解

试下贴段代码：
```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, 世界")
}
```
对了，缩进代码块建议换成围栏代码块：

    缩进代码块太隐晦了
    也没法指定编程语言，容易导致代码高亮失效
    所以建议大家用 ``` 围栏代码块
试下围栏代码块匹配场景：
````markdown
围栏代码块只要开头的 ` 和结束的 ` 数量匹配即可，这样可以实现在围栏代码块中显示围栏代码块：
```
这里只有 3 个 `，所以不会匹配markdown代码块结束
```
下面匹配到就真的结束了。
````
以上块级内容都挤在一坨了，插入合理的空行也很有必要。


但是过多的空行分段也不好啊，用来分段的话一个空行就够了。



接下来让我们试试稍微复杂点的场景，比如列表项包含多个段落的情况：
1. 列表项中的第一段

   这里是第二个段落，贴段代码：
   ```markdown
   要成为Markdown程序员并不容易，同理PPT架构师也是。
   注意代码块中的中西文间并没有插入空格。
   ```
   这里是最后一段了。
1. 整个有序列表是“松散”的：列表项内容要用 `<p>` 标签

最后，我们试下对 GFM 的格式化支持：

|col1|col2  |       col3   |
---           |---------------|--
col1 without left pipe      |   this is col2   | col3 without right pipe
                                 ||need align cell|

**以上就是为什么我们需要Markdown Format，而且是带中西文自动空格的格式化。**
````````

Formatted:

````````markdown
# ATX 标题也有可能需要格式化的

一个简短的段落。

## Setext 说实话我不喜欢 Setext 标题

0. 有序列表可以从 0 开始
1. 应该自增序号的
2. 对齐对齐对齐

我们再来看看另一个有序列表。

1. 没空行的情况下序号要从 1 开始才能打断段落开始一个新列表
2. 虽然乱序不影响渲染
3. 但是随意写序号容易引起误解

试下贴段代码：

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, 世界")
}
```

对了，缩进代码块建议换成围栏代码块：

```
缩进代码块太隐晦了
也没法指定编程语言，容易导致代码高亮失效
所以建议大家用 ``` 围栏代码块
```

试下围栏代码块匹配场景：

````markdown
围栏代码块只要开头的 ` 和结束的 ` 数量匹配即可，这样可以实现在围栏代码块中显示围栏代码块：
```
这里只有 3 个 `，所以不会匹配markdown代码块结束
```
下面匹配到就真的结束了。
````

以上块级内容都挤在一坨了，插入合理的空行也很有必要。

但是过多的空行分段也不好啊，用来分段的话一个空行就够了。

接下来让我们试试稍微复杂点的场景，比如列表项包含多个段落的情况：

1. 列表项中的第一段

   这里是第二个段落，贴段代码：

   ```markdown
   要成为Markdown程序员并不容易，同理PPT架构师也是。
   注意代码块中的中西文间并没有插入空格。
   ```

   这里是最后一段了。
2. 整个有序列表是“松散”的：列表项内容要用 `<p>` 标签

最后，我们试下对 GFM 的格式化支持：

| col1                   | col2            | col3                    |
| ---------------------- | --------------- | ----------------------- |
| col1 without left pipe | this is col2    | col3 without right pipe |
|                        | need align cell |                         |

**以上就是为什么我们需要 Markdown Format，而且是带中西文自动空格的格式化。**
````````

</details>

## ✍️ Terminology revision

Markdown: 

```markdown
Doing open source projects on github is a very happy thing, please don't spell Github as `github`!

In particular, this should never happen in your resume:

> Proficient in using JAVA, Javascript, GIT, have a certain understanding of android, ios development, proficient in using Mysql, postgresql database.
```

after fixing:

```markdown
Doing open source projects on GitHub is a very happy thing, please don't spell Github as `github`!

In particular, this should never happen in your resume:

> Proficient in using Java, JavaScript, Git, have a certain understanding of Android, iOS development, proficient in using MySQL, PostgreSQL database.
```

## ⚡ Performance

Please see [Golang markdown engine performance benchmark](https://ld246.com/article/1574570835061).

## 💪 Robustness

Lute carries all Markdown processing on [LianDi](https://ld246.com), processes millions of parsing and rendering requests every day, and runs stably.

## 🔒 Safety

Lute does not implement [Disallowed Raw HTML (extension)](https://github.github.com/gfm/#disallowed-raw-html-extension-) in GFM, because the extension still has certain vulnerabilities `<input>`).
It is recommended to use other libraries (such as [bluemonday](https://github.com/microcosm-cc/bluemonday)) for HTML security filtering, so that it can better adapt to the application scenario.

## 🛠️ Usages

There are three ways to use Lute:

1. Backend: Introduce `github.com/88250/lute` package in Go language
2. Backend: Start Lute as an HTTP service process for other processes to call, please refer to [here](https://github.com/88250/lute-http)
3. Front end: Introduce lute.min.js in the js directory, support Node.js

### Go

Introduce the Lute library:

```shell
go get -u github.com/88250/lute
```

Working example of minimization:

```go
package main

import (
	"fmt"

	"github.com/88250/lute"
)

func main() {
	luteEngine := lute.New() // GFM support and Chinese context optimization have been enabled by default
	html := luteEngine.MarkdownStr("demo", "**Lute** - A structured markdown engine.")
	fmt.Println(html)
	// <p><strong>Lute</strong> - A structured Markdown engine.</p>
}
```

About code block syntax highlighting:

* The external style sheet is used by default, and the theme is github.css. You can copy the style file from the chroma-styles directory to the project and import it
* You can specify highlight-related parameters such as whether to enable inline styles, line numbers, and themes through `lutenEngine.SetCodeSyntaxHighlightXXX ()`

### JavaScript

For a simple example, please refer to the demo in the JavaScript directory. For the complete usage of the front-end editor, please refer to [Demo in Vditor](https://github.com/Vanessa219/vditor/tree/master/demo)

![Vditor](https://b3logfile.com/file/2020/02/%E6%88%AA%E5%9B%BE%E4%B8%93%E7%94%A8-ef21ef12.png)

Some details:

1. lute.js has no built-in syntax highlighting feature
2. The size of lute.js after compilation is ~3.5MB, the size after regular GZip compression is ~500KB

## 📜 Documentation

* [Interpretation of CommonMark specifications](https://ld246.com/article/1566893557720)
* [Lute Implementation Postscript](https://ld246.com/article/1567062979327)
* [Markdown parsing and Markdown AST](https://ld246.com/article/1587637426085)

## 🏘️ Community

* [Forum](https://ld246.com/tag/lute)
* [Issues](https://github.com/88250/lute/issues/new/choose)

## 📄 License

Lute uses the [Mulan Permissive Software License，Version 2](http://license.coscl.org.cn/MulanPSL2) open source license.

## 🙏 Acknowledgement

* [commonmark.js](https://github.com/commonmark/commonmark.js): CommonMark parser and renderer in JavaScript
* [goldmark](https://github.com/yuin/goldmark)：A markdown parser written in Go
* [golang-commonmark](https://gitlab.com/golang-commonmark/markdown): A CommonMark-compliant markdown parser and renderer in Go
* [Chroma](https://github.com/alecthomas/chroma): A general purpose syntax highlighter in pure Go
* [中文文案排版指北](https://github.com/sparanoid/chinese-copywriting-guidelines): Chinese copywriting guidelines for better written communication
* [GopherJS](https://github.com/gopherjs/gopherjs): A compiler from Go to JavaScript for running Go code in a browser
