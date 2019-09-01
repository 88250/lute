// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// fixTermTypo 修正 node 下文本节点中出现的术语拼写问题。
func (t *Tree) fixTermTypo(node *Node) {
	if nil == node {
		return
	}

	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			text := fromItems(child.tokens)
			text = fixTermTypo0(text)
			child.tokens = toItems(text)
		} else {
			t.fixTermTypo(child) // 递归处理子节点
		}
		child = next
	}
}

func fixTermTypo0(str string) string {
	// 鸣谢 https://github.com/studygolang/autocorrect

	strs := strings.FieldsFunc(str, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})
	for _, s := range strs {
		if s[0] >= utf8.RuneSelf {
			// 术语仅由 ASCII 字符组成
			continue
		}

		for from, to := range terms {
			if strings.EqualFold(s, from) {
				re := regexp.MustCompile("(?i)" + from) // TODO: 不要用正则
				str = re.ReplaceAllString(str, to)
			}
		}
	}

	return str
}

// terms 定义了术语字典，用于术语拼写修正。
// TODO: 考虑提供接口支持开发者添加
var terms = map[string]string{
	"jdbc":          "JDBC",
	"mariadb":       "MariaDB",
	"ipfs":          "IPFS",
	"saas":          "SaaS",
	"paas":          "PaaS",
	"iaas":          "IaaS",
	"ioc":           "IoC",
	"freemarker":    "FreeMarker",
	"ruby":          "Ruby",
	"mri":           "MRI",
	"rails":         "Rails",
	"mina":          "Mina",
	"puppet":        "Puppet",
	"vagrant":       "Vagrant",
	"chef":          "Chef",
	"nodejs":        "Node.js",
	"npm":           "NPM",
	"beego":         "Beego",
	"gin":           "Gin",
	"iris":          "Iris",
	"php":           "PHP",
	"ssh":           "SSH",
	"web":           "Web",
	"api":           "API",
	"css":           "CSS",
	"html":          "HTML",
	"json":          "JSON",
	"jsonp":         "JSONP",
	"xml":           "XML",
	"yaml":          "YAML",
	"yml":           "YAML",
	"ini":           "INI",
	"csv":           "CSV",
	"soap":          "SOAP",
	"ajax":          "AJAX",
	"messagepack":   "MessagePack",
	"javascript":    "JavaScript",
	"java":          "Java",
	"jsp":           "JSP",
	"asp.net":       "ASP.NET",
	".net":          ".NET",
	"restful":       "RESTFul",
	"orm":           "ORM",
	"oauth":         "OAuth",
	"markdown":      "Markdown",
	"facebook":      "Facebook",
	"github":        "GitHub",
	"gist":          "Gist",
	"heroku":        "Heroku",
	"stackoverflow": "Stack Overflow",
	"stackexchange": "Stack Exchange",
	"twitter":       "Twitter",
	"youtube":       "YouTube",
	"dynamodb":      "DynamoDB",
	"mysql":         "MySQL",
	"postgresql":    "PostgreSQL",
	"sqlite":        "SQLite",
	"memcached":     "Memcached",
	"mongodb":       "MongoDB",
	"redis":         "Redis",
	"elasticsearch": "Elasticsearch",
	"solr":          "Solr",
	"sphinx":        "Sphinx",
	"linux":         "Linux",
	"mac":           "Mac",
	"osx":           "OS X",
	"ubuntu":        "Ubuntu",
	"centos":        "CentOS",
	"centos7":       "CentOS7",
	"redhat":        "RedHat",
	"gitlab":        "GitLab",
	"jquery":        "jQuery",
	"angularjs":     "AngularJS",
	"ffmpeg":        "FFMPEG",
	"git":           "Git",
	"svn":           "SVN",
	"vim":           "VIM",
	"emacs":         "Emacs",
	"sublime":       "Sublime",
	"virtualbox":    "VirtualBox",
	"safari":        "Safari",
	"chrome":        "Chrome",
	"ie":            "IE",
	"firefox":       "Firefox",
	"iterm":         "iTerm",
	"iterm2":        "iTerm2",
	"iwork":         "iWork",
	"itunes":        "iTunes",
	"iphoto":        "iPhoto",
	"ibook":         "iBook",
	"imessage":      "iMessage",
	"photoshop":     "Photoshop",
	"excel":         "Excel",
	"powerpoint":    "PowerPoint",
	"ios":           "iOS",
	"iphone":        "iPhone",
	"ipad":          "iPad",
	"android":       "Android",
	"imac":          "iMac",
	"macbook":       "MacBook",
	"vps":           "VPS",
	"vpn":           "VPN",
	"arm":           "ARM",
	"cpu":           "CPU",
	"spring":        "Spring",
	"springboot":    "SpringBoot",
	"springcloud":   "SpringCloud",
	"qq":            "QQ",
}
