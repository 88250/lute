package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestProtyleVideoControls(t *testing.T) {
	l := lute.New()
	l.SetProtyleWYSIWYG(true)
	l.SetSanitize(true)
	tests := []struct {
		name, input, expected string
	}{
		{"missing", `<video src="assets/video.mp4"></video>`, `<video src="assets/video.mp4" controls="controls"></video>`},
		{"existing", `<video controls="controls" src="assets/video.mp4"></video>`, `<video controls="controls" src="assets/video.mp4"></video>`},
		{"boolean", `<video controls src="assets/video.mp4"></video>`, `<video controls="" src="assets/video.mp4"></video>`},
		{"attributes", `<video src="https://example.com/video.mp4?a=1&amp;b=2" poster="assets/poster.png" muted loop></video>`, `<video src="https://example.com/video.mp4?a=1&amp;b=2" poster="assets/poster.png" muted="" loop="" controls="controls"></video>`},
		{"source", `<video><source src="assets/video.mp4" type="video/mp4"></video>`, `<video controls="controls"><source src="assets/video.mp4" type="video/mp4"/></video>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, dom := range map[string]string{
				"markdown":       l.Md2BlockDOM(tt.input, false),
				"html":           l.HTML2BlockDOM(tt.input),
				"existing-block": l.SpinBlockDOM(`<div data-node-id="20260918120000-abcdefg" data-type="NodeVideo" class="iframe"><div class="iframe-content">` + tt.input + `</div></div>`),
			} {
				expected := tt.expected
				if "source" == tt.name && "html" == name {
					expected = `<video src="assets/video.mp4" controls="controls"></video>`
				}
				if !strings.Contains(dom, expected) {
					t.Errorf("%s: expected %s in %s", name, expected, dom)
				}
				spun := l.SpinBlockDOM(dom)
				if !strings.Contains(spun, expected) || strings.Count(spun, " controls=") != 1 {
					t.Errorf("%s: controls or attributes changed after round trip: %s", name, spun)
				}
			}
		})
	}
}

func TestProtyleAudioControls(t *testing.T) {
	l := lute.New()
	l.SetProtyleWYSIWYG(true)
	l.SetSanitize(true)
	tests := []struct {
		name, input, expected string
	}{
		{"missing", `<audio src="assets/audio.mp3"></audio>`, `<audio src="assets/audio.mp3" controls="controls"></audio>`},
		{"existing", `<audio controls="controls" src="assets/audio.mp3"></audio>`, `<audio controls="controls" src="assets/audio.mp3"></audio>`},
		{"boolean", `<audio controls src="assets/audio.mp3"></audio>`, `<audio controls="" src="assets/audio.mp3"></audio>`},
		{"attributes", `<audio src="https://example.com/audio.mp3?a=1&amp;b=2" preload="none" muted loop></audio>`, `<audio src="https://example.com/audio.mp3?a=1&amp;b=2" preload="none" muted="" loop="" controls="controls"></audio>`},
		{"source", `<audio><source src="assets/audio.mp3" type="audio/mpeg"></audio>`, `<audio controls="controls"><source src="assets/audio.mp3" type="audio/mpeg"/></audio>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, dom := range map[string]string{
				"markdown":       l.Md2BlockDOM(tt.input, false),
				"html":           l.HTML2BlockDOM(tt.input),
				"existing-block": l.SpinBlockDOM(`<div data-node-id="20260918120000-abcdefg" data-type="NodeAudio" class="iframe"><div class="iframe-content">` + tt.input + `</div></div>`),
			} {
				if !strings.Contains(dom, tt.expected) {
					t.Errorf("%s: expected %s in %s", name, tt.expected, dom)
				}
				spun := l.SpinBlockDOM(dom)
				if !strings.Contains(spun, tt.expected) || strings.Count(spun, " controls=") != 1 {
					t.Errorf("%s: controls or attributes changed after round trip: %s", name, spun)
				}
			}
		})
	}
}
