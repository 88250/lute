package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/88250/lute/parse"
)

var parsedTree *parse.Tree

// BenchmarkParse 单独测量解析，并在每次迭代恢复可被解析器修改的输入。
func BenchmarkParse(b *testing.B) {
	for _, size := range []int{256, 1024, 4096} {
		for _, name := range []string{"FalseTable", "NoPipe", "AutoLink", "CRLF", "NUL", "References", "TaskList", "Emoji"} {
			b.Run(fmt.Sprintf("%s/%d", name, size), func(b *testing.B) {
				o := parse.NewOptions()
				o.Emoji, o.GFMAutoLink = false, false
				var input string
				switch name {
				case "FalseTable":
					input = strings.Repeat("alpha | beta gamma delta\n", size)
				case "NoPipe":
					input = strings.Repeat("alpha beta gamma delta\n", size)
				case "AutoLink":
					input = strings.Repeat("a", size*16) + "\n"
					o.GFMAutoLink = true
				case "CRLF":
					input = "```text\r\n" + strings.Repeat("alpha beta gamma delta\r\n", size*16) + "```\r\n"
				case "NUL":
					input = "```text\n" + strings.Repeat("a\x00", size*16) + "\n```\n"
				case "References":
					var sb strings.Builder
					for i := 0; i < size; i++ {
						fmt.Fprintf(&sb, "[label%d]: /url\n", i)
					}
					sb.WriteString("\n")
					for i := 0; i < size; i++ {
						fmt.Fprintf(&sb, "[label%d] ", i)
					}
					input = sb.String() + "\n"
				case "TaskList":
					input = strings.Repeat("- [ ] alpha **beta** and `code`\n", size)
				case "Emoji":
					input = strings.Repeat("plain text without emoji aliases\n\n", size)
					o.Emoji = true
				}
				benchmarkParse(b, []byte(input), o)
			})
		}
	}
}

func BenchmarkParseSpec(b *testing.B) {
	input, err := os.ReadFile(spec + ".md")
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParse(b, input, parse.NewOptions())
}

func BenchmarkParseTaskModes(b *testing.B) {
	input := []byte(strings.Repeat("- [ ] alpha **beta** and `code`\n", 500))
	for _, mode := range []string{"Markdown", "Protyle", "Vditor"} {
		b.Run(mode, func(b *testing.B) {
			o := parse.NewOptions()
			o.ProtyleWYSIWYG = mode == "Protyle"
			o.VditorWYSIWYG = mode == "Vditor"
			o.KramdownBlockIAL, o.KramdownSpanIAL = o.ProtyleWYSIWYG, o.ProtyleWYSIWYG
			benchmarkParse(b, input, o)
		})
	}
}

func benchmarkParse(b *testing.B, input []byte, options *parse.Options) {
	buffer := make([]byte, len(input), len(input)+8)
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buffer, input)
		parsedTree = parse.Parse("", buffer, options)
	}
}
