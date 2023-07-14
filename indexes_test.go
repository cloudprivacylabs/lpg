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

func TestBtreeNodeIndex(t *testing.T) {
	g := NewGraph()
	g.index.NodePropertyIndex("index", g, BtreeIndex)
	labels := []string{"a", "b", "c", "d", "e", "f"}
	data := make(map[string]struct{})
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			g.NewNode([]string{l}, map[string]interface{}{"index": fmt.Sprint(i)})
			data[fmt.Sprintf("%s:%d", l, i)] = struct{}{}
		}
	}
	itr := g.index.GetIteratorForNodeProperty("index", "0")
	if size := itr.MaxSize(); size != 6 {
		t.Errorf("Expecting 6, got %d", size)
	}
	foundLabel := make(map[string]struct{})
	for itr.Next() {
		foundLabel[itr.Node().GetLabels().Slice()[0]] = struct{}{}
	}
	for _, l := range labels {
		if _, found := foundLabel[l]; !found {
			t.Errorf("Not found: %s", l)
		}
	}
}
func TestHashNodeIndex(t *testing.T) {
	g := NewGraph()
	g.index.NodePropertyIndex("index", g, HashIndex)
	labels := []string{"a", "b", "c", "d", "e", "f"}
	data := make(map[string]struct{})
	for _, l := range labels {
		for i := 0; i < 10; i++ {
			g.NewNode([]string{l}, map[string]interface{}{"index": fmt.Sprint(i)})
			data[fmt.Sprintf("%s:%d", l, i)] = struct{}{}
		}
	}
	itr := g.index.GetIteratorForNodeProperty("index", "0")
	if size := itr.MaxSize(); size != 6 {
		t.Errorf("Expecting 6, got %d", size)
	}
	foundLabel := make(map[string]struct{})
	for itr.Next() {
		foundLabel[itr.Node().GetLabels().Slice()[0]] = struct{}{}
	}
	for _, l := range labels {
		if _, found := foundLabel[l]; !found {
			t.Errorf("Not found: %s", l)
		}
	}
}
