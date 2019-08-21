// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

// Lute 描述了 Lute 引擎的顶层使用入口。
type Lute struct {
	options
}

// New 创建一个新的 Lute 引擎，默认启用所有 GFM 支持、代码块语法高亮。
func New(opts ...option) (ret *Lute) {
	ret = &Lute{}
	GFM(true)(ret)
	CodeSyntaxHighlight(true)(ret)

	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

// Markdown 将 markdown 文本字符数组处理为相应的 html 字符数组。name 参数用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (lute *Lute) Markdown(name string, markdown []byte) (html []byte, err error) {
	var tree *Tree
	tree, err = parse(name, markdown, lute.options)
	if nil != err {
		return
	}

	renderer := newHTMLRenderer(lute.options)
	html, err = tree.render(renderer)
	if nil != err {
		return
	}
	return
}

// MarkdownStr 接受 string 类型的 markdown 后直接调用 Markdown 进行处理。
func (lute *Lute) MarkdownStr(name, markdown string) (html string, err error) {
	var htmlBytes []byte
	htmlBytes, err = lute.Markdown(name, toItems(markdown))
	if nil != err {
		return
	}

	html = fromItems(htmlBytes)
	return
}

// GFM 设置是否打开所有 GFM 支持。
func GFM(b bool) option {
	return func(lute *Lute) {
		lute.GFMTable = b
		lute.GFMTaskListItem = b
		lute.GFMStrikethrough = b
		lute.GFMAutoLink = b
	}
}

// GFMTable 设置是否打开“GFM 表”支持。
func GFMTable(b bool) option {
	return func(lute *Lute) {
		lute.GFMTable = b
	}
}

// GFMTaskListItem 设置是否打开“GFM 任务列表项”支持。
func GFMTaskListItem(b bool) option {
	return func(lute *Lute) {
		lute.GFMTaskListItem = b
	}
}

// GFMStrikethrough 设置是否打开“GFM 删除线”支持。
func GFMStrikethrough(b bool) option {
	return func(lute *Lute) {
		lute.GFMStrikethrough = b
	}
}

// GFMAutoLink 设置是否打开“GFM 自动链接”支持。
func GFMAutoLink(b bool) option {
	return func(lute *Lute) {
		lute.GFMAutoLink = b
	}
}

// CodeSyntaxHighlight 设置是否对代码块进行语法高亮。
func CodeSyntaxHighlight(b bool) option {
	return func(lute *Lute) {
		lute.CodeSyntaxHighlight = b
	}
}

// options 描述了一些列解析和渲染选项。
type options struct {
	GFMTable         bool // 是否处理 GFM 表
	GFMTaskListItem  bool // 是否处理 GFM 任务列表项
	GFMStrikethrough bool // 是否处理 GFM 删除线
	GFMAutoLink      bool // 是否处理 GFM 自动链接

	CodeSyntaxHighlight bool // 是否对代码块进行语法高亮
}

// option 描述了解析渲染选项设置函数签名。
type option func(lute *Lute)
