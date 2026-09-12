package parse

import (
	"bytes"
	"testing"
)

func TestFileAnnotationRefRejectsNestedOpenersWithoutConsumingInput(t *testing.T) {
	context := &Context{}
	tokens := bytes.Repeat([]byte("<"), 65536)
	passed, remains, id := context.parseFileAnnotationRefID(tokens)
	if len(passed) != 0 || len(id) != 0 || len(remains) != len(tokens) {
		t.Fatal("invalid annotation consumed input after a nested opener")
	}
}
