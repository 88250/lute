# Lute

> 千呼万唤始出来，犹抱琵琶半遮面。转轴拨弦三两声，未成曲调先有情。

## 简介

[Lute](https://github.com/b3log/lute) 是一款结构化的 Markdown 引擎，完整实现了最新的 [CommonMark 规范](https://commonmark.org)，对中文语境支持更好，并具备语法检查、格式化等功能。

## 背景

<details>
<summary>太长不看。</summary>
<br>

之前我一直在使用其他 Markdown 处理器，他们或多或少都有些“瑕疵”：

* 对标准规范的支持不一致
* 对“怪异”的文本处理非常耗时，甚至挂死
* **对中文支持不够好**

Lute 的目标是构建一个结构化的 Markdown 引擎。所谓的“结构化”指的是从输入的 MD 文本构建抽象语法树，通过操作树来进行格式化、HTML 输出等。

关于如何实现一个 Markdown 处理器，网上众说纷纭。有的人说 Markdown 适合用正则解析，因为文法规则太简单；也有的人说 Markdown 可以用编译原理来处理，正则太难维护。我赞同后者，因为正则确实太难维护而且运行效率较低。
</details>

## 特性

* 完整实现最新版 CommonMark 规范
* 更好地支持中文语境
* 自动格式化、Lint 
* 可扩展语法树节点类型以实现自定义输出
* 内置缓存以提升性能

## 性能对比

TBD

## 文档

TBD

* [《提问的智慧》精读注解版](https://hacpai.com/article/1536377163156)
* Lute 使用指南
* CommonMark 规范要点解读
* Lute 实现后记

## 社区

* [讨论区](https://hacpai.com/tag/lute)
* [报告问题](https://github.com/b3log/lute/issues/new/choose)

## 授权

Lute 使用 [Apache License, Version 2](https://www.apache.org/licenses/LICENSE-2.0) 开源协议。

## 鸣谢

Lute 的诞生离不开以下开源项目，在此对这些项目的贡献者致以最崇高的敬意！

* [commonmark.js](https://github.com/commonmark/commonmark.js)：该项目是 CommonMark 官方参考实现的 JavaScript 版，Lute 参考了其解析器实现部分
* [mdast](https://github.com/syntax-tree/mdast)：该项目介绍了一种 Markdown 抽象语法树结构的表现形式，Lute 的 AST 在初始设计阶段参考了该项目
* [goldmark](https://github.com/yuin/goldmark)：另一款用 golang 写的 Markdown 解析器，Lute 参考了其树遍历实现部分

---

## 开源项目推荐

* 如果你需要集成一个浏览器端的 Markdown 编辑器，可以考虑使用 [Vditor](https://github.com/b3log/vditor)
* 如果你需要搭建一个个人博客系统，可以考虑使用 [Solo](https://github.com/b3log/solo)
* 如果你需要搭建一个社区平台，可以考虑使用 [Sym](https://github.com/b3log/symphony)
* 欢迎加入我们的小众开源社区，详情请看[这里](https://hacpai.com/article/1463025124998)
