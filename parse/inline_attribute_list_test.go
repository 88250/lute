package parse

import (
	"reflect"
	"testing"
)

func TestMergeIALPreservingOrder(t *testing.T) {
	tests := []struct {
		name string
		dst  [][]string
		src  [][]string
		want [][]string
	}{
		{
			// 覆盖已有属性但不改变其位置；新增属性按 src 顺序追加，值会先做反转义。
			name: "updates_existing_attrs_and_appends_new_attrs",
			dst:  [][]string{{"id", "old"}, {"class", "foo&amp;bar"}},
			src:  [][]string{{"id", "new"}, {"updated", "20260405"}, {"title", "a&amp;b"}},
			want: [][]string{{"id", "new"}, {"class", "foo&bar"}, {"updated", "20260405"}, {"title", "a&b"}},
		},
		{
			// 对应 SpinBlockDOM case 262：已有块级属性顺序保持不变，新的 IAL 属性只能追加到末尾。
			name: "appends_new_attr_after_existing_block_attrs",
			dst:  [][]string{{"data-node-id", "20250824233004-3qfd5gf"}, {"data-node-index", "1"}, {"data-type", "NodeBlockquote"}, {"class", "bq"}, {"updated", "20250824233509"}},
			src:  [][]string{{"custom-b", "info"}},
			want: [][]string{{"data-node-id", "20250824233004-3qfd5gf"}, {"data-node-index", "1"}, {"data-type", "NodeBlockquote"}, {"class", "bq"}, {"updated", "20250824233509"}, {"custom-b", "info"}},
		},
		{
			// 与旧逻辑的 map 语义对齐：重复 key 需要折叠，最后一个值生效，
			// 但输出顺序保持稳定，不再随机漂移。
			name: "collapses_duplicate_attrs_with_stable_order",
			dst:  [][]string{{"id", "old"}, {"id", "shadow"}, {"class", "foo"}},
			src:  [][]string{{"class", "bar"}, {"id", "new"}},
			want: [][]string{{"id", "new"}, {"class", "bar"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeIALPreservingOrder(tt.dst, tt.src)
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("unexpected merged ial, want %v, got %v", tt.want, got)
			}
		})
	}
}
