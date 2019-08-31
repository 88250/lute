// Lute - A structured markdown engine.
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
	"os"

	"github.com/b3log/gulu"
	"github.com/b3log/lute"
	"github.com/valyala/fasthttp"
)

var logger = gulu.Log.NewLogger(os.Stdout)

// handleMarkdown2HTML 处理 Markdown 转 HTML。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是处理好的 HTML。
func handleMarkdown2HTML(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()
	html, err := engine.Markdown("", body)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		logger.Errorf("markdown text [%s]\n", body)
		return
	}
	ctx.SetBody(html)
}

// handleMarkdownFormat 处理 Markdown 格式化。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是格式化好的 Markdown 文本。
func handleMarkdownFormat(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()
	formatted, err := engine.Format("", body)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		logger.Errorf("markdown text [%s]\n", body)
		return
	}
	ctx.SetBody(formatted)
}

// handle 处理请求分发。
func handle(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/", "":
		handleMarkdown2HTML(ctx)
	case "/format":
		handleMarkdownFormat(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

// Lute 的 HTTP Server 入口点。
func main() {
	gulu.Log.SetLevel("debug")

	addr := "127.0.0.1:8249"
	logger.Infof("booting Lute HTTP Server on [%s]", addr)
	server := &fasthttp.Server{
		Handler:            handle,
		MaxRequestBodySize: 1024 * 1024 * 2, // 2MB
	}
	err := server.ListenAndServe(addr)
	if nil != err {
		logger.Fatalf("booting Lute HTTP server failed: %s", err.Error())
	}
}
