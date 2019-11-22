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
	"bytes"
	"unicode/utf8"
)

// fixTermTypo 修正文本节点 textNode 中出现的术语拼写问题。
func (r *BaseRenderer) fixTermTypo(textNode *Node) {
	textNode.tokens = r.fixTermTypo0(textNode.tokens)
}

func (r *BaseRenderer) fixTermTypo0(tokens []byte) []byte {
	length := len(tokens)
	var token byte
	var i, j, k, l int
	var before, after byte
	var originalTerm []byte
	for ; i < length; i++ {
		token = tokens[i]
		if isNotTerm(token) {
			continue
		}
		if 1 <= i {
			before = tokens[i-1]
			if !isNotTerm(before) {
				// 前一个字节必须是非术语，否则无法分隔
				continue
			}
		}
		if isASCIIPunct(before) {
			// 比如术语前面如果是 . 则不进行修正，因为可能是链接
			// 比如 test.html 虽然不能识别为自动链接，但是也不能进行修正
			continue
		}

		for j = i; j < length; j++ {
			after = tokens[j]
			if isNotTerm(after) || itemDot == after {
				break
			}
		}
		if isASCIIPunct(after) {
			// 比如术语后面如果是 . 则不进行修正，因为可能是链接
			// 比如 github.com 虽然不能识别为自动链接，但是也不能进行修正
			continue
		}

		originalTerm = bytes.ToLower(tokens[i:j])
		if to, ok := r.option.Terms[bytesToStr(originalTerm)]; ok {
			l = 0
			for k = i; k < j; k++ {
				tokens[k] = to[l]
				l++
			}
		}
	}

	return tokens
}

func isNotTerm(token byte) bool {
	return token >= utf8.RuneSelf || isWhitespace(token) || isASCIIPunct(token)
}

func newTerms() (ret map[string]string) {
	ret = make(map[string]string, len(terms))
	for k, v := range terms {
		ret[k] = v
	}
	return
}

// terms 定义了术语字典，用于术语拼写修正。Key 必须是全小写的。
var terms = map[string]string{
	"jetty":         "Jetty",
	"tomcat":        "Tomcat",
	"jdbc":          "JDBC",
	"mariadb":       "MariaDB",
	"ipfs":          "IPFS",
	"saas":          "SaaS",
	"paas":          "PaaS",
	"iaas":          "IaaS",
	"ioc":           "IoC",
	"freemarker":    "FreeMarker",
	"ruby":          "Ruby",
	"rails":         "Rails",
	"mina":          "Mina",
	"puppet":        "Puppet",
	"vagrant":       "Vagrant",
	"chef":          "Chef",
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
	"csv":           "CSV",
	"soap":          "SOAP",
	"ajax":          "AJAX",
	"messagepack":   "MessagePack",
	"javascript":    "JavaScript",
	"java":          "Java",
	"jsp":           "JSP",
	"restful":       "RESTFul",
	"graphql":       "GraphQL",
	"gorm":          "GORM",
	"orm":           "ORM",
	"oauth":         "OAuth",
	"facebook":      "Facebook",
	"github":        "GitHub",
	"gist":          "Gist",
	"heroku":        "Heroku",
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
	"b3log":         "B3log",
	"hacpai":        "HacPai",
	"lute":          "Lute",
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
	"sprintmvc":     "SpringMVC",
	"mybatis":       "MyBatis",
	"qq":            "QQ",
	"sql":           "SQL",
}
