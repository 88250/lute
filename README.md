<p align = "center">
<img alt="Lute" src="https://user-images.githubusercontent.com/873584/65254932-99b17300-db2f-11e9-815a-b0b63e812b89.png">
<br><br>
一款结构化的 Markdown 引擎，为未来而构建<br><br>
<em>
千呼万唤始出来&emsp;犹抱琵琶半遮面<br>
转轴拨弦三两声&emsp;未成曲调先有情
</em>
<br><br>
<a title="Build Status" target="_blank" href="https://travis-ci.org/b3log/lute"><img src="https://img.shields.io/travis/b3log/lute.svg?style=flat-square"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/b3log/lute"><img src="https://goreportcard.com/badge/github.com/b3log/lute?style=flat-square"></a>
<a title="Coverage Status" target="_blank" href="https://coveralls.io/repos/github/b3log/lute/badge.svg?branch=master"><img src="https://img.shields.io/coveralls/github/b3log/lute.svg?style=flat-square&color=CC9933"></a>
<a title="Code Size" target="_blank" href="https://github.com/b3log/lute"><img src="https://img.shields.io/github/languages/code-size/b3log/lute.svg?style=flat-square"></a>
<a title="MulanPSL" target="_blank" href="https://github.com/b3log/lute/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MulanPSL-orange.svg?style=flat-square"></a>
<br>
<a title="GitHub Commits" target="_blank" href="https://github.com/b3log/lute/commits/master"><img src="https://img.shields.io/github/commit-activity/m/b3log/lute.svg?style=flat-square"></a>
<a title="Last Commit" target="_blank" href="https://github.com/b3log/lute/commits/master"><img src="https://img.shields.io/github/last-commit/b3log/lute.svg?style=flat-square&color=FF9900"></a>
<a title="GitHub Pull Requests" target="_blank" href="https://github.com/b3log/lute/pulls"><img src="https://img.shields.io/github/issues-pr-closed/b3log/lute.svg?style=flat-square&color=FF9966"></a>
<a title="Hits" target="_blank" href="https://github.com/b3log/hits"><img src="https://hits.b3log.org/b3log/lute.svg"></a>
<br><br>
<a title="GitHub Watchers" target="_blank" href="https://github.com/b3log/lute/watchers"><img src="https://img.shields.io/github/watchers/b3log/lute.svg?label=Watchers&style=social"></a>&nbsp;&nbsp;
<a title="GitHub Stars" target="_blank" href="https://github.com/b3log/lute/stargazers"><img src="https://img.shields.io/github/stars/b3log/lute.svg?label=Stars&style=social"></a>&nbsp;&nbsp;
<a title="GitHub Forks" target="_blank" href="https://github.com/b3log/lute/network/members"><img src="https://img.shields.io/github/forks/b3log/lute.svg?label=Forks&style=social"></a>&nbsp;&nbsp;
<a title="Author GitHub Followers" target="_blank" href="https://github.com/88250"><img src="https://img.shields.io/github/followers/88250.svg?label=Followers&style=social"></a>
</p>

## 💡 简介

[Lute](https://github.com/b3log/lute) 是一款结构化的 Markdown 引擎，完整实现了最新的 [GFM](https://github.github.com/gfm/)/[CommonMark](https://commonmark.org) 规范，对中文语境支持更好。

## 📽️ 背景

之前我一直在使用其他 Markdown 引擎，它们或多或少都有些“瑕疵”：

* 对标准规范的支持不一致
* 对“怪异”文本处理非常耗时，甚至挂死
* 对中文支持不够好

Lute 的目标是构建一个结构化的 Markdown 引擎，实现 GFM/CM 规范并对中文提供更好的支持。所谓的“结构化”指的是从输入的 MD 文本构建抽象语法树，通过操作树来进行 HTML 输出、原文格式化等。
实现规范是为了保证 Markdown 渲染不存在二义性，让同一份 Markdown 文本可以在实现规范的 Markdown 引擎处理后得到一样的结果，这一点非常重要。

实现规范的引擎并不多，我想试试看自己能不能写上一个，这也是 Lute 的动机之一。关于如何实现一个 Markdown 引擎，网上众说纷纭：

* 有的人说 Markdown 适合用正则解析，因为文法规则太简单
* 也有的人说 Markdown 可以用编译原理来处理，正则太难维护

我赞同后者，因为正则确实太难维护而且运行效率较低。最重要的原因是符合 GFM/CM 规范的 Markdown 引擎的核心解析算法不可能用正则写出来，因为规范定义的规则实在是太复杂了。

最后，还有一个很重要的动机就是 B3log 开源社区需要一款自己的 Markdown 引擎：

* [Solo](https://github.com/b3log/solo)、[Pipe](https://github.com/b3log/pipe)、[Sym](https://github.com/b3log/symphony) 需要效果统一的 Markdown 渲染，并且性能非常重要
* [Vditor](https://github.com/b3log/vditor) 需要一款结构化的引擎作为支撑以实现下一代的 Markdown 编辑器

## ✨  特性

* 实现最新版 GFM/CM 规范
* 零正则，非常快
* 代码块语法高亮
* 更好地支持中文语境
* Emoji 解析
* 支持 Markdown 格式化

## 🗃 案例

* [黑客派](https://hacpai.com)，可以到[这里](https://hacpai.com/guide/markdown)进行在线测试

## 🇨🇳 中文语境优化

* 自动链接识别加强
* 在中西文间自动插入空格
* 术语拼写修正

## ♍ 格式化

格式化功能可将“不整洁”的 Markdown 文本格式化为统一风格，在需要公共编辑的场景下，统一的排版风格能让大家更容易协作。

<details>
<summary>点此展开格式化示例。</summary>
<br>
Markdown 原文：

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
    多以建议大家用 ``` 围栏代码块
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

|表格列a|表格列b|       表格列c   |
:---           |:---------------:|--:
第1列开头不要竖线      |   第2列   |第3列结尾不要竖线
                                 ||这个表格看得我眼都花了|

**以上就是为什么我们需要Markdown Format，而且是带中西文自动空格的格式化。**
````````

格式化后：

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
多以建议大家用 ``` 围栏代码块
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

|表格列 a|表格列 b|表格列 c|
|:---|:---:|---:|
|第 1 列开头不要竖线|第 2 列|第 3 列结尾不要竖线|
||这个表格看得我眼都花了||

**以上就是为什么我们需要 Markdown Format，而且是带中西文自动空格的格式化。**

````````

</details>

## ✍️ 术语修正

Markdown 原文：

```markdown
在github上做开源项目是一件很开心的事情，请不要把Github拼写成`github`哦！

特别是简历中千万不要出现这样的情况：

> 熟练使用JAVA、Javascript、GIT，对android、ios开发有一定了解，熟练使用Mysql、postgresql数据库。
```

修正后：

```markdown
在 GitHub 上做开源项目是一件很开心的事情，请不要把 GitHub 拼写成`github`哦！

特别是简历中千万不要出现这样的情况：

> 熟练使用 Java、JavaScript、Git，对 Android、iOS 开发有一定了解，熟练使用 MySQL、PostgreSQL 数据库。
```

## ⚡ 性能

1. 主要对比了 4 款 Go 实现的 Markdown 引擎：[Lute](https://github.com/b3log/lute)、[golang-commonmark](https://gitlab.com/golang-commonmark/markdown)、[goldmark](https://github.com/yuin/goldmark) 和 [Blackfriday](https://github.com/russross/blackfriday)
2. 均已关闭 Typographer 相关配置，因为在排版优化上各库的功能差异较大，对比性能意义不大
3. 均已开启 GFM 支持，因为 GFM 在实际使用场景下是必备的，仅测试 CommonMark 支持意义不大 
4. Lute 在多核平台上性能优势比较明显，因为 Lute 使用了并行解析算法
5. Blackfriday 没有实现 GFM 所以性能看上去更好一些
6. 基准测试数据 [*CommonMark Spec*](https://github.com/commonmark/commonmark-spec-web/blob/gh-pages/0.29/spec.txt) ~197K，参数 `-test.cpu 2,4,8 -test.benchmem`

```text
BenchmarkLute-2                      200           6058789 ns/op         5356988 B/op      24767 allocs/op
BenchmarkLute-4                      300           4517919 ns/op         5349500 B/op      24755 allocs/op
BenchmarkLute-8                      300           4258608 ns/op         5348686 B/op      24744 allocs/op
BenchmarkGolangCommonMark-2          300           4916845 ns/op         2914352 B/op      18423 allocs/op
BenchmarkGolangCommonMark-4          300           4807139 ns/op         2923699 B/op      18424 allocs/op
BenchmarkGolangCommonMark-8          300           4827190 ns/op         2930485 B/op      18425 allocs/op
BenchmarkGoldMark-2                  300           5894211 ns/op         2308201 B/op      15367 allocs/op
BenchmarkGoldMark-4                  300           5724676 ns/op         2315634 B/op      15368 allocs/op
BenchmarkGoldMark-8                  300           5771226 ns/op         2314451 B/op      15368 allocs/op
BenchmarkBlackFriday-2               300           3982679 ns/op         3296658 B/op      20046 allocs/op
BenchmarkBlackFriday-4               500           3734009 ns/op         3304971 B/op      20048 allocs/op
BenchmarkBlackFriday-8               500           3736069 ns/op         3313467 B/op      20049 allocs/op
```

另外，JavaScript 写的 [markdown-it](https://github.com/markdown-it/markdown-it) 循环渲染 300 次，平均每次调用耗时 ~10ms，耗时大致是 golang 实现的两倍，测试代码[见此](https://github.com/b3log/markdown-http/blob/master/benchmark/test.js)。

## 💪 健壮性

Lute 承载了[黑客派](https://hacpai.com)上的所有 Markdown 处理，每天处理数十万请求，运行表现稳定。

## 🔒 安全性

Lute 没有实现实现 GFM 中的 [Disallowed Raw HTML (extension)](https://github.github.com/gfm/#disallowed-raw-html-extension-)，因为该扩展还是存在一定漏洞（比如没有处理 `<input>`）。
建议通过其他库（比如 [bluemonday](https://github.com/microcosm-cc/bluemonday)）来进行 HTML 安全过滤，这样也能更好地适配应用场景。

## 🛠️ 使用

有三种方式使用 Lute：
        
1. 后端：用 Go 语言的话引入 `github.com/b3log/lute` 包即可 
2. 后端：将 Lute 启动为一个 HTTP 服务进程供其他进程调用，具体请参考[这里](https://github.com/b3log/lute-http)
3. 前端：引入 js 目录下的 lute.min.js 即可，支持 Node.js 

### Go

引入 Lute 库：
```shell
go get -u github.com/b3log/lute
```

最小化可工作示例：

```go
package main

import (
	"fmt"

	"github.com/b3log/lute"
)

func main() {
	luteEngine := lute.New() // 默认已经启用 GFM 支持以及中文语境优化
	html, err := luteEngine.MarkdownStr("demo", "**Lute** - A structured markdown engine.")
	if nil != err {
		panic(err)
	}
	fmt.Println(html)
	// <p><strong>Lute</strong> - A structured Markdown engine.</p>
}
```

关于代码块语法高亮：

* 默认使用外部样式表，主题为 github.css，可从 chroma-styles 目录下拷贝该样式文件到项目中引入
* 可通过 `lutenEngine.SetCodeSyntaxHighlightXXX()` 来指定高亮相关参数，比如是否启用内联样式、行号以及主题

### JavaScript
 
简单示例可参考 js 目录下的 demo，结合前端编辑器的完整用法请参考 [Vditor 中的示例](https://github.com/b3log/vditor/tree/master/demo)。

![Vditor](https://user-images.githubusercontent.com/970828/64320104-624fac00-cff0-11e9-8727-0ad51a6f71c0.png)

一些细节：

1. lute.js 没有内置语法高亮特性
2. lute.js 编译压缩后大小为 ~800KB，通过 `brotli -o lute.min.js.br lute.min.js` 压缩后大小 ~110KB，常规 GZip 压缩后大小 ~150KB

## 📜 文档

* [《提问的智慧》精读注解版](https://hacpai.com/article/1536377163156)
* [CommonMark 规范要点解读](https://hacpai.com/article/1566893557720)
* [Lute 实现后记](https://hacpai.com/article/1567062979327)

## 🏘️ 社区

* [讨论区](https://hacpai.com/tag/lute)
* [报告问题](https://github.com/b3log/lute/issues/new/choose)

## 📄 授权

Lute 使用 [木兰宽松许可证, 第1版](http://license.coscl.org.cn/MulanPSL) 开源协议。

## 🙏 鸣谢

Lute 的诞生离不开以下开源项目，在此对这些项目的贡献者们致敬！

* [commonmark.js](https://github.com/commonmark/commonmark.js)：该项目是 CommonMark 官方参考实现的 JavaScript 版，Lute 参考了其解析器实现部分
* [mdast](https://github.com/syntax-tree/mdast)：该项目介绍了一种 Markdown 抽象语法树结构的表现形式，Lute 的 AST 在初始设计阶段参考了该项目
* [goldmark](https://github.com/yuin/goldmark)：另一款用 golang 写的 Markdown 引擎，Lute 参考了其树遍历实现部分
* [golang-commonmark](https://gitlab.com/golang-commonmark/markdown)：另一款用 golang 写的 Markdown 引擎，Lute 参考了其 URL 编码以及 HTML 转义算法
* [Chroma](https://github.com/alecthomas/chroma)：用 golang 写的语法高亮引擎
* [中文文案排版指北](https://github.com/sparanoid/chinese-copywriting-guidelines)：统一中文文案、排版的相关用法，降低团队成员之间的沟通成本，增强网站气质
* [autocorrect](https://github.com/studygolang/autocorrect)：自动给中英文之间加入空格、术语拼写修正
* [GopherJS](https://github.com/gopherjs/gopherjs)：将 Go 代码编译成 JavaScript 代码

---

## 👍 开源项目推荐

* 如果你需要集成一个浏览器端的 Markdown 编辑器，可以考虑使用 [Vditor](https://github.com/b3log/vditor)
* 如果你需要搭建一个个人博客系统，可以考虑使用 [Solo](https://github.com/b3log/solo)
* 如果你需要搭建一个社区平台，可以考虑使用 [Sym](https://github.com/b3log/symphony)
* 欢迎加入我们的小众开源社区，详情请看[这里](https://hacpai.com/article/1463025124998)
