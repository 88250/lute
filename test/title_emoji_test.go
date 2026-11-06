package test

import (
	"strings"
	"testing"
)

func TestTabsTitleCustomEmojiRoundTrip(t *testing.T) {
	l := tabsEngine()
	l.PutEmojis(map[string]string{"custom": "/emojis/custom.png", "vector": "/emojis/vector.svg"})
	for _, title := range []string{":custom:", "12:custom:", ":custom::vector:", "A:vector:B"} {
		t.Run(title, func(t *testing.T) {
			dom := l.Md2BlockDOM("::: tabs\n@tab "+title+"\n\nBody\n\n:::\n", false)
			count := strings.Count(dom, `class="emoji"`)
			if count == 0 {
				t.Fatal("expected custom emoji in the rendered title")
			}
			for round := 0; round < 3; round++ {
				_, items := tabNodes(l.BlockDOM2Tree(dom).Root)
				if len(items) != 1 || items[0].TabItemTitle != title {
					t.Fatalf("round %d changed the title: %#v", round, items)
				}
				dom = l.SpinBlockDOM(dom)
				if strings.Count(dom, `class="emoji"`) != count || strings.Contains(dom, `data-type="img"`) {
					t.Fatalf("round %d changed custom emoji into ordinary images: %s", round, dom)
				}
			}
		})
	}
}

func TestCalloutTitleCustomEmojiRoundTrip(t *testing.T) {
	l := newCalloutLute()
	l.PutEmojis(map[string]string{"custom": "/emojis/custom.png", "vector": "/emojis/vector.svg"})
	for _, title := range []string{":custom:", "12:custom:", ":custom::vector:", "A:vector:B"} {
		t.Run(title, func(t *testing.T) {
			dom := l.Md2BlockDOM("> [!NOTE] ![callout-icon](https://example.com/icon.png) "+title+"\n> Body", false)
			count := strings.Count(dom, `class="emoji"`)
			if count == 0 {
				t.Fatal("expected custom emoji in the rendered title")
			}
			for round := 0; round < 3; round++ {
				tree := l.BlockDOM2Tree(dom)
				if tree.Root.FirstChild.CalloutTitle != title {
					t.Fatalf("round %d changed the title: %q", round, tree.Root.FirstChild.CalloutTitle)
				}
				dom = l.SpinBlockDOM(dom)
				if strings.Count(dom, `class="emoji"`) != count || strings.Contains(dom, `data-type="img"`) {
					t.Fatalf("round %d changed custom emoji into ordinary images: %s", round, dom)
				}
			}
		})
	}
}
