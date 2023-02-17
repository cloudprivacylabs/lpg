package lpg

import (
	"fmt"
	"testing"
)

func TestNodeMap(t *testing.T) {
	m := NewNodeMap()
	labels := [][]string{{"a"}, {"b", "c", "d"}, {"e", "f"}}
	data := make(map[string]struct{})
	id := 0
	st := &stringTable{}
	st.init()
	idx := st.allocate("index")
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			node := &Node{labels: NewStringSet(l...), id: id}
			id++
			node.properties = make(properties)
			node.properties[idx] = i
			m.Add(node)
			data[fmt.Sprintf("%d:%d", len(l), i)] = struct{}{}
		}
	}
	// itr: 30 items
	itr := m.Iterator()
	found := make(map[string]struct{})
	for itr.Next() {
		node := itr.Node()
		found[fmt.Sprintf("%d:%d", node.labels.Len(), node.properties[idx])] = struct{}{}
	}
	if len(found) != len(data) {
		t.Errorf("found: %v", found)
	}

	// Label-based iteration
	for _, label := range labels {
		itr = m.IteratorAllLabels(NewStringSet(label...))
		found = make(map[string]struct{})
		for itr.Next() {
			node := itr.Node()
			if !node.labels.HasAll(label...) {
				t.Errorf("Expecting %v got %+v", label, node)
			}
			found[fmt.Sprint(node.properties[idx])] = struct{}{}
		}
		if len(found) != 10 {
			t.Errorf("10 entries were expected, got %v", found)
		}
	}
}
