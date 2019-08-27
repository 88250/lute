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
)

// fixTermTypo 修正 str 中出现的术语拼写问题。
func fixTermTypo(str string) string {
	// 鸣谢 https://github.com/studygolang/autocorrect

	for from, to := range terms {
		re := regexp.MustCompile("(?i)" + from)
		str = re.ReplaceAllString(str, to)
	}
	return str
}

// terms 定义了术语字典，用于术语拼写修正。
var terms = map[string]string{
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
