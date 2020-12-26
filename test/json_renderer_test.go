// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"fmt"
	"testing"

	"github.com/88250/lute"
)

var JSONRendererTest = "```mindmap\n- 教程\n- 语法指导\n  - 普通内容\n  - 提及用户\n  - 表情符号 Emoji\n    - 一些表情例子\n  - 大标题 - Heading 3\n    - Heading 4\n      - Heading 5\n        - Heading 6\n  - 图片\n  - 代码块\n    - 普通\n    - 语法高亮支持\n      - 演示 Go 代码高亮\n      - 演示 Java 高亮\n  - 有序、无序、任务列表\n    - 无序列表\n    - 有序列表\n    - 任务列表\n  - 表格\n  - 隐藏细节\n  - 段落\n  - 链接引用\n  - 数学公式\n  - 脑图\n  - 流程图\n  - 时序图\n  - 甘特图\n  - 图表\n  - 五线谱\n  - Graphviz\n  - 多媒体\n  - 脚注\n- 快捷键\n```\n"

func TestJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.KramdownIAL = false

	jsonStr := luteEngine.RenderJSON(JSONRendererTest)
	fmt.Println(jsonStr)
}
