package helper

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name   string
		listA  []string
		listB  []string
		result struct {
			removed []string
			added   []string
		}
	}{
		{
			name:   "b removed and e added",
			listA:  []string{"a", "b", "c", "d"},
			listB:  []string{"a", "c", "d", "e"},
			result: struct{ removed, added []string }{removed: []string{"b"}, added: []string{"e"}},
		},
		{
			name:   "No Differences",
			listA:  []string{"a", "b", "c"},
			listB:  []string{"a", "b", "c"},
			result: struct{ removed, added []string }{removed: nil, added: nil},
		},
		{
			name:   "Items Only in List B",
			listA:  []string{"a"},
			listB:  []string{"a", "b", "c"},
			result: struct{ removed, added []string }{removed: nil, added: []string{"b", "c"}},
		},
		{
			name:   "Empty Lists",
			listA:  []string{},
			listB:  []string{},
			result: struct{ removed, added []string }{removed: nil, added: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removed, added := Diff(tt.listA, tt.listB)

			if !reflect.DeepEqual(removed, tt.result.removed) {
				t.Errorf("Removed items mismatch. Expected %v, got %v", tt.result.removed, removed)
			}

			if !reflect.DeepEqual(added, tt.result.added) {
				t.Errorf("Added items mismatch. Expected %v, got %v", tt.result.added, added)
			}
		})
	}
}
