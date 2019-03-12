# Lute

> 千呼万唤始出来，犹抱琵琶半遮面。转轴拨弦三两声，未成曲调先有情。

## 简介

[Lute](https://github.com/b3log/lute) 是一款结构化的 Markdown 处理引擎，具备语法报错、格式化等功能。

## 背景

<details>
<summary>太长不看。</summary>
<br>

关于如果实现一个 Markdown 处理器，网上众说纷纭。有的人说 Markdown 适合用正则解析，因为文法规则太简单；也有的人说 Markdown 可以用编译原理来处理，正则太难维护。

我赞同后者。之前我一直在使用其他 Markdown 处理器，他们或多或少都有些“瑕疵”：

* 对标准规范的支持不一致
* 对“怪异”的文本处理非常耗时，甚至挂死
* **对中文支持不好**

Lute 的目标是构建一个结构化的 Markdown 引擎。所谓的“结构化”指的是从输入的 MD 文本构建抽象语法树，通过操作树来进行格式化、HTML 输出等。

除了结构化，追求高性能也是很重要的目标。
</details>

## 特性

* 完整实现 [CommonMark 规范](https://spec.commonmark.org)
* 完整实现对中文的支持

## 词法分析


## 语法分析

语法树结构按照 [mdast 规范](https://github.com/syntax-tree/mdast)进行实现。

## 目标代码生成

## 鸣谢

* golang template：很多代码偷师自该包
* [mdast](https://github.com/syntax-tree/mdast)：Markdown 语法树规范
