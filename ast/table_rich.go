package ast

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const TableCellRichSpec = 1
const TableCellRichFormat = "kramdown"
const TableCellRichAttribute = "data-sy-table-cell-rich"
const TableCellRichIAL = "table-cell-rich"

// TableCellRich 保存单元格的富文本源，格式版本独立于文档规范版本。
type TableCellRich struct {
	Spec    int    `json:"spec"`
	Format  string `json:"format"`
	Content string `json:"content"`
	Raw     string `json:"-"`
	Invalid bool   `json:"-"`
}

// Validate 检查载荷版本，损坏的内部载荷保留在 Raw 中，供调用方拒绝写入。
func (rich *TableCellRich) Validate() error {
	if nil == rich {
		return nil
	}
	if rich.Invalid || "" != rich.Raw || TableCellRichSpec != rich.Spec || TableCellRichFormat != rich.Format {
		return fmt.Errorf("unsupported table cell rich text format [%d, %s]", rich.Spec, rich.Format)
	}
	if strings.ContainsRune(rich.Content, 0) || !utf8.ValidString(rich.Content) {
		return fmt.Errorf("table cell rich text contains invalid text")
	}
	return nil
}

// Encode 为内部 Markdown 和 BlockDOM 编码，避免换行、引号和管道符破坏单元格边界。
func (rich *TableCellRich) Encode() string {
	if nil == rich {
		return ""
	}
	if rich.Invalid || "" != rich.Raw {
		return rich.Raw
	}
	data, _ := json.Marshal(rich)
	return base64.RawURLEncoding.EncodeToString(data)
}

// DecodeTableCellRich 保留无法解码的原始载荷，禁止将它当作普通单元格静默覆盖。
func DecodeTableCellRich(encoded string) *TableCellRich {
	rich := &TableCellRich{Raw: encoded, Invalid: true}
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if nil != err || !utf8.Valid(data) {
		return rich
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(data, &fields); nil != err || nil == fields {
		return rich
	}
	if len(fields) != 3 {
		return rich
	}
	for _, key := range []string{"spec", "format", "content"} {
		if len(fields[key]) == 0 || bytes.Equal(bytes.TrimSpace(fields[key]), []byte("null")) {
			return rich
		}
	}
	if err = json.Unmarshal(data, rich); nil != err {
		return &TableCellRich{Raw: encoded, Invalid: true}
	}
	rich.Raw = ""
	rich.Invalid = false
	return rich
}
