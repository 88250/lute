package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestHTMLCodeBlockBanner(t *testing.T) {
	tests := []struct {
		name, html, markdown string
	}{
		{"plain", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre>echo hello</pre></div>`, "```bash\necho hello\n```\n"},
		{"highlight and controls", `<div class="extra md-code-block"><div class="md-code-block-banner extra"><svg><title>Icon</title></svg><span> bash </span><button>复制代码</button><div role="button"><span>Copy</span></div></div><pre><span class="token">echo</span> &lt;a&gt; &amp; b
  next

last</pre></div>`, "```bash\necho <a> & b\n  next\n\nlast\n```\n"},
		{"existing code language", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre><code class="language-python">pass</code></pre></div>`, "```python\npass\n```\n"},
		{"existing pre language", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre class="language-python">pass</pre></div>`, "```python\npass\n```\n"},
		{"existing bare language", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre class="python">pass</pre></div>`, "```python\npass\n```\n"},
		{"custom button", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span><div class="ds-icon-button"><span>复制</span></div></div><pre>echo x</pre></div>`, "```bash\necho x\n```\n"},
		{"existing data language", `<div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre data-language="python">pass</pre></div>`, "```python\npass\n```\n"},
		{"code without language", `<div class="md-code-block"><div class="md-code-block-banner"><span>C++</span></div><pre><code>int x;</code></pre></div>`, "```C++\nint x;\n```\n"},
		{"multiple blocks", `<p>before</p><div class="md-code-block"><div class="md-code-block-banner"><span>bash</span></div><pre>echo x</pre></div><div class="md-code-block"><div class="md-code-block-banner"><span>go</span></div><pre>package main</pre></div><p>after</p>`, "before\n\n```bash\necho x\n```\n\n\n```go\npackage main\n```\n\n\nafter\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := lute.New()
			if got := engine.HTML2Md(tt.html); got != tt.markdown {
				t.Fatalf("expected %q, got %q", tt.markdown, got)
			}
			// 同时验证编辑器使用的转换入口，确保语言和正文在往返转换后保留。
			want := strings.ReplaceAll(tt.markdown, "```\n\n\n", "```\n\n")
			if got := engine.BlockDOM2StdMd(engine.HTML2BlockDOM(tt.html)); got != want {
				t.Fatalf("BlockDOM round trip: expected %q, got %q", want, got)
			}
		})
	}
}

func TestHTMLCodeBlockBannerPreservesUnrecognizedContent(t *testing.T) {
	for _, body := range []string{
		`<div class="md-code-block-banner"><span>bash script</span></div><pre>echo x</pre>`,
		`<div class="md-code-block-banner"><span>bash</span><span>python</span></div><pre>echo x</pre>`,
		`<div class="md-code-block-banner"><span>Copy</span></div><pre>echo x</pre>`,
		`<div class="md-code-block-banner"><span>unknown_language</span></div><pre>echo x</pre>`,
		`<div class="md-code-block-banner"><span>bash</span></div><p>echo x</p>`,
		`<div class="md-code-block-banner"><span>bash</span></div><pre>x</pre><pre>y</pre>`,
		`<div class="md-code-block-banner"><span>bash</span></div><div class="md-code-block-banner">python</div><pre>x</pre>`,
		`<p>bash</p><pre>echo x</pre>`,
		`<pre>echo x</pre>`,
	} {
		engine := lute.New()
		want := engine.HTML2Md("<div>" + body + "</div>")
		if got := engine.HTML2Md(`<div class="md-code-block">` + body + `</div>`); got != want {
			t.Errorf("content changed for %s: expected %q, got %q", body, want, got)
		}
	}
	engine := lute.New()
	html := `<div class="md-code-block-other"><div class="md-code-block-banner"><span>bash</span></div><pre>echo x</pre></div>`
	if got := engine.HTML2Md(html); !strings.HasPrefix(got, "bash\n\n```\n") {
		t.Fatalf("unrelated container changed: %q", got)
	}
}
