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

package main

import (
	"os"

	"github.com/b3log/gulu"
	"github.com/b3log/lute"
	"github.com/valyala/fasthttp"
)

// handleMarkdown2HTML 处理 Markdown 转 HTML。
// POST 请求 Body 传入 Markdown 原文；响应 Body 即处理好的 HTML。
func handleMarkdown2HTML(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()
	html, err := engine.Markdown("", body)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	ctx.SetBody(html)
}

// handle 处理请求分发。
func handle(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/", "":
		handleMarkdown2HTML(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

// 日志记录器。
var logger *gulu.Logger

// Lute 的 HTTP Server 入口点。
func main() {
	gulu.Log.SetLevel("debug")
	logger = gulu.Log.NewLogger(os.Stdout)

	err := fasthttp.ListenAndServe("127.0.0.1:8250", handle)
	if nil != err {
		logger.Fatalf("Lute HTTP server 启动失败，原因：%s", err.Error())
	}
}
