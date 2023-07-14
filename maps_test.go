// Copyright 2021 Cloud Privacy Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
