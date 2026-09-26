* [Golang markdown 引擎性能基准测试](https://ld246.com/article/1574570835061)

纯 Markdown 解析基准（不包含 HTML 渲染）：

```sh
go test ./benchmark -run '^$' -bench '^BenchmarkParse' -benchmem -cpu=1
```

在仓库根目录执行，覆盖不同规模的非表格段落、自动链接、CRLF、NUL、引用链接、任务列表和 Emoji，并包含 CommonMark 规范文档及 Markdown、Protyle、Vditor 模式对照。每次解析前恢复输入缓冲区，避免输入被修改后影响后续测量。
