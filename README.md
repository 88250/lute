# Lute

> 千呼万唤始出来，犹抱琵琶半遮面。转轴拨弦三两声，未成曲调先有情。

## 💡 简介

[Lute](https://github.com/b3log/lute) 是一款结构化的 Markdown 引擎，完整实现了最新的 [GFM](https://github.github.com/gfm/) / [CommonMark](https://commonmark.org) 规范，对中文语境支持更好。

## 📽️ 背景

之前我一直在使用其他 Markdown 引擎，他们或多或少都有些“瑕疵”：

* 对标准规范的支持不一致
* 对“怪异”文本处理非常耗时，甚至挂死

Lute 的目标是构建一个结构化的 Markdown 引擎，实现 GFM/CommonMark 规范。所谓的“结构化”指的是从输入的 MD 文本构建抽象语法树，通过操作树来进行 HTML 输出、原文格式化等。
支持 GFM/CM 规范则是为了保证 Markdown 渲染不存在二义性，让同一份 Markdown 文本可以在实现这两个规范的 Markdown 引擎处理后得到一样的结果，我觉得这一点非常重要。

实现规范的引擎并不多，我想试试看自己能不能写上一个，这也是 Lute 的动机之一。关于如何实现一个 Markdown 引擎，网上众说纷纭：

* 有的人说 Markdown 适合用正则解析，因为文法规则太简单
* 也有的人说 Markdown 可以用编译原理来处理，正则太难维护

我赞同后者，因为正则确实太难维护而且运行效率较低。最重要的原因是符合 GFM/CM 规范的 Markdown 引擎的核心解析算法不大可能用正则写出来，因为规范定义的规则实在是太复杂了。

最后，还有一个很重要的动机就是 B3log 开源社区需要一款自己的 Markdown 引擎：

* 社区项目 [Solo](https://github.com/b3log/solo)、[Pipe](https://github.com/b3log/pipe)、[Sym](https://github.com/b3log/symphony) 需要效果统一的 Markdown 渲染，并且性能非常重要
* 社区项目 [Vditor](https://github.com/b3log/vditor) 需要一款结构化的引擎作为支撑，实现下一代的 Markdown 编辑器，为未来而构建

## ✨  特性

* 完整实现最新版 GFM/CM 规范
* 非常快
* 代码块语法高亮（待实现）
* 更好地支持中文语境（待实现）
* 支持原文格式化（待实现）
* 可扩展语法树节点（待实现）

## ⚡ 性能

在相同机器上，用相同的测试数据（[CommonMark 规范文档](https://github.com/commonmark/commonmark-spec-web/blob/gh-pages/0.29/spec.txt) ~197K）跑基准测试结果如下。
目前看来在实现 CommonMark 规范的前提下，Lute、goldmark 和 golang-commonmark 的性能差距不大。

### [Lute](https://github.com/b3log/lute)

```
BenchmarkLute-2   	     300	   4458074 ns/op	 2767406 B/op	   22450 allocs/op
BenchmarkLute-4   	     300	   4428252 ns/op	 2767664 B/op	   22450 allocs/op
BenchmarkLute-8   	     300	   4591147 ns/op	 2768186 B/op	   22450 allocs/op
```

### [goldmark](https://github.com/yuin/goldmark)

```
BenchmarkGoldMark-2   	     300	   4757269 ns/op	 2110741 B/op	   13901 allocs/op
BenchmarkGoldMark-4   	     300	   4790514 ns/op	 2113679 B/op	   13902 allocs/op
BenchmarkGoldMark-8   	     300	   4860327 ns/op	 2114696 B/op	   13902 allocs/op
```

### [golang-commonmark](https://gitlab.com/golang-commonmark/markdown)

```
BenchmarkGolangCommonMark-2   	     300	   5099691 ns/op	 2973258 B/op	   18827 allocs/op
BenchmarkGolangCommonMark-4   	     300	   5083059 ns/op	 2973794 B/op	   18828 allocs/op
BenchmarkGolangCommonMark-8   	     300	   5103111 ns/op	 2974818 B/op	   18828 allocs/op
```

### [Blackfriday](https://github.com/russross/blackfriday)

```
BenchmarkBlackFriday-2   	     500	   3875623 ns/op	 3318457 B/op	   20052 allocs/op
BenchmarkBlackFriday-4   	     500	   3783871 ns/op	 3334775 B/op	   20056 allocs/op
BenchmarkBlackFriday-8   	     500	   3917515 ns/op	 3341045 B/op	   20058 allocs/op
```

Blackfriday 没有实现 CommonMark 所以性能好一些。

### [markdown-it](https://github.com/markdown-it/markdown-it)

markdown-it 是 JavaScript 写的，它同样实现了 CommonMark 规范。循环渲染 300 次，平均每次调用耗时 9285933ns（9.2ms），耗时大致是 golang 实现的两倍。

## 📜 文档

TBD

* [《提问的智慧》精读注解版](https://hacpai.com/article/1536377163156)
* Lute 使用指南
* CommonMark 规范要点解读
* Lute 实现后记

目前所有文档（包括代码注释）都是中文的，欢迎外文好的同学帮忙进行国际化，不胜感激！

## 🏘️ 社区

* [讨论区](https://hacpai.com/tag/lute)
* [报告问题](https://github.com/b3log/lute/issues/new/choose)

## 📄 授权

Lute 使用 [Apache License, Version 2](https://www.apache.org/licenses/LICENSE-2.0) 开源协议。

## 🙏 鸣谢

Lute 的诞生离不开以下开源项目，在此对这些项目的贡献者们致以最崇高的敬意！

* [commonmark.js](https://github.com/commonmark/commonmark.js)：该项目是 CommonMark 官方参考实现的 JavaScript 版，Lute 参考了其解析器实现部分
* [mdast](https://github.com/syntax-tree/mdast)：该项目介绍了一种 Markdown 抽象语法树结构的表现形式，Lute 的 AST 在初始设计阶段参考了该项目
* [goldmark](https://github.com/yuin/goldmark)：另一款用 golang 写的 Markdown 引擎，Lute 参考了其树遍历实现部分
* [golang-commonmark](https://gitlab.com/golang-commonmark/markdown)：另一款用 golang 写的 Markdown 引擎，Lute 参考了其 URL 编码算法

---

## 👍 开源项目推荐

* 如果你需要集成一个浏览器端的 Markdown 编辑器，可以考虑使用 [Vditor](https://github.com/b3log/vditor)
* 如果你需要搭建一个个人博客系统，可以考虑使用 [Solo](https://github.com/b3log/solo)
* 如果你需要搭建一个社区平台，可以考虑使用 [Sym](https://github.com/b3log/symphony)
* 欢迎加入我们的小众开源社区，详情请看[这里](https://hacpai.com/article/1463025124998)
