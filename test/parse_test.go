// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"os"
	"testing"

	"github.com/88250/lute"
)

type parseTest struct {
	name string
	from string
	to   string
}

var parseTests = []parseTest{
	//
	//	{"link ref node", "[foo]\n\n[foo]: bar", "<p><a href=\"bar\">foo</a></p>\n"},
	//	{"details", `<details>
	//<summary>foo</summary>
	//
	//* bar
	//
	//</details>`, "<details>\n<summary>foo</summary>\n<ul>\n<li>bar</li>\n</ul>\n</details>\n"},
	//
	//	// 链接引用定义问题 https://github.com/88250/lute/issues/3
	//	{"#3", "[foo][]bar\n\n[foo]: /url \"title\"\n", "<p><a href=\"/url\" title=\"title\">foo</a>bar</p>\n"},
	//
	//	{"#177", "[link](/u(ri\n)\n", "<p>[link](/u(ri\n)</p>\n"},
	//	{"#177", "[link](/u(ri )\n", "<p>[link](/u(ri )</p>\n"},

	// commonmark spec inline-level cases

	// 0.31
	{"0.31-spec626", "foo <!--> foo -->\n\nfoo <!---> foo -->\n", "<p>foo <!--> foo --&gt;</p>\n<p>foo <!---> foo --&gt;</p>\n"},
	{"0.31-spec625", "foo <!-- this is a --\ncomment - with hyphens -->\n", "<p>foo <!-- this is a --\ncomment - with hyphens --></p>\n"},

	// Below 0.31
	{"spec637", "`code \nspan`\n", "<p><code>code  span</code></p>\n"},
	{"spec626", "foo <![CDATA[>&<]]>\n", "<p>foo <![CDATA[>&<]]></p>\n"},
	{"spec625", "foo <!ELEMENT br EMPTY>\n", "<p>foo <!ELEMENT br EMPTY></p>\n"},
	{"spec624", "foo <?php echo $a; ?>\n", "<p>foo <?php echo $a; ?></p>\n"},
	{"spec621", "foo <!-- this is a\ncomment - with hyphen -->\n", "<p>foo <!-- this is a\ncomment - with hyphen --></p>\n"},
	{"spec616", "<a href=\"hi'> <a href=hi'>\n", "<p>&lt;a href=&quot;hi'&gt; &lt;a href=hi'&gt;</p>\n"},
	{"spec613", "Foo <responsive-image src=\"foo.jpg\" />\n", "<p>Foo <responsive-image src=\"foo.jpg\" /></p>\n"},
	{"spec612", "<a foo=\"bar\" bam = 'baz <em>\"</em>'\n_boolean zoop:33=zoop:33 />\n", "<p><a foo=\"bar\" bam = 'baz <em>\"</em>'\n_boolean zoop:33=zoop:33 /></p>\n"},
	{"spec611", "<a  /><b2\ndata=\"foo\" >\n", "<p><a  /><b2\ndata=\"foo\" ></p>\n"},
	{"spec604", "< http://foo.bar >\n", "<p>&lt; http://foo.bar &gt;</p>\n"},
	{"spec600", "<foo@bar.example.com>\n", "<p><a href=\"mailto:foo@bar.example.com\">foo@bar.example.com</a></p>\n"},
	{"spec594", "<a+b+c:d>\n", "<p><a href=\"a+b+c:d\">a+b+c:d</a></p>\n"},
	{"spec592", "<irc://foo.bar:2233/baz>\n", "<p><a href=\"irc://foo.bar:2233/baz\">irc://foo.bar:2233/baz</a></p>\n"},
	{"spec591", "<http://foo.bar.baz/test?q=hello&id=22&boolean>\n", "<p><a href=\"http://foo.bar.baz/test?q=hello&amp;id=22&amp;boolean\">http://foo.bar.baz/test?q=hello&amp;id=22&amp;boolean</a></p>\n"},
	{"spec588", "!\\[foo]\n\n[foo]: /url \"title\"\n", "<p>![foo]</p>\n"},
	{"spec569", "![foo *bar*]\n\n[foo *bar*]: train.jpg \"train & tracks\"\n", "<p><img src=\"train.jpg\" alt=\"foo bar\" title=\"train &amp; tracks\" /></p>\n"},
	{"spec568", "![foo](/url \"title\")\n", "<p><img src=\"/url\" alt=\"foo\" title=\"title\" /></p>\n"},
	{"spec552", "[foo] \n[]\n\n[foo]: /url \"title\"\n", "<p><a href=\"/url\" title=\"title\">foo</a>\n[]</p>\n"},
	{"spec549", "[foo][]\n\n[foo]: /url \"title\"\n", "<p><a href=\"/url\" title=\"title\">foo</a></p>\n"},
	{"spec548", "[\n ]\n\n[\n ]: /uri\n", "<p>[\n]</p>\n<p>[\n]: /uri</p>\n"},
	{"spec547", "[]\n\n[]: /uri\n", "<p>[]</p>\n<p>[]: /uri</p>\n"},
	{"spec537", "[Foo\n  bar]: /url\n\n[Baz][Foo bar]\n", "<p><a href=\"/url\">Baz</a></p>\n"},
	{"spec536", "[Толпой][Толпой] is a Russian word.\n\n[ТОЛПОЙ]: /url\n", "<p><a href=\"/url\">Толпой</a> is a Russian word.</p>\n"},
	{"spec529", "[foo *bar [baz][ref]*][ref]\n\n[ref]: /uri\n", "<p>[foo <em>bar <a href=\"/uri\">baz</a></em>]<a href=\"/uri\">ref</a></p>\n"},
	{"spec523", "[foo][bar]\n\n[bar]: /url \"title\"\n", "<p><a href=\"/url\" title=\"title\">foo</a></p>\n"},
	{"spec520", "[foo <bar attr=\"](baz)\">\n", "<p>[foo <bar attr=\"](baz)\"></p>\n"},
	{"spec516", "![[[foo](uri1)](uri2)](uri3)\n", "<p><img src=\"uri3\" alt=\"[foo](uri2)\" /></p>\n"},
	{"spec513", "[![moon](moon.jpg)](/uri)\n", "<p><a href=\"/uri\"><img src=\"moon.jpg\" alt=\"moon\" /></a></p>\n"},
	{"spec512", "[link *foo **bar** `#`*](/uri)\n", "<p><a href=\"/uri\">link <em>foo <strong>bar</strong> <code>#</code></em></a></p>\n"},
	{"spec506", "[link](   /uri\n  \"title\"  )\n", "<p><a href=\"/uri\" title=\"title\">link</a></p>\n"},
	{"spec503", "[link](/url \"title\")\n", "<p><a href=\"/url%C2%A0%22title%22\">link</a></p>\n"},
	{"spec501", "[link](/url \"title\")\n[link](/url 'title')\n[link](/url (title))\n", "<p><a href=\"/url\" title=\"title\">link</a>\n<a href=\"/url\" title=\"title\">link</a>\n<a href=\"/url\" title=\"title\">link</a></p>\n"},
	{"spec491", "[a](<b)c\n[a](<b)c>\n[a](<b>c)\n", "<p>[a](&lt;b)c\n[a](&lt;b)c&gt;\n[a](<b>c)</p>\n"},
	{"spec490", "[link](<foo\\>)\n", "<p>[link](&lt;foo&gt;)</p>\n"},
	{"spec488", "[link](<foo\nbar>)\n", "<p>[link](<foo\nbar>)</p>\n"},
	{"spec484", "[link](<>)\n", "<p><a href=\"\">link</a></p>\n"},
	{"spec483", "[link]()\n", "<p><a href=\"\">link</a></p>\n"},
	{"spec479", "**a<http://foo.bar/?q=**>\n", "<p>**a<a href=\"http://foo.bar/?q=**\">http://foo.bar/?q=**</a></p>\n"},
	{"spec416", "foo******bar*********baz\n", "<p>foo<strong><strong><strong>bar</strong></strong></strong>***baz</p>\n"},
	{"spec408", "__foo_ bar_\n", "<p><em><em>foo</em> bar</em></p>\n"},
	{"spec403", "*foo [bar](/url)*\n", "<p><em>foo <a href=\"/url\">bar</a></em></p>\n"},
	{"spec353", "* a *\n", "<p>* a *</p>\n"},
	{"spec352", "a*\"foo\"*\n", "<p>a*&quot;foo&quot;*</p>\n"},
	{"spec333", "` b `\n", "<p><code> b </code></p>\n"},
	{"spec330", "` `` `\n", "<p><code>``</code></p>\n"},
	{"spec314", "&nbsp &x; &#; &#x;\n&#987654321;\n&#abcdef0;\n&ThisIsNotDefined; &hi?;\n", "<p>&amp;nbsp &amp;x; &amp;#; &amp;#x;\n&amp;#987654321;\n&amp;#abcdef0;\n&amp;ThisIsNotDefined; &amp;hi?;</p>\n"},
	{"spec312", "&#35; &#1234; &#992; &#0;\n", "<p># Ӓ Ϡ �</p>\n"},
	{"spec311", "&nbsp; &amp; &copy; &AElig; &Dcaron;\n&frac34; &HilbertSpace; &DifferentialD;\n&ClockwiseContourIntegral; &ngE;\n", "<p>  &amp; © Æ Ď\n¾ ℋ ⅆ\n∲ ≧̸</p>\n"},
	{"spec310", "``` foo\\+bar\nfoo\n```\n", "<pre><code class=\"language-foo+bar\">foo\n</code></pre>\n"},
	{"spec308", "[foo](/bar\\* \"ti\\*tle\")\n", "<p><a href=\"/bar*\" title=\"ti*tle\">foo</a></p>\n"},
	{"spec306", "<http://example.com?find=\\*>\n", "<p><a href=\"http://example.com?find=%5C*\">http://example.com?find=\\*</a></p>\n"},

	// commonmark spec block-level cases

	{"spec295", "* foo\n  * bar\n\n  baz\n", "<ul>\n<li>\n<p>foo</p>\n<ul>\n<li>bar</li>\n</ul>\n<p>baz</p>\n</li>\n</ul>\n"},
	{"spec289", "- a\n  - b\n\n    c\n- d\n", "<ul>\n<li>a\n<ul>\n<li>\n<p>b</p>\n<p>c</p>\n</li>\n</ul>\n</li>\n<li>d</li>\n</ul>\n"},
	{"spec287", "- a\n- b\n\n  [ref]: /url\n- d\n", "<ul>\n<li>\n<p>a</p>\n</li>\n<li>\n<p>b</p>\n</li>\n<li>\n<p>d</p>\n</li>\n</ul>\n"},
	{"spec285", "* a\n*\n\n* c\n", "<ul>\n<li>\n<p>a</p>\n</li>\n<li></li>\n<li>\n<p>c</p>\n</li>\n</ul>\n"},
	{"spec283", "1. a\n\n  2. b\n\n    3. c\n", "<ol>\n<li>\n<p>a</p>\n</li>\n<li>\n<p>b</p>\n</li>\n</ol>\n<pre><code>3. c\n</code></pre>\n"},
	{"spec282", "- a\n - b\n  - c\n   - d\n    - e\n", "<ul>\n<li>a</li>\n<li>b</li>\n<li>c</li>\n<li>d\n- e</li>\n</ul>\n"},
	{"spec278", "- foo\n- bar\n\n<!-- -->\n\n- baz\n- bim\n", "<ul>\n<li>foo</li>\n<li>bar</li>\n</ul>\n<!-- -->\n<ul>\n<li>baz</li>\n<li>bim</li>\n</ul>\n"},
	{"spec277", "- foo\n  - bar\n    - baz\n\n\n      bim\n", "<ul>\n<li>foo\n<ul>\n<li>bar\n<ul>\n<li>\n<p>baz</p>\n<p>bim</p>\n</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>\n"},
	{"spec276", "- foo\n\n- bar\n\n\n- baz\n", "<ul>\n<li>\n<p>foo</p>\n</li>\n<li>\n<p>bar</p>\n</li>\n<li>\n<p>baz</p>\n</li>\n</ul>\n"},
	{"spec275", "The number of windows in my house is\n1.  The number of doors is 6.\n", "<p>The number of windows in my house is</p>\n<ol>\n<li>The number of doors is 6.</li>\n</ol>\n"},
	{"spec274", "The number of windows in my house is\n14.  The number of doors is 6.\n", "<p>The number of windows in my house is\n14.  The number of doors is 6.</p>\n"},
	{"spec272", "1. foo\n2. bar\n3) baz\n", "<ol>\n<li>foo</li>\n<li>bar</li>\n</ol>\n<ol start=\"3\">\n<li>baz</li>\n</ol>\n"},
	{"spec270", "- # Foo\n- Bar\n  ---\n  baz\n", "<ul>\n<li>\n<h1>Foo</h1>\n</li>\n<li>\n<h2>Bar</h2>\nbaz</li>\n</ul>\n"},
	{"spec265", "- foo\n - bar\n  - baz\n   - boo\n", "<ul>\n<li>foo</li>\n<li>bar</li>\n<li>baz</li>\n<li>boo</li>\n</ul>\n"},
	{"spec262", "> 1. > Blockquote\ncontinued here.\n", "<blockquote>\n<ol>\n<li>\n<blockquote>\n<p>Blockquote\ncontinued here.</p>\n</blockquote>\n</li>\n</ol>\n</blockquote>\n"},
	{"spec255", "foo\n*\n\nfoo\n1.\n", "<p>foo\n*</p>\n<p>foo\n1.</p>\n"},
	{"spec253", "1. foo\n2.\n3. bar\n", "<ol>\n<li>foo</li>\n<li></li>\n<li>bar</li>\n</ol>\n"},
	{"spec251", "- foo\n-\n- bar\n", "<ul>\n<li>foo</li>\n<li></li>\n<li>bar</li>\n</ul>\n"},
	{"spec250", "-\n\n  foo\n", "<ul>\n<li></li>\n</ul>\n<p>foo</p>\n"},
	{"spec249", "-   \n  foo\n", "<ul>\n<li>foo</li>\n</ul>\n"},
	{"spec248", "-\n  foo\n-\n  ```\n  bar\n  ```\n-\n      baz\n", "<ul>\n<li>foo</li>\n<li>\n<pre><code>bar\n</code></pre>\n</li>\n<li>\n<pre><code>baz\n</code></pre>\n</li>\n</ul>\n"},
	{"spec243", "1.     indented code\n\n   paragraph\n\n       more code\n", "<ol>\n<li>\n<pre><code>indented code\n</code></pre>\n<p>paragraph</p>\n<pre><code>more code\n</code></pre>\n</li>\n</ol>\n"},
	{"spec241", "  10.  foo\n\n           bar\n", "<ol start=\"10\">\n<li>\n<p>foo</p>\n<pre><code>bar\n</code></pre>\n</li>\n</ol>\n"},
	{"spec235", "123456789. ok\n", "<ol start=\"123456789\">\n<li>ok</li>\n</ol>\n"},
	{"spec234", "- Foo\n\n      bar\n\n\n      baz\n", "<ul>\n<li>\n<p>Foo</p>\n<pre><code>bar\n\n\nbaz\n</code></pre>\n</li>\n</ul>\n"},
	{"spec233", "1.  foo\n\n    ```\n    bar\n    ```\n\n    baz\n\n    > bam\n", "<ol>\n<li>\n<p>foo</p>\n<pre><code>bar\n</code></pre>\n<p>baz</p>\n<blockquote>\n<p>bam</p>\n</blockquote>\n</li>\n</ol>\n"},
	{"spec229", "   > > 1.  one\n>>\n>>     two\n", "<blockquote>\n<blockquote>\n<ol>\n<li>\n<p>one</p>\n<p>two</p>\n</li>\n</ol>\n</blockquote>\n</blockquote>\n"},
	{"spec227", " -    one\n\n     two\n", "<ul>\n<li>one</li>\n</ul>\n<pre><code> two\n</code></pre>\n"},
	{"spec224", "1.  A paragraph\n    with two lines.\n\n        indented code\n\n    > A block quote.\n", "<ol>\n<li>\n<p>A paragraph\nwith two lines.</p>\n<pre><code>indented code\n</code></pre>\n<blockquote>\n<p>A block quote.</p>\n</blockquote>\n</li>\n</ol>\n"},
	{"spec222", ">     code\n\n>    not code\n", "<blockquote>\n<pre><code>code\n</code></pre>\n</blockquote>\n<blockquote>\n<p>not code</p>\n</blockquote>\n"},
	{"spec221", ">>> foo\n> bar\n>>baz\n", "<blockquote>\n<blockquote>\n<blockquote>\n<p>foo\nbar\nbaz</p>\n</blockquote>\n</blockquote>\n</blockquote>\n"},
	{"spec219", "> bar\n>\nbaz\n", "<blockquote>\n<p>bar</p>\n</blockquote>\n<p>baz</p>\n"},
	{"spec214", "> foo\n>\n> bar\n", "<blockquote>\n<p>foo</p>\n<p>bar</p>\n</blockquote>\n"},
	{"spec211", ">\n> foo\n>  \n", "<blockquote>\n<p>foo</p>\n</blockquote>\n"},
	{"spec210", ">\n>  \n> \n", "<blockquote>\n</blockquote>\n"},
	{"spec207", "> ```\nfoo\n```\n", "<blockquote>\n<pre><code></code></pre>\n</blockquote>\n<p>foo</p>\n<pre><code></code></pre>\n"},
	{"spec206", ">     foo\n    bar\n", "<blockquote>\n<pre><code>foo\n</code></pre>\n</blockquote>\n<pre><code>bar\n</code></pre>\n"},
	{"spec205", "> - foo\n- bar\n", "<blockquote>\n<ul>\n<li>foo</li>\n</ul>\n</blockquote>\n<ul>\n<li>bar</li>\n</ul>\n"},
	{"spec200", "   > # Foo\n   > bar\n > baz\n", "<blockquote>\n<h1>Foo</h1>\n<p>bar\nbaz</p>\n</blockquote>\n"},
	{"spec198", "> # Foo\n> bar\n> baz\n", "<blockquote>\n<h1>Foo</h1>\n<p>bar\nbaz</p>\n</blockquote>\n"},
	{"spec197", "  \n\naaa\n  \n\n# aaa\n\n  \n", "<p>aaa</p>\n<h1>aaa</h1>\n"},
	{"spec196", "aaa     \nbbb     \n", "<p>aaa<br />\nbbb</p>\n"},
	{"spec187", "[foo]\n\n> [foo]: /url\n", "<p><a href=\"/url\">foo</a></p>\n<blockquote>\n</blockquote>\n"},
	{"spec186", "[foo]: /foo-url \"foo\"\n[bar]: /bar-url\n  \"bar\"\n[baz]: /baz-url\n\n[foo],\n[bar],\n[baz]\n", "<p><a href=\"/foo-url\" title=\"foo\">foo</a>,\n<a href=\"/bar-url\" title=\"bar\">bar</a>,\n<a href=\"/baz-url\">baz</a></p>\n"},
	{"spec185", "[foo]: /url\n===\n[foo]\n", "<p>===\n<a href=\"/url\">foo</a></p>\n"},
	{"spec183", "# [Foo]\n[foo]: /url\n> bar\n", "<h1><a href=\"/url\">Foo</a></h1>\n<blockquote>\n<p>bar</p>\n</blockquote>\n"},
	{"spec179", "[foo]: /url\n\"title\" ok\n", "<p>&quot;title&quot; ok</p>\n"},
	{"spec178", "[foo]: /url \"title\" ok\n", "<p>[foo]: /url &quot;title&quot; ok</p>\n"},
	{"spec177", "[\nfoo\n]: /url\nbar\n", "<p>bar</p>\n"},
	{"spec175", "[ΑΓΩ]: /φου\n\n[αγω]\n", "<p><a href=\"/%CF%86%CE%BF%CF%85\">αγω</a></p>\n"},
	{"spec173", "[foo]\n\n[foo]: first\n[foo]: second\n", "<p><a href=\"first\">foo</a></p>\n"},
	{"spec171", "[foo]: /url\\bar\\*baz \"foo\\\"bar\\baz\"\n\n[foo]\n", "<p><a href=\"/url%5Cbar*baz\" title=\"foo&quot;bar\\baz\">foo</a></p>\n"},
	{"spec170", "[foo]: <bar>(baz)\n\n[foo]\n", "<p>[foo]: <bar>(baz)</p>\n<p>[foo]</p>\n"},
	{"spec168", "[foo]:\n\n[foo]\n", "<p>[foo]:</p>\n<p>[foo]</p>\n"},
	{"spec167", "[foo]:\n/url\n\n[foo]\n", "<p><a href=\"/url\">foo</a></p>\n"},
	{"spec166", "[foo]: /url 'title\n\nwith blank line'\n\n[foo]\n", "<p>[foo]: /url 'title</p>\n<p>with blank line'</p>\n<p>[foo]</p>\n"},
	{"spec165", "[foo]: /url '\ntitle\nline1\nline2\n'\n\n[foo]\n", "<p><a href=\"/url\" title=\"\ntitle\nline1\nline2\n\">foo</a></p>\n"},
	{"spec164", "[Foo bar]:\n<my url>\n'title'\n\n[Foo bar]\n", "<p><a href=\"my%20url\" title=\"title\">Foo bar</a></p>\n"},
	{"spec163", "[Foo*bar\\]]:my_(url) 'title (with parens)'\n\n[Foo*bar\\]]\n", "<p><a href=\"my_(url)\" title=\"title (with parens)\">Foo*bar]</a></p>\n"},
	{"spec162", "   [foo]: \n      /url  \n           'the title'  \n\n[foo]\n", "<p><a href=\"/url\" title=\"the title\">foo</a></p>\n"},
	{"spec161", "[foo]: /url \"title\"\n\n[foo]\n", "<p><a href=\"/url\" title=\"title\">foo</a></p>\n"},
	{"spec160", "<table>\n\n  <tr>\n\n    <td>\n      Hi\n    </td>\n\n  </tr>\n\n</table>\n", "<table>\n  <tr>\n<pre><code>&lt;td&gt;\n  Hi\n&lt;/td&gt;\n</code></pre>\n  </tr>\n</table>\n"},
	{"spec156", "Foo\n<a href=\"bar\">\nbaz\n", "<p>Foo\n<a href=\"bar\">\nbaz</p>\n"},
	{"spec151", "<![CDATA[\nfunction matchwo(a,b)\n{\n  if (a < b && a < 0) then {\n    return 1;\n\n  } else {\n\n    return 0;\n  }\n}\n]]>\nokay\n", "<![CDATA[\nfunction matchwo(a,b)\n{\n  if (a < b && a < 0) then {\n    return 1;\n\n  } else {\n\n    return 0;\n  }\n}\n]]>\n<p>okay</p>\n"},
	{"spec149", "<?php\n\n  echo '>';\n\n?>\nokay\n", "<?php\n\n  echo '>';\n\n?>\n<p>okay</p>\n"},
	{"spec148", "<!-- Foo\n\nbar\n   baz -->\nokay\n", "<!-- Foo\n\nbar\n   baz -->\n<p>okay</p>\n"},
	{"spec146", "<!-- foo -->*bar*\n*baz*\n", "<!-- foo -->*bar*\n<p><em>baz</em></p>\n"},
	{"spec145", "<style>p{color:red;}</style>\n*foo*\n", "<style>p{color:red;}</style>\n<p><em>foo</em></p>\n"},
	{"spec144", "- <div>\n- foo\n", "<ul>\n<li>\n<div>\n</li>\n<li>foo</li>\n</ul>\n"},
	{"spec143", "> <div>\n> foo\n\nbar\n", "<blockquote>\n<div>\nfoo\n</blockquote>\n<p>bar</p>\n"},
	{"spec142", "<style\n  type=\"text/css\">\n\nfoo\n", "<style\n  type=\"text/css\">\n\nfoo\n"},
	{"spec141", "<style\n  type=\"text/css\">\nh1 {color:red;}\n\np {color:blue;}\n</style>\nokay\n", "<style\n  type=\"text/css\">\nh1 {color:red;}\n\np {color:blue;}\n</style>\n<p>okay</p>\n"},
	{"spec139", "<pre language=\"haskell\"><code>\nimport Text.HTML.TagSoup\n\nmain :: IO ()\nmain = print $ parseTags tags\n</code></pre>\nokay\n", "<pre language=\"haskell\"><code>\nimport Text.HTML.TagSoup\n\nmain :: IO ()\nmain = print $ parseTags tags\n</code></pre>\n<p>okay</p>\n"},
	{"spec135", "</ins>\n*bar*\n", "</ins>\n*bar*\n"},
	{"spec132", "<a href=\"foo\">\n*bar*\n</a>\n", "<a href=\"foo\">\n*bar*\n</a>\n"},
	{"spec120", " <div>\n  *hello*\n         <foo><a>\n", " <div>\n  *hello*\n         <foo><a>\n"},
	{"spec118", "<table><tr><td>\n<pre>\n**Hello**,\n\n_world_.\n</pre>\n</td></tr></table>\n", "<table><tr><td>\n<pre>\n**Hello**,\n<p><em>world</em>.\n</pre></p>\n</td></tr></table>\n"},
	{"spec117", "```\n``` aaa\n```\n", "<pre><code>``` aaa\n</code></pre>\n"},
	{"spec116", "~~~ aa ``` ~~~\nfoo\n~~~\n", "<pre><code class=\"language-aa\">foo\n</code></pre>\n"},
	{"spec115", "``` aa ```\nfoo\n", "<p><code>aa</code>\nfoo</p>\n"},
	{"spec113", "~~~~    ruby startline=3 $%@#$\ndef foo(x)\n  return 3\nend\n~~~~~~~\n", "<pre><code class=\"language-ruby\">def foo(x)\n  return 3\nend\n</code></pre>\n"},
	{"spec112", "```ruby\ndef foo(x)\n  return 3\nend\n```\n", "<pre><code class=\"language-ruby\">def foo(x)\n  return 3\nend\n</code></pre>\n"},
	{"spec110", "foo\n```\nbar\n```\nbaz\n", "<p>foo</p>\n<pre><code>bar\n</code></pre>\n<p>baz</p>\n"},
	{"spec108", "``` ```\naaa\n", "<p><code> </code>\naaa</p>\n"},
	{"spec103", "   ```\n   aaa\n    aaa\n  aaa\n   ```\n", "<pre><code>aaa\n aaa\naaa\n</code></pre>\n"},
	{"spec101", " ```\n aaa\naaa\n```\n", "<pre><code>aaa\naaa\n</code></pre>\n"},
	{"spec100", "```\n```\n", "<pre><code></code></pre>\n"},
	{"spec98", "> ```\n> aaa\n\nbbb\n", "<blockquote>\n<pre><code>aaa\n</code></pre>\n</blockquote>\n<p>bbb</p>\n"},
	{"spec97", "`````\n\n```\naaa\n", "<pre><code>\n```\naaa\n</code></pre>\n"},
	{"spec96", "```\n", "<pre><code></code></pre>\n"},
	{"spec91", "``\nfoo\n``\n", "<p><code>foo</code></p>\n"},
	{"spec90", "~~~\n<\n >\n~~~\n", "<pre><code>&lt;\n &gt;\n</code></pre>\n"},
	{"spec89", "```\n<\n >\n```\n", "<pre><code>&lt;\n &gt;\n</code></pre>\n"},
	{"spec88", "    foo  \n", "<pre><code>foo  \n</code></pre>\n"},
	{"spec87", "\n    \n    foo\n    \n\n", "<pre><code>foo\n</code></pre>\n"},
	{"spec81", "    chunk1\n\n    chunk2\n  \n \n \n    chunk3\n", "<pre><code>chunk1\n\nchunk2\n\n\n\nchunk3\n</code></pre>\n"},
	{"spec79", "1.  foo\n\n    - bar\n", "<ol>\n<li>\n<p>foo</p>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ol>\n"},
	{"spec77", "    a simple\n      indented code block\n", "<pre><code>a simple\n  indented code block\n</code></pre>\n"},
	{"spec72", "\\> foo\n------\n", "<h2>&gt; foo</h2>\n"},
	{"spec63", "> foo\nbar\n===\n", "<blockquote>\n<p>foo\nbar\n===</p>\n</blockquote>\n"},
	{"spec62", "> Foo\n---\n", "<blockquote>\n<p>Foo</p>\n</blockquote>\n<hr />\n"},
	{"spec61", "`Foo\n----\n`\n\n<a title=\"a lot\n---\nof dashes\"/>\n", "<h2>`Foo</h2>\n<p>`</p>\n<h2>&lt;a title=&quot;a lot</h2>\n<p>of dashes&quot;/&gt;</p>\n"},
	{"spec60", "Foo\\\n----\n", "<h2>Foo\\</h2>\n"},
	{"spec58", "Foo\n= =\n\nFoo\n--- -\n", "<p>Foo\n= =</p>\n<p>Foo</p>\n<hr />\n"},
	{"spec57", "Foo\n    ---\n", "<p>Foo\n---</p>\n"},
	{"spec56", "Foo\n   ----      \n", "<h2>Foo</h2>\n"},
	{"spec55", "    Foo\n    ---\n\n    Foo\n---\n", "<pre><code>Foo\n---\n\nFoo\n</code></pre>\n<hr />\n"},
	{"spec53", "Foo\n-------------------------\n\nFoo\n=\n", "<h2>Foo</h2>\n<h1>Foo</h1>\n"},
	{"spec51", "Foo *bar\nbaz*\n====\n", "<h1>Foo <em>bar\nbaz</em></h1>\n"},
	{"spec49", "## \n#\n### ###\n", "<h2></h2>\n<h1></h1>\n<h3></h3>\n"},
	{"spec41", "## foo ##\n  ###   bar    ###\n", "<h2>foo</h2>\n<h3>bar</h3>\n"},
	{"spec38", " ### foo\n  ## foo\n   # foo\n", "<h3>foo</h3>\n<h2>foo</h2>\n<h1>foo</h1>\n"},
	{"spec37", "#                  foo                     \n", "<h1>foo</h1>\n"},
	{"spec36", "# foo *bar* \\*baz\\*\n", "<h1>foo <em>bar</em> *baz*</h1>\n"},
	{"spec35", "\\## foo\n", "<p>## foo</p>\n"},
	{"spec34", "#5 bolt\n\n#hashtag\n", "<p>#5 bolt</p>\n<p>#hashtag</p>\n"},
	{"spec32", "# foo\n## foo\n### foo\n#### foo\n##### foo\n###### foo\n", "<h1>foo</h1>\n<h2>foo</h2>\n<h3>foo</h3>\n<h4>foo</h4>\n<h5>foo</h5>\n<h6>foo</h6>\n"},
	{"spec30", "* Foo\n* * *\n* Bar\n", "<ul>\n<li>Foo</li>\n</ul>\n<hr />\n<ul>\n<li>Bar</li>\n</ul>\n"},
	{"spec29", "Foo\n---\nbar\n", "<h2>Foo</h2>\n<p>bar</p>\n"},
	{"spec27", "- foo\n***\n- bar\n", "<ul>\n<li>foo</li>\n</ul>\n<hr />\n<ul>\n<li>bar</li>\n</ul>\n"},
	{"spec26", " *-*\n", "<p><em>-</em></p>\n"},
	{"spec25", "_ _ _ _ a\n\na------\n\n---a---\n", "<p>_ _ _ _ a</p>\n<p>a------</p>\n<p>---a---</p>\n"},
	{"spec19", "Foo\n    ***\n", "<p>Foo\n***</p>\n"},
	{"spec18", "    ***\n", "<pre><code>***\n</code></pre>\n"},
	{"spec16", "--\n**\n__\n", "<p>--\n**\n__</p>\n"},
	{"spec14", "+++\n", "<p>+++</p>\n"},
	{"spec13", "***\n---\n___\n", "<hr />\n<hr />\n<hr />\n"},
	{"spec12", "- `one\n- two`\n", "<ul>\n<li>`one</li>\n<li>two`</li>\n</ul>\n"},
	{"spec11", "*\t*\t*\t\n", "<hr />\n"},
	{"spec10", "#\tFoo\n", "<h1>Foo</h1>\n"},
	{"spec9", " - foo\n   - bar\n\t - baz\n", "<ul>\n<li>foo\n<ul>\n<li>bar\n<ul>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>\n"},
	{"spec8", "    foo\n\tbar\n", "<pre><code>foo\nbar\n</code></pre>\n"},
	{"spec7", "-\t\tfoo\n", "<ul>\n<li>\n<pre><code>  foo\n</code></pre>\n</li>\n</ul>\n"},
	{"spec6", ">\t\tfoo\n", "<blockquote>\n<pre><code>  foo\n</code></pre>\n</blockquote>\n"},
	{"spec5", "- foo\n\n\t\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<pre><code>  bar\n</code></pre>\n</li>\n</ul>\n"},
	{"spec4", "  - foo\n\n\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<p>bar</p>\n</li>\n</ul>\n"},
	{"spec3", "    a\ta\n    ὐ\ta\n", "<pre><code>a\ta\nὐ\ta\n</code></pre>\n"},
	{"spec2", "  \tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},
	{"spec1", "\tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},

	// some simple cases

	{"simple16", "lu\u0000te", "<p>lu\uFFFDte</p>\n"},
	{"simple15", "[lute]: /url 'not a title", "<p>[lute]: /url 'not a title</p>\n"},
	{"simple14", "\\[not a lute](/url)", "<p>[not a lute](/url)</p>\n"},
	{"simple13", "- lu\n  - te", "<ul>\n<li>lu\n<ul>\n<li>te</li>\n</ul>\n</li>\n</ul>\n"},
	{"simple12", "`l*ut*e", "<p>`l<em>ut</em>e</p>\n"},
	{"simple11", "`lu\nte`", "<p><code>lu te</code></p>\n"},
	{"simple10", "lu\n\nte", "<p>lu</p>\n<p>te</p>\n"},
	{"simple9", "* lute", "<ul>\n<li>lute</li>\n</ul>\n"},
	{"simple8", "# lute", "<h1>lute</h1>\n"},
	{"simple7", "> lute", "<blockquote>\n<p>lute</p>\n</blockquote>\n"},
	{"simple6", "l**ut**e", "<p>l<strong>ut</strong>e</p>\n"},
	{"simple5", "l*ut*e", "<p>l<em>ut</em>e</p>\n"},
	{"simple4", "    lute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple3", "\tlute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple2", "l`ut`e", "<p>l<code>ut</code>e</p>\n"},
	{"simple1", "lute", "<p>lute</p>\n"},
	{"simple0", "", ""},
}

func TestParse(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMTaskListItem = false
	luteEngine.ParseOptions.GFMTable = false
	luteEngine.ParseOptions.GFMAutoLink = false
	luteEngine.ParseOptions.GFMStrikethrough = false
	luteEngine.RenderOptions.SoftBreak2HardBreak = false
	luteEngine.RenderOptions.CodeSyntaxHighlight = false
	luteEngine.ParseOptions.HeadingID = false
	luteEngine.RenderOptions.HeadingID = false
	luteEngine.RenderOptions.AutoSpace = false
	luteEngine.RenderOptions.FixTermTypo = false
	luteEngine.ParseOptions.Emoji = false

	for _, test := range parseTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

func TestCase1(t *testing.T) {
	bytes, err := os.ReadFile("case1.md")
	if nil != err {
		t.Fatalf("read case failed: %s", err)
	}

	luteEngine := lute.New()
	html := luteEngine.Markdown("case1.md", bytes)
	t.Log(string(html))
}
