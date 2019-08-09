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

// Table 描述了表节点结构。
type Table struct {
	*BaseNode
	Align int // 0：默认对齐，1：左对齐，2：居中对齐，3：右对齐
}

// TableRow 描述了表行节点结构。
type TableRow struct {
	*BaseNode
}

// TableCell 描述了表格节点结构。
type TableCell struct {
}
