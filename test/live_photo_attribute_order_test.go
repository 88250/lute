package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestLivePhotoCustomAttributesStableRoundTrip(t *testing.T) {
	l := lute.New()
	l.SetProtyleWYSIWYG(true)
	l.SetKramdownIAL(true)
	const source = `![first](assets/a.heic){: custom-sy-live-photo="1" custom-sy-live-photo-video="assets/new.mov?x=1&y=2"} ![second](assets/b.heic){: custom-sy-live-photo="1" custom-sy-live-photo-video="assets/other.mov"} assets/old.mov`
	dom := l.Md2BlockDOM(source, false)
	want := l.BlockDOM2Md(dom)
	if !strings.Contains(want, `custom-sy-live-photo-video="assets/new.mov?x=1&amp;y=2"`) ||
		!strings.Contains(want, `custom-sy-live-photo-video="assets/other.mov"`) ||
		!strings.Contains(want, `} assets/old.mov`) {
		t.Fatalf("lost live-photo metadata or literal text: %s", want)
	}
	for i := 0; i < 200; i++ {
		if got := l.BlockDOM2Md(dom); got != want {
			t.Fatalf("same DOM serialized differently on iteration %d:\n%s\nwant:\n%s", i, got, want)
		}
		if got := l.BlockDOM2Md(l.Md2BlockDOM(want, false)); got != want {
			t.Fatalf("Markdown round trip changed attributes on iteration %d:\n%s\nwant:\n%s", i, got, want)
		}
	}
}
