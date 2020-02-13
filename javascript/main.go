// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package main

import (
	"github.com/88250/lute"
	"github.com/gopherjs/gopherjs/js"
)

func New(formatRenderer, vditorRenderer *js.Object) *js.Object {
	engine := lute.New()
	registerRenderer(engine, formatRenderer)
	registerRenderer(engine, vditorRenderer)
	return js.MakeWrapper(engine)
}

func registerRenderer(engine *lute.Lute, renderer *js.Object) {
	switch renderer.Interface().(type) {
	case map[string]interface{}:
		break
	default:
		return
	}

	renderFuncs := renderer.Interface().(map[string]interface{})
	for funcName, _ := range renderFuncs {
		nodeType := "Node" + funcName[len("render"):]
		engine.FormatRendererFuncs[lute.Str2NodeType(nodeType)] = func(n *lute.Node, entering bool) (status lute.WalkStatus, err error) {
			walkStatus := renderer.Call(funcName, js.MakeWrapper(n), entering).Int()
			return lute.WalkStatus(walkStatus), nil
		}
	}
}

func main() {
	js.Global.Set("Lute", map[string]interface{}{
		"Version":          lute.Version,
		"New":              New,
		"WalkStop":         lute.WalkStop,
		"WalkSkipChildren": lute.WalkSkipChildren,
		"WalkContinue":     lute.WalkContinue,
	})
}
