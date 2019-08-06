# Lute

> 千呼万唤始出来，犹抱琵琶半遮面。转轴拨弦三两声，未成曲调先有情。

## 💡 简介

[Lute](https://github.com/b3log/lute) 是一款结构化的 Markdown 引擎，完整实现了最新的 [GFM](https://github.github.com/gfm/) / [CommonMark](https://commonmark.org)规范，对中文语境支持更好。

## 📽️ 背景

<details>
<summary>太长不看。</summary>
<br>

之前我一直在使用其他 Markdown 引擎，他们或多或少都有些“瑕疵”：

* 对标准规范的支持不一致
* 对“怪异”文本处理非常耗时，甚至挂死

Lute 的目标是构建一个结构化的 Markdown 引擎，实现 GFM/CommonMark 规范。所谓的“结构化”指的是从输入的 MD 文本构建抽象语法树，通过操作树来进行 HTML 输出、原文格式化等。
支持 GFM/CM 规范则是为了保证没有二义性，让同一份 Markdown 文本可以在实现这两个规范的 Markdown 引擎处理后得到一样的结果，这一点非常重要。

实现规范的 Markdown 引擎并不多，我想试试看自己能不能写上一个，这也是 Lute 的动机之一。关于如何实现一个 Markdown 引擎，网上众说纷纭：

* 有的人说 Markdown 适合用正则解析，因为文法规则太简单
* 也有的人说 Markdown 可以用编译原理来处理，正则太难维护

我赞同后者，因为正则确实太难维护而且运行效率较低。最重要的原因是符合 GFM/CM 规范的 Markdown 引擎的核心解析算法是不可能用正则写出来的，因为规范定义的规则实在是太复杂了。

暂时抛开实现方式，回到“结构化”这一点上。结构化的目的并不只是为了优美，它的意义在于为 [Vditor](https://github.com/b3log/vditor) 提供良好的数据结构，让 Vditor 实现所见即所得的特性提供有力支撑。

最终，我们会将 Vditor 打造为下一代的 Markdown 编辑器，为未来而构建。

</details>

## ✨  特性

* [x] 完整实现最新版 CommonMark 规范
* [ ] 更好地支持中文语境
* [ ] 格式化
* [ ] 可扩展语法树节点类型以实现自定义输出

## ⚡ 性能

以下是对 [CommonMark 规范文档](https://github.com/commonmark/commonmark-spec-web/blob/gh-pages/0.29/spec.txt)（~198K，9700 行）跑基准测试的结果：

```
BenchmarkLute-2   	     200	   6392891 ns/op	 4281718 B/op	   51869 allocs/op
BenchmarkLute-4   	     200	   6238305 ns/op	 4282277 B/op	   51870 allocs/op
BenchmarkLute-8   	     200	   6492781 ns/op	 4283295 B/op	   51871 allocs/op
```

Lute 在性能方面还有很大优化空间，目标是做到至少和 [goldmark](https://github.com/yuin/goldmark) 一样快（不得不说，goldmark 真的很快）。

```
BenchmarkGoldMark-2   	     300	   4724041 ns/op	 2110378 B/op	   13901 allocs/op
BenchmarkGoldMark-4   	     300	   4817211 ns/op	 2113808 B/op	   13902 allocs/op
BenchmarkGoldMark-8   	     300	   4860328 ns/op	 2114412 B/op	   13902 allocs/op
```

## 📜 文档

TBD

* [《提问的智慧》精读注解版](https://hacpai.com/article/1536377163156)
* Lute 使用指南
* CommonMark 规范要点解读
* Lute 实现后记

## 🏘️ 社区

* [讨论区](https://hacpai.com/tag/lute)
* [报告问题](https://github.com/b3log/lute/issues/new/choose)

## 📄 授权

Lute 使用 [Apache License, Version 2](https://www.apache.org/licenses/LICENSE-2.0) 开源协议。

## 🙏 鸣谢

Lute 的诞生离不开以下开源项目，在此对这些项目的贡献者们致以最崇高的敬意！

* [commonmark.js](https://github.com/commonmark/commonmark.js)：该项目是 CommonMark 官方参考实现的 JavaScript 版，Lute 参考了其解析器实现部分
* [mdast](https://github.com/syntax-tree/mdast)：该项目介绍了一种 Markdown 抽象语法树结构的表现形式，Lute 的 AST 在初始设计阶段参考了该项目
* [goldmark](https://github.com/yuin/goldmark)：另一款用 golang 写的 Markdown 解析器，Lute 参考了其树遍历实现部分

---

## 👍 开源项目推荐

* 如果你需要集成一个浏览器端的 Markdown 编辑器，可以考虑使用 [Vditor](https://github.com/b3log/vditor)
* 如果你需要搭建一个个人博客系统，可以考虑使用 [Solo](https://github.com/b3log/solo)
* 如果你需要搭建一个社区平台，可以考虑使用 [Sym](https://github.com/b3log/symphony)
* 欢迎加入我们的小众开源社区，详情请看[这里](https://hacpai.com/article/1463025124998)
