// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

// tagfilter 将一些标签的 < 替换为 &lt;
func (r *Renderer) tagfilter(tokens items) items {
	//tokens = bytes.ReplaceAll(tokens, []byte("<xmp>"), []byte("&lt;xmp>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<title>"), []byte("&lt;title>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<style>"), []byte("&lt;style>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<script>"), []byte("&lt;script>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<iframe>"), []byte("&lt;iframe>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<noembed>"), []byte("&lt;noembed>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<textarea>"), []byte("&lt;textarea>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<noframes>"), []byte("&lt;noframes>"))
	//tokens = bytes.ReplaceAll(tokens, []byte("<plaintext>"), []byte("&lt;plaintext>"))

	length := len(tokens)
	var i int
	var token byte
	for ; i < length; i++ {
		token = tokens[i]
		if itemLess != token {
			continue
		}

		if i < length-6 {
			token1 := tokens[i+1]
			token2 := tokens[i+2]
			token3 := tokens[i+3]
			token4 := tokens[i+4]
			token5 := tokens[i+5]
			token6 := tokens[i+6]

			if ('t' == token1 || 'T' == token1) &&
				('i' == token2 || 'I' == token2) &&
				('t' == token3 || 'T' == token3) &&
				('l' == token4 || 'L' == token4) &&
				('e' == token5 || 'E' == token5) &&
				(itemGreater == token6) {
				tokens = append(tokens, 0, 0, 0)
				copy(tokens[i+3:], tokens[i:])
				tokens[i] = '&'
				tokens[i+1] = 'l'
				tokens[i+2] = 't'
				tokens[i+3] = ';'
				i += 3
				length += 3
			}

		}
	}

	return tokens
}
