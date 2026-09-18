package test

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/parse"
)

// 数据来自 https://www.unicode.org/Public/17.0.0/emoji/emoji-test.txt 中 E17.0 的 fully-qualified 序列。
//
//go:embed emoji_unicode17.json
var unicode17EmojiJSON []byte

func TestUnicode17Emoji(t *testing.T) {
	var entries []struct {
		Alias string `json:"alias"`
		Emoji string `json:"emoji"`
	}
	if err := json.Unmarshal(unicode17EmojiJSON, &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 163 {
		t.Fatalf("expected 163 sequences, got %d", len(entries))
	}
	engine := lute.New()
	engine.SetCallout(true)
	for _, entry := range entries {
		t.Run(entry.Alias, func(t *testing.T) {
			if got := engine.GetEmojis()[entry.Alias]; got != entry.Emoji {
				t.Fatalf("mapping: expected %q, got %q", entry.Emoji, got)
			}
			if got := parse.EmojiUnicodeAlias[entry.Emoji]; got != entry.Alias {
				t.Fatalf("reverse mapping: expected %q, got %q", entry.Alias, got)
			}
			for _, input := range []string{entry.Emoji, ":" + entry.Alias + ":"} {
				if got := engine.MarkdownStr("", input); got != "<p>"+entry.Emoji+"</p>\n" {
					t.Fatalf("inline %q: %s", input, got)
				}
				got := engine.MarkdownStr("", "> [!NOTE] "+input+" Title\n> Content\n")
				want := `<span class="callout-icon">` + entry.Emoji + `</span><span class="callout-title">Title</span>`
				if !strings.Contains(got, want) {
					t.Fatalf("callout %q: %s", input, got)
				}
			}
		})
	}
}
