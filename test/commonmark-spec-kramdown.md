# Introduction
{: id="20210408153138-5wxl0jy"}

## What is Markdown?
{: id="20210408153138-1yrt53j"}

Markdown is a plain text format for writing structured documents,
based on conventions for indicating formatting in email
and usenet posts.  It was developed by John Gruber (with
help from Aaron Swartz) and released in 2004 in the form of a
[syntax description](http://daringfireball.net/projects/markdown/syntax)
and a Perl script (`Markdown.pl`) for converting Markdown to
HTML.  In the next decade, dozens of implementations were
developed in many languages.  Some extended the original
Markdown syntax with conventions for footnotes, tables, and
other document elements.  Some allowed Markdown documents to be
rendered in formats other than HTML.  Websites like Reddit,
StackOverflow, and GitHub had millions of people using Markdown.
And Markdown started to be used beyond the web, to author books,
articles, slide shows, letters, and lecture notes.
{: id="20210408153138-unkdwim"}

What distinguishes Markdown from many other lightweight markup
syntaxes, which are often easier to write, is its readability.
As Gruber writes:
{: id="20210408153138-3kdk9kn"}

> The overriding design goal for Markdown's formatting syntax is
> to make it as readable as possible. The idea is that a
> Markdown-formatted document should be publishable as-is, as
> plain text, without looking like it's been marked up with tags
> or formatting instructions.
> ([http://daringfireball.net/projects/markdown/](http://daringfireball.net/projects/markdown/))
> {: id="20210408153138-b9zg01c"}
{: id="20210408153138-13s3vfi"}

The point can be illustrated by comparing a sample of
[AsciiDoc](http://www.methods.co.nz/asciidoc/) with
an equivalent sample of Markdown.  Here is a sample of
AsciiDoc from the AsciiDoc manual:
{: id="20210408153138-slpj2nu"}

```
1. List item one.
+
List item one continued with a second paragraph followed by an
Indented block.
+
.................
$ ls *.sh
$ mv *.sh ~/tmp
.................
+
List item continued with a third paragraph.

2. List item two continued with an open block.
+
--
This paragraph is part of the preceding list item.

a. This list is nested and does not require explicit item
continuation.
+
This paragraph is part of the preceding list item.

b. List item b.

This paragraph belongs to item two of the outer list.
--
```
{: id="20210408153138-lrl5a80"}

And here is the equivalent in Markdown:
{: id="20210408153138-e1y76ke"}

```
1.  List item one.

    List item one continued with a second paragraph followed by an
    Indented block.

        $ ls *.sh
        $ mv *.sh ~/tmp

    List item continued with a third paragraph.

2.  List item two continued with an open block.

    This paragraph is part of the preceding list item.

    1. This list is nested and does not require explicit item continuation.

       This paragraph is part of the preceding list item.

    2. List item b.

    This paragraph belongs to item two of the outer list.
```
{: id="20210408153138-5w021y2"}

The AsciiDoc version is, arguably, easier to write. You don't need
to worry about indentation.  But the Markdown version is much easier
to read.  The nesting of list items is apparent to the eye in the
source, not just in the processed document.
{: id="20210408153138-yjcb2ml"}

## Why is a spec needed?
{: id="20210408153138-bjuhko8"}

John Gruber's [canonical description of Markdown's
syntax](http://daringfireball.net/projects/markdown/syntax)
does not specify the syntax unambiguously.  Here are some examples of
questions it does not answer:
{: id="20210408153138-uht8rj8"}

1. {: id="20210408153137-759mgxi"}How much indentation is needed for a sublist?  The spec says that
   continuation paragraphs need to be indented four spaces, but is
   not fully explicit about sublists.  It is natural to think that
   they, too, must be indented four spaces, but `Markdown.pl` does
   not require that.  This is hardly a "corner case," and divergences
   between implementations on this issue often lead to surprises for
   users in real documents. (See [this comment by John
   Gruber](http://article.gmane.org/gmane.text.markdown.general/1997).)
   {: id="20210408153138-x73y6l5"}
2. {: id="20210408153137-69vy34k"}Is a blank line needed before a block quote or heading?
   Most implementations do not require the blank line.  However,
   this can lead to unexpected results in hard-wrapped text, and
   also to ambiguities in parsing (note that some implementations
   put the heading inside the blockquote, while others do not).
   (John Gruber has also spoken [in favor of requiring the blank
   lines](http://article.gmane.org/gmane.text.markdown.general/2146).)
   {: id="20210408153138-i5uw92o"}
3. {: id="20210408153137-ps9d04h"}Is a blank line needed before an indented code block?
   (`Markdown.pl` requires it, but this is not mentioned in the
   documentation, and some implementations do not require it.)
   {: id="20210408153138-8q7ml5l"}
   ```markdown
   paragraph
       code?
   ```
   {: id="20210408153138-19b6c8r"}
4. {: id="20210408153137-vcurqjl"}What is the exact rule for determining when list items get
   wrapped in `<p>` tags?  Can a list be partially "loose" and partially
   "tight"?  What should we do with a list like this?
   {: id="20210408153138-5l6waxg"}
   ```markdown
   1. one

   2. two
   3. three
   ```
   {: id="20210408153138-bn1koqa"}
   Or this?
   {: id="20210408153138-26kmi53"}
   ```markdown
   1.  one
       - a

       - b
   2.  two
   ```
   {: id="20210408153138-y1nyuau"}
   (There are some relevant comments by John Gruber
   [here](http://article.gmane.org/gmane.text.markdown.general/2554).)
   {: id="20210408153138-pua6mjd"}
5. {: id="20210408153137-e1tzayu"}Can list markers be indented?  Can ordered list markers be right-aligned?
   {: id="20210408153138-awtbjca"}
   ```markdown
    8. item 1
    9. item 2
   10. item 2a
   ```
   {: id="20210408153138-xl6yj2p"}
6. {: id="20210408153137-fxudwvh"}Is this one list with a thematic break in its second item,
   or two lists separated by a thematic break?
   {: id="20210408153138-x82y814"}
   ```markdown
   * a
   * * * * *
   * b
   ```
   {: id="20210408153138-aw1a3qd"}
7. {: id="20210408153137-vf9e6a2"}When list markers change from numbers to bullets, do we have
   two lists or one?  (The Markdown syntax description suggests two,
   but the perl scripts and many other implementations produce one.)
   {: id="20210408153138-vkp8o7u"}
   ```markdown
   1. fee
   2. fie
   -  foe
   -  fum
   ```
   {: id="20210408153138-1uy37w8"}
8. {: id="20210408153137-487qtsz"}What are the precedence rules for the markers of inline structure?
   For example, is the following a valid link, or does the code span
   take precedence ?
   {: id="20210408153138-6lg15rf"}
   ```markdown
   [a backtick (`)](/url) and [another backtick (`)](/url).
   ```
   {: id="20210408153138-y0ad2vf"}
9. {: id="20210408153137-d1z44e0"}What are the precedence rules for markers of emphasis and strong
   emphasis?  For example, how should the following be parsed?
   {: id="20210408153138-3s9tyti"}
   ```markdown
   *foo *bar* baz*
   ```
   {: id="20210408153138-ip72o7d"}
10. {: id="20210408153137-e84bfzk"}What are the precedence rules between block-level and inline-level
    structure?  For example, how should the following be parsed?
    {: id="20210408153138-k9cbb7k"}
    ```markdown
    - `a long code span can contain a hyphen like this
      - and it can screw things up`
    ```
    {: id="20210408153138-2s7xarj"}
11. {: id="20210408153137-7myour3"}Can list items include section headings?  (`Markdown.pl` does not
    allow this, but does allow blockquotes to include headings.)
    {: id="20210408153138-d1msyxr"}
    ```markdown
    - # Heading
    ```
    {: id="20210408153138-s8recga"}
12. {: id="20210408153137-v5pg9m0"}Can list items be empty?
    {: id="20210408153138-ps9w26n"}
    ```markdown
    * a
    *
    * b
    ```
    {: id="20210408153138-c85hydx"}
13. {: id="20210408153137-zvb96gk"}Can link references be defined inside block quotes or list items?
    {: id="20210408153138-bdo61aq"}
    ```markdown
    > Blockquote [foo].
    >
    > [foo]: /url
    ```
    {: id="20210408153138-zc2yzad"}
14. {: id="20210408153137-cpl8tnf"}If there are multiple definitions for the same reference, which takes
    precedence?
    {: id="20210408153138-pdprid9"}
    ```markdown
    [foo]: /url1
    [foo]: /url2

    [foo][]
    ```
    {: id="20210408153138-whl6etu"}
{: id="20210408153138-omod3yn"}

In the absence of a spec, early implementers consulted `Markdown.pl`
to resolve these ambiguities.  But `Markdown.pl` was quite buggy, and
gave manifestly bad results in many cases, so it was not a
satisfactory replacement for a spec.
{: id="20210408153138-znsov61"}

Because there is no unambiguous spec, implementations have diverged
considerably.  As a result, users are often surprised to find that
a document that renders one way on one system (say, a GitHub wiki)
renders differently on another (say, converting to docbook using
pandoc).  To make matters worse, because nothing in Markdown counts
as a "syntax error," the divergence often isn't discovered right away.
{: id="20210408153138-pbahvc9"}

## About this document
{: id="20210408153138-w5q4v7x"}

This document attempts to specify Markdown syntax unambiguously.
It contains many examples with side-by-side Markdown and
HTML.  These are intended to double as conformance tests.  An
accompanying script `spec_tests.py` can be used to run the tests
against any Markdown program:
{: id="20210408153138-1bu56iv"}

```
python test/spec_tests.py --spec spec.txt --program PROGRAM
```
{: id="20210408153138-tic8puq"}

Since this document describes how Markdown is to be parsed into
an abstract syntax tree, it would have made sense to use an abstract
representation of the syntax tree instead of HTML.  But HTML is capable
of representing the structural distinctions we need to make, and the
choice of HTML for the tests makes it possible to run the tests against
an implementation without writing an abstract syntax tree renderer.
{: id="20210408153138-al10t1b"}

This document is generated from a text file, `spec.txt`, written
in Markdown with a small extension for the side-by-side tests.
The script `tools/makespec.py` can be used to convert `spec.txt` into
HTML or CommonMark (which can then be converted into other formats).
{: id="20210408153138-tgut4ir"}

In the examples, the `→` character is used to represent tabs.
{: id="20210408153138-mr7kkzt"}

# Preliminaries
{: id="20210408153138-wsq00kb"}

## Characters and lines
{: id="20210408153138-wg18myj"}

Any sequence of [characters] is a valid CommonMark
document.
{: id="20210408153138-6xro8bq"}

A [character](@) is a Unicode code point.  Although some
code points (for example, combining accents) do not correspond to
characters in an intuitive sense, all code points count as characters
for purposes of this spec.
{: id="20210408153138-e3yng5i"}

This spec does not specify an encoding; it thinks of lines as composed
of [characters] rather than bytes.  A conforming parser may be limited
to a certain encoding.
{: id="20210408153138-caopxh7"}

A [line](@) is a sequence of zero or more [characters]
other than newline (`U+000A`) or carriage return (`U+000D`),
followed by a [line ending] or by the end of file.
{: id="20210408153138-9mud4dq"}

A [line ending](@) is a newline (`U+000A`), a carriage return
(`U+000D`) not followed by a newline, or a carriage return and a
following newline.
{: id="20210408153138-spkrw1r"}

A line containing no characters, or a line containing only spaces
(`U+0020`) or tabs (`U+0009`), is called a [blank line](@).
{: id="20210408153138-krayrvp"}

The following definitions of character classes will be used in this spec:
{: id="20210408153138-h0rnxtg"}

A [whitespace character](@) is a space
(`U+0020`), tab (`U+0009`), newline (`U+000A`), line tabulation (`U+000B`),
form feed (`U+000C`), or carriage return (`U+000D`).
{: id="20210408153138-rd7zibf"}

[Whitespace](@) is a sequence of one or more [whitespace
characters].
{: id="20210408153138-ddq7l8i"}

A [Unicode whitespace character](@) is
any code point in the Unicode `Zs` general category, or a tab (`U+0009`),
carriage return (`U+000D`), newline (`U+000A`), or form feed
(`U+000C`).
{: id="20210408153138-h3r0gf2"}

[Unicode whitespace](@) is a sequence of one
or more [Unicode whitespace characters].
{: id="20210408153138-z8ko4ko"}

A [space](@) is `U+0020`.
{: id="20210408153138-uggzd3e"}

A [non-whitespace character](@) is any character
that is not a [whitespace character].
{: id="20210408153138-3n4f3vg"}

An [ASCII punctuation character](@)
is `!`, `"`, `#`, `$`, `%`, `&`, `'`, `(`, `)`,
`*`, `+`, `,`, `-`, `.`, `/` (U+0021–2F),
`:`, `;`, `<`, `=`, `>`, `?`, `@` (U+003A–0040),
`[`, `\`, `]`, `^`, `_`, `` ` `` (U+005B–0060),
`{`, `|`, `}`, or `~` (U+007B–007E).
{: id="20210408153138-tbj2squ"}

A [punctuation character](@) is an [ASCII
punctuation character] or anything in
the general Unicode categories  `Pc`, `Pd`, `Pe`, `Pf`, `Pi`, `Po`, or `Ps`.
{: id="20210408153138-agxvvtl"}

## Tabs
{: id="20210408153138-1a2atux"}

Tabs in lines are not expanded to [spaces].  However,
in contexts where whitespace helps to define block structure,
tabs behave as if they were replaced by spaces with a tab stop
of 4 characters.
{: id="20210408153138-1fme5ba"}

Thus, for example, a tab can be used instead of four spaces
in an indented code block.  (Note, however, that internal
tabs are passed through as literal tabs, not expanded to
spaces.)
{: id="20210408153138-fj5uvt8"}

````````````````````````````````example
→foo→baz→→bim
.
<pre><code>foo→baz→→bim
</code></pre>
````````````````````````````````
{: id="20210408153138-cnmit9r"}

````````````````````````````````example
  →foo→baz→→bim
.
<pre><code>foo→baz→→bim
</code></pre>
````````````````````````````````
{: id="20210408153138-us8kcd1"}

````````````````````````````````example
    a→a
    ὐ→a
.
<pre><code>a→a
ὐ→a
</code></pre>
````````````````````````````````
{: id="20210408153138-i2lnfm0"}

In the following example, a continuation paragraph of a list
item is indented with a tab; this has exactly the same effect
as indentation with four spaces would:
{: id="20210408153138-qyliixs"}

````````````````````````````````example
  - foo

→bar
.
<ul>
<li>
<p>foo</p>
<p>bar</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-2mg6ai6"}

````````````````````````````````example
- foo

→→bar
.
<ul>
<li>
<p>foo</p>
<pre><code>  bar
</code></pre>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-4ycbxiv"}

Normally the `>` that begins a block quote may be followed
optionally by a space, which is not considered part of the
content.  In the following case `>` is followed by a tab,
which is treated as if it were expanded into three spaces.
Since one of these spaces is considered part of the
delimiter, `foo` is considered to be indented six spaces
inside the block quote context, so we get an indented
code block starting with two spaces.
{: id="20210408153138-3crlj1u"}

````````````````````````````````example
>→→foo
.
<blockquote>
<pre><code>  foo
</code></pre>
</blockquote>
````````````````````````````````
{: id="20210408153138-5e7vqcj"}

````````````````````````````````example
-→→foo
.
<ul>
<li>
<pre><code>  foo
</code></pre>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-wdkxljm"}

````````````````````````````````example
    foo
→bar
.
<pre><code>foo
bar
</code></pre>
````````````````````````````````
{: id="20210408153138-ebz7nqd"}

````````````````````````````````example
 - foo
   - bar
→ - baz
.
<ul>
<li>foo
<ul>
<li>bar
<ul>
<li>baz</li>
</ul>
</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-atgon9x"}

````````````````````````````````example
#→Foo
.
<h1>Foo</h1>
````````````````````````````````
{: id="20210408153138-49prcgj"}

````````````````````````````````example
*→*→*→
.
<hr />
````````````````````````````````
{: id="20210408153138-dqdpbvq"}

## Insecure characters
{: id="20210408153138-uoke57m"}

For security reasons, the Unicode character `U+0000` must be replaced
with the REPLACEMENT CHARACTER (`U+FFFD`).
{: id="20210408153138-uw5w4se"}

# Blocks and inlines
{: id="20210408153138-9fbvb6c"}

We can think of a document as a sequence of
[blocks](@)---structural elements like paragraphs, block
quotations, lists, headings, rules, and code blocks.  Some blocks (like
block quotes and list items) contain other blocks; others (like
headings and paragraphs) contain [inline](@) content---text,
links, emphasized text, images, code spans, and so on.
{: id="20210408153138-yc7ay90"}

## Precedence
{: id="20210408153138-qzqnejk"}

Indicators of block structure always take precedence over indicators
of inline structure.  So, for example, the following is a list with
two items, not a list with one item containing a code span:
{: id="20210408153138-xvutdsk"}

````````````````````````````````example
- `one
- two`
.
<ul>
<li>`one</li>
<li>two`</li>
</ul>
````````````````````````````````
{: id="20210408153138-f25btr4"}

This means that parsing can proceed in two steps:  first, the block
structure of the document can be discerned; second, text lines inside
paragraphs, headings, and other block constructs can be parsed for inline
structure.  The second step requires information about link reference
definitions that will be available only at the end of the first
step.  Note that the first step requires processing lines in sequence,
but the second can be parallelized, since the inline parsing of
one block element does not affect the inline parsing of any other.
{: id="20210408153138-e9ntafj"}

## Container blocks and leaf blocks
{: id="20210408153138-v3igaai"}

We can divide blocks into two types:
[container blocks](@),
which can contain other blocks, and [leaf blocks](@),
which cannot.
{: id="20210408153138-ud61hxz"}

# Leaf blocks
{: id="20210408153138-45o4o1f"}

This section describes the different kinds of leaf block that make up a
Markdown document.
{: id="20210408153138-dj6389t"}

## Thematic breaks
{: id="20210408153138-ihcvppg"}

A line consisting of 0-3 spaces of indentation, followed by a sequence
of three or more matching `-`, `_`, or `*` characters, each followed
optionally by any number of spaces or tabs, forms a
[thematic break](@).
{: id="20210408153138-nt3b3de"}

````````````````````````````````example
***
---
___
.
<hr />
<hr />
<hr />
````````````````````````````````
{: id="20210408153138-hnihvdo"}

Wrong characters:
{: id="20210408153138-qsp9wpr"}

````````````````````````````````example
+++
.
<p>+++</p>
````````````````````````````````
{: id="20210408153138-8wzvebw"}

````````````````````````````````example
===
.
<p>===</p>
````````````````````````````````
{: id="20210408153138-ar3ccnd"}

Not enough characters:
{: id="20210408153138-jbzxijl"}

````````````````````````````````example
--
**
__
.
<p>--
**
__</p>
````````````````````````````````
{: id="20210408153138-78ruobd"}

One to three spaces indent are allowed:
{: id="20210408153138-kzmu067"}

````````````````````````````````example
 ***
  ***
   ***
.
<hr />
<hr />
<hr />
````````````````````````````````
{: id="20210408153138-qe06zrn"}

Four spaces is too many:
{: id="20210408153138-i7bsb1x"}

````````````````````````````````example
    ***
.
<pre><code>***
</code></pre>
````````````````````````````````
{: id="20210408153138-b2r63jo"}

````````````````````````````````example
Foo
    ***
.
<p>Foo
***</p>
````````````````````````````````
{: id="20210408153138-f37upqb"}

More than three characters may be used:
{: id="20210408153138-wvykz9u"}

````````````````````````````````example
_____________________________________
.
<hr />
````````````````````````````````
{: id="20210408153138-w5w0tdj"}

Spaces are allowed between the characters:
{: id="20210408153138-rxhc7g0"}

````````````````````````````````example
 - - -
.
<hr />
````````````````````````````````
{: id="20210408153138-qq0vamm"}

````````````````````````````````example
 **  * ** * ** * **
.
<hr />
````````````````````````````````
{: id="20210408153138-ifve2i3"}

````````````````````````````````example
-     -      -      -
.
<hr />
````````````````````````````````
{: id="20210408153138-ilbjbqq"}

Spaces are allowed at the end:
{: id="20210408153138-ssaybem"}

````````````````````````````````example
- - - -    
.
<hr />
````````````````````````````````
{: id="20210408153138-p4m2460"}

However, no other characters may occur in the line:
{: id="20210408153138-ligww61"}

````````````````````````````````example
_ _ _ _ a

a------

---a---
.
<p>_ _ _ _ a</p>
<p>a------</p>
<p>---a---</p>
````````````````````````````````
{: id="20210408153138-bcox982"}

It is required that all of the [non-whitespace characters] be the same.
So, this is not a thematic break:
{: id="20210408153138-rjd5qgd"}

````````````````````````````````example
 *-*
.
<p><em>-</em></p>
````````````````````````````````
{: id="20210408153138-sbily1w"}

Thematic breaks do not need blank lines before or after:
{: id="20210408153138-ix58mka"}

````````````````````````````````example
- foo
***
- bar
.
<ul>
<li>foo</li>
</ul>
<hr />
<ul>
<li>bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-p81d9uo"}

Thematic breaks can interrupt a paragraph:
{: id="20210408153138-b821mn5"}

````````````````````````````````example
Foo
***
bar
.
<p>Foo</p>
<hr />
<p>bar</p>
````````````````````````````````
{: id="20210408153138-uyx9hvd"}

If a line of dashes that meets the above conditions for being a
thematic break could also be interpreted as the underline of a [setext
heading], the interpretation as a
[setext heading] takes precedence. Thus, for example,
this is a setext heading, not a paragraph followed by a thematic break:
{: id="20210408153138-lngezov"}

````````````````````````````````example
Foo
---
bar
.
<h2>Foo</h2>
<p>bar</p>
````````````````````````````````
{: id="20210408153138-d9yaeyi"}

When both a thematic break and a list item are possible
interpretations of a line, the thematic break takes precedence:
{: id="20210408153138-jnvnxyn"}

````````````````````````````````example
* Foo
* * *
* Bar
.
<ul>
<li>Foo</li>
</ul>
<hr />
<ul>
<li>Bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-sciwdv8"}

If you want a thematic break in a list item, use a different bullet:
{: id="20210408153138-qzqvehy"}

````````````````````````````````example
- Foo
- * * *
.
<ul>
<li>Foo</li>
<li>
<hr />
</li>
</ul>
````````````````````````````````
{: id="20210408153138-c3oq1bk"}

## ATX headings
{: id="20210408153138-cqjiawd"}

An [ATX heading](@)
consists of a string of characters, parsed as inline content, between an
opening sequence of 1--6 unescaped `#` characters and an optional
closing sequence of any number of unescaped `#` characters.
The opening sequence of `#` characters must be followed by a
[space] or by the end of line. The optional closing sequence of `#`s must be
preceded by a [space] and may be followed by spaces only.  The opening
`#` character may be indented 0-3 spaces.  The raw contents of the
heading are stripped of leading and trailing spaces before being parsed
as inline content.  The heading level is equal to the number of `#`
characters in the opening sequence.
{: id="20210408153138-an7ekoh"}

Simple headings:
{: id="20210408153138-og1ta1n"}

````````````````````````````````example
# foo
## foo
### foo
#### foo
##### foo
###### foo
.
<h1>foo</h1>
<h2>foo</h2>
<h3>foo</h3>
<h4>foo</h4>
<h5>foo</h5>
<h6>foo</h6>
````````````````````````````````
{: id="20210408153138-gym45p4"}

More than six `#` characters is not a heading:
{: id="20210408153138-p8b3qp3"}

````````````````````````````````example
####### foo
.
<p>####### foo</p>
````````````````````````````````
{: id="20210408153138-nddz0qv"}

At least one space is required between the `#` characters and the
heading's contents, unless the heading is empty.  Note that many
implementations currently do not require the space.  However, the
space was required by the
[original ATX implementation](http://www.aaronsw.com/2002/atx/atx.py),
and it helps prevent things like the following from being parsed as
headings:
{: id="20210408153138-fy96ffm"}

````````````````````````````````example
#5 bolt

#hashtag
.
<p>#5 bolt</p>
<p>#hashtag</p>
````````````````````````````````
{: id="20210408153138-6ivjybe"}

This is not a heading, because the first `#` is escaped:
{: id="20210408153138-zar0o8k"}

````````````````````````````````example
\## foo
.
<p>## foo</p>
````````````````````````````````
{: id="20210408153138-s62t887"}

Contents are parsed as inlines:
{: id="20210408153138-nom872l"}

````````````````````````````````example
# foo *bar* \*baz\*
.
<h1>foo <em>bar</em> *baz*</h1>
````````````````````````````````
{: id="20210408153138-l4vhlq6"}

Leading and trailing [whitespace] is ignored in parsing inline content:
{: id="20210408153138-kihxhp2"}

````````````````````````````````example
#                  foo                     
.
<h1>foo</h1>
````````````````````````````````
{: id="20210408153138-gbndhi5"}

One to three spaces indentation are allowed:
{: id="20210408153138-ibtdq84"}

````````````````````````````````example
 ### foo
  ## foo
   # foo
.
<h3>foo</h3>
<h2>foo</h2>
<h1>foo</h1>
````````````````````````````````
{: id="20210408153138-8vepags"}

Four spaces are too much:
{: id="20210408153138-t9fvb9e"}

````````````````````````````````example
    # foo
.
<pre><code># foo
</code></pre>
````````````````````````````````
{: id="20210408153138-c9ayng3"}

````````````````````````````````example
foo
    # bar
.
<p>foo
# bar</p>
````````````````````````````````
{: id="20210408153138-twi522x"}

A closing sequence of `#` characters is optional:
{: id="20210408153138-hrbnmhm"}

````````````````````````````````example
## foo ##
  ###   bar    ###
.
<h2>foo</h2>
<h3>bar</h3>
````````````````````````````````
{: id="20210408153138-dw35evz"}

It need not be the same length as the opening sequence:
{: id="20210408153138-tlofp1k"}

````````````````````````````````example
# foo ##################################
##### foo ##
.
<h1>foo</h1>
<h5>foo</h5>
````````````````````````````````
{: id="20210408153138-8dt0bqd"}

Spaces are allowed after the closing sequence:
{: id="20210408153138-vuvxn52"}

````````````````````````````````example
### foo ###     
.
<h3>foo</h3>
````````````````````````````````
{: id="20210408153138-h4m3sd5"}

A sequence of `#` characters with anything but [spaces] following it
is not a closing sequence, but counts as part of the contents of the
heading:
{: id="20210408153138-d2sxdsq"}

````````````````````````````````example
### foo ### b
.
<h3>foo ### b</h3>
````````````````````````````````
{: id="20210408153138-7cax8av"}

The closing sequence must be preceded by a space:
{: id="20210408153138-krcln0c"}

````````````````````````````````example
# foo#
.
<h1>foo#</h1>
````````````````````````````````
{: id="20210408153138-6l056po"}

Backslash-escaped `#` characters do not count as part
of the closing sequence:
{: id="20210408153138-aqrlc1f"}

````````````````````````````````example
### foo \###
## foo #\##
# foo \#
.
<h3>foo ###</h3>
<h2>foo ###</h2>
<h1>foo #</h1>
````````````````````````````````
{: id="20210408153138-hwq0rn1"}

ATX headings need not be separated from surrounding content by blank
lines, and they can interrupt paragraphs:
{: id="20210408153138-ylipii9"}

````````````````````````````````example
****
## foo
****
.
<hr />
<h2>foo</h2>
<hr />
````````````````````````````````
{: id="20210408153138-c00i24g"}

````````````````````````````````example
Foo bar
# baz
Bar foo
.
<p>Foo bar</p>
<h1>baz</h1>
<p>Bar foo</p>
````````````````````````````````
{: id="20210408153138-3zpmzk0"}

ATX headings can be empty:
{: id="20210408153138-qaft30m"}

````````````````````````````````example
## 
#
### ###
.
<h2></h2>
<h1></h1>
<h3></h3>
````````````````````````````````
{: id="20210408153138-tr7dai5"}

## Setext headings
{: id="20210408153138-xqwcxb7"}

A [setext heading](@) consists of one or more
lines of text, each containing at least one [non-whitespace
character], with no more than 3 spaces indentation, followed by
a [setext heading underline].  The lines of text must be such
that, were they not followed by the setext heading underline,
they would be interpreted as a paragraph:  they cannot be
interpretable as a [code fence], [ATX heading][ATX headings],
[block quote][block quotes], [thematic break][thematic breaks],
[list item][list items], or [HTML block][HTML blocks].
{: id="20210408153138-xt4bhrx"}

A [setext heading underline](@) is a sequence of
`=` characters or a sequence of `-` characters, with no more than 3
spaces indentation and any number of trailing spaces.  If a line
containing a single `-` can be interpreted as an
empty [list items], it should be interpreted this way
and not as a [setext heading underline].
{: id="20210408153138-8kqwqfh"}

The heading is a level 1 heading if `=` characters are used in
the [setext heading underline], and a level 2 heading if `-`
characters are used.  The contents of the heading are the result
of parsing the preceding lines of text as CommonMark inline
content.
{: id="20210408153138-f6lk6w8"}

In general, a setext heading need not be preceded or followed by a
blank line.  However, it cannot interrupt a paragraph, so when a
setext heading comes after a paragraph, a blank line is needed between
them.
{: id="20210408153138-07jcxo3"}

Simple examples:
{: id="20210408153138-jjvcuh9"}

````````````````````````````````example
Foo *bar*
=========

Foo *bar*
---------
.
<h1>Foo <em>bar</em></h1>
<h2>Foo <em>bar</em></h2>
````````````````````````````````
{: id="20210408153138-h03hnen"}

The content of the header may span more than one line:
{: id="20210408153138-noz2w58"}

````````````````````````````````example
Foo *bar
baz*
====
.
<h1>Foo <em>bar
baz</em></h1>
````````````````````````````````
{: id="20210408153138-4qy8ath"}

The contents are the result of parsing the headings's raw
content as inlines.  The heading's raw content is formed by
concatenating the lines and removing initial and final
[whitespace].
{: id="20210408153138-fhc551v"}

````````````````````````````````example
  Foo *bar
baz*→
====
.
<h1>Foo <em>bar
baz</em></h1>
````````````````````````````````
{: id="20210408153138-i0eh9wd"}

The underlining can be any length:
{: id="20210408153138-tur6xad"}

````````````````````````````````example
Foo
-------------------------

Foo
=
.
<h2>Foo</h2>
<h1>Foo</h1>
````````````````````````````````
{: id="20210408153138-gs3agsq"}

The heading content can be indented up to three spaces, and need
not line up with the underlining:
{: id="20210408153138-0vtvwww"}

````````````````````````````````example
   Foo
---

  Foo
-----

  Foo
  ===
.
<h2>Foo</h2>
<h2>Foo</h2>
<h1>Foo</h1>
````````````````````````````````
{: id="20210408153138-gya96aw"}

Four spaces indent is too much:
{: id="20210408153138-hcwd0dq"}

````````````````````````````````example
    Foo
    ---

    Foo
---
.
<pre><code>Foo
---

Foo
</code></pre>
<hr />
````````````````````````````````
{: id="20210408153138-iw393b4"}

The setext heading underline can be indented up to three spaces, and
may have trailing spaces:
{: id="20210408153138-st2hvwh"}

````````````````````````````````example
Foo
   ----      
.
<h2>Foo</h2>
````````````````````````````````
{: id="20210408153138-w3084xa"}

Four spaces is too much:
{: id="20210408153138-88ir97h"}

````````````````````````````````example
Foo
    ---
.
<p>Foo
---</p>
````````````````````````````````
{: id="20210408153138-ivf0ax1"}

The setext heading underline cannot contain internal spaces:
{: id="20210408153138-ne1y5sy"}

````````````````````````````````example
Foo
= =

Foo
--- -
.
<p>Foo
= =</p>
<p>Foo</p>
<hr />
````````````````````````````````
{: id="20210408153138-ne0dxfk"}

Trailing spaces in the content line do not cause a line break:
{: id="20210408153138-bhh6y65"}

````````````````````````````````example
Foo  
-----
.
<h2>Foo</h2>
````````````````````````````````
{: id="20210408153138-5luc2y4"}

Nor does a backslash at the end:
{: id="20210408153138-zpyipzb"}

````````````````````````````````example
Foo\
----
.
<h2>Foo\</h2>
````````````````````````````````
{: id="20210408153138-dcwprnv"}

Since indicators of block structure take precedence over
indicators of inline structure, the following are setext headings:
{: id="20210408153138-zhvd6ye"}

````````````````````````````````example
`Foo
----
`

<a title="a lot
---
of dashes"/>
.
<h2>`Foo</h2>
<p>`</p>
<h2>&lt;a title=&quot;a lot</h2>
<p>of dashes&quot;/&gt;</p>
````````````````````````````````
{: id="20210408153138-bn8uees"}

The setext heading underline cannot be a [lazy continuation
line] in a list item or block quote:
{: id="20210408153138-wnfw8dr"}

````````````````````````````````example
> Foo
---
.
<blockquote>
<p>Foo</p>
</blockquote>
<hr />
````````````````````````````````
{: id="20210408153138-mywowxf"}

````````````````````````````````example
> foo
bar
===
.
<blockquote>
<p>foo
bar
===</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-u4eswmk"}

````````````````````````````````example
- Foo
---
.
<ul>
<li>Foo</li>
</ul>
<hr />
````````````````````````````````
{: id="20210408153138-0crdppk"}

A blank line is needed between a paragraph and a following
setext heading, since otherwise the paragraph becomes part
of the heading's content:
{: id="20210408153138-2gz80vs"}

````````````````````````````````example
Foo
Bar
---
.
<h2>Foo
Bar</h2>
````````````````````````````````
{: id="20210408153138-zvj5zgs"}

But in general a blank line is not required before or after
setext headings:
{: id="20210408153138-c42mvja"}

````````````````````````````````example
---
Foo
---
Bar
---
Baz
.
<hr />
<h2>Foo</h2>
<h2>Bar</h2>
<p>Baz</p>
````````````````````````````````
{: id="20210408153138-v80i5vo"}

Setext headings cannot be empty:
{: id="20210408153138-ddi6v2o"}

````````````````````````````````example

====
.
<p>====</p>
````````````````````````````````
{: id="20210408153138-8v6b2fe"}

Setext heading text lines must not be interpretable as block
constructs other than paragraphs.  So, the line of dashes
in these examples gets interpreted as a thematic break:
{: id="20210408153138-h1sjgu8"}

````````````````````````````````example
---
---
.
<hr />
<hr />
````````````````````````````````
{: id="20210408153138-lkt33g0"}

````````````````````````````````example
- foo
-----
.
<ul>
<li>foo</li>
</ul>
<hr />
````````````````````````````````
{: id="20210408153138-4j2m7bf"}

````````````````````````````````example
    foo
---
.
<pre><code>foo
</code></pre>
<hr />
````````````````````````````````
{: id="20210408153138-n905le7"}

````````````````````````````````example
> foo
-----
.
<blockquote>
<p>foo</p>
</blockquote>
<hr />
````````````````````````````````
{: id="20210408153138-9zvqs88"}

If you want a heading with `> foo` as its literal text, you can
use backslash escapes:
{: id="20210408153138-miixysn"}

````````````````````````````````example
\> foo
------
.
<h2>&gt; foo</h2>
````````````````````````````````
{: id="20210408153138-3k8vpxm"}

**Compatibility note:**  Most existing Markdown implementations
do not allow the text of setext headings to span multiple lines.
But there is no consensus about how to interpret
{: id="20210408153138-s5k7lu2"}

```markdown
Foo
bar
---
baz
```
{: id="20210408153138-6jlc4rz"}

One can find four different interpretations:
{: id="20210408153138-43w63hk"}

1. {: id="20210408153137-iw9iew5"}paragraph "Foo", heading "bar", paragraph "baz"
   {: id="20210408153138-7azoosa"}
2. {: id="20210408153137-wedd4ag"}paragraph "Foo bar", thematic break, paragraph "baz"
   {: id="20210408153138-qhjc08k"}
3. {: id="20210408153137-gkyksr9"}paragraph "Foo bar --- baz"
   {: id="20210408153138-d04p4q3"}
4. {: id="20210408153137-68zmriz"}heading "Foo bar", paragraph "baz"
   {: id="20210408153138-7iqhkg0"}
{: id="20210408153138-z1yoxfn"}

We find interpretation 4 most natural, and interpretation 4
increases the expressive power of CommonMark, by allowing
multiline headings.  Authors who want interpretation 1 can
put a blank line after the first paragraph:
{: id="20210408153138-tm7tatp"}

````````````````````````````````example
Foo

bar
---
baz
.
<p>Foo</p>
<h2>bar</h2>
<p>baz</p>
````````````````````````````````
{: id="20210408153138-06zoxpn"}

Authors who want interpretation 2 can put blank lines around
the thematic break,
{: id="20210408153138-wisdhh5"}

````````````````````````````````example
Foo
bar

---

baz
.
<p>Foo
bar</p>
<hr />
<p>baz</p>
````````````````````````````````
{: id="20210408153138-vgluxaz"}

or use a thematic break that cannot count as a [setext heading
underline], such as
{: id="20210408153138-9y69ddi"}

````````````````````````````````example
Foo
bar
* * *
baz
.
<p>Foo
bar</p>
<hr />
<p>baz</p>
````````````````````````````````
{: id="20210408153138-z1tzkve"}

Authors who want interpretation 3 can use backslash escapes:
{: id="20210408153138-gefdj6d"}

````````````````````````````````example
Foo
bar
\---
baz
.
<p>Foo
bar
---
baz</p>
````````````````````````````````
{: id="20210408153138-0tx5p6f"}

## Indented code blocks
{: id="20210408153138-x582uxe"}

An [indented code block](@) is composed of one or more
[indented chunks] separated by blank lines.
An [indented chunk](@) is a sequence of non-blank lines,
each indented four or more spaces. The contents of the code block are
the literal contents of the lines, including trailing
[line endings], minus four spaces of indentation.
An indented code block has no [info string].
{: id="20210408153138-s9i7v46"}

An indented code block cannot interrupt a paragraph, so there must be
a blank line between a paragraph and a following indented code block.
(A blank line is not needed, however, between a code block and a following
paragraph.)
{: id="20210408153138-dsm60oc"}

````````````````````````````````example
    a simple
      indented code block
.
<pre><code>a simple
  indented code block
</code></pre>
````````````````````````````````
{: id="20210408153138-cgt4l1w"}

If there is any ambiguity between an interpretation of indentation
as a code block and as indicating that material belongs to a [list
item][list items], the list item interpretation takes precedence:
{: id="20210408153138-kcv0ohd"}

````````````````````````````````example
  - foo

    bar
.
<ul>
<li>
<p>foo</p>
<p>bar</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-w8puzxs"}

````````````````````````````````example
1.  foo

    - bar
.
<ol>
<li>
<p>foo</p>
<ul>
<li>bar</li>
</ul>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-cbtj2xu"}

The contents of a code block are literal text, and do not get parsed
as Markdown:
{: id="20210408153138-vt4hk5x"}

````````````````````````````````example
    <a/>
    *hi*

    - one
.
<pre><code>&lt;a/&gt;
*hi*

- one
</code></pre>
````````````````````````````````
{: id="20210408153138-birz8q7"}

Here we have three chunks separated by blank lines:
{: id="20210408153138-t7bfv8c"}

````````````````````````````````example
    chunk1

    chunk2
  
 
 
    chunk3
.
<pre><code>chunk1

chunk2



chunk3
</code></pre>
````````````````````````````````
{: id="20210408153138-wmegqtv"}

Any initial spaces beyond four will be included in the content, even
in interior blank lines:
{: id="20210408153138-yl12j2d"}

````````````````````````````````example
    chunk1
      
      chunk2
.
<pre><code>chunk1
  
  chunk2
</code></pre>
````````````````````````````````
{: id="20210408153138-dsvin7l"}

An indented code block cannot interrupt a paragraph.  (This
allows hanging indents and the like.)
{: id="20210408153138-7ii10zh"}

````````````````````````````````example
Foo
    bar

.
<p>Foo
bar</p>
````````````````````````````````
{: id="20210408153138-v02osf7"}

However, any non-blank line with fewer than four leading spaces ends
the code block immediately.  So a paragraph may occur immediately
after indented code:
{: id="20210408153138-oils4qc"}

````````````````````````````````example
    foo
bar
.
<pre><code>foo
</code></pre>
<p>bar</p>
````````````````````````````````
{: id="20210408153138-kc6qj7k"}

And indented code can occur immediately before and after other kinds of
blocks:
{: id="20210408153138-swpyygh"}

````````````````````````````````example
# Heading
    foo
Heading
------
    foo
----
.
<h1>Heading</h1>
<pre><code>foo
</code></pre>
<h2>Heading</h2>
<pre><code>foo
</code></pre>
<hr />
````````````````````````````````
{: id="20210408153138-1bledii"}

The first line can be indented more than four spaces:
{: id="20210408153138-4fxit3i"}

````````````````````````````````example
        foo
    bar
.
<pre><code>    foo
bar
</code></pre>
````````````````````````````````
{: id="20210408153138-yfkfsae"}

Blank lines preceding or following an indented code block
are not included in it:
{: id="20210408153138-dqrsk1t"}

````````````````````````````````example

    
    foo
    

.
<pre><code>foo
</code></pre>
````````````````````````````````
{: id="20210408153138-idy5fu3"}

Trailing spaces are included in the code block's content:
{: id="20210408153138-e5eyxuj"}

````````````````````````````````example
    foo  
.
<pre><code>foo  
</code></pre>
````````````````````````````````
{: id="20210408153138-3221ffv"}

## Fenced code blocks
{: id="20210408153138-qmz81lk"}

A [code fence](@) is a sequence
of at least three consecutive backtick characters (`` ` ``) or
tildes (`~`).  (Tildes and backticks cannot be mixed.)
A [fenced code block](@)
begins with a code fence, indented no more than three spaces.
{: id="20210408153138-t36mldy"}

The line with the opening code fence may optionally contain some text
following the code fence; this is trimmed of leading and trailing
whitespace and called the [info string](@). If the [info string] comes
after a backtick fence, it may not contain any backtick
characters.  (The reason for this restriction is that otherwise
some inline code would be incorrectly interpreted as the
beginning of a fenced code block.)
{: id="20210408153138-odftvy6"}

The content of the code block consists of all subsequent lines, until
a closing [code fence] of the same type as the code block
began with (backticks or tildes), and with at least as many backticks
or tildes as the opening code fence.  If the leading code fence is
indented N spaces, then up to N spaces of indentation are removed from
each line of the content (if present).  (If a content line is not
indented, it is preserved unchanged.  If it is indented less than N
spaces, all of the indentation is removed.)
{: id="20210408153138-hrgzyfi"}

The closing code fence may be indented up to three spaces, and may be
followed only by spaces, which are ignored.  If the end of the
containing block (or document) is reached and no closing code fence
has been found, the code block contains all of the lines after the
opening code fence until the end of the containing block (or
document).  (An alternative spec would require backtracking in the
event that a closing code fence is not found.  But this makes parsing
much less efficient, and there seems to be no real down side to the
behavior described here.)
{: id="20210408153138-nir4d0c"}

A fenced code block may interrupt a paragraph, and does not require
a blank line either before or after.
{: id="20210408153138-8kbo036"}

The content of a code fence is treated as literal text, not parsed
as inlines.  The first word of the [info string] is typically used to
specify the language of the code sample, and rendered in the `class`
attribute of the `code` tag.  However, this spec does not mandate any
particular treatment of the [info string].
{: id="20210408153138-mtkdomk"}

Here is a simple example with backticks:
{: id="20210408153138-pz20v0j"}

````````````````````````````````example
```
<
 >
```
.
<pre><code>&lt;
 &gt;
</code></pre>
````````````````````````````````
{: id="20210408153138-kx9w4nd"}

With tildes:
{: id="20210408153138-erhliip"}

````````````````````````````````example
~~~
<
 >
~~~
.
<pre><code>&lt;
 &gt;
</code></pre>
````````````````````````````````
{: id="20210408153138-u6lbxx3"}

Fewer than three backticks is not enough:
{: id="20210408153138-xdxf42w"}

````````````````````````````````example
``
foo
``
.
<p><code>foo</code></p>
````````````````````````````````
{: id="20210408153138-pe2b0cy"}

The closing code fence must use the same character as the opening
fence:
{: id="20210408153138-0v08rd4"}

````````````````````````````````example
```
aaa
~~~
```
.
<pre><code>aaa
~~~
</code></pre>
````````````````````````````````
{: id="20210408153138-dkl2w4w"}

````````````````````````````````example
~~~
aaa
```
~~~
.
<pre><code>aaa
```
</code></pre>
````````````````````````````````
{: id="20210408153138-fsoacs9"}

The closing code fence must be at least as long as the opening fence:
{: id="20210408153138-msf8u6w"}

````````````````````````````````example
````
aaa
```
``````
.
<pre><code>aaa
```
</code></pre>
````````````````````````````````
{: id="20210408153138-43f6toq"}

````````````````````````````````example
~~~~
aaa
~~~
~~~~
.
<pre><code>aaa
~~~
</code></pre>
````````````````````````````````
{: id="20210408153138-ab635ws"}

Unclosed code blocks are closed by the end of the document
(or the enclosing [block quote][block quotes] or [list item][list items]):
{: id="20210408153138-55p8nlt"}

````````````````````````````````example
```
.
<pre><code></code></pre>
````````````````````````````````
{: id="20210408153138-yhumtrf"}

````````````````````````````````example
`````

```
aaa
.
<pre><code>
```
aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-2y009ji"}

````````````````````````````````example
> ```
> aaa

bbb
.
<blockquote>
<pre><code>aaa
</code></pre>
</blockquote>
<p>bbb</p>
````````````````````````````````
{: id="20210408153138-21vkz7u"}

A code block can have all empty lines as its content:
{: id="20210408153138-hjfr7ko"}

````````````````````````````````example
```

  
```
.
<pre><code>
  
</code></pre>
````````````````````````````````
{: id="20210408153138-h46klyr"}

A code block can be empty:
{: id="20210408153138-n9grk5f"}

````````````````````````````````example
```
```
.
<pre><code></code></pre>
````````````````````````````````
{: id="20210408153138-gncbod9"}

Fences can be indented.  If the opening fence is indented,
content lines will have equivalent opening indentation removed,
if present:
{: id="20210408153138-9le6lgu"}

````````````````````````````````example
 ```
 aaa
aaa
```
.
<pre><code>aaa
aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-91n9rhk"}

````````````````````````````````example
  ```
aaa
  aaa
aaa
  ```
.
<pre><code>aaa
aaa
aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-uljv120"}

````````````````````````````````example
   ```
   aaa
    aaa
  aaa
   ```
.
<pre><code>aaa
 aaa
aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-m7tdorf"}

Four spaces indentation produces an indented code block:
{: id="20210408153138-xerb97z"}

````````````````````````````````example
    ```
    aaa
    ```
.
<pre><code>```
aaa
```
</code></pre>
````````````````````````````````
{: id="20210408153138-ivh9orp"}

Closing fences may be indented by 0-3 spaces, and their indentation
need not match that of the opening fence:
{: id="20210408153138-tunjbsj"}

````````````````````````````````example
```
aaa
  ```
.
<pre><code>aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-oqpmqs9"}

````````````````````````````````example
   ```
aaa
  ```
.
<pre><code>aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-l3ci1n7"}

This is not a closing fence, because it is indented 4 spaces:
{: id="20210408153138-81t7ax6"}

````````````````````````````````example
```
aaa
    ```
.
<pre><code>aaa
    ```
</code></pre>
````````````````````````````````
{: id="20210408153138-n1f6s85"}

Code fences (opening and closing) cannot contain internal spaces:
{: id="20210408153138-g4a7k3n"}

````````````````````````````````example
``` ```
aaa
.
<p><code> </code>
aaa</p>
````````````````````````````````
{: id="20210408153138-9c690dv"}

````````````````````````````````example
~~~~~~
aaa
~~~ ~~
.
<pre><code>aaa
~~~ ~~
</code></pre>
````````````````````````````````
{: id="20210408153138-scgsr88"}

Fenced code blocks can interrupt paragraphs, and can be followed
directly by paragraphs, without a blank line between:
{: id="20210408153138-xqtbhrm"}

````````````````````````````````example
foo
```
bar
```
baz
.
<p>foo</p>
<pre><code>bar
</code></pre>
<p>baz</p>
````````````````````````````````
{: id="20210408153138-nm83bv7"}

Other blocks can also occur before and after fenced code blocks
without an intervening blank line:
{: id="20210408153138-ofpibsl"}

````````````````````````````````example
foo
---
~~~
bar
~~~
# baz
.
<h2>foo</h2>
<pre><code>bar
</code></pre>
<h1>baz</h1>
````````````````````````````````
{: id="20210408153138-xnfndwh"}

An [info string] can be provided after the opening code fence.
Although this spec doesn't mandate any particular treatment of
the info string, the first word is typically used to specify
the language of the code block. In HTML output, the language is
normally indicated by adding a class to the `code` element consisting
of `language-` followed by the language name.
{: id="20210408153138-1f00g5w"}

````````````````````````````````example
```ruby
def foo(x)
  return 3
end
```
.
<pre><code class="language-ruby">def foo(x)
  return 3
end
</code></pre>
````````````````````````````````
{: id="20210408153138-0wdlow9"}

````````````````````````````````example
~~~~    ruby startline=3 $%@#$
def foo(x)
  return 3
end
~~~~~~~
.
<pre><code class="language-ruby">def foo(x)
  return 3
end
</code></pre>
````````````````````````````````
{: id="20210408153138-cgketyn"}

````````````````````````````````example
````;
````
.
<pre><code class="language-;"></code></pre>
````````````````````````````````
{: id="20210408153138-gob99nh"}

[Info strings] for backtick code blocks cannot contain backticks:
{: id="20210408153138-i09gauk"}

````````````````````````````````example
``` aa ```
foo
.
<p><code>aa</code>
foo</p>
````````````````````````````````
{: id="20210408153138-msry4yk"}

[Info strings] for tilde code blocks can contain backticks and tildes:
{: id="20210408153138-yhsabg4"}

````````````````````````````````example
~~~ aa ``` ~~~
foo
~~~
.
<pre><code class="language-aa">foo
</code></pre>
````````````````````````````````
{: id="20210408153138-uqapkqs"}

Closing code fences cannot have [info strings]:
{: id="20210408153138-9sbon38"}

````````````````````````````````example
```
``` aaa
```
.
<pre><code>``` aaa
</code></pre>
````````````````````````````````
{: id="20210408153138-sbpa7j5"}

## HTML blocks
{: id="20210408153138-cm3o8sl"}

An [HTML block](@) is a group of lines that is treated
as raw HTML (and will not be escaped in HTML output).
{: id="20210408153138-mm2epqw"}

There are seven kinds of [HTML block], which can be defined by their
start and end conditions.  The block begins with a line that meets a
[start condition](@) (after up to three spaces optional indentation).
It ends with the first subsequent line that meets a matching [end
condition](@), or the last line of the document, or the last line of
the [container block](#container-blocks) containing the current HTML
block, if no line is encountered that meets the [end condition].  If
the first line meets both the [start condition] and the [end
condition], the block will contain just that line.
{: id="20210408153138-2wt3uwg"}

1. {: id="20210408153137-2ok5xx1"}**Start condition:**  line begins with the string `<script`,
   `<pre`, or `<style` (case-insensitive), followed by whitespace,
   the string `>`, or the end of the line.
   **End condition:**  line contains an end tag
   `</script>`, `</pre>`, or `</style>` (case-insensitive; it
   need not match the start tag).
   {: id="20210408153138-y6k8xdj"}
2. {: id="20210408153137-xykq8gg"}**Start condition:** line begins with the string `<!--`.
   **End condition:**  line contains the string `-->`.
   {: id="20210408153138-7i9csq1"}
3. {: id="20210408153137-4mhtoh0"}**Start condition:** line begins with the string `<?`.
   **End condition:** line contains the string `?>`.
   {: id="20210408153138-6dqh4bp"}
4. {: id="20210408153137-tdsqe9i"}**Start condition:** line begins with the string `<!`
   followed by an uppercase ASCII letter.
   **End condition:** line contains the character `>`.
   {: id="20210408153138-wcffr34"}
5. {: id="20210408153137-qyg4ybn"}**Start condition:**  line begins with the string
   `<![CDATA[`.
   **End condition:** line contains the string `]]>`.
   {: id="20210408153138-5h8vh1a"}
6. {: id="20210408153137-bbu83n6"}**Start condition:** line begins the string `<` or `</`
   followed by one of the strings (case-insensitive) `address`,
   `article`, `aside`, `base`, `basefont`, `blockquote`, `body`,
   `caption`, `center`, `col`, `colgroup`, `dd`, `details`, `dialog`,
   `dir`, `div`, `dl`, `dt`, `fieldset`, `figcaption`, `figure`,
   `footer`, `form`, `frame`, `frameset`,
   `h1`, `h2`, `h3`, `h4`, `h5`, `h6`, `head`, `header`, `hr`,
   `html`, `iframe`, `legend`, `li`, `link`, `main`, `menu`, `menuitem`,
   `nav`, `noframes`, `ol`, `optgroup`, `option`, `p`, `param`,
   `section`, `source`, `summary`, `table`, `tbody`, `td`,
   `tfoot`, `th`, `thead`, `title`, `tr`, `track`, `ul`, followed
   by [whitespace], the end of the line, the string `>`, or
   the string `/>`.
   **End condition:** line is followed by a [blank line].
   {: id="20210408153138-z67m371"}
7. {: id="20210408153137-kfhx4o3"}**Start condition:**  line begins with a complete [open tag]
   (with any [tag name] other than `script`,
   `style`, or `pre`) or a complete [closing tag],
   followed only by [whitespace] or the end of the line.
   **End condition:** line is followed by a [blank line].
   {: id="20210408153138-0xo3bnb"}
{: id="20210408153138-i59bwrk"}

HTML blocks continue until they are closed by their appropriate
[end condition], or the last line of the document or other [container
block](#container-blocks).  This means any HTML **within an HTML
block** that might otherwise be recognised as a start condition will
be ignored by the parser and passed through as-is, without changing
the parser's state.
{: id="20210408153138-ur9e4j8"}

For instance, `<pre>` within a HTML block started by `<table>` will not affect
the parser state; as the HTML block was started in by start condition 6, it
will end at any blank line. This can be surprising:
{: id="20210408153138-pltw9k2"}

````````````````````````````````example
<table><tr><td>
<pre>
**Hello**,

_world_.
</pre>
</td></tr></table>
.
<table><tr><td>
<pre>
**Hello**,
<p><em>world</em>.
</pre></p>
</td></tr></table>
````````````````````````````````
{: id="20210408153138-uabzm30"}

In this case, the HTML block is terminated by the newline — the `**Hello**`
text remains verbatim — and regular parsing resumes, with a paragraph,
emphasised `world` and inline and block HTML following.
{: id="20210408153138-hvkprt3"}

All types of [HTML blocks] except type 7 may interrupt
a paragraph.  Blocks of type 7 may not interrupt a paragraph.
(This restriction is intended to prevent unwanted interpretation
of long tags inside a wrapped paragraph as starting HTML blocks.)
{: id="20210408153138-uj5867c"}

Some simple examples follow.  Here are some basic HTML blocks
of type 6:
{: id="20210408153138-00joqd1"}

````````````````````````````````example
<table>
  <tr>
    <td>
           hi
    </td>
  </tr>
</table>

okay.
.
<table>
  <tr>
    <td>
           hi
    </td>
  </tr>
</table>
<p>okay.</p>
````````````````````````````````
{: id="20210408153138-d14jdxu"}

````````````````````````````````example
 <div>
  *hello*
         <foo><a>
.
 <div>
  *hello*
         <foo><a>
````````````````````````````````
{: id="20210408153138-fjyu9w7"}

A block can also start with a closing tag:
{: id="20210408153138-7s6fc1f"}

````````````````````````````````example
</div>
*foo*
.
</div>
*foo*
````````````````````````````````
{: id="20210408153138-wt0dhid"}

Here we have two HTML blocks with a Markdown paragraph between them:
{: id="20210408153138-vgxbqnt"}

````````````````````````````````example
<DIV CLASS="foo">

*Markdown*

</DIV>
.
<DIV CLASS="foo">
<p><em>Markdown</em></p>
</DIV>
````````````````````````````````
{: id="20210408153138-rm4nnd6"}

The tag on the first line can be partial, as long
as it is split where there would be whitespace:
{: id="20210408153138-m8f7qwh"}

````````````````````````````````example
<div id="foo"
  class="bar">
</div>
.
<div id="foo"
  class="bar">
</div>
````````````````````````````````
{: id="20210408153138-714zzfp"}

````````````````````````````````example
<div id="foo" class="bar
  baz">
</div>
.
<div id="foo" class="bar
  baz">
</div>
````````````````````````````````
{: id="20210408153138-pawkp2s"}

An open tag need not be closed:
{: id="20210408153138-abrwzjs"}

````````````````````````````````example
<div>
*foo*

*bar*
.
<div>
*foo*
<p><em>bar</em></p>
````````````````````````````````
{: id="20210408153138-yknywc2"}

A partial tag need not even be completed (garbage
in, garbage out):
{: id="20210408153138-elqfdlk"}

````````````````````````````````example
<div id="foo"
*hi*
.
<div id="foo"
*hi*
````````````````````````````````
{: id="20210408153138-pqmlp9r"}

````````````````````````````````example
<div class
foo
.
<div class
foo
````````````````````````````````
{: id="20210408153138-k0hlyq2"}

The initial tag doesn't even need to be a valid
tag, as long as it starts like one:
{: id="20210408153138-4nqzvff"}

````````````````````````````````example
<div *???-&&&-<---
*foo*
.
<div *???-&&&-<---
*foo*
````````````````````````````````
{: id="20210408153138-53i1k7r"}

In type 6 blocks, the initial tag need not be on a line by
itself:
{: id="20210408153138-ud0kzpo"}

````````````````````````````````example
<div><a href="bar">*foo*</a></div>
.
<div><a href="bar">*foo*</a></div>
````````````````````````````````
{: id="20210408153138-8rs2hm2"}

````````````````````````````````example
<table><tr><td>
foo
</td></tr></table>
.
<table><tr><td>
foo
</td></tr></table>
````````````````````````````````
{: id="20210408153138-lmwvjtd"}

Everything until the next blank line or end of document
gets included in the HTML block.  So, in the following
example, what looks like a Markdown code block
is actually part of the HTML block, which continues until a blank
line or the end of the document is reached:
{: id="20210408153138-s0ev25i"}

````````````````````````````````example
<div></div>
``` c
int x = 33;
```
.
<div></div>
``` c
int x = 33;
```
````````````````````````````````
{: id="20210408153138-rm2ijw8"}

To start an [HTML block] with a tag that is *not* in the
list of block-level tags in (6), you must put the tag by
itself on the first line (and it must be complete):
{: id="20210408153138-3fjmyea"}

````````````````````````````````example
<a href="foo">
*bar*
</a>
.
<a href="foo">
*bar*
</a>
````````````````````````````````
{: id="20210408153138-kdpeyjo"}

In type 7 blocks, the [tag name] can be anything:
{: id="20210408153138-m9x88ad"}

````````````````````````````````example
<Warning>
*bar*
</Warning>
.
<Warning>
*bar*
</Warning>
````````````````````````````````
{: id="20210408153138-an2sqws"}

````````````````````````````````example
<i class="foo">
*bar*
</i>
.
<i class="foo">
*bar*
</i>
````````````````````````````````
{: id="20210408153138-u6wf5eo"}

````````````````````````````````example
</ins>
*bar*
.
</ins>
*bar*
````````````````````````````````
{: id="20210408153138-1suh0so"}

These rules are designed to allow us to work with tags that
can function as either block-level or inline-level tags.
The `<del>` tag is a nice example.  We can surround content with
`<del>` tags in three different ways.  In this case, we get a raw
HTML block, because the `<del>` tag is on a line by itself:
{: id="20210408153138-0jgfxgm"}

````````````````````````````````example
<del>
*foo*
</del>
.
<del>
*foo*
</del>
````````````````````````````````
{: id="20210408153138-cb7kfug"}

In this case, we get a raw HTML block that just includes
the `<del>` tag (because it ends with the following blank
line).  So the contents get interpreted as CommonMark:
{: id="20210408153138-nk2peom"}

````````````````````````````````example
<del>

*foo*

</del>
.
<del>
<p><em>foo</em></p>
</del>
````````````````````````````````
{: id="20210408153138-kiaexe8"}

Finally, in this case, the `<del>` tags are interpreted
as [raw HTML] *inside* the CommonMark paragraph.  (Because
the tag is not on a line by itself, we get inline HTML
rather than an [HTML block].)
{: id="20210408153138-548ip2b"}

````````````````````````````````example
<del>*foo*</del>
.
<p><del><em>foo</em></del></p>
````````````````````````````````
{: id="20210408153138-lunx73t"}

HTML tags designed to contain literal content
(`script`, `style`, `pre`), comments, processing instructions,
and declarations are treated somewhat differently.
Instead of ending at the first blank line, these blocks
end at the first line containing a corresponding end tag.
As a result, these blocks can contain blank lines:
{: id="20210408153138-7l40nv4"}

A pre tag (type 1):
{: id="20210408153138-kz7o889"}

````````````````````````````````example
<pre language="haskell"><code>
import Text.HTML.TagSoup

main :: IO ()
main = print $ parseTags tags
</code></pre>
okay
.
<pre language="haskell"><code>
import Text.HTML.TagSoup

main :: IO ()
main = print $ parseTags tags
</code></pre>
<p>okay</p>
````````````````````````````````
{: id="20210408153138-enbgrw5"}

A script tag (type 1):
{: id="20210408153138-nmzfzma"}

````````````````````````````````example
<script type="text/javascript">
// JavaScript example

document.getElementById("demo").innerHTML = "Hello JavaScript!";
</script>
okay
.
<script type="text/javascript">
// JavaScript example

document.getElementById("demo").innerHTML = "Hello JavaScript!";
</script>
<p>okay</p>
````````````````````````````````
{: id="20210408153138-ftv56u2"}

A style tag (type 1):
{: id="20210408153138-5ihzk1c"}

````````````````````````````````example
<style
  type="text/css">
h1 {color:red;}

p {color:blue;}
</style>
okay
.
<style
  type="text/css">
h1 {color:red;}

p {color:blue;}
</style>
<p>okay</p>
````````````````````````````````
{: id="20210408153138-qsv8vba"}

If there is no matching end tag, the block will end at the
end of the document (or the enclosing [block quote][block quotes]
or [list item][list items]):
{: id="20210408153138-2h96v20"}

````````````````````````````````example
<style
  type="text/css">

foo
.
<style
  type="text/css">

foo
````````````````````````````````
{: id="20210408153138-pt0zr07"}

````````````````````````````````example
> <div>
> foo

bar
.
<blockquote>
<div>
foo
</blockquote>
<p>bar</p>
````````````````````````````````
{: id="20210408153138-hprz1i4"}

````````````````````````````````example
- <div>
- foo
.
<ul>
<li>
<div>
</li>
<li>foo</li>
</ul>
````````````````````````````````
{: id="20210408153138-9ffjh8l"}

The end tag can occur on the same line as the start tag:
{: id="20210408153138-17wvd64"}

````````````````````````````````example
<style>p{color:red;}</style>
*foo*
.
<style>p{color:red;}</style>
<p><em>foo</em></p>
````````````````````````````````
{: id="20210408153138-lc9ovjl"}

````````````````````````````````example
<!-- foo -->*bar*
*baz*
.
<!-- foo -->*bar*
<p><em>baz</em></p>
````````````````````````````````
{: id="20210408153138-423c5oa"}

Note that anything on the last line after the
end tag will be included in the [HTML block]:
{: id="20210408153138-daq152c"}

````````````````````````````````example
<script>
foo
</script>1. *bar*
.
<script>
foo
</script>1. *bar*
````````````````````````````````
{: id="20210408153138-nl9v9x9"}

A comment (type 2):
{: id="20210408153138-ec4hvnb"}

````````````````````````````````example
<!-- Foo

bar
   baz -->
okay
.
<!-- Foo

bar
   baz -->
<p>okay</p>
````````````````````````````````
{: id="20210408153138-guhta0c"}

A processing instruction (type 3):
{: id="20210408153138-p4m8l2j"}

````````````````````````````````example
<?php

  echo '>';

?>
okay
.
<?php

  echo '>';

?>
<p>okay</p>
````````````````````````````````
{: id="20210408153138-pwzk9cw"}

A declaration (type 4):
{: id="20210408153138-4l0s8uy"}

````````````````````````````````example
<!DOCTYPE html>
.
<!DOCTYPE html>
````````````````````````````````
{: id="20210408153138-12845pt"}

CDATA (type 5):
{: id="20210408153138-0wiyds4"}

````````````````````````````````example
<![CDATA[
function matchwo(a,b)
{
  if (a < b && a < 0) then {
    return 1;

  } else {

    return 0;
  }
}
]]>
okay
.
<![CDATA[
function matchwo(a,b)
{
  if (a < b && a < 0) then {
    return 1;

  } else {

    return 0;
  }
}
]]>
<p>okay</p>
````````````````````````````````
{: id="20210408153138-stzkazv"}

The opening tag can be indented 1-3 spaces, but not 4:
{: id="20210408153138-jgygw7j"}

````````````````````````````````example
  <!-- foo -->

    <!-- foo -->
.
  <!-- foo -->
<pre><code>&lt;!-- foo --&gt;
</code></pre>
````````````````````````````````
{: id="20210408153138-89weusc"}

````````````````````````````````example
  <div>

    <div>
.
  <div>
<pre><code>&lt;div&gt;
</code></pre>
````````````````````````````````
{: id="20210408153138-w0ml7ot"}

An HTML block of types 1--6 can interrupt a paragraph, and need not be
preceded by a blank line.
{: id="20210408153138-9m6ikkc"}

````````````````````````````````example
Foo
<div>
bar
</div>
.
<p>Foo</p>
<div>
bar
</div>
````````````````````````````````
{: id="20210408153138-443hylb"}

However, a following blank line is needed, except at the end of
a document, and except for blocks of types 1--5, [above][HTML
block]:
{: id="20210408153138-wtxwlx8"}

````````````````````````````````example
<div>
bar
</div>
*foo*
.
<div>
bar
</div>
*foo*
````````````````````````````````
{: id="20210408153138-rn26e5j"}

HTML blocks of type 7 cannot interrupt a paragraph:
{: id="20210408153138-9ddfw76"}

````````````````````````````````example
Foo
<a href="bar">
baz
.
<p>Foo
<a href="bar">
baz</p>
````````````````````````````````
{: id="20210408153138-jg5maxb"}

This rule differs from John Gruber's original Markdown syntax
specification, which says:
{: id="20210408153138-yzx1shj"}

> The only restrictions are that block-level HTML elements —
> e.g. `<div>`, `<table>`, `<pre>`, `<p>`, etc. — must be separated from
> surrounding content by blank lines, and the start and end tags of the
> block should not be indented with tabs or spaces.
> {: id="20210408153138-pzvit7o"}
{: id="20210408153138-d6ygxuy"}

In some ways Gruber's rule is more restrictive than the one given
here:
{: id="20210408153138-zio3swl"}

- {: id="20210408153137-zbq17tw"}It requires that an HTML block be preceded by a blank line.
  {: id="20210408153138-ez2utwv"}
- {: id="20210408153137-kf9yu6v"}It does not allow the start tag to be indented.
  {: id="20210408153138-rv5293z"}
- {: id="20210408153137-ecv4swz"}It requires a matching end tag, which it also does not allow to
  be indented.
  {: id="20210408153138-j195icm"}
{: id="20210408153138-zzmhgvn"}

Most Markdown implementations (including some of Gruber's own) do not
respect all of these restrictions.
{: id="20210408153138-52ekgad"}

There is one respect, however, in which Gruber's rule is more liberal
than the one given here, since it allows blank lines to occur inside
an HTML block.  There are two reasons for disallowing them here.
First, it removes the need to parse balanced tags, which is
expensive and can require backtracking from the end of the document
if no matching end tag is found. Second, it provides a very simple
and flexible way of including Markdown content inside HTML tags:
simply separate the Markdown from the HTML using blank lines:
{: id="20210408153138-34iv7ye"}

Compare:
{: id="20210408153138-v1qn2x4"}

````````````````````````````````example
<div>

*Emphasized* text.

</div>
.
<div>
<p><em>Emphasized</em> text.</p>
</div>
````````````````````````````````
{: id="20210408153138-b3ytewk"}

````````````````````````````````example
<div>
*Emphasized* text.
</div>
.
<div>
*Emphasized* text.
</div>
````````````````````````````````
{: id="20210408153138-pt30b16"}

Some Markdown implementations have adopted a convention of
interpreting content inside tags as text if the open tag has
the attribute `markdown=1`.  The rule given above seems a simpler and
more elegant way of achieving the same expressive power, which is also
much simpler to parse.
{: id="20210408153138-ft2428o"}

The main potential drawback is that one can no longer paste HTML
blocks into Markdown documents with 100% reliability.  However,
*in most cases* this will work fine, because the blank lines in
HTML are usually followed by HTML block tags.  For example:
{: id="20210408153138-j3ict1p"}

````````````````````````````````example
<table>

<tr>

<td>
Hi
</td>

</tr>

</table>
.
<table>
<tr>
<td>
Hi
</td>
</tr>
</table>
````````````````````````````````
{: id="20210408153138-xam9gxf"}

There are problems, however, if the inner tags are indented
*and* separated by spaces, as then they will be interpreted as
an indented code block:
{: id="20210408153138-tiu998m"}

````````````````````````````````example
<table>

  <tr>

    <td>
      Hi
    </td>

  </tr>

</table>
.
<table>
  <tr>
<pre><code>&lt;td&gt;
  Hi
&lt;/td&gt;
</code></pre>
  </tr>
</table>
````````````````````````````````
{: id="20210408153138-jxop4tp"}

Fortunately, blank lines are usually not necessary and can be
deleted.  The exception is inside `<pre>` tags, but as described
[above][HTML blocks], raw HTML blocks starting with `<pre>`
*can* contain blank lines.
{: id="20210408153138-ws7854a"}

## Link reference definitions
{: id="20210408153138-g0fa8u1"}

A [link reference definition](@)
consists of a [link label], indented up to three spaces, followed
by a colon (`:`), optional [whitespace] (including up to one
[line ending]), a [link destination],
optional [whitespace] (including up to one
[line ending]), and an optional [link
title], which if it is present must be separated
from the [link destination] by [whitespace].
No further [non-whitespace characters] may occur on the line.
{: id="20210408153138-mpnlc5k"}

A [link reference definition]
does not correspond to a structural element of a document.  Instead, it
defines a label which can be used in [reference links]
and reference-style [images] elsewhere in the document.  [Link
reference definitions] can come either before or after the links that use
them.
{: id="20210408153138-xbti93i"}

````````````````````````````````example
[foo]: /url "title"

[foo]
.
<p><a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-3dkk8la"}

````````````````````````````````example
   [foo]: 
      /url  
           'the title'  

[foo]
.
<p><a href="/url" title="the title">foo</a></p>
````````````````````````````````
{: id="20210408153138-sf00wnb"}

````````````````````````````````example
[Foo*bar\]]:my_(url) 'title (with parens)'

[Foo*bar\]]
.
<p><a href="my_(url)" title="title (with parens)">Foo*bar]</a></p>
````````````````````````````````
{: id="20210408153138-v7thoo6"}

````````````````````````````````example
[Foo bar]:
<my url>
'title'

[Foo bar]
.
<p><a href="my%20url" title="title">Foo bar</a></p>
````````````````````````````````
{: id="20210408153138-dewy9uo"}

The title may extend over multiple lines:
{: id="20210408153138-zetc4mw"}

````````````````````````````````example
[foo]: /url '
title
line1
line2
'

[foo]
.
<p><a href="/url" title="
title
line1
line2
">foo</a></p>
````````````````````````````````
{: id="20210408153138-n056c3k"}

However, it may not contain a [blank line]:
{: id="20210408153138-bln4sol"}

````````````````````````````````example
[foo]: /url 'title

with blank line'

[foo]
.
<p>[foo]: /url 'title</p>
<p>with blank line'</p>
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-ujvd8fc"}

The title may be omitted:
{: id="20210408153138-hj3qec9"}

````````````````````````````````example
[foo]:
/url

[foo]
.
<p><a href="/url">foo</a></p>
````````````````````````````````
{: id="20210408153138-myig2wo"}

The link destination may not be omitted:
{: id="20210408153138-s5pcyka"}

````````````````````````````````example
[foo]:

[foo]
.
<p>[foo]:</p>
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-1xlyiv3"}

However, an empty link destination may be specified using
angle brackets:
{: id="20210408153138-cg7xqvo"}

````````````````````````````````example
[foo]: <>

[foo]
.
<p><a href="">foo</a></p>
````````````````````````````````
{: id="20210408153138-syihr91"}

The title must be separated from the link destination by
whitespace:
{: id="20210408153138-fecjv62"}

````````````````````````````````example
[foo]: <bar>(baz)

[foo]
.
<p>[foo]: <bar>(baz)</p>
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-zo2bam4"}

Both title and destination can contain backslash escapes
and literal backslashes:
{: id="20210408153138-29teion"}

````````````````````````````````example
[foo]: /url\bar\*baz "foo\"bar\baz"

[foo]
.
<p><a href="/url%5Cbar*baz" title="foo&quot;bar\baz">foo</a></p>
````````````````````````````````
{: id="20210408153138-llmkb0a"}

A link can come before its corresponding definition:
{: id="20210408153138-hrrbmra"}

````````````````````````````````example
[foo]

[foo]: url
.
<p><a href="url">foo</a></p>
````````````````````````````````
{: id="20210408153138-h8547ga"}

If there are several matching definitions, the first one takes
precedence:
{: id="20210408153138-pbd8tn4"}

````````````````````````````````example
[foo]

[foo]: first
[foo]: second
.
<p><a href="first">foo</a></p>
````````````````````````````````
{: id="20210408153138-ejbqvn2"}

As noted in the section on [Links], matching of labels is
case-insensitive (see [matches]).
{: id="20210408153138-1d0udc9"}

````````````````````````````````example
[FOO]: /url

[Foo]
.
<p><a href="/url">Foo</a></p>
````````````````````````````````
{: id="20210408153138-0elp83c"}

````````````````````````````````example
[ΑΓΩ]: /φου

[αγω]
.
<p><a href="/%CF%86%CE%BF%CF%85">αγω</a></p>
````````````````````````````````
{: id="20210408153138-5mlad44"}

Here is a link reference definition with no corresponding link.
It contributes nothing to the document.
{: id="20210408153138-6p4ibgn"}

````````````````````````````````example
[foo]: /url
.
````````````````````````````````
{: id="20210408153138-9a5hz0v"}

Here is another one:
{: id="20210408153138-9rb78af"}

````````````````````````````````example
[
foo
]: /url
bar
.
<p>bar</p>
````````````````````````````````
{: id="20210408153138-jij7768"}

This is not a link reference definition, because there are
[non-whitespace characters] after the title:
{: id="20210408153138-6xqkadh"}

````````````````````````````````example
[foo]: /url "title" ok
.
<p>[foo]: /url &quot;title&quot; ok</p>
````````````````````````````````
{: id="20210408153138-d2yc4k8"}

This is a link reference definition, but it has no title:
{: id="20210408153138-x7glcy0"}

````````````````````````````````example
[foo]: /url
"title" ok
.
<p>&quot;title&quot; ok</p>
````````````````````````````````
{: id="20210408153138-rgfur7k"}

This is not a link reference definition, because it is indented
four spaces:
{: id="20210408153138-qdsoptg"}

````````````````````````````````example
    [foo]: /url "title"

[foo]
.
<pre><code>[foo]: /url &quot;title&quot;
</code></pre>
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-0j6b06i"}

This is not a link reference definition, because it occurs inside
a code block:
{: id="20210408153138-1xx06pu"}

````````````````````````````````example
```
[foo]: /url
```

[foo]
.
<pre><code>[foo]: /url
</code></pre>
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-hcsa05p"}

A [link reference definition] cannot interrupt a paragraph.
{: id="20210408153138-9fxtvmj"}

````````````````````````````````example
Foo
[bar]: /baz

[bar]
.
<p>Foo
[bar]: /baz</p>
<p>[bar]</p>
````````````````````````````````
{: id="20210408153138-j9ezabq"}

However, it can directly follow other block elements, such as headings
and thematic breaks, and it need not be followed by a blank line.
{: id="20210408153138-ia5sz09"}

````````````````````````````````example
# [Foo]
[foo]: /url
> bar
.
<h1><a href="/url">Foo</a></h1>
<blockquote>
<p>bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-jkhquar"}

````````````````````````````````example
[foo]: /url
bar
===
[foo]
.
<h1>bar</h1>
<p><a href="/url">foo</a></p>
````````````````````````````````
{: id="20210408153138-cx1pyor"}

````````````````````````````````example
[foo]: /url
===
[foo]
.
<p>===
<a href="/url">foo</a></p>
````````````````````````````````
{: id="20210408153138-y33ck2m"}

Several [link reference definitions]
can occur one after another, without intervening blank lines.
{: id="20210408153138-20d9fyy"}

````````````````````````````````example
[foo]: /foo-url "foo"
[bar]: /bar-url
  "bar"
[baz]: /baz-url

[foo],
[bar],
[baz]
.
<p><a href="/foo-url" title="foo">foo</a>,
<a href="/bar-url" title="bar">bar</a>,
<a href="/baz-url">baz</a></p>
````````````````````````````````
{: id="20210408153138-3urg4xs"}

[Link reference definitions] can occur
inside block containers, like lists and block quotations.  They
affect the entire document, not just the container in which they
are defined:
{: id="20210408153138-koiytrm"}

````````````````````````````````example
[foo]

> [foo]: /url
.
<p><a href="/url">foo</a></p>
<blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-br65rto"}

Whether something is a [link reference definition] is
independent of whether the link reference it defines is
used in the document.  Thus, for example, the following
document contains just a link reference definition, and
no visible content:
{: id="20210408153138-x42dugs"}

````````````````````````````````example
[foo]: /url
.
````````````````````````````````
{: id="20210408153138-xmku4jo"}

## Paragraphs
{: id="20210408153138-z463gz4"}

A sequence of non-blank lines that cannot be interpreted as other
kinds of blocks forms a [paragraph](@).
The contents of the paragraph are the result of parsing the
paragraph's raw content as inlines.  The paragraph's raw content
is formed by concatenating the lines and removing initial and final
[whitespace].
{: id="20210408153138-ubdrrb5"}

A simple example with two paragraphs:
{: id="20210408153138-i7ngsp0"}

````````````````````````````````example
aaa

bbb
.
<p>aaa</p>
<p>bbb</p>
````````````````````````````````
{: id="20210408153138-zhhivhu"}

Paragraphs can contain multiple lines, but no blank lines:
{: id="20210408153138-upnyudo"}

````````````````````````````````example
aaa
bbb

ccc
ddd
.
<p>aaa
bbb</p>
<p>ccc
ddd</p>
````````````````````````````````
{: id="20210408153138-ncmal43"}

Multiple blank lines between paragraph have no effect:
{: id="20210408153138-70tjje9"}

````````````````````````````````example
aaa


bbb
.
<p>aaa</p>
<p>bbb</p>
````````````````````````````````
{: id="20210408153138-pspy1y9"}

Leading spaces are skipped:
{: id="20210408153138-jvrupe0"}

````````````````````````````````example
  aaa
 bbb
.
<p>aaa
bbb</p>
````````````````````````````````
{: id="20210408153138-9ziyrgm"}

Lines after the first may be indented any amount, since indented
code blocks cannot interrupt paragraphs.
{: id="20210408153138-ljnywiz"}

````````````````````````````````example
aaa
             bbb
                                       ccc
.
<p>aaa
bbb
ccc</p>
````````````````````````````````
{: id="20210408153138-nqfhx9m"}

However, the first line may be indented at most three spaces,
or an indented code block will be triggered:
{: id="20210408153138-pgsybtm"}

````````````````````````````````example
   aaa
bbb
.
<p>aaa
bbb</p>
````````````````````````````````
{: id="20210408153138-sqq0i3g"}

````````````````````````````````example
    aaa
bbb
.
<pre><code>aaa
</code></pre>
<p>bbb</p>
````````````````````````````````
{: id="20210408153138-9s77pdr"}

Final spaces are stripped before inline parsing, so a paragraph
that ends with two or more spaces will not end with a [hard line
break]:
{: id="20210408153138-1h5md5c"}

````````````````````````````````example
aaa     
bbb     
.
<p>aaa<br />
bbb</p>
````````````````````````````````
{: id="20210408153138-d14p93s"}

## Blank lines
{: id="20210408153138-yxf8w65"}

[Blank lines] between block-level elements are ignored,
except for the role they play in determining whether a [list]
is [tight] or [loose].
{: id="20210408153138-57i52cx"}

Blank lines at the beginning and end of the document are also ignored.
{: id="20210408153138-l7pn12m"}

````````````````````````````````example
  

aaa
  

# aaa

  
.
<p>aaa</p>
<h1>aaa</h1>
````````````````````````````````
{: id="20210408153138-q95dxpe"}

# Container blocks
{: id="20210408153138-pr08yfm"}

A [container block](#container-blocks) is a block that has other
blocks as its contents.  There are two basic kinds of container blocks:
[block quotes] and [list items].
[Lists] are meta-containers for [list items].
{: id="20210408153138-39f3hhr"}

We define the syntax for container blocks recursively.  The general
form of the definition is:
{: id="20210408153138-p9hqio8"}

> If X is a sequence of blocks, then the result of
> transforming X in such-and-such a way is a container of type Y
> with these blocks as its content.
> {: id="20210408153138-ayu4qr1"}
{: id="20210408153138-5tssj0c"}

So, we explain what counts as a block quote or list item by explaining
how these can be *generated* from their contents. This should suffice
to define the syntax, although it does not give a recipe for *parsing*
these constructions.  (A recipe is provided below in the section entitled
[A parsing strategy](#appendix-a-parsing-strategy).)
{: id="20210408153138-e5kdanz"}

## Block quotes
{: id="20210408153138-eeig5i2"}

A [block quote marker](@)
consists of 0-3 spaces of initial indent, plus (a) the character `>` together
with a following space, or (b) a single character `>` not followed by a space.
{: id="20210408153138-tfv7nhq"}

The following rules define [block quotes]:
{: id="20210408153138-bkyrda3"}

1. {: id="20210408153137-l7xawf6"}**Basic case.**  If a string of lines *Ls* constitute a sequence
   of blocks *Bs*, then the result of prepending a [block quote
   marker] to the beginning of each line in *Ls*
   is a [block quote](#block-quotes) containing *Bs*.
   {: id="20210408153138-4cosk7g"}
2. {: id="20210408153137-lhzhql0"}**Laziness.**  If a string of lines *Ls* constitute a [block
   quote](#block-quotes) with contents *Bs*, then the result of deleting
   the initial [block quote marker] from one or
   more lines in which the next [non-whitespace character] after the [block
   quote marker] is [paragraph continuation
   text] is a block quote with *Bs* as its content.
   [Paragraph continuation text](@) is text
   that will be parsed as part of the content of a paragraph, but does
   not occur at the beginning of the paragraph.
   {: id="20210408153138-i3jbwrp"}
3. {: id="20210408153137-g2czbac"}**Consecutiveness.**  A document cannot contain two [block
   quotes] in a row unless there is a [blank line] between them.
   {: id="20210408153138-hgo8akw"}
{: id="20210408153138-sn1y9y7"}

Nothing else counts as a [block quote](#block-quotes).
{: id="20210408153138-192q4kz"}

Here is a simple example:
{: id="20210408153138-64wi04h"}

````````````````````````````````example
> # Foo
> bar
> baz
.
<blockquote>
<h1>Foo</h1>
<p>bar
baz</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-ojph1x5"}

The spaces after the `>` characters can be omitted:
{: id="20210408153138-gtinjqw"}

````````````````````````````````example
># Foo
>bar
> baz
.
<blockquote>
<h1>Foo</h1>
<p>bar
baz</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-7wilk2a"}

The `>` characters can be indented 1-3 spaces:
{: id="20210408153138-5utwi0i"}

````````````````````````````````example
   > # Foo
   > bar
 > baz
.
<blockquote>
<h1>Foo</h1>
<p>bar
baz</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-uy8l4uq"}

Four spaces gives us a code block:
{: id="20210408153138-4vmg7vk"}

````````````````````````````````example
    > # Foo
    > bar
    > baz
.
<pre><code>&gt; # Foo
&gt; bar
&gt; baz
</code></pre>
````````````````````````````````
{: id="20210408153138-j0svha0"}

The Laziness clause allows us to omit the `>` before
[paragraph continuation text]:
{: id="20210408153138-pmq8keg"}

````````````````````````````````example
> # Foo
> bar
baz
.
<blockquote>
<h1>Foo</h1>
<p>bar
baz</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-164ow5d"}

A block quote can contain some lazy and some non-lazy
continuation lines:
{: id="20210408153138-t3rpmmi"}

````````````````````````````````example
> bar
baz
> foo
.
<blockquote>
<p>bar
baz
foo</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-joa0hqd"}

Laziness only applies to lines that would have been continuations of
paragraphs had they been prepended with [block quote markers].
For example, the `> ` cannot be omitted in the second line of
{: id="20210408153138-8ijwkwc"}

```markdown
> foo
> ---
```
{: id="20210408153138-5268rop"}

without changing the meaning:
{: id="20210408153138-d7wk50l"}

````````````````````````````````example
> foo
---
.
<blockquote>
<p>foo</p>
</blockquote>
<hr />
````````````````````````````````
{: id="20210408153138-ltv2s8b"}

Similarly, if we omit the `> ` in the second line of
{: id="20210408153138-xopx3f7"}

```markdown
> - foo
> - bar
```
{: id="20210408153138-bfrwq8u"}

then the block quote ends after the first line:
{: id="20210408153138-htojklr"}

````````````````````````````````example
> - foo
- bar
.
<blockquote>
<ul>
<li>foo</li>
</ul>
</blockquote>
<ul>
<li>bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-llt16w4"}

For the same reason, we can't omit the `> ` in front of
subsequent lines of an indented or fenced code block:
{: id="20210408153138-6s6f9ra"}

````````````````````````````````example
>     foo
    bar
.
<blockquote>
<pre><code>foo
</code></pre>
</blockquote>
<pre><code>bar
</code></pre>
````````````````````````````````
{: id="20210408153138-ffr5dej"}

````````````````````````````````example
> ```
foo
```
.
<blockquote>
<pre><code></code></pre>
</blockquote>
<p>foo</p>
<pre><code></code></pre>
````````````````````````````````
{: id="20210408153138-0hgr27y"}

Note that in the following case, we have a [lazy
continuation line]:
{: id="20210408153138-olo1vmw"}

````````````````````````````````example
> foo
    - bar
.
<blockquote>
<p>foo
- bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-756qqxy"}

To see why, note that in
{: id="20210408153138-wsjadl6"}

```markdown
> foo
>     - bar
```
{: id="20210408153138-f4dxkt3"}

the `- bar` is indented too far to start a list, and can't
be an indented code block because indented code blocks cannot
interrupt paragraphs, so it is [paragraph continuation text].
{: id="20210408153138-19kltk4"}

A block quote can be empty:
{: id="20210408153138-ziuywft"}

````````````````````````````````example
>
.
<blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-m6kvsrz"}

````````````````````````````````example
>
>  
> 
.
<blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-fjjewxo"}

A block quote can have initial or final blank lines:
{: id="20210408153138-obftke6"}

````````````````````````````````example
>
> foo
>  
.
<blockquote>
<p>foo</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-q6bavli"}

A blank line always separates block quotes:
{: id="20210408153138-f38ise0"}

````````````````````````````````example
> foo

> bar
.
<blockquote>
<p>foo</p>
</blockquote>
<blockquote>
<p>bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-zzkmgwe"}

(Most current Markdown implementations, including John Gruber's
original `Markdown.pl`, will parse this example as a single block quote
with two paragraphs.  But it seems better to allow the author to decide
whether two block quotes or one are wanted.)
{: id="20210408153138-ac55u4j"}

Consecutiveness means that if we put these block quotes together,
we get a single block quote:
{: id="20210408153138-hieu2iz"}

````````````````````````````````example
> foo
> bar
.
<blockquote>
<p>foo
bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-tyu5isk"}

To get a block quote with two paragraphs, use:
{: id="20210408153138-40z65si"}

````````````````````````````````example
> foo
>
> bar
.
<blockquote>
<p>foo</p>
<p>bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-pl3evj2"}

Block quotes can interrupt paragraphs:
{: id="20210408153138-s2a9qn4"}

````````````````````````````````example
foo
> bar
.
<p>foo</p>
<blockquote>
<p>bar</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-lzjpfyc"}

In general, blank lines are not needed before or after block
quotes:
{: id="20210408153138-3ensksl"}

````````````````````````````````example
> aaa
***
> bbb
.
<blockquote>
<p>aaa</p>
</blockquote>
<hr />
<blockquote>
<p>bbb</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-drddf35"}

However, because of laziness, a blank line is needed between
a block quote and a following paragraph:
{: id="20210408153138-luxgljk"}

````````````````````````````````example
> bar
baz
.
<blockquote>
<p>bar
baz</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-edd80nb"}

````````````````````````````````example
> bar

baz
.
<blockquote>
<p>bar</p>
</blockquote>
<p>baz</p>
````````````````````````````````
{: id="20210408153138-mfxj631"}

````````````````````````````````example
> bar
>
baz
.
<blockquote>
<p>bar</p>
</blockquote>
<p>baz</p>
````````````````````````````````
{: id="20210408153138-680c35b"}

It is a consequence of the Laziness rule that any number
of initial `>`s may be omitted on a continuation line of a
nested block quote:
{: id="20210408153138-qb35omb"}

````````````````````````````````example
> > > foo
bar
.
<blockquote>
<blockquote>
<blockquote>
<p>foo
bar</p>
</blockquote>
</blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-5e8pg9t"}

````````````````````````````````example
>>> foo
> bar
>>baz
.
<blockquote>
<blockquote>
<blockquote>
<p>foo
bar
baz</p>
</blockquote>
</blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-ssun0dq"}

When including an indented code block in a block quote,
remember that the [block quote marker] includes
both the `>` and a following space.  So *five spaces* are needed after
the `>`:
{: id="20210408153138-iiu6e26"}

````````````````````````````````example
>     code

>    not code
.
<blockquote>
<pre><code>code
</code></pre>
</blockquote>
<blockquote>
<p>not code</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-1cde28r"}

## List items
{: id="20210408153138-g4u1h8i"}

A [list marker](@) is a
[bullet list marker] or an [ordered list marker].
{: id="20210408153138-e0uyyxz"}

A [bullet list marker](@)
is a `-`, `+`, or `*` character.
{: id="20210408153138-lca43z2"}

An [ordered list marker](@)
is a sequence of 1--9 arabic digits (`0-9`), followed by either a
`.` character or a `)` character.  (The reason for the length
limit is that with 10 digits we start seeing integer overflows
in some browsers.)
{: id="20210408153138-p4cf5md"}

The following rules define [list items]:
{: id="20210408153138-ttzudjx"}

1. {: id="20210408153137-ab10tle"}**Basic case.**  If a sequence of lines *Ls* constitute a sequence of
   blocks *Bs* starting with a [non-whitespace character], and *M* is a
   list marker of width *W* followed by 1 ≤ *N* ≤ 4 spaces, then the result
   of prepending *M* and the following spaces to the first line of
   *Ls*, and indenting subsequent lines of *Ls* by *W + N* spaces, is a
   list item with *Bs* as its contents.  The type of the list item
   (bullet or ordered) is determined by the type of its list marker.
   If the list item is ordered, then it is also assigned a start
   number, based on the ordered list marker.
   {: id="20210408153138-lwxjp1b"}

   Exceptions:
   {: id="20210408153138-wfnakup"}

   1. {: id="20210408153137-segzl8h"}When the first list item in a [list] interrupts
      a paragraph---that is, when it starts on a line that would
      otherwise count as [paragraph continuation text]---then (a)
      the lines *Ls* must not begin with a blank line, and (b) if
      the list item is ordered, the start number must be 1.
      {: id="20210408153138-iaap10c"}
   2. {: id="20210408153137-kevtppc"}If any line is a [thematic break][thematic breaks] then
      that line is not a list item.
      {: id="20210408153138-cs963s4"}
   {: id="20210408153138-p6yeuzo"}
{: id="20210408153138-kv356c2"}

For example, let *Ls* be the lines
{: id="20210408153138-egla54m"}

````````````````````````````````example
A paragraph
with two lines.

    indented code

> A block quote.
.
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
````````````````````````````````
{: id="20210408153138-bb2q6sl"}

And let *M* be the marker `1.`, and *N* = 2.  Then rule #1# says
that the following is an ordered list item with start number 1,
and the same contents as *Ls*:
{: id="20210408153138-yn04asc"}

````````````````````````````````example
1.  A paragraph
    with two lines.

        indented code

    > A block quote.
.
<ol>
<li>
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-budtfj0"}

The most important thing to notice is that the position of
the text after the list marker determines how much indentation
is needed in subsequent blocks in the list item.  If the list
marker takes up two spaces, and there are three spaces between
the list marker and the next [non-whitespace character], then blocks
must be indented five spaces in order to fall under the list
item.
{: id="20210408153138-g1k5af2"}

Here are some examples showing how far content must be indented to be
put under the list item:
{: id="20210408153138-kymb7sl"}

````````````````````````````````example
- one

 two
.
<ul>
<li>one</li>
</ul>
<p>two</p>
````````````````````````````````
{: id="20210408153138-dm8b75j"}

````````````````````````````````example
- one

  two
.
<ul>
<li>
<p>one</p>
<p>two</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-5k3iohy"}

````````````````````````````````example
 -    one

     two
.
<ul>
<li>one</li>
</ul>
<pre><code> two
</code></pre>
````````````````````````````````
{: id="20210408153138-q6wdc9x"}

````````````````````````````````example
 -    one

      two
.
<ul>
<li>
<p>one</p>
<p>two</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-5do8spd"}

It is tempting to think of this in terms of columns:  the continuation
blocks must be indented at least to the column of the first
[non-whitespace character] after the list marker. However, that is not quite right.
The spaces after the list marker determine how much relative indentation
is needed.  Which column this indentation reaches will depend on
how the list item is embedded in other constructions, as shown by
this example:
{: id="20210408153138-f6cdbzj"}

````````````````````````````````example
   > > 1.  one
>>
>>     two
.
<blockquote>
<blockquote>
<ol>
<li>
<p>one</p>
<p>two</p>
</li>
</ol>
</blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-ofnohtz"}

Here `two` occurs in the same column as the list marker `1.`,
but is actually contained in the list item, because there is
sufficient indentation after the last containing blockquote marker.
{: id="20210408153138-cg4w1s8"}

The converse is also possible.  In the following example, the word `two`
occurs far to the right of the initial text of the list item, `one`, but
it is not considered part of the list item, because it is not indented
far enough past the blockquote marker:
{: id="20210408153138-5vv4dwe"}

````````````````````````````````example
>>- one
>>
  >  > two
.
<blockquote>
<blockquote>
<ul>
<li>one</li>
</ul>
<p>two</p>
</blockquote>
</blockquote>
````````````````````````````````
{: id="20210408153138-0p37lo6"}

Note that at least one space is needed between the list marker and
any following content, so these are not list items:
{: id="20210408153138-086hj6q"}

````````````````````````````````example
-one

2.two
.
<p>-one</p>
<p>2.two</p>
````````````````````````````````
{: id="20210408153138-mwwcvgm"}

A list item may contain blocks that are separated by more than
one blank line.
{: id="20210408153138-8c79dbe"}

````````````````````````````````example
- foo


  bar
.
<ul>
<li>
<p>foo</p>
<p>bar</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-lgs1dyo"}

A list item may contain any kind of block:
{: id="20210408153138-zc7us8e"}

````````````````````````````````example
1.  foo

    ```
    bar
    ```

    baz

    > bam
.
<ol>
<li>
<p>foo</p>
<pre><code>bar
</code></pre>
<p>baz</p>
<blockquote>
<p>bam</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-nez6fvx"}

A list item that contains an indented code block will preserve
empty lines within the code block verbatim.
{: id="20210408153138-geea9a0"}

````````````````````````````````example
- Foo

      bar


      baz
.
<ul>
<li>
<p>Foo</p>
<pre><code>bar


baz
</code></pre>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-fovmt1m"}

Note that ordered list start numbers must be nine digits or less:
{: id="20210408153138-dk9fwhk"}

````````````````````````````````example
123456789. ok
.
<ol start="123456789">
<li>ok</li>
</ol>
````````````````````````````````
{: id="20210408153138-l1gsbgs"}

````````````````````````````````example
1234567890. not ok
.
<p>1234567890. not ok</p>
````````````````````````````````
{: id="20210408153138-1m332rs"}

A start number may begin with 0s:
{: id="20210408153138-7u4efy7"}

````````````````````````````````example
0. ok
.
<ol start="0">
<li>ok</li>
</ol>
````````````````````````````````
{: id="20210408153138-0e14bsw"}

````````````````````````````````example
003. ok
.
<ol start="3">
<li>ok</li>
</ol>
````````````````````````````````
{: id="20210408153138-fwolxjk"}

A start number may not be negative:
{: id="20210408153138-q68j865"}

````````````````````````````````example
-1. not ok
.
<p>-1. not ok</p>
````````````````````````````````
{: id="20210408153138-jyi49fe"}

2. {: id="20210408153137-xnjwgl9"}**Item starting with indented code.**  If a sequence of lines *Ls*
   constitute a sequence of blocks *Bs* starting with an indented code
   block, and *M* is a list marker of width *W* followed by
   one space, then the result of prepending *M* and the following
   space to the first line of *Ls*, and indenting subsequent lines of
   *Ls* by *W + 1* spaces, is a list item with *Bs* as its contents.
   If a line is empty, then it need not be indented.  The type of the
   list item (bullet or ordered) is determined by the type of its list
   marker.  If the list item is ordered, then it is also assigned a
   start number, based on the ordered list marker.
   {: id="20210408153138-4715xi6"}
{: id="20210408153138-linjsep"}

An indented code block will have to be indented four spaces beyond
the edge of the region where text will be included in the list item.
In the following case that is 6 spaces:
{: id="20210408153138-2qgiv1h"}

````````````````````````````````example
- foo

      bar
.
<ul>
<li>
<p>foo</p>
<pre><code>bar
</code></pre>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-l24gtz0"}

And in this case it is 11 spaces:
{: id="20210408153138-m2xcmpx"}

````````````````````````````````example
  10.  foo

           bar
.
<ol start="10">
<li>
<p>foo</p>
<pre><code>bar
</code></pre>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-m5vs922"}

If the *first* block in the list item is an indented code block,
then by rule #2,# the contents must be indented *one* space after the
list marker:
{: id="20210408153138-k0mfsrw"}

````````````````````````````````example
    indented code

paragraph

    more code
.
<pre><code>indented code
</code></pre>
<p>paragraph</p>
<pre><code>more code
</code></pre>
````````````````````````````````
{: id="20210408153138-elr4kwj"}

````````````````````````````````example
1.     indented code

   paragraph

       more code
.
<ol>
<li>
<pre><code>indented code
</code></pre>
<p>paragraph</p>
<pre><code>more code
</code></pre>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-u3vf3ls"}

Note that an additional space indent is interpreted as space
inside the code block:
{: id="20210408153138-tdypqgl"}

````````````````````````````````example
1.      indented code

   paragraph

       more code
.
<ol>
<li>
<pre><code> indented code
</code></pre>
<p>paragraph</p>
<pre><code>more code
</code></pre>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-sdqf84o"}

Note that rules #1# and #2# only apply to two cases:  (a) cases
in which the lines to be included in a list item begin with a
[non-whitespace character], and (b) cases in which
they begin with an indented code
block.  In a case like the following, where the first block begins with
a three-space indent, the rules do not allow us to form a list item by
indenting the whole thing and prepending a list marker:
{: id="20210408153138-blagqfy"}

````````````````````````````````example
   foo

bar
.
<p>foo</p>
<p>bar</p>
````````````````````````````````
{: id="20210408153138-cnzmq18"}

````````````````````````````````example
-    foo

  bar
.
<ul>
<li>foo</li>
</ul>
<p>bar</p>
````````````````````````````````
{: id="20210408153138-2u5vwzt"}

This is not a significant restriction, because when a block begins
with 1-3 spaces indent, the indentation can always be removed without
a change in interpretation, allowing rule #1# to be applied.  So, in
the above case:
{: id="20210408153138-j409ckv"}

````````````````````````````````example
-  foo

   bar
.
<ul>
<li>
<p>foo</p>
<p>bar</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-kb5vau9"}

3. {: id="20210408153137-4omfqe7"}**Item starting with a blank line.**  If a sequence of lines *Ls*
   starting with a single [blank line] constitute a (possibly empty)
   sequence of blocks *Bs*, not separated from each other by more than
   one blank line, and *M* is a list marker of width *W*,
   then the result of prepending *M* to the first line of *Ls*, and
   indenting subsequent lines of *Ls* by *W + 1* spaces, is a list
   item with *Bs* as its contents.
   If a line is empty, then it need not be indented.  The type of the
   list item (bullet or ordered) is determined by the type of its list
   marker.  If the list item is ordered, then it is also assigned a
   start number, based on the ordered list marker.
   {: id="20210408153138-vtskoao"}
{: id="20210408153138-3ndg0wf"}

Here are some list items that start with a blank line but are not empty:
{: id="20210408153138-zu3w5ry"}

````````````````````````````````example
-
  foo
-
  ```
  bar
  ```
-
      baz
.
<ul>
<li>foo</li>
<li>
<pre><code>bar
</code></pre>
</li>
<li>
<pre><code>baz
</code></pre>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-eqmfbjn"}

When the list item starts with a blank line, the number of spaces
following the list marker doesn't change the required indentation:
{: id="20210408153138-w5611kl"}

````````````````````````````````example
-   
  foo
.
<ul>
<li>foo</li>
</ul>
````````````````````````````````
{: id="20210408153138-gxel22o"}

A list item can begin with at most one blank line.
In the following example, `foo` is not part of the list
item:
{: id="20210408153138-pa5vz28"}

````````````````````````````````example
-

  foo
.
<ul>
<li></li>
</ul>
<p>foo</p>
````````````````````````````````
{: id="20210408153138-67x2wmn"}

Here is an empty bullet list item:
{: id="20210408153138-ocx34lw"}

````````````````````````````````example
- foo
-
- bar
.
<ul>
<li>foo</li>
<li></li>
<li>bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-71pdjg3"}

It does not matter whether there are spaces following the [list marker]:
{: id="20210408153138-b9ynheq"}

````````````````````````````````example
- foo
-   
- bar
.
<ul>
<li>foo</li>
<li></li>
<li>bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-b8roiv2"}

Here is an empty ordered list item:
{: id="20210408153138-8f8elel"}

````````````````````````````````example
1. foo
2.
3. bar
.
<ol>
<li>foo</li>
<li></li>
<li>bar</li>
</ol>
````````````````````````````````
{: id="20210408153138-1antk57"}

A list may start or end with an empty list item:
{: id="20210408153138-k1752hy"}

````````````````````````````````example
*
.
<ul>
<li></li>
</ul>
````````````````````````````````
{: id="20210408153138-65vpozd"}

However, an empty list item cannot interrupt a paragraph:
{: id="20210408153138-m1d667h"}

````````````````````````````````example
foo
*

foo
1.
.
<p>foo
*</p>
<p>foo
1.</p>
````````````````````````````````
{: id="20210408153138-ahpskkd"}

4. {: id="20210408153137-d61ddio"}**Indentation.**  If a sequence of lines *Ls* constitutes a list item
   according to rule #1,# #2,# or #3,# then the result of indenting each line
   of *Ls* by 1-3 spaces (the same for each line) also constitutes a
   list item with the same contents and attributes.  If a line is
   empty, then it need not be indented.
   {: id="20210408153138-ji7emkg"}
{: id="20210408153138-lsex7i5"}

Indented one space:
{: id="20210408153138-i3fb7no"}

````````````````````````````````example
 1.  A paragraph
     with two lines.

         indented code

     > A block quote.
.
<ol>
<li>
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-tlyxqyr"}

Indented two spaces:
{: id="20210408153138-6ove1rc"}

````````````````````````````````example
  1.  A paragraph
      with two lines.

          indented code

      > A block quote.
.
<ol>
<li>
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-ckgeuar"}

Indented three spaces:
{: id="20210408153138-ejgiyeb"}

````````````````````````````````example
   1.  A paragraph
       with two lines.

           indented code

       > A block quote.
.
<ol>
<li>
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-h6myfkg"}

Four spaces indent gives a code block:
{: id="20210408153138-6ga0poo"}

````````````````````````````````example
    1.  A paragraph
        with two lines.

            indented code

        > A block quote.
.
<pre><code>1.  A paragraph
    with two lines.

        indented code

    &gt; A block quote.
</code></pre>
````````````````````````````````
{: id="20210408153138-7e2lzxf"}

5. {: id="20210408153137-4gvs2ra"}**Laziness.**  If a string of lines *Ls* constitute a [list
   item](#list-items) with contents *Bs*, then the result of deleting
   some or all of the indentation from one or more lines in which the
   next [non-whitespace character] after the indentation is
   [paragraph continuation text] is a
   list item with the same contents and attributes.  The unindented
   lines are called
   [lazy continuation line](@)s.
   {: id="20210408153138-k7tehv2"}
{: id="20210408153138-694q79d"}

Here is an example with [lazy continuation lines]:
{: id="20210408153138-z9grpnm"}

````````````````````````````````example
  1.  A paragraph
with two lines.

          indented code

      > A block quote.
.
<ol>
<li>
<p>A paragraph
with two lines.</p>
<pre><code>indented code
</code></pre>
<blockquote>
<p>A block quote.</p>
</blockquote>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-by4lvh3"}

Indentation can be partially deleted:
{: id="20210408153138-vjui905"}

````````````````````````````````example
  1.  A paragraph
    with two lines.
.
<ol>
<li>A paragraph
with two lines.</li>
</ol>
````````````````````````````````
{: id="20210408153138-hdu5vvp"}

These examples show how laziness can work in nested structures:
{: id="20210408153138-4hnyirk"}

````````````````````````````````example
> 1. > Blockquote
continued here.
.
<blockquote>
<ol>
<li>
<blockquote>
<p>Blockquote
continued here.</p>
</blockquote>
</li>
</ol>
</blockquote>
````````````````````````````````
{: id="20210408153138-fxfd83e"}

````````````````````````````````example
> 1. > Blockquote
> continued here.
.
<blockquote>
<ol>
<li>
<blockquote>
<p>Blockquote
continued here.</p>
</blockquote>
</li>
</ol>
</blockquote>
````````````````````````````````
{: id="20210408153138-fagwbfq"}

6. {: id="20210408153137-si6kg7d"}**That's all.** Nothing that is not counted as a list item by rules
   #1--5# counts as a [list item](#list-items).
   {: id="20210408153138-1qc1dxs"}
{: id="20210408153138-n81yr8b"}

The rules for sublists follow from the general rules
[above][List items].  A sublist must be indented the same number
of spaces a paragraph would need to be in order to be included
in the list item.
{: id="20210408153138-j0mzvhm"}

So, in this case we need two spaces indent:
{: id="20210408153138-ntovhc6"}

````````````````````````````````example
- foo
  - bar
    - baz
      - boo
.
<ul>
<li>foo
<ul>
<li>bar
<ul>
<li>baz
<ul>
<li>boo</li>
</ul>
</li>
</ul>
</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-7qzuxvs"}

One is not enough:
{: id="20210408153138-4k2xt08"}

````````````````````````````````example
- foo
 - bar
  - baz
   - boo
.
<ul>
<li>foo</li>
<li>bar</li>
<li>baz</li>
<li>boo</li>
</ul>
````````````````````````````````
{: id="20210408153138-04cr7rd"}

Here we need four, because the list marker is wider:
{: id="20210408153138-6survcc"}

````````````````````````````````example
10) foo
    - bar
.
<ol start="10">
<li>foo
<ul>
<li>bar</li>
</ul>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-hr3eql1"}

Three is not enough:
{: id="20210408153138-49vx6yx"}

````````````````````````````````example
10) foo
   - bar
.
<ol start="10">
<li>foo</li>
</ol>
<ul>
<li>bar</li>
</ul>
````````````````````````````````
{: id="20210408153138-ylbb81t"}

A list may be the first block in a list item:
{: id="20210408153138-t47d7rn"}

````````````````````````````````example
- - foo
.
<ul>
<li>
<ul>
<li>foo</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-q7xj1z3"}

````````````````````````````````example
1. - 2. foo
.
<ol>
<li>
<ul>
<li>
<ol start="2">
<li>foo</li>
</ol>
</li>
</ul>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-er9erq7"}

A list item can contain a heading:
{: id="20210408153138-1tvdqg8"}

````````````````````````````````example
- # Foo
- Bar
  ---
  baz
.
<ul>
<li>
<h1>Foo</h1>
</li>
<li>
<h2>Bar</h2>
baz</li>
</ul>
````````````````````````````````
{: id="20210408153138-vcvryh8"}

### Motivation
{: id="20210408153138-6va76nt"}

John Gruber's Markdown spec says the following about list items:
{: id="20210408153138-g2n4l5s"}

1. {: id="20210408153137-6a7klw7"}"List markers typically start at the left margin, but may be indented
   by up to three spaces. List markers must be followed by one or more
   spaces or a tab."
   {: id="20210408153138-dkayxzv"}
2. {: id="20210408153137-xk30ztn"}"To make lists look nice, you can wrap items with hanging indents....
   But if you don't want to, you don't have to."
   {: id="20210408153138-yryec9b"}
3. {: id="20210408153137-xo9iwrt"}"List items may consist of multiple paragraphs. Each subsequent
   paragraph in a list item must be indented by either 4 spaces or one
   tab."
   {: id="20210408153138-2qaigpy"}
4. {: id="20210408153137-wsfc99u"}"It looks nice if you indent every line of the subsequent paragraphs,
   but here again, Markdown will allow you to be lazy."
   {: id="20210408153138-f3e9bg8"}
5. {: id="20210408153137-c2q7pe7"}"To put a blockquote within a list item, the blockquote's `>`
   delimiters need to be indented."
   {: id="20210408153138-zikid5w"}
6. {: id="20210408153137-ogz02fa"}"To put a code block within a list item, the code block needs to be
   indented twice — 8 spaces or two tabs."
   {: id="20210408153138-aovf8x7"}
{: id="20210408153138-z932n0k"}

These rules specify that a paragraph under a list item must be indented
four spaces (presumably, from the left margin, rather than the start of
the list marker, but this is not said), and that code under a list item
must be indented eight spaces instead of the usual four.  They also say
that a block quote must be indented, but not by how much; however, the
example given has four spaces indentation.  Although nothing is said
about other kinds of block-level content, it is certainly reasonable to
infer that *all* block elements under a list item, including other
lists, must be indented four spaces.  This principle has been called the
*four-space rule*.
{: id="20210408153138-dxgm6dh"}

The four-space rule is clear and principled, and if the reference
implementation `Markdown.pl` had followed it, it probably would have
become the standard.  However, `Markdown.pl` allowed paragraphs and
sublists to start with only two spaces indentation, at least on the
outer level.  Worse, its behavior was inconsistent: a sublist of an
outer-level list needed two spaces indentation, but a sublist of this
sublist needed three spaces.  It is not surprising, then, that different
implementations of Markdown have developed very different rules for
determining what comes under a list item.  (Pandoc and python-Markdown,
for example, stuck with Gruber's syntax description and the four-space
rule, while discount, redcarpet, marked, PHP Markdown, and others
followed `Markdown.pl`'s behavior more closely.)
{: id="20210408153138-e52p962"}

Unfortunately, given the divergences between implementations, there
is no way to give a spec for list items that will be guaranteed not
to break any existing documents.  However, the spec given here should
correctly handle lists formatted with either the four-space rule or
the more forgiving `Markdown.pl` behavior, provided they are laid out
in a way that is natural for a human to read.
{: id="20210408153138-44okdbs"}

The strategy here is to let the width and indentation of the list marker
determine the indentation necessary for blocks to fall under the list
item, rather than having a fixed and arbitrary number.  The writer can
think of the body of the list item as a unit which gets indented to the
right enough to fit the list marker (and any indentation on the list
marker).  (The laziness rule, #5,# then allows continuation lines to be
unindented if needed.)
{: id="20210408153138-f9ihnlp"}

This rule is superior, we claim, to any rule requiring a fixed level of
indentation from the margin.  The four-space rule is clear but
unnatural. It is quite unintuitive that
{: id="20210408153138-hzopxnk"}

```markdown
- foo

  bar

  - baz
```
{: id="20210408153138-xklp9uz"}

should be parsed as two lists with an intervening paragraph,
{: id="20210408153138-ejgnrsz"}

```html
<ul>
<li>foo</li>
</ul>
<p>bar</p>
<ul>
<li>baz</li>
</ul>
```
{: id="20210408153138-yjv31wq"}

as the four-space rule demands, rather than a single list,
{: id="20210408153138-k94qgtr"}

```html
<ul>
<li>
<p>foo</p>
<p>bar</p>
<ul>
<li>baz</li>
</ul>
</li>
</ul>
```
{: id="20210408153138-ost2pai"}

The choice of four spaces is arbitrary.  It can be learned, but it is
not likely to be guessed, and it trips up beginners regularly.
{: id="20210408153138-98hes36"}

Would it help to adopt a two-space rule?  The problem is that such
a rule, together with the rule allowing 1--3 spaces indentation of the
initial list marker, allows text that is indented *less than* the
original list marker to be included in the list item. For example,
`Markdown.pl` parses
{: id="20210408153138-r59qqjt"}

```markdown
   - one

  two
```
{: id="20210408153138-qjyjeg6"}

as a single list item, with `two` a continuation paragraph:
{: id="20210408153138-jd8cf81"}

```html
<ul>
<li>
<p>one</p>
<p>two</p>
</li>
</ul>
```
{: id="20210408153138-p6hqjzh"}

and similarly
{: id="20210408153138-tcuyg2j"}

```markdown
>   - one
>
>  two
```
{: id="20210408153138-1ucllhw"}

as
{: id="20210408153138-8fx7gv9"}

```html
<blockquote>
<ul>
<li>
<p>one</p>
<p>two</p>
</li>
</ul>
</blockquote>
```
{: id="20210408153138-wmwg8ev"}

This is extremely unintuitive.
{: id="20210408153138-16g781k"}

Rather than requiring a fixed indent from the margin, we could require
a fixed indent (say, two spaces, or even one space) from the list marker (which
may itself be indented).  This proposal would remove the last anomaly
discussed.  Unlike the spec presented above, it would count the following
as a list item with a subparagraph, even though the paragraph `bar`
is not indented as far as the first paragraph `foo`:
{: id="20210408153138-grsww9q"}

```markdown
 10. foo

   bar  
```
{: id="20210408153138-w5byukd"}

Arguably this text does read like a list item with `bar` as a subparagraph,
which may count in favor of the proposal.  However, on this proposal indented
code would have to be indented six spaces after the list marker.  And this
would break a lot of existing Markdown, which has the pattern:
{: id="20210408153138-dd5javc"}

```markdown
1.  foo

        indented code
```
{: id="20210408153138-b8dorof"}

where the code is indented eight spaces.  The spec above, by contrast, will
parse this text as expected, since the code block's indentation is measured
from the beginning of `foo`.
{: id="20210408153138-97sw8bm"}

The one case that needs special treatment is a list item that *starts*
with indented code.  How much indentation is required in that case, since
we don't have a "first paragraph" to measure from?  Rule #2# simply stipulates
that in such cases, we require one space indentation from the list marker
(and then the normal four spaces for the indented code).  This will match the
four-space rule in cases where the list marker plus its initial indentation
takes four spaces (a common case), but diverge in other cases.
{: id="20210408153138-dbyzdt6"}

## Lists
{: id="20210408153138-oysjwp3"}

A [list](@) is a sequence of one or more
list items [of the same type].  The list items
may be separated by any number of blank lines.
{: id="20210408153138-nt5bwu3"}

Two list items are [of the same type](@)
if they begin with a [list marker] of the same type.
Two list markers are of the
same type if (a) they are bullet list markers using the same character
(`-`, `+`, or `*`) or (b) they are ordered list numbers with the same
delimiter (either `.` or `)`).
{: id="20210408153138-dfd34jj"}

A list is an [ordered list](@)
if its constituent list items begin with
[ordered list markers], and a
[bullet list](@) if its constituent list
items begin with [bullet list markers].
{: id="20210408153138-skm2qnw"}

The [start number](@)
of an [ordered list] is determined by the list number of
its initial list item.  The numbers of subsequent list items are
disregarded.
{: id="20210408153138-an301it"}

A list is [loose](@) if any of its constituent
list items are separated by blank lines, or if any of its constituent
list items directly contain two block-level elements with a blank line
between them.  Otherwise a list is [tight](@).
(The difference in HTML output is that paragraphs in a loose list are
wrapped in `<p>` tags, while paragraphs in a tight list are not.)
{: id="20210408153138-c4cd01q"}

Changing the bullet or ordered list delimiter starts a new list:
{: id="20210408153138-m1pv0vm"}

````````````````````````````````example
- foo
- bar
+ baz
.
<ul>
<li>foo</li>
<li>bar</li>
</ul>
<ul>
<li>baz</li>
</ul>
````````````````````````````````
{: id="20210408153138-he5qjwn"}

````````````````````````````````example
1. foo
2. bar
3) baz
.
<ol>
<li>foo</li>
<li>bar</li>
</ol>
<ol start="3">
<li>baz</li>
</ol>
````````````````````````````````
{: id="20210408153138-tv80lba"}

In CommonMark, a list can interrupt a paragraph. That is,
no blank line is needed to separate a paragraph from a following
list:
{: id="20210408153138-7kymxrx"}

````````````````````````````````example
Foo
- bar
- baz
.
<p>Foo</p>
<ul>
<li>bar</li>
<li>baz</li>
</ul>
````````````````````````````````
{: id="20210408153138-dxmbjue"}

`Markdown.pl` does not allow this, through fear of triggering a list
via a numeral in a hard-wrapped line:
{: id="20210408153138-e2b35wd"}

```markdown
The number of windows in my house is
14.  The number of doors is 6.
```
{: id="20210408153138-3n5q57e"}

Oddly, though, `Markdown.pl` *does* allow a blockquote to
interrupt a paragraph, even though the same considerations might
apply.
{: id="20210408153138-a889aqs"}

In CommonMark, we do allow lists to interrupt paragraphs, for
two reasons.  First, it is natural and not uncommon for people
to start lists without blank lines:
{: id="20210408153138-00xutiu"}

```markdown
I need to buy
- new shoes
- a coat
- a plane ticket
```
{: id="20210408153138-701xlsf"}

Second, we are attracted to a
{: id="20210408153138-u8t8qc7"}

> [principle of uniformity](@):
> if a chunk of text has a certain
> meaning, it will continue to have the same meaning when put into a
> container block (such as a list item or blockquote).
> {: id="20210408153138-voezf29"}
{: id="20210408153138-mnly9d9"}

(Indeed, the spec for [list items] and [block quotes] presupposes
this principle.) This principle implies that if
{: id="20210408153138-00ug5sb"}

```markdown
  * I need to buy
    - new shoes
    - a coat
    - a plane ticket
```
{: id="20210408153138-2nbsis4"}

is a list item containing a paragraph followed by a nested sublist,
as all Markdown implementations agree it is (though the paragraph
may be rendered without `<p>` tags, since the list is "tight"),
then
{: id="20210408153138-6db8yv2"}

```markdown
I need to buy
- new shoes
- a coat
- a plane ticket
```
{: id="20210408153138-vdd3nf8"}

by itself should be a paragraph followed by a nested sublist.
{: id="20210408153138-wxmy9bi"}

Since it is well established Markdown practice to allow lists to
interrupt paragraphs inside list items, the [principle of
uniformity] requires us to allow this outside list items as
well.  ([reStructuredText](http://docutils.sourceforge.net/rst.html)
takes a different approach, requiring blank lines before lists
even inside other list items.)
{: id="20210408153138-49u0wvz"}

In order to solve of unwanted lists in paragraphs with
hard-wrapped numerals, we allow only lists starting with `1` to
interrupt paragraphs.  Thus,
{: id="20210408153138-rauettr"}

````````````````````````````````example
The number of windows in my house is
14.  The number of doors is 6.
.
<p>The number of windows in my house is
14.  The number of doors is 6.</p>
````````````````````````````````
{: id="20210408153138-r4ezk41"}

We may still get an unintended result in cases like
{: id="20210408153138-03u5mae"}

````````````````````````````````example
The number of windows in my house is
1.  The number of doors is 6.
.
<p>The number of windows in my house is</p>
<ol>
<li>The number of doors is 6.</li>
</ol>
````````````````````````````````
{: id="20210408153138-ch2k87z"}

but this rule should prevent most spurious list captures.
{: id="20210408153138-ns8ydj0"}

There can be any number of blank lines between items:
{: id="20210408153138-iqyupxi"}

````````````````````````````````example
- foo

- bar


- baz
.
<ul>
<li>
<p>foo</p>
</li>
<li>
<p>bar</p>
</li>
<li>
<p>baz</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-xubm42f"}

````````````````````````````````example
- foo
  - bar
    - baz


      bim
.
<ul>
<li>foo
<ul>
<li>bar
<ul>
<li>
<p>baz</p>
<p>bim</p>
</li>
</ul>
</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-l22s3c7"}

To separate consecutive lists of the same type, or to separate a
list from an indented code block that would otherwise be parsed
as a subparagraph of the final list item, you can insert a blank HTML
comment:
{: id="20210408153138-ufhtnvt"}

````````````````````````````````example
- foo
- bar

<!-- -->

- baz
- bim
.
<ul>
<li>foo</li>
<li>bar</li>
</ul>
<!-- -->
<ul>
<li>baz</li>
<li>bim</li>
</ul>
````````````````````````````````
{: id="20210408153138-1ppse6d"}

````````````````````````````````example
-   foo

    notcode

-   foo

<!-- -->

    code
.
<ul>
<li>
<p>foo</p>
<p>notcode</p>
</li>
<li>
<p>foo</p>
</li>
</ul>
<!-- -->
<pre><code>code
</code></pre>
````````````````````````````````
{: id="20210408153138-zzn56pr"}

List items need not be indented to the same level.  The following
list items will be treated as items at the same list level,
since none is indented enough to belong to the previous list
item:
{: id="20210408153138-ixun8vb"}

````````````````````````````````example
- a
 - b
  - c
   - d
  - e
 - f
- g
.
<ul>
<li>a</li>
<li>b</li>
<li>c</li>
<li>d</li>
<li>e</li>
<li>f</li>
<li>g</li>
</ul>
````````````````````````````````
{: id="20210408153138-bevzqkd"}

````````````````````````````````example
1. a

  2. b

   3. c
.
<ol>
<li>
<p>a</p>
</li>
<li>
<p>b</p>
</li>
<li>
<p>c</p>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-ph0m566"}

Note, however, that list items may not be indented more than
three spaces.  Here `- e` is treated as a paragraph continuation
line, because it is indented more than three spaces:
{: id="20210408153138-3k0ms2k"}

````````````````````````````````example
- a
 - b
  - c
   - d
    - e
.
<ul>
<li>a</li>
<li>b</li>
<li>c</li>
<li>d
- e</li>
</ul>
````````````````````````````````
{: id="20210408153138-susde5u"}

And here, `3. c` is treated as in indented code block,
because it is indented four spaces and preceded by a
blank line.
{: id="20210408153138-m3xt39u"}

````````````````````````````````example
1. a

  2. b

    3. c
.
<ol>
<li>
<p>a</p>
</li>
<li>
<p>b</p>
</li>
</ol>
<pre><code>3. c
</code></pre>
````````````````````````````````
{: id="20210408153138-a10j90u"}

This is a loose list, because there is a blank line between
two of the list items:
{: id="20210408153138-bk4ijqb"}

````````````````````````````````example
- a
- b

- c
.
<ul>
<li>
<p>a</p>
</li>
<li>
<p>b</p>
</li>
<li>
<p>c</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-lot5rcp"}

So is this, with a empty second item:
{: id="20210408153138-n5l04jy"}

````````````````````````````````example
* a
*

* c
.
<ul>
<li>
<p>a</p>
</li>
<li></li>
<li>
<p>c</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-zvkiiwy"}

These are loose lists, even though there is no space between the items,
because one of the items directly contains two block-level elements
with a blank line between them:
{: id="20210408153138-ry6yy8o"}

````````````````````````````````example
- a
- b

  c
- d
.
<ul>
<li>
<p>a</p>
</li>
<li>
<p>b</p>
<p>c</p>
</li>
<li>
<p>d</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-utqcrbe"}

````````````````````````````````example
- a
- b

  [ref]: /url
- d
.
<ul>
<li>
<p>a</p>
</li>
<li>
<p>b</p>
</li>
<li>
<p>d</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-cikrhmd"}

This is a tight list, because the blank lines are in a code block:
{: id="20210408153138-06sq3xw"}

````````````````````````````````example
- a
- ```
  b


  ```
- c
.
<ul>
<li>a</li>
<li>
<pre><code>b


</code></pre>
</li>
<li>c</li>
</ul>
````````````````````````````````
{: id="20210408153138-45aykkg"}

This is a tight list, because the blank line is between two
paragraphs of a sublist.  So the sublist is loose while
the outer list is tight:
{: id="20210408153138-b6kmuc1"}

````````````````````````````````example
- a
  - b

    c
- d
.
<ul>
<li>a
<ul>
<li>
<p>b</p>
<p>c</p>
</li>
</ul>
</li>
<li>d</li>
</ul>
````````````````````````````````
{: id="20210408153138-tboqbn9"}

This is a tight list, because the blank line is inside the
block quote:
{: id="20210408153138-g9enbcq"}

````````````````````````````````example
* a
  > b
  >
* c
.
<ul>
<li>a
<blockquote>
<p>b</p>
</blockquote>
</li>
<li>c</li>
</ul>
````````````````````````````````
{: id="20210408153138-bdhxshp"}

This list is tight, because the consecutive block elements
are not separated by blank lines:
{: id="20210408153138-r7nea5p"}

````````````````````````````````example
- a
  > b
  ```
  c
  ```
- d
.
<ul>
<li>a
<blockquote>
<p>b</p>
</blockquote>
<pre><code>c
</code></pre>
</li>
<li>d</li>
</ul>
````````````````````````````````
{: id="20210408153138-bdle4ni"}

A single-paragraph list is tight:
{: id="20210408153138-40b3kiz"}

````````````````````````````````example
- a
.
<ul>
<li>a</li>
</ul>
````````````````````````````````
{: id="20210408153138-dcf459z"}

````````````````````````````````example
- a
  - b
.
<ul>
<li>a
<ul>
<li>b</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-rjpvcf5"}

This list is loose, because of the blank line between the
two block elements in the list item:
{: id="20210408153138-16qnvj3"}

````````````````````````````````example
1. ```
   foo
   ```

   bar
.
<ol>
<li>
<pre><code>foo
</code></pre>
<p>bar</p>
</li>
</ol>
````````````````````````````````
{: id="20210408153138-4smailt"}

Here the outer list is loose, the inner list tight:
{: id="20210408153138-ykyiuhf"}

````````````````````````````````example
* foo
  * bar

  baz
.
<ul>
<li>
<p>foo</p>
<ul>
<li>bar</li>
</ul>
<p>baz</p>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-0i7p91h"}

````````````````````````````````example
- a
  - b
  - c

- d
  - e
  - f
.
<ul>
<li>
<p>a</p>
<ul>
<li>b</li>
<li>c</li>
</ul>
</li>
<li>
<p>d</p>
<ul>
<li>e</li>
<li>f</li>
</ul>
</li>
</ul>
````````````````````````````````
{: id="20210408153138-5kidbss"}

# Inlines
{: id="20210408153138-cit5svo"}

Inlines are parsed sequentially from the beginning of the character
stream to the end (left to right, in left-to-right languages).
Thus, for example, in
{: id="20210408153138-zajjyb9"}

````````````````````````````````example
`hi`lo`
.
<p><code>hi</code>lo`</p>
````````````````````````````````
{: id="20210408153138-omz0utz"}

`hi` is parsed as code, leaving the backtick at the end as a literal
backtick.
{: id="20210408153138-dm4d7bh"}

## Backslash escapes
{: id="20210408153138-gnrnsiu"}

Any ASCII punctuation character may be backslash-escaped:
{: id="20210408153138-8m0mftj"}

````````````````````````````````example
\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\`\{\|\}\~
.
<p>!&quot;#$%&amp;'()*+,-./:;&lt;=&gt;?@[\]^_`{|}~</p>
````````````````````````````````
{: id="20210408153138-9bx27cj"}

Backslashes before other characters are treated as literal
backslashes:
{: id="20210408153138-6fuzrvg"}

````````````````````````````````example
\→\A\a\ \3\φ\«
.
<p>\→\A\a\ \3\φ\«</p>
````````````````````````````````
{: id="20210408153138-yit3yln"}

Escaped characters are treated as regular characters and do
not have their usual Markdown meanings:
{: id="20210408153138-ptksgee"}

````````````````````````````````example
\*not emphasized*
\<br/> not a tag
\[not a link](/foo)
\`not code`
1\. not a list
\* not a list
\# not a heading
\[foo]: /url "not a reference"
\&ouml; not a character entity
.
<p>*not emphasized*
&lt;br/&gt; not a tag
[not a link](/foo)
`not code`
1. not a list
* not a list
# not a heading
[foo]: /url &quot;not a reference&quot;
&amp;ouml; not a character entity</p>
````````````````````````````````
{: id="20210408153138-w3h70re"}

If a backslash is itself escaped, the following character is not:
{: id="20210408153138-88bar4s"}

````````````````````````````````example
\\*emphasis*
.
<p>\<em>emphasis</em></p>
````````````````````````````````
{: id="20210408153138-6xmnspc"}

A backslash at the end of the line is a [hard line break]:
{: id="20210408153138-u0rxgwr"}

````````````````````````````````example
foo\
bar
.
<p>foo<br />
bar</p>
````````````````````````````````
{: id="20210408153138-t0b0twb"}

Backslash escapes do not work in code blocks, code spans, autolinks, or
raw HTML:
{: id="20210408153138-94mzp2q"}

````````````````````````````````example
`` \[\` ``
.
<p><code>\[\`</code></p>
````````````````````````````````
{: id="20210408153138-c94tntu"}

````````````````````````````````example
    \[\]
.
<pre><code>\[\]
</code></pre>
````````````````````````````````
{: id="20210408153138-hy9qat4"}

````````````````````````````````example
~~~
\[\]
~~~
.
<pre><code>\[\]
</code></pre>
````````````````````````````````
{: id="20210408153138-rjiu9yb"}

````````````````````````````````example
<http://example.com?find=\*>
.
<p><a href="http://example.com?find=%5C*">http://example.com?find=\*</a></p>
````````````````````````````````
{: id="20210408153138-hxizbzg"}

````````````````````````````````example
<a href="/bar\/)">
.
<a href="/bar\/)">
````````````````````````````````
{: id="20210408153138-w4phvnd"}

But they work in all other contexts, including URLs and link titles,
link references, and [info strings] in [fenced code blocks]:
{: id="20210408153138-d3s8loe"}

````````````````````````````````example
[foo](/bar\* "ti\*tle")
.
<p><a href="/bar*" title="ti*tle">foo</a></p>
````````````````````````````````
{: id="20210408153138-8j6cr5i"}

````````````````````````````````example
[foo]

[foo]: /bar\* "ti\*tle"
.
<p><a href="/bar*" title="ti*tle">foo</a></p>
````````````````````````````````
{: id="20210408153138-kz3zcg6"}

````````````````````````````````example
``` foo\+bar
foo
```
.
<pre><code class="language-foo+bar">foo
</code></pre>
````````````````````````````````
{: id="20210408153138-weqr6kt"}

## Entity and numeric character references
{: id="20210408153138-4csvduw"}

Valid HTML entity references and numeric character references
can be used in place of the corresponding Unicode character,
with the following exceptions:
{: id="20210408153138-vuex70g"}

- {: id="20210408153137-3s0h86w"}Entity and character references are not recognized in code
  blocks and code spans.
  {: id="20210408153138-nt6051b"}
- {: id="20210408153137-3d8nsy9"}Entity and character references cannot stand in place of
  special characters that define structural elements in
  CommonMark.  For example, although `&#42;` can be used
  in place of a literal `*` character, `&#42;` cannot replace
  `*` in emphasis delimiters, bullet list markers, or thematic
  breaks.
  {: id="20210408153138-pezfri0"}
{: id="20210408153138-wvlooza"}

Conforming CommonMark parsers need not store information about
whether a particular character was represented in the source
using a Unicode character or an entity reference.
{: id="20210408153138-0q8c7q6"}

[Entity references](@) consist of `&` + any of the valid
HTML5 entity names + `;`. The
document [https://html.spec.whatwg.org/multipage/entities.json](https://html.spec.whatwg.org/multipage/entities.json)
is used as an authoritative source for the valid entity
references and their corresponding code points.
{: id="20210408153138-x1lo42s"}

````````````````````````````````example
&nbsp; &amp; &copy; &AElig; &Dcaron;
&frac34; &HilbertSpace; &DifferentialD;
&ClockwiseContourIntegral; &ngE;
.
<p>  &amp; © Æ Ď
¾ ℋ ⅆ
∲ ≧̸</p>
````````````````````````````````
{: id="20210408153138-f5a7ojc"}

[Decimal numeric character
references](@)
consist of `&#` + a string of 1--7 arabic digits + `;`. A
numeric character reference is parsed as the corresponding
Unicode character. Invalid Unicode code points will be replaced by
the REPLACEMENT CHARACTER (`U+FFFD`).  For security reasons,
the code point `U+0000` will also be replaced by `U+FFFD`.
{: id="20210408153138-3g19hgk"}

````````````````````````````````example
&#35; &#1234; &#992; &#0;
.
<p># Ӓ Ϡ �</p>
````````````````````````````````
{: id="20210408153138-f67jepu"}

[Hexadecimal numeric character
references](@) consist of `&#` +
either `X` or `x` + a string of 1-6 hexadecimal digits + `;`.
They too are parsed as the corresponding Unicode character (this
time specified with a hexadecimal numeral instead of decimal).
{: id="20210408153138-rcljizk"}

````````````````````````````````example
&#X22; &#XD06; &#xcab;
.
<p>&quot; ആ ಫ</p>
````````````````````````````````
{: id="20210408153138-pnp00ie"}

Here are some nonentities:
{: id="20210408153138-jw486ta"}

````````````````````````````````example
&nbsp &x; &#; &#x;
&#987654321;
&#abcdef0;
&ThisIsNotDefined; &hi?;
.
<p>&amp;nbsp &amp;x; &amp;#; &amp;#x;
&amp;#987654321;
&amp;#abcdef0;
&amp;ThisIsNotDefined; &amp;hi?;</p>
````````````````````````````````
{: id="20210408153138-hohgqt7"}

Although HTML5 does accept some entity references
without a trailing semicolon (such as `&copy`), these are not
recognized here, because it makes the grammar too ambiguous:
{: id="20210408153138-gnzy7le"}

````````````````````````````````example
&copy
.
<p>&amp;copy</p>
````````````````````````````````
{: id="20210408153138-jxv2gk2"}

Strings that are not on the list of HTML5 named entities are not
recognized as entity references either:
{: id="20210408153138-pl7j453"}

````````````````````````````````example
&MadeUpEntity;
.
<p>&amp;MadeUpEntity;</p>
````````````````````````````````
{: id="20210408153138-uktefd1"}

Entity and numeric character references are recognized in any
context besides code spans or code blocks, including
URLs, [link titles], and [fenced code block][] [info strings]:
{: id="20210408153138-0gk8r4h"}

````````````````````````````````example
<a href="&ouml;&ouml;.html">
.
<a href="&ouml;&ouml;.html">
````````````````````````````````
{: id="20210408153138-1rf343j"}

````````````````````````````````example
[foo](/f&ouml;&ouml; "f&ouml;&ouml;")
.
<p><a href="/f%C3%B6%C3%B6" title="föö">foo</a></p>
````````````````````````````````
{: id="20210408153138-920od9j"}

````````````````````````````````example
[foo]

[foo]: /f&ouml;&ouml; "f&ouml;&ouml;"
.
<p><a href="/f%C3%B6%C3%B6" title="föö">foo</a></p>
````````````````````````````````
{: id="20210408153138-qjej9db"}

````````````````````````````````example
``` f&ouml;&ouml;
foo
```
.
<pre><code class="language-föö">foo
</code></pre>
````````````````````````````````
{: id="20210408153138-37ky6bh"}

Entity and numeric character references are treated as literal
text in code spans and code blocks:
{: id="20210408153138-95w9kxl"}

````````````````````````````````example
`f&ouml;&ouml;`
.
<p><code>f&amp;ouml;&amp;ouml;</code></p>
````````````````````````````````
{: id="20210408153138-4ug5p4w"}

````````````````````````````````example
    f&ouml;f&ouml;
.
<pre><code>f&amp;ouml;f&amp;ouml;
</code></pre>
````````````````````````````````
{: id="20210408153138-6cvjyto"}

Entity and numeric character references cannot be used
in place of symbols indicating structure in CommonMark
documents.
{: id="20210408153138-uhhbgsx"}

````````````````````````````````example
&#42;foo&#42;
*foo*
.
<p>*foo*
<em>foo</em></p>
````````````````````````````````
{: id="20210408153138-wd1qlnz"}

````````````````````````````````example
&#42; foo

* foo
.
<p>* foo</p>
<ul>
<li>foo</li>
</ul>
````````````````````````````````
{: id="20210408153138-1z5i7ww"}

````````````````````````````````example
foo&#10;&#10;bar
.
<p>foo

bar</p>
````````````````````````````````
{: id="20210408153138-g4glpal"}

````````````````````````````````example
&#9;foo
.
<p>→foo</p>
````````````````````````````````
{: id="20210408153138-sr96qz0"}

````````````````````````````````example
[a](url &quot;tit&quot;)
.
<p>[a](url &quot;tit&quot;)</p>
````````````````````````````````
{: id="20210408153138-atmmr17"}

## Code spans
{: id="20210408153138-ccxgb0g"}

A [backtick string](@)
is a string of one or more backtick characters (`` ` ``) that is neither
preceded nor followed by a backtick.
{: id="20210408153138-jne9uvn"}

A [code span](@) begins with a backtick string and ends with
a backtick string of equal length.  The contents of the code span are
the characters between the two backtick strings, normalized in the
following ways:
{: id="20210408153138-mh1d8w3"}

- {: id="20210408153137-ynqs5xc"}First, [line endings] are converted to [spaces].
  {: id="20210408153138-72r7bir"}
- {: id="20210408153137-hfvjlzi"}If the resulting string both begins *and* ends with a [space]
  character, but does not consist entirely of [space]
  characters, a single [space] character is removed from the
  front and back.  This allows you to include code that begins
  or ends with backtick characters, which must be separated by
  whitespace from the opening or closing backtick strings.
  {: id="20210408153138-7ycfdwe"}
{: id="20210408153138-jwb0bny"}

This is a simple code span:
{: id="20210408153138-m46h7l4"}

````````````````````````````````example
`foo`
.
<p><code>foo</code></p>
````````````````````````````````
{: id="20210408153138-zl1rvla"}

Here two backticks are used, because the code contains a backtick.
This example also illustrates stripping of a single leading and
trailing space:
{: id="20210408153138-vpcg3db"}

````````````````````````````````example
`` foo ` bar ``
.
<p><code>foo ` bar</code></p>
````````````````````````````````
{: id="20210408153138-n1vetow"}

This example shows the motivation for stripping leading and trailing
spaces:
{: id="20210408153138-xhvm9if"}

````````````````````````````````example
` `` `
.
<p><code>``</code></p>
````````````````````````````````
{: id="20210408153138-124dmj0"}

Note that only *one* space is stripped:
{: id="20210408153138-wt5dste"}

````````````````````````````````example
`  ``  `
.
<p><code> `` </code></p>
````````````````````````````````
{: id="20210408153138-7y9g89o"}

The stripping only happens if the space is on both
sides of the string:
{: id="20210408153138-dlkvaeo"}

````````````````````````````````example
` a`
.
<p><code> a</code></p>
````````````````````````````````
{: id="20210408153138-kkbjnhx"}

Only [spaces], and not [unicode whitespace] in general, are
stripped in this way:
{: id="20210408153138-6jrl68v"}

````````````````````````````````example
` b `
.
<p><code> b </code></p>
````````````````````````````````
{: id="20210408153138-68nwx6k"}

No stripping occurs if the code span contains only spaces:
{: id="20210408153138-ah7iwjh"}

````````````````````````````````example
` `
`  `
.
<p><code> </code>
<code>  </code></p>
````````````````````````````````
{: id="20210408153138-6knvmb8"}

[Line endings] are treated like spaces:
{: id="20210408153138-asis6vq"}

````````````````````````````````example
``
foo
bar  
baz
``
.
<p><code>foo bar   baz</code></p>
````````````````````````````````
{: id="20210408153138-78tbfwk"}

````````````````````````````````example
``
foo 
``
.
<p><code>foo </code></p>
````````````````````````````````
{: id="20210408153138-l6mq5n4"}

Interior spaces are not collapsed:
{: id="20210408153138-awjzu3t"}

````````````````````````````````example
`foo   bar 
baz`
.
<p><code>foo   bar  baz</code></p>
````````````````````````````````
{: id="20210408153138-amwkrky"}

Note that browsers will typically collapse consecutive spaces
when rendering `<code>` elements, so it is recommended that
the following CSS be used:
{: id="20210408153138-x95ykzf"}

```
code{white-space: pre-wrap;}
```
{: id="20210408153138-ts4fnra"}

Note that backslash escapes do not work in code spans. All backslashes
are treated literally:
{: id="20210408153138-eqlviz7"}

````````````````````````````````example
`foo\`bar`
.
<p><code>foo\</code>bar`</p>
````````````````````````````````
{: id="20210408153138-4edqw82"}

Backslash escapes are never needed, because one can always choose a
string of *n* backtick characters as delimiters, where the code does
not contain any strings of exactly *n* backtick characters.
{: id="20210408153138-rrlwou6"}

````````````````````````````````example
``foo`bar``
.
<p><code>foo`bar</code></p>
````````````````````````````````
{: id="20210408153138-t1uufm2"}

````````````````````````````````example
` foo `` bar `
.
<p><code>foo `` bar</code></p>
````````````````````````````````
{: id="20210408153138-ujln4zp"}

Code span backticks have higher precedence than any other inline
constructs except HTML tags and autolinks.  Thus, for example, this is
not parsed as emphasized text, since the second `*` is part of a code
span:
{: id="20210408153138-v3vjoc0"}

````````````````````````````````example
*foo`*`
.
<p>*foo<code>*</code></p>
````````````````````````````````
{: id="20210408153138-p1usozq"}

And this is not parsed as a link:
{: id="20210408153138-urkf62b"}

````````````````````````````````example
[not a `link](/foo`)
.
<p>[not a <code>link](/foo</code>)</p>
````````````````````````````````
{: id="20210408153138-ssb2hce"}

Code spans, HTML tags, and autolinks have the same precedence.
Thus, this is code:
{: id="20210408153138-4wd3y66"}

````````````````````````````````example
`<a href="`">`
.
<p><code>&lt;a href=&quot;</code>&quot;&gt;`</p>
````````````````````````````````
{: id="20210408153138-txqrdn8"}

But this is an HTML tag:
{: id="20210408153138-u3ax5wr"}

````````````````````````````````example
<a href="`">`
.
<p><a href="`">`</p>
````````````````````````````````
{: id="20210408153138-nvj1vmq"}

And this is code:
{: id="20210408153138-6klryt1"}

````````````````````````````````example
`<http://foo.bar.`baz>`
.
<p><code>&lt;http://foo.bar.</code>baz&gt;`</p>
````````````````````````````````
{: id="20210408153138-4yznqs6"}

But this is an autolink:
{: id="20210408153138-ba9we8p"}

````````````````````````````````example
<http://foo.bar.`baz>`
.
<p><a href="http://foo.bar.%60baz">http://foo.bar.`baz</a>`</p>
````````````````````````````````
{: id="20210408153138-5qyat1y"}

When a backtick string is not closed by a matching backtick string,
we just have literal backticks:
{: id="20210408153138-ppc1uwp"}

````````````````````````````````example
```foo``
.
<p>```foo``</p>
````````````````````````````````
{: id="20210408153138-ftmp344"}

````````````````````````````````example
`foo
.
<p>`foo</p>
````````````````````````````````
{: id="20210408153138-ulvrf8y"}

The following case also illustrates the need for opening and
closing backtick strings to be equal in length:
{: id="20210408153138-amczxbk"}

````````````````````````````````example
`foo``bar``
.
<p>`foo<code>bar</code></p>
````````````````````````````````
{: id="20210408153138-cmu2vcj"}

## Emphasis and strong emphasis
{: id="20210408153138-57wsfla"}

John Gruber's original [Markdown syntax
description](http://daringfireball.net/projects/markdown/syntax#em) says:
{: id="20210408153138-xp675mt"}

> Markdown treats asterisks (`*`) and underscores (`_`) as indicators of
> emphasis. Text wrapped with one `*` or `_` will be wrapped with an HTML
> `<em>` tag; double `*`'s or `_`'s will be wrapped with an HTML `<strong>`
> tag.
> {: id="20210408153138-k0qcxju"}
{: id="20210408153138-u7w8w83"}

This is enough for most users, but these rules leave much undecided,
especially when it comes to nested emphasis.  The original
`Markdown.pl` test suite makes it clear that triple `***` and
`___` delimiters can be used for strong emphasis, and most
implementations have also allowed the following patterns:
{: id="20210408153138-6a2hnhs"}

```markdown
***strong emph***
***strong** in emph*
***emph* in strong**
**in strong *emph***
*in emph **strong***
```
{: id="20210408153138-h2oa7s0"}

The following patterns are less widely supported, but the intent
is clear and they are useful (especially in contexts like bibliography
entries):
{: id="20210408153138-okbl8ag"}

```markdown
*emph *with emph* in it*
**strong **with strong** in it**
```
{: id="20210408153138-rwl23sl"}

Many implementations have also restricted intraword emphasis to
the `*` forms, to avoid unwanted emphasis in words containing
internal underscores.  (It is best practice to put these in code
spans, but users often do not.)
{: id="20210408153138-1wm8xxh"}

```markdown
internal emphasis: foo*bar*baz
no emphasis: foo_bar_baz
```
{: id="20210408153138-jbyn36s"}

The rules given below capture all of these patterns, while allowing
for efficient parsing strategies that do not backtrack.
{: id="20210408153138-gefhd27"}

First, some definitions.  A [delimiter run](@) is either
a sequence of one or more `*` characters that is not preceded or
followed by a non-backslash-escaped `*` character, or a sequence
of one or more `_` characters that is not preceded or followed by
a non-backslash-escaped `_` character.
{: id="20210408153138-ma78cks"}

A [left-flanking delimiter run](@) is
a [delimiter run] that is (1) not followed by [Unicode whitespace],
and either (2a) not followed by a [punctuation character], or
(2b) followed by a [punctuation character] and
preceded by [Unicode whitespace] or a [punctuation character].
For purposes of this definition, the beginning and the end of
the line count as Unicode whitespace.
{: id="20210408153138-xkgjn90"}

A [right-flanking delimiter run](@) is
a [delimiter run] that is (1) not preceded by [Unicode whitespace],
and either (2a) not preceded by a [punctuation character], or
(2b) preceded by a [punctuation character] and
followed by [Unicode whitespace] or a [punctuation character].
For purposes of this definition, the beginning and the end of
the line count as Unicode whitespace.
{: id="20210408153138-btxpdxf"}

Here are some examples of delimiter runs.
{: id="20210408153138-1ljq2qu"}

- {: id="20210408153137-trojyfa"}left-flanking but not right-flanking:
  {: id="20210408153138-5hc5x73"}

  ```
  ***abc
    _abc
  **"abc"
   _"abc"
  ```
  {: id="20210408153138-axr33pu"}
- {: id="20210408153137-jva6zx4"}right-flanking but not left-flanking:
  {: id="20210408153138-ztaaxvh"}
  ```
   abc***
   abc_
  "abc"**
  "abc"_
  ```
  {: id="20210408153138-yd2g11z"}
- {: id="20210408153137-c5f7wfn"}Both left and right-flanking:
  {: id="20210408153138-i5rlzc6"}
  ```
   abc***def
  "abc"_"def"
  ```
  {: id="20210408153138-gb7ydf0"}
- {: id="20210408153137-tnqc6pu"}Neither left nor right-flanking:
  {: id="20210408153138-9tqvxxb"}
  ```
  abc *** def
  a _ b
  ```
  {: id="20210408153138-p60ww6d"}
{: id="20210408153138-uiqdlhy"}

(The idea of distinguishing left-flanking and right-flanking
delimiter runs based on the character before and the character
after comes from Roopesh Chander's
[vfmd](http://www.vfmd.org/vfmd-spec/specification/#procedure-for-identifying-emphasis-tags).
vfmd uses the terminology "emphasis indicator string" instead of "delimiter
run," and its rules for distinguishing left- and right-flanking runs
are a bit more complex than the ones given here.)
{: id="20210408153138-5qi9tts"}

The following rules define emphasis and strong emphasis:
{: id="20210408153138-emlv7cn"}

1. {: id="20210408153137-bzzzod4"}A single `*` character [can open emphasis](@)
   iff (if and only if) it is part of a [left-flanking delimiter run].
   {: id="20210408153138-x65uvwm"}
2. {: id="20210408153137-0sv0ftz"}A single `_` character [can open emphasis] iff
   it is part of a [left-flanking delimiter run]
   and either (a) not part of a [right-flanking delimiter run]
   or (b) part of a [right-flanking delimiter run]
   preceded by punctuation.
   {: id="20210408153138-476fswm"}
3. {: id="20210408153137-fcze3h5"}A single `*` character [can close emphasis](@)
   iff it is part of a [right-flanking delimiter run].
   {: id="20210408153138-pauonxa"}
4. {: id="20210408153137-hcudpxj"}A single `_` character [can close emphasis] iff
   it is part of a [right-flanking delimiter run]
   and either (a) not part of a [left-flanking delimiter run]
   or (b) part of a [left-flanking delimiter run]
   followed by punctuation.
   {: id="20210408153138-dwg2164"}
5. {: id="20210408153137-nuq3gck"}A double `**` [can open strong emphasis](@)
   iff it is part of a [left-flanking delimiter run].
   {: id="20210408153138-57n75de"}
6. {: id="20210408153137-5q1jc8w"}A double `__` [can open strong emphasis] iff
   it is part of a [left-flanking delimiter run]
   and either (a) not part of a [right-flanking delimiter run]
   or (b) part of a [right-flanking delimiter run]
   preceded by punctuation.
   {: id="20210408153138-jrgeiat"}
7. {: id="20210408153137-uwofi4d"}A double `**` [can close strong emphasis](@)
   iff it is part of a [right-flanking delimiter run].
   {: id="20210408153138-xnqoi9o"}
8. {: id="20210408153137-ao76yzf"}A double `__` [can close strong emphasis] iff
   it is part of a [right-flanking delimiter run]
   and either (a) not part of a [left-flanking delimiter run]
   or (b) part of a [left-flanking delimiter run]
   followed by punctuation.
   {: id="20210408153138-kya4h2j"}
9. {: id="20210408153137-pr6q2bw"}Emphasis begins with a delimiter that [can open emphasis] and ends
   with a delimiter that [can close emphasis], and that uses the same
   character (`_` or `*`) as the opening delimiter.  The
   opening and closing delimiters must belong to separate
   [delimiter runs].  If one of the delimiters can both
   open and close emphasis, then the sum of the lengths of the
   delimiter runs containing the opening and closing delimiters
   must not be a multiple of 3 unless both lengths are
   multiples of 3.
   {: id="20210408153138-llhqe2m"}
10. {: id="20210408153137-p224xxa"}Strong emphasis begins with a delimiter that
    [can open strong emphasis] and ends with a delimiter that
    [can close strong emphasis], and that uses the same character
    (`_` or `*`) as the opening delimiter.  The
    opening and closing delimiters must belong to separate
    [delimiter runs].  If one of the delimiters can both open
    and close strong emphasis, then the sum of the lengths of
    the delimiter runs containing the opening and closing
    delimiters must not be a multiple of 3 unless both lengths
    are multiples of 3.
    {: id="20210408153138-axrg69y"}
11. {: id="20210408153137-qr98cfh"}A literal `*` character cannot occur at the beginning or end of
    `*`-delimited emphasis or `**`-delimited strong emphasis, unless it
    is backslash-escaped.
    {: id="20210408153138-nyu7ve3"}
12. {: id="20210408153137-6zpjxc7"}A literal `_` character cannot occur at the beginning or end of
    `_`-delimited emphasis or `__`-delimited strong emphasis, unless it
    is backslash-escaped.
    {: id="20210408153138-f8csh1d"}
{: id="20210408153138-bqg2o3f"}

Where rules 1--12 above are compatible with multiple parsings,
the following principles resolve ambiguity:
{: id="20210408153138-q6ht9dx"}

13. {: id="20210408153137-zp29q4z"}The number of nestings should be minimized. Thus, for example,
    an interpretation `<strong>...</strong>` is always preferred to
    `<em><em>...</em></em>`.
    {: id="20210408153138-vyell5w"}
14. {: id="20210408153137-rv0rlvu"}An interpretation `<em><strong>...</strong></em>` is always
    preferred to `<strong><em>...</em></strong>`.
    {: id="20210408153138-ies0fwj"}
15. {: id="20210408153137-llsn5z3"}When two potential emphasis or strong emphasis spans overlap,
    so that the second begins before the first ends and ends after
    the first ends, the first takes precedence. Thus, for example,
    `*foo _bar* baz_` is parsed as `<em>foo _bar</em> baz_` rather
    than `*foo <em>bar* baz</em>`.
    {: id="20210408153138-pw5ksun"}
16. {: id="20210408153137-g34ke7f"}When there are two potential emphasis or strong emphasis spans
    with the same closing delimiter, the shorter one (the one that
    opens later) takes precedence. Thus, for example,
    `**foo **bar baz**` is parsed as `**foo <strong>bar baz</strong>`
    rather than `<strong>foo **bar baz</strong>`.
    {: id="20210408153138-dpl22yr"}
17. {: id="20210408153137-oh361vy"}Inline code spans, links, images, and HTML tags group more tightly
    than emphasis.  So, when there is a choice between an interpretation
    that contains one of these elements and one that does not, the
    former always wins.  Thus, for example, `*[foo*](bar)` is
    parsed as `*<a href="bar">foo*</a>` rather than as
    `<em>[foo</em>](bar)`.
    {: id="20210408153138-wiww8bv"}
{: id="20210408153138-os4tzcv"}

These rules can be illustrated through a series of examples.
{: id="20210408153138-a2fmlmc"}

Rule 1:
{: id="20210408153138-n0lnzkp"}

````````````````````````````````example
*foo bar*
.
<p><em>foo bar</em></p>
````````````````````````````````
{: id="20210408153138-54yc9y1"}

This is not emphasis, because the opening `*` is followed by
whitespace, and hence not part of a [left-flanking delimiter run]:
{: id="20210408153138-ay8rj46"}

````````````````````````````````example
a * foo bar*
.
<p>a * foo bar*</p>
````````````````````````````````
{: id="20210408153138-iupsgxw"}

This is not emphasis, because the opening `*` is preceded
by an alphanumeric and followed by punctuation, and hence
not part of a [left-flanking delimiter run]:
{: id="20210408153138-c7viisr"}

````````````````````````````````example
a*"foo"*
.
<p>a*&quot;foo&quot;*</p>
````````````````````````````````
{: id="20210408153138-1vd9r8q"}

Unicode nonbreaking spaces count as whitespace, too:
{: id="20210408153138-3bxdp3j"}

````````````````````````````````example
* a *
.
<p>* a *</p>
````````````````````````````````
{: id="20210408153138-ct2sntw"}

Intraword emphasis with `*` is permitted:
{: id="20210408153138-ik0hn8b"}

````````````````````````````````example
foo*bar*
.
<p>foo<em>bar</em></p>
````````````````````````````````
{: id="20210408153138-76utr25"}

````````````````````````````````example
5*6*78
.
<p>5<em>6</em>78</p>
````````````````````````````````
{: id="20210408153138-cf2rcwg"}

Rule 2:
{: id="20210408153138-mcfiiey"}

````````````````````````````````example
_foo bar_
.
<p><em>foo bar</em></p>
````````````````````````````````
{: id="20210408153138-imx8qn2"}

This is not emphasis, because the opening `_` is followed by
whitespace:
{: id="20210408153138-883s7f3"}

````````````````````````````````example
_ foo bar_
.
<p>_ foo bar_</p>
````````````````````````````````
{: id="20210408153138-ijqoxd3"}

This is not emphasis, because the opening `_` is preceded
by an alphanumeric and followed by punctuation:
{: id="20210408153138-rwh6ktc"}

````````````````````````````````example
a_"foo"_
.
<p>a_&quot;foo&quot;_</p>
````````````````````````````````
{: id="20210408153138-getyzdq"}

Emphasis with `_` is not allowed inside words:
{: id="20210408153138-88zdnw3"}

````````````````````````````````example
foo_bar_
.
<p>foo_bar_</p>
````````````````````````````````
{: id="20210408153138-2sp8mxk"}

````````````````````````````````example
5_6_78
.
<p>5_6_78</p>
````````````````````````````````
{: id="20210408153138-h289ne3"}

````````````````````````````````example
пристаням_стремятся_
.
<p>пристаням_стремятся_</p>
````````````````````````````````
{: id="20210408153138-fu10oft"}

Here `_` does not generate emphasis, because the first delimiter run
is right-flanking and the second left-flanking:
{: id="20210408153138-qufhfik"}

````````````````````````````````example
aa_"bb"_cc
.
<p>aa_&quot;bb&quot;_cc</p>
````````````````````````````````
{: id="20210408153138-4pbrwse"}

This is emphasis, even though the opening delimiter is
both left- and right-flanking, because it is preceded by
punctuation:
{: id="20210408153138-qxh9b86"}

````````````````````````````````example
foo-_(bar)_
.
<p>foo-<em>(bar)</em></p>
````````````````````````````````
{: id="20210408153138-kqakyvr"}

Rule 3:
{: id="20210408153138-nsmg0r7"}

This is not emphasis, because the closing delimiter does
not match the opening delimiter:
{: id="20210408153138-26y61dm"}

````````````````````````````````example
_foo*
.
<p>_foo*</p>
````````````````````````````````
{: id="20210408153138-pem2bl9"}

This is not emphasis, because the closing `*` is preceded by
whitespace:
{: id="20210408153138-6fvaltr"}

````````````````````````````````example
*foo bar *
.
<p>*foo bar *</p>
````````````````````````````````
{: id="20210408153138-o25us7o"}

A newline also counts as whitespace:
{: id="20210408153138-ozxlxyo"}

````````````````````````````````example
*foo bar
*
.
<p>*foo bar
*</p>
````````````````````````````````
{: id="20210408153138-1omt08l"}

This is not emphasis, because the second `*` is
preceded by punctuation and followed by an alphanumeric
(hence it is not part of a [right-flanking delimiter run]:
{: id="20210408153138-dmlbweb"}

````````````````````````````````example
*(*foo)
.
<p>*(*foo)</p>
````````````````````````````````
{: id="20210408153138-bkjzggt"}

The point of this restriction is more easily appreciated
with this example:
{: id="20210408153138-6er8kcf"}

````````````````````````````````example
*(*foo*)*
.
<p><em>(<em>foo</em>)</em></p>
````````````````````````````````
{: id="20210408153138-rnq670n"}

Intraword emphasis with `*` is allowed:
{: id="20210408153138-55ewcii"}

````````````````````````````````example
*foo*bar
.
<p><em>foo</em>bar</p>
````````````````````````````````
{: id="20210408153138-h02vcs8"}

Rule 4:
{: id="20210408153138-d8jow45"}

This is not emphasis, because the closing `_` is preceded by
whitespace:
{: id="20210408153138-8s66pvm"}

````````````````````````````````example
_foo bar _
.
<p>_foo bar _</p>
````````````````````````````````
{: id="20210408153138-7mkn91g"}

This is not emphasis, because the second `_` is
preceded by punctuation and followed by an alphanumeric:
{: id="20210408153138-6v4j3p3"}

````````````````````````````````example
_(_foo)
.
<p>_(_foo)</p>
````````````````````````````````
{: id="20210408153138-lsb6aon"}

This is emphasis within emphasis:
{: id="20210408153138-sai7358"}

````````````````````````````````example
_(_foo_)_
.
<p><em>(<em>foo</em>)</em></p>
````````````````````````````````
{: id="20210408153138-uirhqvy"}

Intraword emphasis is disallowed for `_`:
{: id="20210408153138-xbawfia"}

````````````````````````````````example
_foo_bar
.
<p>_foo_bar</p>
````````````````````````````````
{: id="20210408153138-dkg8tbz"}

````````````````````````````````example
_пристаням_стремятся
.
<p>_пристаням_стремятся</p>
````````````````````````````````
{: id="20210408153138-acagzzd"}

````````````````````````````````example
_foo_bar_baz_
.
<p><em>foo_bar_baz</em></p>
````````````````````````````````
{: id="20210408153138-z4hgycm"}

This is emphasis, even though the closing delimiter is
both left- and right-flanking, because it is followed by
punctuation:
{: id="20210408153138-1swov2k"}

````````````````````````````````example
_(bar)_.
.
<p><em>(bar)</em>.</p>
````````````````````````````````
{: id="20210408153138-sjyqfrf"}

Rule 5:
{: id="20210408153138-vlfioqd"}

````````````````````````````````example
**foo bar**
.
<p><strong>foo bar</strong></p>
````````````````````````````````
{: id="20210408153138-fhx7dx2"}

This is not strong emphasis, because the opening delimiter is
followed by whitespace:
{: id="20210408153138-me9c2lv"}

````````````````````````````````example
** foo bar**
.
<p>** foo bar**</p>
````````````````````````````````
{: id="20210408153138-0cq3kel"}

This is not strong emphasis, because the opening `**` is preceded
by an alphanumeric and followed by punctuation, and hence
not part of a [left-flanking delimiter run]:
{: id="20210408153138-dykln8y"}

````````````````````````````````example
a**"foo"**
.
<p>a**&quot;foo&quot;**</p>
````````````````````````````````
{: id="20210408153138-aj2rhwq"}

Intraword strong emphasis with `**` is permitted:
{: id="20210408153138-j1yn62k"}

````````````````````````````````example
foo**bar**
.
<p>foo<strong>bar</strong></p>
````````````````````````````````
{: id="20210408153138-hiwvain"}

Rule 6:
{: id="20210408153138-ohs7e4o"}

````````````````````````````````example
__foo bar__
.
<p><strong>foo bar</strong></p>
````````````````````````````````
{: id="20210408153138-8mljb42"}

This is not strong emphasis, because the opening delimiter is
followed by whitespace:
{: id="20210408153138-2euw37y"}

````````````````````````````````example
__ foo bar__
.
<p>__ foo bar__</p>
````````````````````````````````
{: id="20210408153138-mildtoj"}

A newline counts as whitespace:
{: id="20210408153138-qcwg26n"}

````````````````````````````````example
__
foo bar__
.
<p>__
foo bar__</p>
````````````````````````````````
{: id="20210408153138-w860zb0"}

This is not strong emphasis, because the opening `__` is preceded
by an alphanumeric and followed by punctuation:
{: id="20210408153138-d8boeu3"}

````````````````````````````````example
a__"foo"__
.
<p>a__&quot;foo&quot;__</p>
````````````````````````````````
{: id="20210408153138-b60vwpl"}

Intraword strong emphasis is forbidden with `__`:
{: id="20210408153138-izwcy3p"}

````````````````````````````````example
foo__bar__
.
<p>foo__bar__</p>
````````````````````````````````
{: id="20210408153138-18ngk3p"}

````````````````````````````````example
5__6__78
.
<p>5__6__78</p>
````````````````````````````````
{: id="20210408153138-pdb8ehq"}

````````````````````````````````example
пристаням__стремятся__
.
<p>пристаням__стремятся__</p>
````````````````````````````````
{: id="20210408153138-8gtpb3z"}

````````````````````````````````example
__foo, __bar__, baz__
.
<p><strong>foo, <strong>bar</strong>, baz</strong></p>
````````````````````````````````
{: id="20210408153138-6zwuh3t"}

This is strong emphasis, even though the opening delimiter is
both left- and right-flanking, because it is preceded by
punctuation:
{: id="20210408153138-h57gz36"}

````````````````````````````````example
foo-__(bar)__
.
<p>foo-<strong>(bar)</strong></p>
````````````````````````````````
{: id="20210408153138-m8g7z3j"}

Rule 7:
{: id="20210408153138-8idfwby"}

This is not strong emphasis, because the closing delimiter is preceded
by whitespace:
{: id="20210408153138-d1j44y6"}

````````````````````````````````example
**foo bar **
.
<p>**foo bar **</p>
````````````````````````````````
{: id="20210408153138-zmfug5c"}

(Nor can it be interpreted as an emphasized `*foo bar *`, because of
Rule 11.)
{: id="20210408153138-d5k9qa3"}

This is not strong emphasis, because the second `**` is
preceded by punctuation and followed by an alphanumeric:
{: id="20210408153138-qreessj"}

````````````````````````````````example
**(**foo)
.
<p>**(**foo)</p>
````````````````````````````````
{: id="20210408153138-vmjrdff"}

The point of this restriction is more easily appreciated
with these examples:
{: id="20210408153138-1e9m5q7"}

````````````````````````````````example
*(**foo**)*
.
<p><em>(<strong>foo</strong>)</em></p>
````````````````````````````````
{: id="20210408153138-pj22cw7"}

````````````````````````````````example
**Gomphocarpus (*Gomphocarpus physocarpus*, syn.
*Asclepias physocarpa*)**
.
<p><strong>Gomphocarpus (<em>Gomphocarpus physocarpus</em>, syn.
<em>Asclepias physocarpa</em>)</strong></p>
````````````````````````````````
{: id="20210408153138-qg3ucku"}

````````````````````````````````example
**foo "*bar*" foo**
.
<p><strong>foo &quot;<em>bar</em>&quot; foo</strong></p>
````````````````````````````````
{: id="20210408153138-qtphjyt"}

Intraword emphasis:
{: id="20210408153138-84ikdli"}

````````````````````````````````example
**foo**bar
.
<p><strong>foo</strong>bar</p>
````````````````````````````````
{: id="20210408153138-g7epur8"}

Rule 8:
{: id="20210408153138-fplsis4"}

This is not strong emphasis, because the closing delimiter is
preceded by whitespace:
{: id="20210408153138-aym6eng"}

````````````````````````````````example
__foo bar __
.
<p>__foo bar __</p>
````````````````````````````````
{: id="20210408153138-6absmj9"}

This is not strong emphasis, because the second `__` is
preceded by punctuation and followed by an alphanumeric:
{: id="20210408153138-qg80rzz"}

````````````````````````````````example
__(__foo)
.
<p>__(__foo)</p>
````````````````````````````````
{: id="20210408153138-a4mmbsb"}

The point of this restriction is more easily appreciated
with this example:
{: id="20210408153138-06mpboo"}

````````````````````````````````example
_(__foo__)_
.
<p><em>(<strong>foo</strong>)</em></p>
````````````````````````````````
{: id="20210408153138-5e6qyfw"}

Intraword strong emphasis is forbidden with `__`:
{: id="20210408153138-4w5w9wl"}

````````````````````````````````example
__foo__bar
.
<p>__foo__bar</p>
````````````````````````````````
{: id="20210408153138-lwxlaho"}

````````````````````````````````example
__пристаням__стремятся
.
<p>__пристаням__стремятся</p>
````````````````````````````````
{: id="20210408153138-gqg5z9n"}

````````````````````````````````example
__foo__bar__baz__
.
<p><strong>foo__bar__baz</strong></p>
````````````````````````````````
{: id="20210408153138-4kfwi5p"}

This is strong emphasis, even though the closing delimiter is
both left- and right-flanking, because it is followed by
punctuation:
{: id="20210408153138-q0exgvq"}

````````````````````````````````example
__(bar)__.
.
<p><strong>(bar)</strong>.</p>
````````````````````````````````
{: id="20210408153138-t61cis1"}

Rule 9:
{: id="20210408153138-ckhc430"}

Any nonempty sequence of inline elements can be the contents of an
emphasized span.
{: id="20210408153138-bcsiwgn"}

````````````````````````````````example
*foo [bar](/url)*
.
<p><em>foo <a href="/url">bar</a></em></p>
````````````````````````````````
{: id="20210408153138-ix8hn55"}

````````````````````````````````example
*foo
bar*
.
<p><em>foo
bar</em></p>
````````````````````````````````
{: id="20210408153138-donf4d8"}

In particular, emphasis and strong emphasis can be nested
inside emphasis:
{: id="20210408153138-vcb4g9y"}

````````````````````````````````example
_foo __bar__ baz_
.
<p><em>foo <strong>bar</strong> baz</em></p>
````````````````````````````````
{: id="20210408153138-ibpb8n1"}

````````````````````````````````example
_foo _bar_ baz_
.
<p><em>foo <em>bar</em> baz</em></p>
````````````````````````````````
{: id="20210408153138-64dznwu"}

````````````````````````````````example
__foo_ bar_
.
<p><em><em>foo</em> bar</em></p>
````````````````````````````````
{: id="20210408153138-hfpw89p"}

````````````````````````````````example
*foo *bar**
.
<p><em>foo <em>bar</em></em></p>
````````````````````````````````
{: id="20210408153138-sfi24kq"}

````````````````````````````````example
*foo **bar** baz*
.
<p><em>foo <strong>bar</strong> baz</em></p>
````````````````````````````````
{: id="20210408153138-akxhps6"}

````````````````````````````````example
*foo**bar**baz*
.
<p><em>foo<strong>bar</strong>baz</em></p>
````````````````````````````````
{: id="20210408153138-pxgngzm"}

Note that in the preceding case, the interpretation
{: id="20210408153138-ezytutp"}

```markdown
<p><em>foo</em><em>bar<em></em>baz</em></p>
```
{: id="20210408153138-yfw1wx5"}

is precluded by the condition that a delimiter that
can both open and close (like the `*` after `foo`)
cannot form emphasis if the sum of the lengths of
the delimiter runs containing the opening and
closing delimiters is a multiple of 3 unless
both lengths are multiples of 3.
{: id="20210408153138-ff65cbg"}

For the same reason, we don't get two consecutive
emphasis sections in this example:
{: id="20210408153138-wv5y872"}

````````````````````````````````example
*foo**bar*
.
<p><em>foo**bar</em></p>
````````````````````````````````
{: id="20210408153138-pu0v3vb"}

The same condition ensures that the following
cases are all strong emphasis nested inside
emphasis, even when the interior spaces are
omitted:
{: id="20210408153138-1mruez5"}

````````````````````````````````example
***foo** bar*
.
<p><em><strong>foo</strong> bar</em></p>
````````````````````````````````
{: id="20210408153138-o527lps"}

````````````````````````````````example
*foo **bar***
.
<p><em>foo <strong>bar</strong></em></p>
````````````````````````````````
{: id="20210408153138-u3iimjx"}

````````````````````````````````example
*foo**bar***
.
<p><em>foo<strong>bar</strong></em></p>
````````````````````````````````
{: id="20210408153138-n5q5ozr"}

When the lengths of the interior closing and opening
delimiter runs are *both* multiples of 3, though,
they can match to create emphasis:
{: id="20210408153138-swc67ri"}

````````````````````````````````example
foo***bar***baz
.
<p>foo<em><strong>bar</strong></em>baz</p>
````````````````````````````````
{: id="20210408153138-f99pvtc"}

````````````````````````````````example
foo******bar*********baz
.
<p>foo<strong><strong><strong>bar</strong></strong></strong>***baz</p>
````````````````````````````````
{: id="20210408153138-oftch2o"}

Indefinite levels of nesting are possible:
{: id="20210408153138-m52nwup"}

````````````````````````````````example
*foo **bar *baz* bim** bop*
.
<p><em>foo <strong>bar <em>baz</em> bim</strong> bop</em></p>
````````````````````````````````
{: id="20210408153138-d6k9mso"}

````````````````````````````````example
*foo [*bar*](/url)*
.
<p><em>foo <a href="/url"><em>bar</em></a></em></p>
````````````````````````````````
{: id="20210408153138-p8x78a6"}

There can be no empty emphasis or strong emphasis:
{: id="20210408153138-oko2d73"}

````````````````````````````````example
** is not an empty emphasis
.
<p>** is not an empty emphasis</p>
````````````````````````````````
{: id="20210408153138-nu28r19"}

````````````````````````````````example
**** is not an empty strong emphasis
.
<p>**** is not an empty strong emphasis</p>
````````````````````````````````
{: id="20210408153138-gor08vu"}

Rule 10:
{: id="20210408153138-0qhflrh"}

Any nonempty sequence of inline elements can be the contents of an
strongly emphasized span.
{: id="20210408153138-c1214rr"}

````````````````````````````````example
**foo [bar](/url)**
.
<p><strong>foo <a href="/url">bar</a></strong></p>
````````````````````````````````
{: id="20210408153138-p3fl98d"}

````````````````````````````````example
**foo
bar**
.
<p><strong>foo
bar</strong></p>
````````````````````````````````
{: id="20210408153138-2lnbniy"}

In particular, emphasis and strong emphasis can be nested
inside strong emphasis:
{: id="20210408153138-4irkagw"}

````````````````````````````````example
__foo _bar_ baz__
.
<p><strong>foo <em>bar</em> baz</strong></p>
````````````````````````````````
{: id="20210408153138-rbx2fyr"}

````````````````````````````````example
__foo __bar__ baz__
.
<p><strong>foo <strong>bar</strong> baz</strong></p>
````````````````````````````````
{: id="20210408153138-wwlb2zi"}

````````````````````````````````example
____foo__ bar__
.
<p><strong><strong>foo</strong> bar</strong></p>
````````````````````````````````
{: id="20210408153138-sq4yuc8"}

````````````````````````````````example
**foo **bar****
.
<p><strong>foo <strong>bar</strong></strong></p>
````````````````````````````````
{: id="20210408153138-suz67dg"}

````````````````````````````````example
**foo *bar* baz**
.
<p><strong>foo <em>bar</em> baz</strong></p>
````````````````````````````````
{: id="20210408153138-i2mtned"}

````````````````````````````````example
**foo*bar*baz**
.
<p><strong>foo<em>bar</em>baz</strong></p>
````````````````````````````````
{: id="20210408153138-h9cto8o"}

````````````````````````````````example
***foo* bar**
.
<p><strong><em>foo</em> bar</strong></p>
````````````````````````````````
{: id="20210408153138-y9b6p8s"}

````````````````````````````````example
**foo *bar***
.
<p><strong>foo <em>bar</em></strong></p>
````````````````````````````````
{: id="20210408153138-ifhs8co"}

Indefinite levels of nesting are possible:
{: id="20210408153138-tydzwwc"}

````````````````````````````````example
**foo *bar **baz**
bim* bop**
.
<p><strong>foo <em>bar <strong>baz</strong>
bim</em> bop</strong></p>
````````````````````````````````
{: id="20210408153138-5oojkq7"}

````````````````````````````````example
**foo [*bar*](/url)**
.
<p><strong>foo <a href="/url"><em>bar</em></a></strong></p>
````````````````````````````````
{: id="20210408153138-bscd5cs"}

There can be no empty emphasis or strong emphasis:
{: id="20210408153138-vyqrxk2"}

````````````````````````````````example
__ is not an empty emphasis
.
<p>__ is not an empty emphasis</p>
````````````````````````````````
{: id="20210408153138-30q5m3d"}

````````````````````````````````example
____ is not an empty strong emphasis
.
<p>____ is not an empty strong emphasis</p>
````````````````````````````````
{: id="20210408153138-rc2dzy2"}

Rule 11:
{: id="20210408153138-lafbx0q"}

````````````````````````````````example
foo ***
.
<p>foo ***</p>
````````````````````````````````
{: id="20210408153138-90qlkcb"}

````````````````````````````````example
foo *\**
.
<p>foo <em>*</em></p>
````````````````````````````````
{: id="20210408153138-qihoj9r"}

````````````````````````````````example
foo *_*
.
<p>foo <em>_</em></p>
````````````````````````````````
{: id="20210408153138-114dqjj"}

````````````````````````````````example
foo *****
.
<p>foo *****</p>
````````````````````````````````
{: id="20210408153138-rmw72ue"}

````````````````````````````````example
foo **\***
.
<p>foo <strong>*</strong></p>
````````````````````````````````
{: id="20210408153138-eaaj28z"}

````````````````````````````````example
foo **_**
.
<p>foo <strong>_</strong></p>
````````````````````````````````
{: id="20210408153138-gf3au4g"}

Note that when delimiters do not match evenly, Rule 11 determines
that the excess literal `*` characters will appear outside of the
emphasis, rather than inside it:
{: id="20210408153138-npvnt7k"}

````````````````````````````````example
**foo*
.
<p>*<em>foo</em></p>
````````````````````````````````
{: id="20210408153138-2cthg4c"}

````````````````````````````````example
*foo**
.
<p><em>foo</em>*</p>
````````````````````````````````
{: id="20210408153138-bhm63hk"}

````````````````````````````````example
***foo**
.
<p>*<strong>foo</strong></p>
````````````````````````````````
{: id="20210408153138-718euvh"}

````````````````````````````````example
****foo*
.
<p>***<em>foo</em></p>
````````````````````````````````
{: id="20210408153138-uhjg6ft"}

````````````````````````````````example
**foo***
.
<p><strong>foo</strong>*</p>
````````````````````````````````
{: id="20210408153138-cffi8lr"}

````````````````````````````````example
*foo****
.
<p><em>foo</em>***</p>
````````````````````````````````
{: id="20210408153138-udjm70d"}

Rule 12:
{: id="20210408153138-qqp5ldw"}

````````````````````````````````example
foo ___
.
<p>foo ___</p>
````````````````````````````````
{: id="20210408153138-2hfyiy1"}

````````````````````````````````example
foo _\__
.
<p>foo <em>_</em></p>
````````````````````````````````
{: id="20210408153138-vt4jbdi"}

````````````````````````````````example
foo _*_
.
<p>foo <em>*</em></p>
````````````````````````````````
{: id="20210408153138-yjk36nk"}

````````````````````````````````example
foo _____
.
<p>foo _____</p>
````````````````````````````````
{: id="20210408153138-7e0qhm1"}

````````````````````````````````example
foo __\___
.
<p>foo <strong>_</strong></p>
````````````````````````````````
{: id="20210408153138-qlu5dmf"}

````````````````````````````````example
foo __*__
.
<p>foo <strong>*</strong></p>
````````````````````````````````
{: id="20210408153138-02ccz6r"}

````````````````````````````````example
__foo_
.
<p>_<em>foo</em></p>
````````````````````````````````
{: id="20210408153138-594n1ny"}

Note that when delimiters do not match evenly, Rule 12 determines
that the excess literal `_` characters will appear outside of the
emphasis, rather than inside it:
{: id="20210408153138-y4xv95c"}

````````````````````````````````example
_foo__
.
<p><em>foo</em>_</p>
````````````````````````````````
{: id="20210408153138-ysdg6h8"}

````````````````````````````````example
___foo__
.
<p>_<strong>foo</strong></p>
````````````````````````````````
{: id="20210408153138-d9uvyn9"}

````````````````````````````````example
____foo_
.
<p>___<em>foo</em></p>
````````````````````````````````
{: id="20210408153138-b6hje3m"}

````````````````````````````````example
__foo___
.
<p><strong>foo</strong>_</p>
````````````````````````````````
{: id="20210408153138-8vtudg6"}

````````````````````````````````example
_foo____
.
<p><em>foo</em>___</p>
````````````````````````````````
{: id="20210408153138-k38rhop"}

Rule 13 implies that if you want emphasis nested directly inside
emphasis, you must use different delimiters:
{: id="20210408153138-utrf6am"}

````````````````````````````````example
**foo**
.
<p><strong>foo</strong></p>
````````````````````````````````
{: id="20210408153138-wtjzcvn"}

````````````````````````````````example
*_foo_*
.
<p><em><em>foo</em></em></p>
````````````````````````````````
{: id="20210408153138-tmdr5fk"}

````````````````````````````````example
__foo__
.
<p><strong>foo</strong></p>
````````````````````````````````
{: id="20210408153138-qockk4w"}

````````````````````````````````example
_*foo*_
.
<p><em><em>foo</em></em></p>
````````````````````````````````
{: id="20210408153138-fcmdb6a"}

However, strong emphasis within strong emphasis is possible without
switching delimiters:
{: id="20210408153138-5kdm77z"}

````````````````````````````````example
****foo****
.
<p><strong><strong>foo</strong></strong></p>
````````````````````````````````
{: id="20210408153138-9g57zdo"}

````````````````````````````````example
____foo____
.
<p><strong><strong>foo</strong></strong></p>
````````````````````````````````
{: id="20210408153138-y41din3"}

Rule 13 can be applied to arbitrarily long sequences of
delimiters:
{: id="20210408153138-0ky3d09"}

````````````````````````````````example
******foo******
.
<p><strong><strong><strong>foo</strong></strong></strong></p>
````````````````````````````````
{: id="20210408153138-74cqim9"}

Rule 14:
{: id="20210408153138-yubcq0t"}

````````````````````````````````example
***foo***
.
<p><em><strong>foo</strong></em></p>
````````````````````````````````
{: id="20210408153138-9c9lo1s"}

````````````````````````````````example
_____foo_____
.
<p><em><strong><strong>foo</strong></strong></em></p>
````````````````````````````````
{: id="20210408153138-79d8y1e"}

Rule 15:
{: id="20210408153138-45b0edm"}

````````````````````````````````example
*foo _bar* baz_
.
<p><em>foo _bar</em> baz_</p>
````````````````````````````````
{: id="20210408153138-u7fwo4l"}

````````````````````````````````example
*foo __bar *baz bim__ bam*
.
<p><em>foo <strong>bar *baz bim</strong> bam</em></p>
````````````````````````````````
{: id="20210408153138-swi763v"}

Rule 16:
{: id="20210408153138-9boc95a"}

````````````````````````````````example
**foo **bar baz**
.
<p>**foo <strong>bar baz</strong></p>
````````````````````````````````
{: id="20210408153138-zlgdebl"}

````````````````````````````````example
*foo *bar baz*
.
<p>*foo <em>bar baz</em></p>
````````````````````````````````
{: id="20210408153138-3ct2p55"}

Rule 17:
{: id="20210408153138-ql4rw5n"}

````````````````````````````````example
*[bar*](/url)
.
<p>*<a href="/url">bar*</a></p>
````````````````````````````````
{: id="20210408153138-y41pzcn"}

````````````````````````````````example
_foo [bar_](/url)
.
<p>_foo <a href="/url">bar_</a></p>
````````````````````````````````
{: id="20210408153138-ufe4n0n"}

````````````````````````````````example
*<img src="foo" title="*"/>
.
<p>*<img src="foo" title="*"/></p>
````````````````````````````````
{: id="20210408153138-ihkzljz"}

````````````````````````````````example
**<a href="**">
.
<p>**<a href="**"></p>
````````````````````````````````
{: id="20210408153138-em2mj6w"}

````````````````````````````````example
__<a href="__">
.
<p>__<a href="__"></p>
````````````````````````````````
{: id="20210408153138-r81nd3y"}

````````````````````````````````example
*a `*`*
.
<p><em>a <code>*</code></em></p>
````````````````````````````````
{: id="20210408153138-1vo2tsh"}

````````````````````````````````example
_a `_`_
.
<p><em>a <code>_</code></em></p>
````````````````````````````````
{: id="20210408153138-82o0z5t"}

````````````````````````````````example
**a<http://foo.bar/?q=**>
.
<p>**a<a href="http://foo.bar/?q=**">http://foo.bar/?q=**</a></p>
````````````````````````````````
{: id="20210408153138-st9pn8e"}

````````````````````````````````example
__a<http://foo.bar/?q=__>
.
<p>__a<a href="http://foo.bar/?q=__">http://foo.bar/?q=__</a></p>
````````````````````````````````
{: id="20210408153138-7xvy5nt"}

## Links
{: id="20210408153138-fnczk0x"}

A link contains [link text] (the visible text), a [link destination]
(the URI that is the link destination), and optionally a [link title].
There are two basic kinds of links in Markdown.  In [inline links] the
destination and title are given immediately after the link text.  In
[reference links] the destination and title are defined elsewhere in
the document.
{: id="20210408153138-nxoablh"}

A [link text](@) consists of a sequence of zero or more
inline elements enclosed by square brackets (`[` and `]`).  The
following rules apply:
{: id="20210408153138-3faanc6"}

- {: id="20210408153137-jhretv9"}Links may not contain other links, at any level of nesting. If
  multiple otherwise valid link definitions appear nested inside each
  other, the inner-most definition is used.
  {: id="20210408153138-inl4oxn"}
- {: id="20210408153137-0covlb4"}Brackets are allowed in the [link text] only if (a) they
  are backslash-escaped or (b) they appear as a matched pair of brackets,
  with an open bracket `[`, a sequence of zero or more inlines, and
  a close bracket `]`.
  {: id="20210408153138-qm012kl"}
- {: id="20210408153137-57xeaun"}Backtick [code spans], [autolinks], and raw [HTML tags] bind more tightly
  than the brackets in link text.  Thus, for example,
  ``[foo`]` `` could not be a link text, since the second `]`
  is part of a code span.
  {: id="20210408153138-497ut28"}
- {: id="20210408153137-dc9u2a6"}The brackets in link text bind more tightly than markers for
  [emphasis and strong emphasis]. Thus, for example, `*[foo*](url)` is a link.
  {: id="20210408153138-zkwkbyd"}
{: id="20210408153138-5oqmj13"}

A [link destination](@) consists of either
{: id="20210408153138-modr696"}

- {: id="20210408153137-4ppa7bt"}a sequence of zero or more characters between an opening `<` and a
  closing `>` that contains no line breaks or unescaped
  `<` or `>` characters, or
  {: id="20210408153138-bj6071j"}
- {: id="20210408153137-klezbr1"}a nonempty sequence of characters that does not start with
  `<`, does not include ASCII space or control characters, and
  includes parentheses only if (a) they are backslash-escaped or
  (b) they are part of a balanced pair of unescaped parentheses.
  (Implementations may impose limits on parentheses nesting to
  avoid performance issues, but at least three levels of nesting
  should be supported.)
  {: id="20210408153138-3fbg0yv"}
{: id="20210408153138-d16iwv3"}

A [link title](@)  consists of either
{: id="20210408153138-2dnonau"}

- {: id="20210408153137-di47sbv"}a sequence of zero or more characters between straight double-quote
  characters (`"`), including a `"` character only if it is
  backslash-escaped, or
  {: id="20210408153138-8uf20hz"}
- {: id="20210408153137-5f416wt"}a sequence of zero or more characters between straight single-quote
  characters (`'`), including a `'` character only if it is
  backslash-escaped, or
  {: id="20210408153138-nbts9og"}
- {: id="20210408153137-lvb46mi"}a sequence of zero or more characters between matching parentheses
  (`(...)`), including a `(` or `)` character only if it is
  backslash-escaped.
  {: id="20210408153138-w8bvw0r"}
{: id="20210408153138-v0czndo"}

Although [link titles] may span multiple lines, they may not contain
a [blank line].
{: id="20210408153138-95ui650"}

An [inline link](@) consists of a [link text] followed immediately
by a left parenthesis `(`, optional [whitespace], an optional
[link destination], an optional [link title] separated from the link
destination by [whitespace], optional [whitespace], and a right
parenthesis `)`. The link's text consists of the inlines contained
in the [link text] (excluding the enclosing square brackets).
The link's URI consists of the link destination, excluding enclosing
`<...>` if present, with backslash-escapes in effect as described
above.  The link's title consists of the link title, excluding its
enclosing delimiters, with backslash-escapes in effect as described
above.
{: id="20210408153138-upzlrzt"}

Here is a simple inline link:
{: id="20210408153138-d0aa4bi"}

````````````````````````````````example
[link](/uri "title")
.
<p><a href="/uri" title="title">link</a></p>
````````````````````````````````
{: id="20210408153138-ae5ga88"}

The title may be omitted:
{: id="20210408153138-bv4zi6r"}

````````````````````````````````example
[link](/uri)
.
<p><a href="/uri">link</a></p>
````````````````````````````````
{: id="20210408153138-958mh9k"}

Both the title and the destination may be omitted:
{: id="20210408153138-s4a29da"}

````````````````````````````````example
[link]()
.
<p><a href="">link</a></p>
````````````````````````````````
{: id="20210408153138-xw1irg0"}

````````````````````````````````example
[link](<>)
.
<p><a href="">link</a></p>
````````````````````````````````
{: id="20210408153138-75t34lu"}

The destination can only contain spaces if it is
enclosed in pointy brackets:
{: id="20210408153138-cekfsry"}

````````````````````````````````example
[link](/my uri)
.
<p>[link](/my uri)</p>
````````````````````````````````
{: id="20210408153138-f33n7ks"}

````````````````````````````````example
[link](</my uri>)
.
<p><a href="/my%20uri">link</a></p>
````````````````````````````````
{: id="20210408153138-zn6w8mo"}

The destination cannot contain line breaks,
even if enclosed in pointy brackets:
{: id="20210408153138-jotv5w3"}

````````````````````````````````example
[link](foo
bar)
.
<p>[link](foo
bar)</p>
````````````````````````````````
{: id="20210408153138-j65fgce"}

````````````````````````````````example
[link](<foo
bar>)
.
<p>[link](<foo
bar>)</p>
````````````````````````````````
{: id="20210408153138-a7pnowc"}

The destination can contain `)` if it is enclosed
in pointy brackets:
{: id="20210408153138-6vtzja6"}

````````````````````````````````example
[a](<b)c>)
.
<p><a href="b)c">a</a></p>
````````````````````````````````
{: id="20210408153138-9pivsop"}

Pointy brackets that enclose links must be unescaped:
{: id="20210408153138-tla8xny"}

````````````````````````````````example
[link](<foo\>)
.
<p>[link](&lt;foo&gt;)</p>
````````````````````````````````
{: id="20210408153138-ilwhif3"}

These are not links, because the opening pointy bracket
is not matched properly:
{: id="20210408153138-139u5jp"}

````````````````````````````````example
[a](<b)c
[a](<b)c>
[a](<b>c)
.
<p>[a](&lt;b)c
[a](&lt;b)c&gt;
[a](<b>c)</p>
````````````````````````````````
{: id="20210408153138-d682egg"}

Parentheses inside the link destination may be escaped:
{: id="20210408153138-74is7o1"}

````````````````````````````````example
[link](\(foo\))
.
<p><a href="(foo)">link</a></p>
````````````````````````````````
{: id="20210408153138-sj6nuep"}

Any number of parentheses are allowed without escaping, as long as they are
balanced:
{: id="20210408153138-l29wz55"}

````````````````````````````````example
[link](foo(and(bar)))
.
<p><a href="foo(and(bar))">link</a></p>
````````````````````````````````
{: id="20210408153138-tlmv7ur"}

However, if you have unbalanced parentheses, you need to escape or use the
`<...>` form:
{: id="20210408153138-qd5e0mn"}

````````````````````````````````example
[link](foo\(and\(bar\))
.
<p><a href="foo(and(bar)">link</a></p>
````````````````````````````````
{: id="20210408153138-7ziv7u4"}

````````````````````````````````example
[link](<foo(and(bar)>)
.
<p><a href="foo(and(bar)">link</a></p>
````````````````````````````````
{: id="20210408153138-lhr44ts"}

Parentheses and other symbols can also be escaped, as usual
in Markdown:
{: id="20210408153138-9dvzl37"}

````````````````````````````````example
[link](foo\)\:)
.
<p><a href="foo):">link</a></p>
````````````````````````````````
{: id="20210408153138-ate4z3n"}

A link can contain fragment identifiers and queries:
{: id="20210408153138-cx94x7a"}

````````````````````````````````example
[link](#fragment)

[link](http://example.com#fragment)

[link](http://example.com?foo=3#frag)
.
<p><a href="#fragment">link</a></p>
<p><a href="http://example.com#fragment">link</a></p>
<p><a href="http://example.com?foo=3#frag">link</a></p>
````````````````````````````````
{: id="20210408153138-jhdborv"}

Note that a backslash before a non-escapable character is
just a backslash:
{: id="20210408153138-bvqtcss"}

````````````````````````````````example
[link](foo\bar)
.
<p><a href="foo%5Cbar">link</a></p>
````````````````````````````````
{: id="20210408153138-mi0txab"}

URL-escaping should be left alone inside the destination, as all
URL-escaped characters are also valid URL characters. Entity and
numerical character references in the destination will be parsed
into the corresponding Unicode code points, as usual.  These may
be optionally URL-escaped when written as HTML, but this spec
does not enforce any particular policy for rendering URLs in
HTML or other formats.  Renderers may make different decisions
about how to escape or normalize URLs in the output.
{: id="20210408153138-x6e2rf4"}

````````````````````````````````example
[link](foo%20b&auml;)
.
<p><a href="foo%20b%C3%A4">link</a></p>
````````````````````````````````
{: id="20210408153138-z1su51m"}

Note that, because titles can often be parsed as destinations,
if you try to omit the destination and keep the title, you'll
get unexpected results:
{: id="20210408153138-b1saf1f"}

````````````````````````````````example
[link]("title")
.
<p><a href="%22title%22">link</a></p>
````````````````````````````````
{: id="20210408153138-8xyl8qd"}

Titles may be in single quotes, double quotes, or parentheses:
{: id="20210408153138-kiv5pyc"}

````````````````````````````````example
[link](/url "title")
[link](/url 'title')
[link](/url (title))
.
<p><a href="/url" title="title">link</a>
<a href="/url" title="title">link</a>
<a href="/url" title="title">link</a></p>
````````````````````````````````
{: id="20210408153138-tpxojwt"}

Backslash escapes and entity and numeric character references
may be used in titles:
{: id="20210408153138-96606zt"}

````````````````````````````````example
[link](/url "title \"&quot;")
.
<p><a href="/url" title="title &quot;&quot;">link</a></p>
````````````````````````````````
{: id="20210408153138-blmbxx2"}

Titles must be separated from the link using a [whitespace].
Other [Unicode whitespace] like non-breaking space doesn't work.
{: id="20210408153138-0k8io8k"}

````````````````````````````````example
[link](/url "title")
.
<p><a href="/url%C2%A0%22title%22">link</a></p>
````````````````````````````````
{: id="20210408153138-ut0l9nm"}

Nested balanced quotes are not allowed without escaping:
{: id="20210408153138-cgjxku3"}

````````````````````````````````example
[link](/url "title "and" title")
.
<p>[link](/url &quot;title &quot;and&quot; title&quot;)</p>
````````````````````````````````
{: id="20210408153138-4obfp91"}

But it is easy to work around this by using a different quote type:
{: id="20210408153138-hoefww6"}

````````````````````````````````example
[link](/url 'title "and" title')
.
<p><a href="/url" title="title &quot;and&quot; title">link</a></p>
````````````````````````````````
{: id="20210408153138-i2y7z9o"}

(Note:  `Markdown.pl` did allow double quotes inside a double-quoted
title, and its test suite included a test demonstrating this.
But it is hard to see a good rationale for the extra complexity this
brings, since there are already many ways---backslash escaping,
entity and numeric character references, or using a different
quote type for the enclosing title---to write titles containing
double quotes.  `Markdown.pl`'s handling of titles has a number
of other strange features.  For example, it allows single-quoted
titles in inline links, but not reference links.  And, in
reference links but not inline links, it allows a title to begin
with `"` and end with `)`.  `Markdown.pl` 1.0.1 even allows
titles with no closing quotation mark, though 1.0.2b8 does not.
It seems preferable to adopt a simple, rational rule that works
the same way in inline links and link reference definitions.)
{: id="20210408153138-4h9t1x9"}

[Whitespace] is allowed around the destination and title:
{: id="20210408153138-3cmu7ym"}

````````````````````````````````example
[link](   /uri
  "title"  )
.
<p><a href="/uri" title="title">link</a></p>
````````````````````````````````
{: id="20210408153138-24upob9"}

But it is not allowed between the link text and the
following parenthesis:
{: id="20210408153138-krqm4v6"}

````````````````````````````````example
[link] (/uri)
.
<p>[link] (/uri)</p>
````````````````````````````````
{: id="20210408153138-jxjneae"}

The link text may contain balanced brackets, but not unbalanced ones,
unless they are escaped:
{: id="20210408153138-tltutr8"}

````````````````````````````````example
[link [foo [bar]]](/uri)
.
<p><a href="/uri">link [foo [bar]]</a></p>
````````````````````````````````
{: id="20210408153138-7lejvuk"}

````````````````````````````````example
[link] bar](/uri)
.
<p>[link] bar](/uri)</p>
````````````````````````````````
{: id="20210408153138-ucd646b"}

````````````````````````````````example
[link [bar](/uri)
.
<p>[link <a href="/uri">bar</a></p>
````````````````````````````````
{: id="20210408153138-2y7ijmg"}

````````````````````````````````example
[link \[bar](/uri)
.
<p><a href="/uri">link [bar</a></p>
````````````````````````````````
{: id="20210408153138-xnvdj4e"}

The link text may contain inline content:
{: id="20210408153138-87uchsr"}

````````````````````````````````example
[link *foo **bar** `#`*](/uri)
.
<p><a href="/uri">link <em>foo <strong>bar</strong> <code>#</code></em></a></p>
````````````````````````````````
{: id="20210408153138-xz3xax8"}

````````````````````````````````example
[![moon](moon.jpg)](/uri)
.
<p><a href="/uri"><img src="moon.jpg" alt="moon" /></a></p>
````````````````````````````````
{: id="20210408153138-h1431tc"}

However, links may not contain other links, at any level of nesting.
{: id="20210408153138-0zbq8eb"}

````````````````````````````````example
[foo [bar](/uri)](/uri)
.
<p>[foo <a href="/uri">bar</a>](/uri)</p>
````````````````````````````````
{: id="20210408153138-5dm3tuc"}

````````````````````````````````example
[foo *[bar [baz](/uri)](/uri)*](/uri)
.
<p>[foo <em>[bar <a href="/uri">baz</a>](/uri)</em>](/uri)</p>
````````````````````````````````
{: id="20210408153138-3by4czk"}

````````````````````````````````example
![[[foo](uri1)](uri2)](uri3)
.
<p><img src="uri3" alt="[foo](uri2)" /></p>
````````````````````````````````
{: id="20210408153138-jpfb3j1"}

These cases illustrate the precedence of link text grouping over
emphasis grouping:
{: id="20210408153138-ft14wlp"}

````````````````````````````````example
*[foo*](/uri)
.
<p>*<a href="/uri">foo*</a></p>
````````````````````````````````
{: id="20210408153138-faa3crl"}

````````````````````````````````example
[foo *bar](baz*)
.
<p><a href="baz*">foo *bar</a></p>
````````````````````````````````
{: id="20210408153138-34kvklq"}

Note that brackets that *aren't* part of links do not take
precedence:
{: id="20210408153138-pos23hq"}

````````````````````````````````example
*foo [bar* baz]
.
<p><em>foo [bar</em> baz]</p>
````````````````````````````````
{: id="20210408153138-elkixdo"}

These cases illustrate the precedence of HTML tags, code spans,
and autolinks over link grouping:
{: id="20210408153138-mnohs4i"}

````````````````````````````````example
[foo <bar attr="](baz)">
.
<p>[foo <bar attr="](baz)"></p>
````````````````````````````````
{: id="20210408153138-6wskk1y"}

````````````````````````````````example
[foo`](/uri)`
.
<p>[foo<code>](/uri)</code></p>
````````````````````````````````
{: id="20210408153138-2zdb9lp"}

````````````````````````````````example
[foo<http://example.com/?search=](uri)>
.
<p>[foo<a href="http://example.com/?search=%5D(uri)">http://example.com/?search=](uri)</a></p>
````````````````````````````````
{: id="20210408153138-tj4a13h"}

There are three kinds of [reference link](@)s:
[full](#full-reference-link), [collapsed](#collapsed-reference-link),
and [shortcut](#shortcut-reference-link).
{: id="20210408153138-aw6yjhz"}

A [full reference link](@)
consists of a [link text] immediately followed by a [link label]
that [matches] a [link reference definition] elsewhere in the document.
{: id="20210408153138-1801qve"}

A [link label](@)  begins with a left bracket (`[`) and ends
with the first right bracket (`]`) that is not backslash-escaped.
Between these brackets there must be at least one [non-whitespace character].
Unescaped square bracket characters are not allowed inside the
opening and closing square brackets of [link labels].  A link
label can have at most 999 characters inside the square
brackets.
{: id="20210408153138-lxiumvn"}

One label [matches](@)
another just in case their normalized forms are equal.  To normalize a
label, strip off the opening and closing brackets,
perform the *Unicode case fold*, strip leading and trailing
[whitespace] and collapse consecutive internal
[whitespace] to a single space.  If there are multiple
matching reference link definitions, the one that comes first in the
document is used.  (It is desirable in such cases to emit a warning.)
{: id="20210408153138-nsc9sr7"}

The contents of the first link label are parsed as inlines, which are
used as the link's text.  The link's URI and title are provided by the
matching [link reference definition].
{: id="20210408153138-w7aoff0"}

Here is a simple example:
{: id="20210408153138-mv9lg12"}

````````````````````````````````example
[foo][bar]

[bar]: /url "title"
.
<p><a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-nsjj777"}

The rules for the [link text] are the same as with
[inline links].  Thus:
{: id="20210408153138-16em904"}

The link text may contain balanced brackets, but not unbalanced ones,
unless they are escaped:
{: id="20210408153138-78tltps"}

````````````````````````````````example
[link [foo [bar]]][ref]

[ref]: /uri
.
<p><a href="/uri">link [foo [bar]]</a></p>
````````````````````````````````
{: id="20210408153138-ulc0ysn"}

````````````````````````````````example
[link \[bar][ref]

[ref]: /uri
.
<p><a href="/uri">link [bar</a></p>
````````````````````````````````
{: id="20210408153138-ahinvm2"}

The link text may contain inline content:
{: id="20210408153138-0nt9jzl"}

````````````````````````````````example
[link *foo **bar** `#`*][ref]

[ref]: /uri
.
<p><a href="/uri">link <em>foo <strong>bar</strong> <code>#</code></em></a></p>
````````````````````````````````
{: id="20210408153138-z41yg8k"}

````````````````````````````````example
[![moon](moon.jpg)][ref]

[ref]: /uri
.
<p><a href="/uri"><img src="moon.jpg" alt="moon" /></a></p>
````````````````````````````````
{: id="20210408153138-45548yo"}

However, links may not contain other links, at any level of nesting.
{: id="20210408153138-8n4j7hs"}

````````````````````````````````example
[foo [bar](/uri)][ref]

[ref]: /uri
.
<p>[foo <a href="/uri">bar</a>]<a href="/uri">ref</a></p>
````````````````````````````````
{: id="20210408153138-i5pqwru"}

````````````````````````````````example
[foo *bar [baz][ref]*][ref]

[ref]: /uri
.
<p>[foo <em>bar <a href="/uri">baz</a></em>]<a href="/uri">ref</a></p>
````````````````````````````````
{: id="20210408153138-geql3d8"}

(In the examples above, we have two [shortcut reference links]
instead of one [full reference link].)
{: id="20210408153138-un70wa9"}

The following cases illustrate the precedence of link text grouping over
emphasis grouping:
{: id="20210408153138-sy1ukpj"}

````````````````````````````````example
*[foo*][ref]

[ref]: /uri
.
<p>*<a href="/uri">foo*</a></p>
````````````````````````````````
{: id="20210408153138-oud3wym"}

````````````````````````````````example
[foo *bar][ref]

[ref]: /uri
.
<p><a href="/uri">foo *bar</a></p>
````````````````````````````````
{: id="20210408153138-m0atbkn"}

These cases illustrate the precedence of HTML tags, code spans,
and autolinks over link grouping:
{: id="20210408153138-2s193o0"}

````````````````````````````````example
[foo <bar attr="][ref]">

[ref]: /uri
.
<p>[foo <bar attr="][ref]"></p>
````````````````````````````````
{: id="20210408153138-ntuanp7"}

````````````````````````````````example
[foo`][ref]`

[ref]: /uri
.
<p>[foo<code>][ref]</code></p>
````````````````````````````````
{: id="20210408153138-108ngik"}

````````````````````````````````example
[foo<http://example.com/?search=][ref]>

[ref]: /uri
.
<p>[foo<a href="http://example.com/?search=%5D%5Bref%5D">http://example.com/?search=][ref]</a></p>
````````````````````````````````
{: id="20210408153138-9l2a47n"}

Matching is case-insensitive:
{: id="20210408153138-1nf6v3o"}

````````````````````````````````example
[foo][BaR]

[bar]: /url "title"
.
<p><a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-vlh8dvp"}

Unicode case fold is used:
{: id="20210408153138-k6h5v7e"}

````````````````````````````````example
[Толпой][Толпой] is a Russian word.

[ТОЛПОЙ]: /url
.
<p><a href="/url">Толпой</a> is a Russian word.</p>
````````````````````````````````
{: id="20210408153138-2g7l9ez"}

Consecutive internal [whitespace] is treated as one space for
purposes of determining matching:
{: id="20210408153138-8rxpijd"}

````````````````````````````````example
[Foo
  bar]: /url

[Baz][Foo bar]
.
<p><a href="/url">Baz</a></p>
````````````````````````````````
{: id="20210408153138-froi0t0"}

No [whitespace] is allowed between the [link text] and the
[link label]:
{: id="20210408153138-1hj85zv"}

````````````````````````````````example
[foo] [bar]

[bar]: /url "title"
.
<p>[foo] <a href="/url" title="title">bar</a></p>
````````````````````````````````
{: id="20210408153138-34wk0pn"}

````````````````````````````````example
[foo]
[bar]

[bar]: /url "title"
.
<p>[foo]
<a href="/url" title="title">bar</a></p>
````````````````````````````````
{: id="20210408153138-ngd9quk"}

This is a departure from John Gruber's original Markdown syntax
description, which explicitly allows whitespace between the link
text and the link label.  It brings reference links in line with
[inline links], which (according to both original Markdown and
this spec) cannot have whitespace after the link text.  More
importantly, it prevents inadvertent capture of consecutive
[shortcut reference links]. If whitespace is allowed between the
link text and the link label, then in the following we will have
a single reference link, not two shortcut reference links, as
intended:
{: id="20210408153138-mkid4vd"}

```markdown
[foo]
[bar]

[foo]: /url1
[bar]: /url2
```
{: id="20210408153138-i0p1nwq"}

(Note that [shortcut reference links] were introduced by Gruber
himself in a beta version of `Markdown.pl`, but never included
in the official syntax description.  Without shortcut reference
links, it is harmless to allow space between the link text and
link label; but once shortcut references are introduced, it is
too dangerous to allow this, as it frequently leads to
unintended results.)
{: id="20210408153138-5yiwfk9"}

When there are multiple matching [link reference definitions],
the first is used:
{: id="20210408153138-fkbwq29"}

````````````````````````````````example
[foo]: /url1

[foo]: /url2

[bar][foo]
.
<p><a href="/url1">bar</a></p>
````````````````````````````````
{: id="20210408153138-p4v2vh7"}

Note that matching is performed on normalized strings, not parsed
inline content.  So the following does not match, even though the
labels define equivalent inline content:
{: id="20210408153138-yykiut6"}

````````````````````````````````example
[bar][foo\!]

[foo!]: /url
.
<p>[bar][foo!]</p>
````````````````````````````````
{: id="20210408153138-ijvkpby"}

[Link labels] cannot contain brackets, unless they are
backslash-escaped:
{: id="20210408153138-mn9okwm"}

````````````````````````````````example
[foo][ref[]

[ref[]: /uri
.
<p>[foo][ref[]</p>
<p>[ref[]: /uri</p>
````````````````````````````````
{: id="20210408153138-s213luu"}

````````````````````````````````example
[foo][ref[bar]]

[ref[bar]]: /uri
.
<p>[foo][ref[bar]]</p>
<p>[ref[bar]]: /uri</p>
````````````````````````````````
{: id="20210408153138-cid1wp8"}

````````````````````````````````example
[[[foo]]]

[[[foo]]]: /url
.
<p>[[[foo]]]</p>
<p>[[[foo]]]: /url</p>
````````````````````````````````
{: id="20210408153138-pcojhi3"}

````````````````````````````````example
[foo][ref\[]

[ref\[]: /uri
.
<p><a href="/uri">foo</a></p>
````````````````````````````````
{: id="20210408153138-3s1hqzw"}

Note that in this example `]` is not backslash-escaped:
{: id="20210408153138-rqhmaoh"}

````````````````````````````````example
[bar\\]: /uri

[bar\\]
.
<p><a href="/uri">bar\</a></p>
````````````````````````````````
{: id="20210408153138-sgr64ek"}

A [link label] must contain at least one [non-whitespace character]:
{: id="20210408153138-cm8rquv"}

````````````````````````````````example
[]

[]: /uri
.
<p>[]</p>
<p>[]: /uri</p>
````````````````````````````````
{: id="20210408153138-sh4e00p"}

````````````````````````````````example
[
 ]

[
 ]: /uri
.
<p>[
]</p>
<p>[
]: /uri</p>
````````````````````````````````
{: id="20210408153138-5e8629t"}

A [collapsed reference link](@)
consists of a [link label] that [matches] a
[link reference definition] elsewhere in the
document, followed by the string `[]`.
The contents of the first link label are parsed as inlines,
which are used as the link's text.  The link's URI and title are
provided by the matching reference link definition.  Thus,
`[foo][]` is equivalent to `[foo][foo]`.
{: id="20210408153138-kzx2j8z"}

````````````````````````````````example
[foo][]

[foo]: /url "title"
.
<p><a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-8kuqre1"}

````````````````````````````````example
[*foo* bar][]

[*foo* bar]: /url "title"
.
<p><a href="/url" title="title"><em>foo</em> bar</a></p>
````````````````````````````````
{: id="20210408153138-5tg9qi4"}

The link labels are case-insensitive:
{: id="20210408153138-v7ctpik"}

````````````````````````````````example
[Foo][]

[foo]: /url "title"
.
<p><a href="/url" title="title">Foo</a></p>
````````````````````````````````
{: id="20210408153138-umkkuyr"}

As with full reference links, [whitespace] is not
allowed between the two sets of brackets:
{: id="20210408153138-0y6dpsv"}

````````````````````````````````example
[foo] 
[]

[foo]: /url "title"
.
<p><a href="/url" title="title">foo</a>
[]</p>
````````````````````````````````
{: id="20210408153138-11retkx"}

A [shortcut reference link](@)
consists of a [link label] that [matches] a
[link reference definition] elsewhere in the
document and is not followed by `[]` or a link label.
The contents of the first link label are parsed as inlines,
which are used as the link's text.  The link's URI and title
are provided by the matching link reference definition.
Thus, `[foo]` is equivalent to `[foo][]`.
{: id="20210408153138-6092nl1"}

````````````````````````````````example
[foo]

[foo]: /url "title"
.
<p><a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-6mxfn8w"}

````````````````````````````````example
[*foo* bar]

[*foo* bar]: /url "title"
.
<p><a href="/url" title="title"><em>foo</em> bar</a></p>
````````````````````````````````
{: id="20210408153138-7qq8xz8"}

````````````````````````````````example
[[*foo* bar]]

[*foo* bar]: /url "title"
.
<p>[<a href="/url" title="title"><em>foo</em> bar</a>]</p>
````````````````````````````````
{: id="20210408153138-2smenjp"}

````````````````````````````````example
[[bar [foo]

[foo]: /url
.
<p>[[bar <a href="/url">foo</a></p>
````````````````````````````````
{: id="20210408153138-uk5bsqf"}

The link labels are case-insensitive:
{: id="20210408153138-u1hr4yk"}

````````````````````````````````example
[Foo]

[foo]: /url "title"
.
<p><a href="/url" title="title">Foo</a></p>
````````````````````````````````
{: id="20210408153138-xga8rtc"}

A space after the link text should be preserved:
{: id="20210408153138-lem9xfg"}

````````````````````````````````example
[foo] bar

[foo]: /url
.
<p><a href="/url">foo</a> bar</p>
````````````````````````````````
{: id="20210408153138-4rhna55"}

If you just want bracketed text, you can backslash-escape the
opening bracket to avoid links:
{: id="20210408153138-659ljzi"}

````````````````````````````````example
\[foo]

[foo]: /url "title"
.
<p>[foo]</p>
````````````````````````````````
{: id="20210408153138-b4820ri"}

Note that this is a link, because a link label ends with the first
following closing bracket:
{: id="20210408153138-cfxumwu"}

````````````````````````````````example
[foo*]: /url

*[foo*]
.
<p>*<a href="/url">foo*</a></p>
````````````````````````````````
{: id="20210408153138-318i9x8"}

Full and compact references take precedence over shortcut
references:
{: id="20210408153138-3nk3xct"}

````````````````````````````````example
[foo][bar]

[foo]: /url1
[bar]: /url2
.
<p><a href="/url2">foo</a></p>
````````````````````````````````
{: id="20210408153138-xktdp2s"}

````````````````````````````````example
[foo][]

[foo]: /url1
.
<p><a href="/url1">foo</a></p>
````````````````````````````````
{: id="20210408153138-weqx72i"}

Inline links also take precedence:
{: id="20210408153138-0xn6pra"}

````````````````````````````````example
[foo]()

[foo]: /url1
.
<p><a href="">foo</a></p>
````````````````````````````````
{: id="20210408153138-yshhytq"}

````````````````````````````````example
[foo](not a link)

[foo]: /url1
.
<p><a href="/url1">foo</a>(not a link)</p>
````````````````````````````````
{: id="20210408153138-jvum0r5"}

In the following case `[bar][baz]` is parsed as a reference,
`[foo]` as normal text:
{: id="20210408153138-a5ocmnq"}

````````````````````````````````example
[foo][bar][baz]

[baz]: /url
.
<p>[foo]<a href="/url">bar</a></p>
````````````````````````````````
{: id="20210408153138-q8muqzm"}

Here, though, `[foo][bar]` is parsed as a reference, since
`[bar]` is defined:
{: id="20210408153138-jcxi048"}

````````````````````````````````example
[foo][bar][baz]

[baz]: /url1
[bar]: /url2
.
<p><a href="/url2">foo</a><a href="/url1">baz</a></p>
````````````````````````````````
{: id="20210408153138-w8kz0bh"}

Here `[foo]` is not parsed as a shortcut reference, because it
is followed by a link label (even though `[bar]` is not defined):
{: id="20210408153138-qvvq2y6"}

````````````````````````````````example
[foo][bar][baz]

[baz]: /url1
[foo]: /url2
.
<p>[foo]<a href="/url1">bar</a></p>
````````````````````````````````
{: id="20210408153138-3604q1m"}

## Images
{: id="20210408153138-jk6e5fp"}

Syntax for images is like the syntax for links, with one
difference. Instead of [link text], we have an
[image description](@).  The rules for this are the
same as for [link text], except that (a) an
image description starts with `![` rather than `[`, and
(b) an image description may contain links.
An image description has inline elements
as its contents.  When an image is rendered to HTML,
this is standardly used as the image's `alt` attribute.
{: id="20210408153138-m3c8usl"}

````````````````````````````````example
![foo](/url "title")
.
<p><img src="/url" alt="foo" title="title" /></p>
````````````````````````````````
{: id="20210408153138-bshe6n6"}

````````````````````````````````example
![foo *bar*]

[foo *bar*]: train.jpg "train & tracks"
.
<p><img src="train.jpg" alt="foo bar" title="train &amp; tracks" /></p>
````````````````````````````````
{: id="20210408153138-nhj0bvn"}

````````````````````````````````example
![foo ![bar](/url)](/url2)
.
<p><img src="/url2" alt="foo bar" /></p>
````````````````````````````````
{: id="20210408153138-paxfg2t"}

````````````````````````````````example
![foo [bar](/url)](/url2)
.
<p><img src="/url2" alt="foo bar" /></p>
````````````````````````````````
{: id="20210408153138-ip1glrd"}

Though this spec is concerned with parsing, not rendering, it is
recommended that in rendering to HTML, only the plain string content
of the [image description] be used.  Note that in
the above example, the alt attribute's value is `foo bar`, not `foo [bar](/url)` or `foo <a href="/url">bar</a>`.  Only the plain string
content is rendered, without formatting.
{: id="20210408153138-ezlsfns"}

````````````````````````````````example
![foo *bar*][]

[foo *bar*]: train.jpg "train & tracks"
.
<p><img src="train.jpg" alt="foo bar" title="train &amp; tracks" /></p>
````````````````````````````````
{: id="20210408153138-nzvkwbx"}

````````````````````````````````example
![foo *bar*][foobar]

[FOOBAR]: train.jpg "train & tracks"
.
<p><img src="train.jpg" alt="foo bar" title="train &amp; tracks" /></p>
````````````````````````````````
{: id="20210408153138-oegwyq2"}

````````````````````````````````example
![foo](train.jpg)
.
<p><img src="train.jpg" alt="foo" /></p>
````````````````````````````````
{: id="20210408153138-l4mhyc2"}

````````````````````````````````example
My ![foo bar](/path/to/train.jpg  "title"   )
.
<p>My <img src="/path/to/train.jpg" alt="foo bar" title="title" /></p>
````````````````````````````````
{: id="20210408153138-kwq6mph"}

````````````````````````````````example
![foo](<url>)
.
<p><img src="url" alt="foo" /></p>
````````````````````````````````
{: id="20210408153138-eyjr4c5"}

````````````````````````````````example
![](/url)
.
<p><img src="/url" alt="" /></p>
````````````````````````````````
{: id="20210408153138-kqq88t7"}

Reference-style:
{: id="20210408153138-ksvdnl5"}

````````````````````````````````example
![foo][bar]

[bar]: /url
.
<p><img src="/url" alt="foo" /></p>
````````````````````````````````
{: id="20210408153138-etuvwsv"}

````````````````````````````````example
![foo][bar]

[BAR]: /url
.
<p><img src="/url" alt="foo" /></p>
````````````````````````````````
{: id="20210408153138-qmm4zxq"}

Collapsed:
{: id="20210408153138-eudn7t4"}

````````````````````````````````example
![foo][]

[foo]: /url "title"
.
<p><img src="/url" alt="foo" title="title" /></p>
````````````````````````````````
{: id="20210408153138-e2k3hi3"}

````````````````````````````````example
![*foo* bar][]

[*foo* bar]: /url "title"
.
<p><img src="/url" alt="foo bar" title="title" /></p>
````````````````````````````````
{: id="20210408153138-paizbfa"}

The labels are case-insensitive:
{: id="20210408153138-px394k4"}

````````````````````````````````example
![Foo][]

[foo]: /url "title"
.
<p><img src="/url" alt="Foo" title="title" /></p>
````````````````````````````````
{: id="20210408153138-tsjb5dp"}

As with reference links, [whitespace] is not allowed
between the two sets of brackets:
{: id="20210408153138-j5pqcxl"}

````````````````````````````````example
![foo] 
[]

[foo]: /url "title"
.
<p><img src="/url" alt="foo" title="title" />
[]</p>
````````````````````````````````
{: id="20210408153138-crrijjq"}

Shortcut:
{: id="20210408153138-pyvpwro"}

````````````````````````````````example
![foo]

[foo]: /url "title"
.
<p><img src="/url" alt="foo" title="title" /></p>
````````````````````````````````
{: id="20210408153138-sn9stk1"}

````````````````````````````````example
![*foo* bar]

[*foo* bar]: /url "title"
.
<p><img src="/url" alt="foo bar" title="title" /></p>
````````````````````````````````
{: id="20210408153138-kf7pdqu"}

Note that link labels cannot contain unescaped brackets:
{: id="20210408153138-oh2ym5z"}

````````````````````````````````example
![[foo]]

[[foo]]: /url "title"
.
<p>![[foo]]</p>
<p>[[foo]]: /url &quot;title&quot;</p>
````````````````````````````````
{: id="20210408153138-1gqr7s4"}

The link labels are case-insensitive:
{: id="20210408153138-qloqitm"}

````````````````````````````````example
![Foo]

[foo]: /url "title"
.
<p><img src="/url" alt="Foo" title="title" /></p>
````````````````````````````````
{: id="20210408153138-2rbqmot"}

If you just want a literal `!` followed by bracketed text, you can
backslash-escape the opening `[`:
{: id="20210408153138-8w7rki6"}

````````````````````````````````example
!\[foo]

[foo]: /url "title"
.
<p>![foo]</p>
````````````````````````````````
{: id="20210408153138-m33e2mx"}

If you want a link after a literal `!`, backslash-escape the
`!`:
{: id="20210408153138-g6tq653"}

````````````````````````````````example
\![foo]

[foo]: /url "title"
.
<p>!<a href="/url" title="title">foo</a></p>
````````````````````````````````
{: id="20210408153138-e1x169n"}

## Autolinks
{: id="20210408153138-5xtug6a"}

[Autolink](@)s are absolute URIs and email addresses inside
`<` and `>`. They are parsed as links, with the URL or email address
as the link label.
{: id="20210408153138-9rs81dj"}

A [URI autolink](@) consists of `<`, followed by an
[absolute URI] followed by `>`.  It is parsed as
a link to the URI, with the URI as the link's label.
{: id="20210408153138-ziafu5o"}

An [absolute URI](@),
for these purposes, consists of a [scheme] followed by a colon (`:`)
followed by zero or more characters other than ASCII
[whitespace] and control characters, `<`, and `>`.  If
the URI includes these characters, they must be percent-encoded
(e.g. `%20` for a space).
{: id="20210408153138-uacbsah"}

For purposes of this spec, a [scheme](@) is any sequence
of 2--32 characters beginning with an ASCII letter and followed
by any combination of ASCII letters, digits, or the symbols plus
("+"), period ("."), or hyphen ("-").
{: id="20210408153138-bab6uqt"}

Here are some valid autolinks:
{: id="20210408153138-yskgmru"}

````````````````````````````````example
<http://foo.bar.baz>
.
<p><a href="http://foo.bar.baz">http://foo.bar.baz</a></p>
````````````````````````````````
{: id="20210408153138-a5xjtdc"}

````````````````````````````````example
<http://foo.bar.baz/test?q=hello&id=22&boolean>
.
<p><a href="http://foo.bar.baz/test?q=hello&amp;id=22&amp;boolean">http://foo.bar.baz/test?q=hello&amp;id=22&amp;boolean</a></p>
````````````````````````````````
{: id="20210408153138-tuymlj1"}

````````````````````````````````example
<irc://foo.bar:2233/baz>
.
<p><a href="irc://foo.bar:2233/baz">irc://foo.bar:2233/baz</a></p>
````````````````````````````````
{: id="20210408153138-8q4ookr"}

Uppercase is also fine:
{: id="20210408153138-zjyecri"}

````````````````````````````````example
<MAILTO:FOO@BAR.BAZ>
.
<p><a href="MAILTO:FOO@BAR.BAZ">MAILTO:FOO@BAR.BAZ</a></p>
````````````````````````````````
{: id="20210408153138-4fgh936"}

Note that many strings that count as [absolute URIs] for
purposes of this spec are not valid URIs, because their
schemes are not registered or because of other problems
with their syntax:
{: id="20210408153138-ddve8ll"}

````````````````````````````````example
<a+b+c:d>
.
<p><a href="a+b+c:d">a+b+c:d</a></p>
````````````````````````````````
{: id="20210408153138-e87khtl"}

````````````````````````````````example
<made-up-scheme://foo,bar>
.
<p><a href="made-up-scheme://foo,bar">made-up-scheme://foo,bar</a></p>
````````````````````````````````
{: id="20210408153138-sm2bw54"}

````````````````````````````````example
<http://../>
.
<p><a href="http://../">http://../</a></p>
````````````````````````````````
{: id="20210408153138-0ua5s61"}

````````````````````````````````example
<localhost:5001/foo>
.
<p><a href="localhost:5001/foo">localhost:5001/foo</a></p>
````````````````````````````````
{: id="20210408153138-c36w8ny"}

Spaces are not allowed in autolinks:
{: id="20210408153138-mr0lbas"}

````````````````````````````````example
<http://foo.bar/baz bim>
.
<p>&lt;http://foo.bar/baz bim&gt;</p>
````````````````````````````````
{: id="20210408153138-2mkdzwz"}

Backslash-escapes do not work inside autolinks:
{: id="20210408153138-i9cx8kt"}

````````````````````````````````example
<http://example.com/\[\>
.
<p><a href="http://example.com/%5C%5B%5C">http://example.com/\[\</a></p>
````````````````````````````````
{: id="20210408153138-oel1rn0"}

An [email autolink](@)
consists of `<`, followed by an [email address],
followed by `>`.  The link's label is the email address,
and the URL is `mailto:` followed by the email address.
{: id="20210408153138-5ifkzov"}

An [email address](@),
for these purposes, is anything that matches
the [non-normative regex from the HTML5
spec](https://html.spec.whatwg.org/multipage/forms.html#e-mail-state-(type=email)):
{: id="20210408153138-1nmmx8u"}

```
/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?
(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
```
{: id="20210408153138-zjh7pu2"}

Examples of email autolinks:
{: id="20210408153138-wqe8lgn"}

````````````````````````````````example
<foo@bar.example.com>
.
<p><a href="mailto:foo@bar.example.com">foo@bar.example.com</a></p>
````````````````````````````````
{: id="20210408153138-sxzxfkh"}

````````````````````````````````example
<foo+special@Bar.baz-bar0.com>
.
<p><a href="mailto:foo+special@Bar.baz-bar0.com">foo+special@Bar.baz-bar0.com</a></p>
````````````````````````````````
{: id="20210408153138-4ijg34e"}

Backslash-escapes do not work inside email autolinks:
{: id="20210408153138-gp5py5d"}

````````````````````````````````example
<foo\+@bar.example.com>
.
<p>&lt;foo+@bar.example.com&gt;</p>
````````````````````````````````
{: id="20210408153138-9bbtnfw"}

These are not autolinks:
{: id="20210408153138-47p14o1"}

````````````````````````````````example
<>
.
<p>&lt;&gt;</p>
````````````````````````````````
{: id="20210408153138-w9n7fxx"}

````````````````````````````````example
< http://foo.bar >
.
<p>&lt; http://foo.bar &gt;</p>
````````````````````````````````
{: id="20210408153138-htveric"}

````````````````````````````````example
<m:abc>
.
<p>&lt;m:abc&gt;</p>
````````````````````````````````
{: id="20210408153138-0mvo2qq"}

````````````````````````````````example
<foo.bar.baz>
.
<p>&lt;foo.bar.baz&gt;</p>
````````````````````````````````
{: id="20210408153138-cx9ocsc"}

````````````````````````````````example
http://example.com
.
<p>http://example.com</p>
````````````````````````````````
{: id="20210408153138-djzvtws"}

````````````````````````````````example
foo@bar.example.com
.
<p>foo@bar.example.com</p>
````````````````````````````````
{: id="20210408153138-jp89ozy"}

## Raw HTML
{: id="20210408153138-0vudrt3"}

Text between `<` and `>` that looks like an HTML tag is parsed as a
raw HTML tag and will be rendered in HTML without escaping.
Tag and attribute names are not limited to current HTML tags,
so custom tags (and even, say, DocBook tags) may be used.
{: id="20210408153138-twjpyys"}

Here is the grammar for tags:
{: id="20210408153138-gyiwvnx"}

A [tag name](@) consists of an ASCII letter
followed by zero or more ASCII letters, digits, or
hyphens (`-`).
{: id="20210408153138-dt3jfs9"}

An [attribute](@) consists of [whitespace],
an [attribute name], and an optional
[attribute value specification].
{: id="20210408153138-9y5hh02"}

An [attribute name](@)
consists of an ASCII letter, `_`, or `:`, followed by zero or more ASCII
letters, digits, `_`, `.`, `:`, or `-`.  (Note:  This is the XML
specification restricted to ASCII.  HTML5 is laxer.)
{: id="20210408153138-og8bg7w"}

An [attribute value specification](@)
consists of optional [whitespace],
a `=` character, optional [whitespace], and an [attribute
value].
{: id="20210408153138-9sxxd8z"}

An [attribute value](@)
consists of an [unquoted attribute value],
a [single-quoted attribute value], or a [double-quoted attribute value].
{: id="20210408153138-te73rx1"}

An [unquoted attribute value](@)
is a nonempty string of characters not
including [whitespace], `"`, `'`, `=`, `<`, `>`, or `` ` ``.
{: id="20210408153138-rx8z3hg"}

A [single-quoted attribute value](@)
consists of `'`, zero or more
characters not including `'`, and a final `'`.
{: id="20210408153138-sofxnvs"}

A [double-quoted attribute value](@)
consists of `"`, zero or more
characters not including `"`, and a final `"`.
{: id="20210408153138-hcliamx"}

An [open tag](@) consists of a `<` character, a [tag name],
zero or more [attributes], optional [whitespace], an optional `/`
character, and a `>` character.
{: id="20210408153138-dm53xzw"}

A [closing tag](@) consists of the string `</`, a
[tag name], optional [whitespace], and the character `>`.
{: id="20210408153138-38x1l0q"}

An [HTML comment](@) consists of `<!--` + *text* + `-->`,
where *text* does not start with `>` or `->`, does not end with `-`,
and does not contain `--`.  (See the
[HTML5 spec](http://www.w3.org/TR/html5/syntax.html#comments).)
{: id="20210408153138-s3e45hl"}

A [processing instruction](@)
consists of the string `<?`, a string
of characters not including the string `?>`, and the string
`?>`.
{: id="20210408153138-7omco65"}

A [declaration](@) consists of the
string `<!`, a name consisting of one or more uppercase ASCII letters,
[whitespace], a string of characters not including the
character `>`, and the character `>`.
{: id="20210408153138-o22vg2m"}

A [CDATA section](@) consists of
the string `<![CDATA[`, a string of characters not including the string
`]]>`, and the string `]]>`.
{: id="20210408153138-05fptni"}

An [HTML tag](@) consists of an [open tag], a [closing tag],
an [HTML comment], a [processing instruction], a [declaration],
or a [CDATA section].
{: id="20210408153138-ygjjovu"}

Here are some simple open tags:
{: id="20210408153138-we6bnya"}

````````````````````````````````example
<a><bab><c2c>
.
<p><a><bab><c2c></p>
````````````````````````````````
{: id="20210408153138-u2qu0pw"}

Empty elements:
{: id="20210408153138-uvcz8of"}

````````````````````````````````example
<a/><b2/>
.
<p><a/><b2/></p>
````````````````````````````````
{: id="20210408153138-rt2ydr0"}

[Whitespace] is allowed:
{: id="20210408153138-mvi6898"}

````````````````````````````````example
<a  /><b2
data="foo" >
.
<p><a  /><b2
data="foo" ></p>
````````````````````````````````
{: id="20210408153138-fkxjvsp"}

With attributes:
{: id="20210408153138-gg0dx11"}

````````````````````````````````example
<a foo="bar" bam = 'baz <em>"</em>'
_boolean zoop:33=zoop:33 />
.
<p><a foo="bar" bam = 'baz <em>"</em>'
_boolean zoop:33=zoop:33 /></p>
````````````````````````````````
{: id="20210408153138-ndp30nf"}

Custom tag names can be used:
{: id="20210408153138-5l5oh5x"}

````````````````````````````````example
Foo <responsive-image src="foo.jpg" />
.
<p>Foo <responsive-image src="foo.jpg" /></p>
````````````````````````````````
{: id="20210408153138-pd3aur9"}

Illegal tag names, not parsed as HTML:
{: id="20210408153138-8tbdmpm"}

````````````````````````````````example
<33> <__>
.
<p>&lt;33&gt; &lt;__&gt;</p>
````````````````````````````````
{: id="20210408153138-u0eyefa"}

Illegal attribute names:
{: id="20210408153138-cnk8bfv"}

````````````````````````````````example
<a h*#ref="hi">
.
<p>&lt;a h*#ref=&quot;hi&quot;&gt;</p>
````````````````````````````````
{: id="20210408153138-uf5li0w"}

Illegal attribute values:
{: id="20210408153138-psbjpkd"}

````````````````````````````````example
<a href="hi'> <a href=hi'>
.
<p>&lt;a href=&quot;hi'&gt; &lt;a href=hi'&gt;</p>
````````````````````````````````
{: id="20210408153138-9ksvc2h"}

Illegal [whitespace]:
{: id="20210408153138-k78e1n5"}

````````````````````````````````example
< a><
foo><bar/ >
<foo bar=baz
bim!bop />
.
<p>&lt; a&gt;&lt;
foo&gt;&lt;bar/ &gt;
&lt;foo bar=baz
bim!bop /&gt;</p>
````````````````````````````````
{: id="20210408153138-qf3cbl7"}

Missing [whitespace]:
{: id="20210408153138-lnqtt8f"}

````````````````````````````````example
<a href='bar'title=title>
.
<p>&lt;a href='bar'title=title&gt;</p>
````````````````````````````````
{: id="20210408153138-njt5ngq"}

Closing tags:
{: id="20210408153138-874dkoa"}

````````````````````````````````example
</a></foo >
.
<p></a></foo ></p>
````````````````````````````````
{: id="20210408153138-wsx02bl"}

Illegal attributes in closing tag:
{: id="20210408153138-b5tfuex"}

````````````````````````````````example
</a href="foo">
.
<p>&lt;/a href=&quot;foo&quot;&gt;</p>
````````````````````````````````
{: id="20210408153138-hdlea5r"}

Comments:
{: id="20210408153138-ssnam2f"}

````````````````````````````````example
foo <!-- this is a
comment - with hyphen -->
.
<p>foo <!-- this is a
comment - with hyphen --></p>
````````````````````````````````
{: id="20210408153138-6d9d1uq"}

````````````````````````````````example
foo <!-- not a comment -- two hyphens -->
.
<p>foo &lt;!-- not a comment -- two hyphens --&gt;</p>
````````````````````````````````
{: id="20210408153138-hklnly2"}

Not comments:
{: id="20210408153138-pqvwo4r"}

````````````````````````````````example
foo <!--> foo -->

foo <!-- foo--->
.
<p>foo &lt;!--&gt; foo --&gt;</p>
<p>foo &lt;!-- foo---&gt;</p>
````````````````````````````````
{: id="20210408153138-8vf2j6y"}

Processing instructions:
{: id="20210408153138-9tkt3v2"}

````````````````````````````````example
foo <?php echo $a; ?>
.
<p>foo <?php echo $a; ?></p>
````````````````````````````````
{: id="20210408153138-pmpmz8m"}

Declarations:
{: id="20210408153138-rxuudta"}

````````````````````````````````example
foo <!ELEMENT br EMPTY>
.
<p>foo <!ELEMENT br EMPTY></p>
````````````````````````````````
{: id="20210408153138-cj2zq4o"}

CDATA sections:
{: id="20210408153138-q0u1ebd"}

````````````````````````````````example
foo <![CDATA[>&<]]>
.
<p>foo <![CDATA[>&<]]></p>
````````````````````````````````
{: id="20210408153138-b31h6wi"}

Entity and numeric character references are preserved in HTML
attributes:
{: id="20210408153138-62egxvu"}

````````````````````````````````example
foo <a href="&ouml;">
.
<p>foo <a href="&ouml;"></p>
````````````````````````````````
{: id="20210408153138-oon7amn"}

Backslash escapes do not work in HTML attributes:
{: id="20210408153138-i6e6g5i"}

````````````````````````````````example
foo <a href="\*">
.
<p>foo <a href="\*"></p>
````````````````````````````````
{: id="20210408153138-ztjgokd"}

````````````````````````````````example
<a href="\"">
.
<p>&lt;a href=&quot;&quot;&quot;&gt;</p>
````````````````````````````````
{: id="20210408153138-mkj88y5"}

## Hard line breaks
{: id="20210408153138-irqxiox"}

A line break (not in a code span or HTML tag) that is preceded
by two or more spaces and does not occur at the end of a block
is parsed as a [hard line break](@) (rendered
in HTML as a `<br />` tag):
{: id="20210408153138-y84of2w"}

````````````````````````````````example
foo  
baz
.
<p>foo<br />
baz</p>
````````````````````````````````
{: id="20210408153138-m8hoax9"}

For a more visible alternative, a backslash before the
[line ending] may be used instead of two spaces:
{: id="20210408153138-a01ngje"}

````````````````````````````````example
foo\
baz
.
<p>foo<br />
baz</p>
````````````````````````````````
{: id="20210408153138-msdt7bs"}

More than two spaces can be used:
{: id="20210408153138-21zpp60"}

````````````````````````````````example
foo       
baz
.
<p>foo<br />
baz</p>
````````````````````````````````
{: id="20210408153138-1owncmr"}

Leading spaces at the beginning of the next line are ignored:
{: id="20210408153138-e380cbh"}

````````````````````````````````example
foo  
     bar
.
<p>foo<br />
bar</p>
````````````````````````````````
{: id="20210408153138-591nlt3"}

````````````````````````````````example
foo\
     bar
.
<p>foo<br />
bar</p>
````````````````````````````````
{: id="20210408153138-kfqrqhj"}

Line breaks can occur inside emphasis, links, and other constructs
that allow inline content:
{: id="20210408153138-0rvyxnn"}

````````````````````````````````example
*foo  
bar*
.
<p><em>foo<br />
bar</em></p>
````````````````````````````````
{: id="20210408153138-37au04b"}

````````````````````````````````example
*foo\
bar*
.
<p><em>foo<br />
bar</em></p>
````````````````````````````````
{: id="20210408153138-xj4bif0"}

Line breaks do not occur inside code spans
{: id="20210408153138-h0sxfg1"}

````````````````````````````````example
`code 
span`
.
<p><code>code  span</code></p>
````````````````````````````````
{: id="20210408153138-eax9dua"}

````````````````````````````````example
`code\
span`
.
<p><code>code\ span</code></p>
````````````````````````````````
{: id="20210408153138-dzwt2su"}

or HTML tags:
{: id="20210408153138-etyc3nn"}

````````````````````````````````example
<a href="foo  
bar">
.
<p><a href="foo  
bar"></p>
````````````````````````````````
{: id="20210408153138-plh0762"}

````````````````````````````````example
<a href="foo\
bar">
.
<p><a href="foo\
bar"></p>
````````````````````````````````
{: id="20210408153138-cmywaaz"}

Hard line breaks are for separating inline content within a block.
Neither syntax for hard line breaks works at the end of a paragraph or
other block element:
{: id="20210408153138-itxo4up"}

````````````````````````````````example
foo\
.
<p>foo\</p>
````````````````````````````````
{: id="20210408153138-s72c5b3"}

````````````````````````````````example
foo  
.
<p>foo</p>
````````````````````````````````
{: id="20210408153138-tpn0x9u"}

````````````````````````````````example
### foo\
.
<h3>foo\</h3>
````````````````````````````````
{: id="20210408153138-c4f81sr"}

````````````````````````````````example
### foo  
.
<h3>foo</h3>
````````````````````````````````
{: id="20210408153138-mp9x99e"}

## Soft line breaks
{: id="20210408153138-7yrm41i"}

A regular line break (not in a code span or HTML tag) that is not
preceded by two or more spaces or a backslash is parsed as a
[softbreak](@).  (A softbreak may be rendered in HTML either as a
[line ending] or as a space. The result will be the same in
browsers. In the examples here, a [line ending] will be used.)
{: id="20210408153138-3yt321q"}

````````````````````````````````example
foo
baz
.
<p>foo
baz</p>
````````````````````````````````
{: id="20210408153138-44yeq80"}

Spaces at the end of the line and beginning of the next line are
removed:
{: id="20210408153138-dgy51ve"}

````````````````````````````````example
foo 
 baz
.
<p>foo
baz</p>
````````````````````````````````
{: id="20210408153138-7tez8m4"}

A conforming parser may render a soft line break in HTML either as a
line break or as a space.
{: id="20210408153138-qpm6ts8"}

A renderer may also provide an option to render soft line breaks
as hard line breaks.
{: id="20210408153138-htyjcfy"}

## Textual content
{: id="20210408153138-l6ncunn"}

Any characters not given an interpretation by the above rules will
be parsed as plain textual content.
{: id="20210408153138-09fjebh"}

````````````````````````````````example
hello $.;'there
.
<p>hello $.;'there</p>
````````````````````````````````
{: id="20210408153138-7a8nebo"}

````````````````````````````````example
Foo χρῆν
.
<p>Foo χρῆν</p>
````````````````````````````````
{: id="20210408153138-4kkdjus"}

Internal spaces are preserved verbatim:
{: id="20210408153138-srq50m8"}

````````````````````````````````example
Multiple     spaces
.
<p>Multiple     spaces</p>
````````````````````````````````
{: id="20210408153138-fgbt82l"}

<!-- END TESTS -->

# Appendix: A parsing strategy
{: id="20210408153138-xyt3u27"}

In this appendix we describe some features of the parsing strategy
used in the CommonMark reference implementations.
{: id="20210408153138-socu7mo"}

## Overview
{: id="20210408153138-a069lr9"}

Parsing has two phases:
{: id="20210408153138-tire2vu"}

1. {: id="20210408153137-ux1o75d"}In the first phase, lines of input are consumed and the block
   structure of the document---its division into paragraphs, block quotes,
   list items, and so on---is constructed.  Text is assigned to these
   blocks but not parsed. Link reference definitions are parsed and a
   map of links is constructed.
   {: id="20210408153138-eaxa8g8"}
2. {: id="20210408153137-aqjkh0v"}In the second phase, the raw text contents of paragraphs and headings
   are parsed into sequences of Markdown inline elements (strings,
   code spans, links, emphasis, and so on), using the map of link
   references constructed in phase 1.
   {: id="20210408153138-yar6ip1"}
{: id="20210408153138-9g8fqfi"}

At each point in processing, the document is represented as a tree of
**blocks**.  The root of the tree is a `document` block.  The `document`
may have any number of other blocks as **children**.  These children
may, in turn, have other blocks as children.  The last child of a block
is normally considered **open**, meaning that subsequent lines of input
can alter its contents.  (Blocks that are not open are **closed**.)
Here, for example, is a possible document tree, with the open blocks
marked by arrows:
{: id="20210408153138-bn5x2cb"}

```tree
-> document
  -> block_quote
       paragraph
         "Lorem ipsum dolor\nsit amet."
    -> list (type=bullet tight=true bullet_char=-)
         list_item
           paragraph
             "Qui *quodsi iracundia*"
      -> list_item
        -> paragraph
             "aliquando id"
```
{: id="20210408153138-n3g00ae"}

## Phase 1: block structure
{: id="20210408153138-4c3ldel"}

Each line that is processed has an effect on this tree.  The line is
analyzed and, depending on its contents, the document may be altered
in one or more of the following ways:
{: id="20210408153138-1qvbleb"}

1. {: id="20210408153137-0u6sxk1"}One or more open blocks may be closed.
   {: id="20210408153138-2hatj57"}
2. {: id="20210408153137-woydadm"}One or more new blocks may be created as children of the
   last open block.
   {: id="20210408153138-vbts7gw"}
3. {: id="20210408153137-lwhmsg1"}Text may be added to the last (deepest) open block remaining
   on the tree.
   {: id="20210408153138-lfiawlb"}
{: id="20210408153138-h0wmt15"}

Once a line has been incorporated into the tree in this way,
it can be discarded, so input can be read in a stream.
{: id="20210408153138-k9tva38"}

For each line, we follow this procedure:
{: id="20210408153138-4elh5od"}

1. {: id="20210408153137-etge74f"}First we iterate through the open blocks, starting with the
   root document, and descending through last children down to the last
   open block.  Each block imposes a condition that the line must satisfy
   if the block is to remain open.  For example, a block quote requires a
   `>` character.  A paragraph requires a non-blank line.
   In this phase we may match all or just some of the open
   blocks.  But we cannot close unmatched blocks yet, because we may have a
   [lazy continuation line].
   {: id="20210408153138-q76o7iq"}
2. {: id="20210408153137-xk3e2od"}Next, after consuming the continuation markers for existing
   blocks, we look for new block starts (e.g. `>` for a block quote).
   If we encounter a new block start, we close any blocks unmatched
   in step 1 before creating the new block as a child of the last
   matched block.
   {: id="20210408153138-y2ixmlm"}
3. {: id="20210408153137-z8o4hb8"}Finally, we look at the remainder of the line (after block
   markers like `>`, list markers, and indentation have been consumed).
   This is text that can be incorporated into the last open
   block (a paragraph, code block, heading, or raw HTML).
   {: id="20210408153138-eqbrna5"}
{: id="20210408153138-st377o6"}

Setext headings are formed when we see a line of a paragraph
that is a [setext heading underline].
{: id="20210408153138-mx295x3"}

Reference link definitions are detected when a paragraph is closed;
the accumulated text lines are parsed to see if they begin with
one or more reference link definitions.  Any remainder becomes a
normal paragraph.
{: id="20210408153138-i5ji7w8"}

We can see how this works by considering how the tree above is
generated by four lines of Markdown:
{: id="20210408153138-o0vqj4t"}

```markdown
> Lorem ipsum dolor
sit amet.
> - Qui *quodsi iracundia*
> - aliquando id
```
{: id="20210408153138-tevzkr1"}

At the outset, our document model is just
{: id="20210408153138-5qr6w0c"}

```tree
-> document
```
{: id="20210408153138-6jyffin"}

The first line of our text,
{: id="20210408153138-chuvzc3"}

```markdown
> Lorem ipsum dolor
```
{: id="20210408153138-14hsclp"}

causes a `block_quote` block to be created as a child of our
open `document` block, and a `paragraph` block as a child of
the `block_quote`.  Then the text is added to the last open
block, the `paragraph`:
{: id="20210408153138-xssbius"}

```tree
-> document
  -> block_quote
    -> paragraph
         "Lorem ipsum dolor"
```
{: id="20210408153138-icpe88m"}

The next line,
{: id="20210408153138-uwf7l94"}

```markdown
sit amet.
```
{: id="20210408153138-0yy30ni"}

is a "lazy continuation" of the open `paragraph`, so it gets added
to the paragraph's text:
{: id="20210408153138-wa38b6m"}

```tree
-> document
  -> block_quote
    -> paragraph
         "Lorem ipsum dolor\nsit amet."
```
{: id="20210408153138-e9rrznl"}

The third line,
{: id="20210408153138-mu08uyz"}

```markdown
> - Qui *quodsi iracundia*
```
{: id="20210408153138-fywhsqo"}

causes the `paragraph` block to be closed, and a new `list` block
opened as a child of the `block_quote`.  A `list_item` is also
added as a child of the `list`, and a `paragraph` as a child of
the `list_item`.  The text is then added to the new `paragraph`:
{: id="20210408153138-de50x0t"}

```tree
-> document
  -> block_quote
       paragraph
         "Lorem ipsum dolor\nsit amet."
    -> list (type=bullet tight=true bullet_char=-)
      -> list_item
        -> paragraph
             "Qui *quodsi iracundia*"
```
{: id="20210408153138-lok6d24"}

The fourth line,
{: id="20210408153138-oimo3y0"}

```markdown
> - aliquando id
```
{: id="20210408153138-n01sjkc"}

causes the `list_item` (and its child the `paragraph`) to be closed,
and a new `list_item` opened up as child of the `list`.  A `paragraph`
is added as a child of the new `list_item`, to contain the text.
We thus obtain the final tree:
{: id="20210408153138-u1p8maj"}

```tree
-> document
  -> block_quote
       paragraph
         "Lorem ipsum dolor\nsit amet."
    -> list (type=bullet tight=true bullet_char=-)
         list_item
           paragraph
             "Qui *quodsi iracundia*"
      -> list_item
        -> paragraph
             "aliquando id"
```
{: id="20210408153138-27o8wsk"}

## Phase 2: inline structure
{: id="20210408153138-andgagc"}

Once all of the input has been parsed, all open blocks are closed.
{: id="20210408153138-5htxgsh"}

We then "walk the tree," visiting every node, and parse raw
string contents of paragraphs and headings as inlines.  At this
point we have seen all the link reference definitions, so we can
resolve reference links as we go.
{: id="20210408153138-iwerikv"}

```tree
document
  block_quote
    paragraph
      str "Lorem ipsum dolor"
      softbreak
      str "sit amet."
    list (type=bullet tight=true bullet_char=-)
      list_item
        paragraph
          str "Qui "
          emph
            str "quodsi iracundia"
      list_item
        paragraph
          str "aliquando id"
```
{: id="20210408153138-p06qv18"}

Notice how the [line ending] in the first paragraph has
been parsed as a `softbreak`, and the asterisks in the first list item
have become an `emph`.
{: id="20210408153138-70mv8gd"}

### An algorithm for parsing nested emphasis and links
{: id="20210408153138-ezhcgrc"}

By far the trickiest part of inline parsing is handling emphasis,
strong emphasis, links, and images.  This is done using the following
algorithm.
{: id="20210408153138-qp9gh38"}

When we're parsing inlines and we hit either
{: id="20210408153138-xn51gqj"}

- {: id="20210408153137-h72mtvy"}a run of `*` or `_` characters, or
  {: id="20210408153138-iqtfmw3"}
- {: id="20210408153137-uvitpus"}a `[` or `![`
  {: id="20210408153138-t84z0d9"}
{: id="20210408153138-dbf15yl"}

we insert a text node with these symbols as its literal content, and we
add a pointer to this text node to the [delimiter stack](@).
{: id="20210408153138-n97x8bz"}

The [delimiter stack] is a doubly linked list.  Each
element contains a pointer to a text node, plus information about
{: id="20210408153138-990ubvs"}

- {: id="20210408153137-yz0h7j6"}the type of delimiter (`[`, `![`, `*`, `_`)
  {: id="20210408153138-4doqve5"}
- {: id="20210408153137-c2vhokf"}the number of delimiters,
  {: id="20210408153138-6v9hwc8"}
- {: id="20210408153137-strh7by"}whether the delimiter is "active" (all are active to start), and
  {: id="20210408153138-an8zw0n"}
- {: id="20210408153137-c321cf9"}whether the delimiter is a potential opener, a potential closer,
  or both (which depends on what sort of characters precede
  and follow the delimiters).
  {: id="20210408153138-825niqg"}
{: id="20210408153138-vmlj59f"}

When we hit a `]` character, we call the *look for link or image*
procedure (see below).
{: id="20210408153138-y6qsozh"}

When we hit the end of the input, we call the *process emphasis*
procedure (see below), with `stack_bottom` = NULL.
{: id="20210408153138-yhjbpsj"}

#### *look for link or image*
{: id="20210408153138-sig5cbi"}

Starting at the top of the delimiter stack, we look backwards
through the stack for an opening `[` or `![` delimiter.
{: id="20210408153138-tqv6x3h"}

- {: id="20210408153137-co4sp2b"}If we don't find one, we return a literal text node `]`.
  {: id="20210408153138-1j49s7v"}
- {: id="20210408153137-gqp37kz"}If we do find one, but it's not *active*, we remove the inactive
  delimiter from the stack, and return a literal text node `]`.
  {: id="20210408153138-pbup4n2"}
- {: id="20210408153137-4q7690u"}If we find one and it's active, then we parse ahead to see if
  we have an inline link/image, reference link/image, compact reference
  link/image, or shortcut reference link/image.
  {: id="20210408153138-1ftro2r"}
  + {: id="20210408153137-a5m16ar"}If we don't, then we remove the opening delimiter from the
    delimiter stack and return a literal text node `]`.
    {: id="20210408153138-2frpetj"}
  + {: id="20210408153137-9jo89ae"}If we do, then
    {: id="20210408153138-tolil6d"}
    * {: id="20210408153137-zds0o4x"}We return a link or image node whose children are the inlines
      after the text node pointed to by the opening delimiter.
      {: id="20210408153138-td774lp"}
    * {: id="20210408153137-ctjw81x"}We run *process emphasis* on these inlines, with the `[` opener
      as `stack_bottom`.
      {: id="20210408153138-j6u0373"}
    * {: id="20210408153137-6noygii"}We remove the opening delimiter.
      {: id="20210408153138-e8488lj"}
    * {: id="20210408153137-xcw86aw"}If we have a link (and not an image), we also set all
      `[` delimiters before the opening delimiter to *inactive*.  (This
      will prevent us from getting links within links.)
      {: id="20210408153138-iy1y34g"}
    {: id="20210408153138-i8u68ii"}
  {: id="20210408153138-lqkzo7j"}
{: id="20210408153138-92ryz2m"}

#### *process emphasis*
{: id="20210408153138-qkssyf4"}

Parameter `stack_bottom` sets a lower bound to how far we
descend in the [delimiter stack].  If it is NULL, we can
go all the way to the bottom.  Otherwise, we stop before
visiting `stack_bottom`.
{: id="20210408153138-9nlgxfw"}

Let `current_position` point to the element on the [delimiter stack]
just above `stack_bottom` (or the first element if `stack_bottom`
is NULL).
{: id="20210408153138-0wlayrb"}

We keep track of the `openers_bottom` for each delimiter
type (`*`, `_`) and each length of the closing delimiter run
(modulo 3).  Initialize this to `stack_bottom`.
{: id="20210408153138-yaseq2v"}

Then we repeat the following until we run out of potential
closers:
{: id="20210408153138-vq8lqz0"}

- {: id="20210408153137-pcfv7wb"}Move `current_position` forward in the delimiter stack (if needed)
  until we find the first potential closer with delimiter `*` or `_`.
  (This will be the potential closer closest
  to the beginning of the input -- the first one in parse order.)
  {: id="20210408153138-kvyfbns"}
- {: id="20210408153137-kripx1a"}Now, look back in the stack (staying above `stack_bottom` and
  the `openers_bottom` for this delimiter type) for the
  first matching potential opener ("matching" means same delimiter).
  {: id="20210408153138-tf4s1jt"}
- {: id="20210408153137-ta8s4k4"}If one is found:
  {: id="20210408153138-5vc2unu"}
  + {: id="20210408153137-9u38z2s"}Figure out whether we have emphasis or strong emphasis:
    if both closer and opener spans have length >= 2, we have
    strong, otherwise regular.
    {: id="20210408153138-xsi5upq"}
  + {: id="20210408153137-bqogwp5"}Insert an emph or strong emph node accordingly, after
    the text node corresponding to the opener.
    {: id="20210408153138-s968w3y"}
  + {: id="20210408153137-5h0xvxu"}Remove any delimiters between the opener and closer from
    the delimiter stack.
    {: id="20210408153138-92hyvpr"}
  + {: id="20210408153137-h63uszz"}Remove 1 (for regular emph) or 2 (for strong emph) delimiters
    from the opening and closing text nodes.  If they become empty
    as a result, remove them and remove the corresponding element
    of the delimiter stack.  If the closing node is removed, reset
    `current_position` to the next element in the stack.
    {: id="20210408153138-vfn36yr"}
  {: id="20210408153138-801o0k2"}
- {: id="20210408153137-ujmu51r"}If none is found:
  {: id="20210408153138-gvqkbh4"}
  + {: id="20210408153137-c26ub9f"}Set `openers_bottom` to the element before `current_position`.
    (We know that there are no openers for this kind of closer up to and
    including this point, so this puts a lower bound on future searches.)
    {: id="20210408153138-55gtgds"}
  + {: id="20210408153137-33y74qy"}If the closer at `current_position` is not a potential opener,
    remove it from the delimiter stack (since we know it can't
    be a closer either).
    {: id="20210408153138-lxswc82"}
  + {: id="20210408153137-mxn5pc3"}Advance `current_position` to the next element in the stack.
    {: id="20210408153138-hdczp4i"}
  {: id="20210408153138-g92mf19"}
{: id="20210408153138-fwihidt"}

After we're done, we remove all delimiters above `stack_bottom` from the
delimiter stack.
{: id="20210408153138-yprun6m"}


{: id="20210408153138-3t979ww" type="doc"}
