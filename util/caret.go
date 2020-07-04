package util

// Caret 插入符 \u2038。
const Caret = "‸"

// CaretTokens 是插入符的字节数组。
var CaretTokens = []byte(Caret)

// CaretReplacement 用于解析过程中临时替换。
const CaretReplacement = "caretreplacement"