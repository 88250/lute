package ast

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

const DocumentTitleSpec = 1
const DocumentTitleFormat = "inline-ast"
const DocumentTitleMaxBytes = 1 << 20
const DocumentTitleMaxNodes = 8192
const DocumentTitleMaxDepth = 64

var ErrInvalidDocumentTitle = errors.New("invalid document title source")

// DocumentTitle 保存规范行级节点数组的 JSON 文本，避免源与临时投影共享可变节点。
type DocumentTitle struct {
	Spec    int    `json:"spec"`
	Format  string `json:"format"`
	Content string `json:"content"`
}

// Validate 验证源版本、大小和 JSON 边界；节点类型与结构由标题片段解析器继续验证。
func (title *DocumentTitle) Validate() error {
	if nil == title || title.Spec != DocumentTitleSpec || title.Format != DocumentTitleFormat ||
		len(title.Content) > DocumentTitleMaxBytes || !utf8.ValidString(title.Content) {
		return ErrInvalidDocumentTitle
	}
	content := bytes.TrimSpace([]byte(title.Content))
	if len(content) == 0 || content[0] != '[' || validateDocumentTitleJSON(content) != nil {
		return ErrInvalidDocumentTitle
	}
	return nil
}

// ParseDocumentTitleJSON 在普通结构体解码前拒绝重复字段、未知字段和不完整的源封装。
func ParseDocumentTitleJSON(data []byte) (*DocumentTitle, error) {
	if len(data) > 6*DocumentTitleMaxBytes || !utf8.Valid(data) || validateDocumentTitleJSON(data) != nil {
		return nil, ErrInvalidDocumentTitle
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(data, &fields) != nil || len(fields) != 3 {
		return nil, ErrInvalidDocumentTitle
	}
	for _, key := range []string{"spec", "format", "content"} {
		if len(fields[key]) == 0 || bytes.Equal(bytes.TrimSpace(fields[key]), []byte("null")) {
			return nil, ErrInvalidDocumentTitle
		}
	}
	ret := &DocumentTitle{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if decoder.Decode(ret) != nil || ret.Validate() != nil {
		return nil, ErrInvalidDocumentTitle
	}
	return ret, nil
}

func validateDocumentTitleJSON(data []byte) error {
	if !validDocumentTitleUnicodeEscapes(data) {
		return ErrInvalidDocumentTitle
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value func(int) error
	value = func(depth int) error {
		if depth > DocumentTitleMaxDepth {
			return ErrInvalidDocumentTitle
		}
		token, err := decoder.Token()
		if nil != err {
			return err
		}
		switch token := token.(type) {
		case nil:
			return ErrInvalidDocumentTitle
		case string:
			if strings.ContainsRune(token, 0) {
				return ErrInvalidDocumentTitle
			}
		case json.Delim:
			switch token {
			case '{':
				keys := map[string]bool{}
				for decoder.More() {
					key, keyErr := decoder.Token()
					name, ok := key.(string)
					if nil != keyErr || !ok || keys[name] || strings.ContainsRune(name, 0) {
						return ErrInvalidDocumentTitle
					}
					keys[name] = true
					if err = value(depth + 1); nil != err {
						return err
					}
				}
				end, endErr := decoder.Token()
				if nil != endErr || end != json.Delim('}') {
					return ErrInvalidDocumentTitle
				}
			case '[':
				for decoder.More() {
					if err = value(depth + 1); nil != err {
						return err
					}
				}
				end, endErr := decoder.Token()
				if nil != endErr || end != json.Delim(']') {
					return ErrInvalidDocumentTitle
				}
			default:
				return ErrInvalidDocumentTitle
			}
		}
		return nil
	}
	if err := value(0); nil != err {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return ErrInvalidDocumentTitle
	}
	return nil
}

// encoding/json replaces unpaired UTF-16 surrogates with U+FFFD. A canonical
// source reader must reject that corruption instead of silently changing text.
func validDocumentTitleUnicodeEscapes(data []byte) bool {
	inString := false
	for i := 0; i < len(data); i++ {
		if data[i] == '"' {
			inString = !inString
			continue
		}
		if !inString || data[i] != '\\' {
			continue
		}
		if i+1 >= len(data) {
			return false
		}
		if data[i+1] != 'u' {
			i++
			continue
		}
		if i+6 > len(data) {
			return false
		}
		value, err := strconv.ParseUint(string(data[i+2:i+6]), 16, 16)
		if err != nil || value >= 0xDC00 && value <= 0xDFFF {
			return false
		}
		if value >= 0xD800 && value <= 0xDBFF {
			if i+12 > len(data) || data[i+6] != '\\' || data[i+7] != 'u' {
				return false
			}
			low, err := strconv.ParseUint(string(data[i+8:i+12]), 16, 16)
			if err != nil || low < 0xDC00 || low > 0xDFFF {
				return false
			}
			i += 6
		}
		i += 5
	}
	return true
}
