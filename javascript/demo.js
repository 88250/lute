require("./lute.min.js")

const lute = Lute.New()

const renderers = {
  renderText: (node, entering) => {
    if (entering) {
      console.log("    render text")
      return [node.Text() + " via Lute", Lute.WalkContinue]
    }
    return ["", Lute.WalkContinue]
  },
  renderString: (node, entering) => {
    entering ? console.log("    start render string") : console.log("    end render string")
    return ["", Lute.WalkContinue]
  },
  renderParagraph: (node, entering) => {
    entering ? console.log("    start render paragraph") : console.log("    end render paragraph")
    return ["", Lute.WalkContinue]
  }
}

lute.SetJSRenderers({
  renderers: {
    Md2HTML: renderers
  },
})

const markdown = "**Markdown**"
console.log("\nmarkdown input:", markdown, "\n")
let result = lute.MarkdownStr("", markdown)
console.log("\nfinal render output:", result)
