package lpg

import (
	"fmt"
	"testing"
)

func TestEdgeMap(t *testing.T) {
	m := NewEdgeMap()
	labels := []string{"a", "b", "c", "d", "e", "f"}
	data := make(map[string]struct{})
	id := 0
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			edge := &Edge{label: l, id: id}
			id++
			edge.properties = make(properties)
			edge.properties["index"] = i
			m.Add(edge)
			data[fmt.Sprintf("%s:%d", l, i)] = struct{}{}
		}
	}
	// itr: 60 items
	itr := m.Iterator()
	found := make(map[string]struct{})
	for itr.Next() {
		edge := itr.Edge()
		found[fmt.Sprintf("%s:%d", edge.label, edge.properties["index"])] = struct{}{}
	}
	if len(found) != len(data) {
		t.Errorf("found: %v", found)
	}

	// Label-based iteration
	for _, label := range labels {
		itr = m.IteratorLabel(label)
		found = make(map[string]struct{})
		for itr.Next() {
			edge := itr.Edge()
			if edge.label != label {
				t.Errorf("Expecting %s got %+v", label, edge)
			}
			found[fmt.Sprint(edge.properties["index"])] = struct{}{}
		}
		if len(found) != 10 {
			t.Errorf("10 entries were expected, got %v", found)
		}
	}

	itr = m.IteratorAnyLabel(NewStringSet("a", "c", "e", "g"))
	found = make(map[string]struct{})
	for itr.Next() {
		edge := itr.Edge()
		if edge.label != "a" && edge.label != "c" && edge.label != "e" {
			t.Errorf("Unexpected label: %s", edge.label)
		}
		found[fmt.Sprintf("%s:%d", edge.label, edge.properties["index"])] = struct{}{}
	}
	if len(found) != 30 {
		t.Errorf("Expecting 30, got %v", found)
	}
}

func TestEdgeMap2(t *testing.T) {
	m := newEdgeMap()
	labels := []string{"a", "b", "c", "d", "e", "f"}
	data := make(map[string]struct{})
	id := 0
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			edge := &Edge{label: l, id: id}
			id++
			edge.properties = make(properties)
			edge.properties["index"] = i
			m.add(edge)
			data[fmt.Sprintf("%s:%d", l, i)] = struct{}{}
		}
	}
	// itr: 60 items
	itr := m.iterator()
	found := make(map[string]struct{})
	for itr.Next() {
		edge := itr.Edge()
		found[fmt.Sprintf("%s:%d", edge.label, edge.properties["index"])] = struct{}{}
	}
	if len(found) != len(data) {
		t.Errorf("found: %v", found)
	}

	// Label-based iteration
	for _, label := range labels {
		itr = m.iteratorLabel(label)
		found = make(map[string]struct{})
		for itr.Next() {
			edge := itr.Edge()
			if edge.label != label {
				t.Errorf("Expecting %s got %+v", label, edge)
			}
			found[fmt.Sprint(edge.properties["index"])] = struct{}{}
		}
		if len(found) != 10 {
			t.Errorf("10 entries were expected, got %v", found)
		}
	}

	itr = m.iteratorAnyLabel(NewStringSet("a", "c", "e", "g"))
	found = make(map[string]struct{})
	for itr.Next() {
		edge := itr.Edge()
		if edge.label != "a" && edge.label != "c" && edge.label != "e" {
			t.Errorf("Unexpected label: %s", edge.label)
		}
		found[fmt.Sprintf("%s:%d", edge.label, edge.properties["index"])] = struct{}{}
	}
	if len(found) != 30 {
		t.Errorf("Expecting 30, got %v", found)
	}
}

func TestNodeMap(t *testing.T) {
	m := NewNodeMap()
	labels := [][]string{{"a"}, {"b", "c", "d"}, {"e", "f"}}
	data := make(map[string]struct{})
	id := 0
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			node := &Node{labels: NewStringSet(l...), id: id}
			id++
			node.properties = make(properties)
			node.properties["index"] = i
			m.Add(node)
			data[fmt.Sprintf("%d:%d", len(l), i)] = struct{}{}
		}
	}
	// itr: 30 items
	itr := m.Iterator()
	found := make(map[string]struct{})
	for itr.Next() {
		node := itr.Node()
		found[fmt.Sprintf("%d:%d", node.labels.Len(), node.properties["index"])] = struct{}{}
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
			found[fmt.Sprint(node.properties["index"])] = struct{}{}
		}
		if len(found) != 10 {
			t.Errorf("10 entries were expected, got %v", found)
		}
	}
}
