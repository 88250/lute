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

func New(options map[string]map[string]*js.Object) *js.Object {
	engine := lute.New()
	registerRenderers(engine, options)
	return js.MakeWrapper(engine)
}

func registerRenderers(engine *lute.Lute, options map[string]map[string]*js.Object) {
	for rendererType, extRenderer := range options["renderers"] {
		switch extRenderer.Interface().(type) {
		case map[string]interface{}:
			break
		default:
			continue
		}

		var rendererFuncs map[lute.NodeType]lute.ExtRendererFunc
		if "HTML2Md" == rendererType {
			rendererFuncs = engine.HTML2MdRendererFuncs
		} else if "HTML2VditorDOM" == rendererType {
			rendererFuncs = engine.HTML2VditorDOMRendererFuncs
		} else {
			continue
		}

		renderFuncs := extRenderer.Interface().(map[string]interface{})
		for funcName, _ := range renderFuncs {
			nodeType := "Node" + funcName[len("render"):]
			rendererFuncs[lute.Str2NodeType(nodeType)] = func(node *lute.Node, entering bool) (string, lute.WalkStatus) {
				nodeType := node.Typ.String()
				funcName = "render" + nodeType[len("Node"):]
				ret := extRenderer.Call(funcName, js.MakeWrapper(node), entering).Interface().([]interface{})
				return ret[0].(string), lute.WalkStatus(ret[1].(float64))
			}
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
