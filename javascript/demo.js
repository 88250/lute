require("./lute.min.js")

const lute = Lute.New()

const renderers = {
  renderText: (node, entering) => {
    console.log("    start render text")
    return [node.Text() + " via Lute", Lute.WalkStop]
  },
  renderStrong: (node, entering) => {
    entering ? console.log("    start render strong") : console.log("    end render strong")
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