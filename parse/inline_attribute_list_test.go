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
			// 锁定当前兼容行为：dst 中若已有重复 key，只更新第一次出现的位置。
			name: "updates_only_first_duplicate_key_in_dst",
			dst:  [][]string{{"id", "old"}, {"id", "shadow"}, {"class", "foo"}},
			src:  [][]string{{"id", "new"}, {"class", "bar"}},
			want: [][]string{{"id", "new"}, {"id", "shadow"}, {"class", "bar"}},
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
