package test

import (
	"os"
	"testing"

	"github.com/alecthomas/chroma/quick"
)

func TestChroma(t *testing.T) {
	java := `
	@RequestProcessing("/")
	public void index(final RequestContext context) {
		context.setRenderer(new SimpleFMRenderer("index.ftl"));
		final Map<String, Object> dataModel = context.getRenderer().getRenderDataModel();
		dataModel.put("greeting", "Hello, Latke!");
	}`
	err := quick.Highlight(os.Stdout, java, "java", "html", "github")
	if nil != err {
		t.Fatalf(err.Error())
	}

}
