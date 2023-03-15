// Copyright 2021 Cloud Privacy Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lpg

import (
	"container/list"
)

// A fastSet is a set of objects with constant-time
// insertion/deletion, with iterator support
type fastSet struct {
	m map[int]*list.Element
	l *list.List
}

func newFastSet() *fastSet {
	return &fastSet{
		m: make(map[int]*list.Element),
		l: list.New(),
	}
}

func (f *fastSet) init() {
	f.m = make(map[int]*list.Element)
	f.l = list.New()
}

func (f fastSet) size() int { return len(f.m) }

// Add a new item. Returns true if added
func (f *fastSet) add(id int, item interface{}) bool {
	_, exists := f.m[id]
	if exists {
		return false
	}
	el := f.l.PushBack(item)
	f.m[id] = el
	return true
}

func (f *fastSet) get(id int) (interface{}, bool) {
	el, ok := f.m[id]
	if !ok {
		return nil, false
	}
	return el.Value, true
}

// Remove an item. Returns true if removed
func (f *fastSet) remove(id int) bool {
	el := f.m[id]
	if el == nil {
		return false
	}
	delete(f.m, id)
	f.l.Remove(el)
	return true
}

func (f fastSet) has(id int) bool {
	_, exists := f.m[id]
	return exists
}

func (f fastSet) iterator() Iterator {
	return &listIterator{next: f.l.Front(), size: f.size()}
}

type NodeSet struct {
	set fastSet
}

func NewNodeSet() *NodeSet {
	nm := &NodeSet{}
	nm.set.init()
	return nm
}

func (set *NodeSet) Add(node *Node) {
	set.set.add(node.id, node)
}

func (set NodeSet) Remove(node *Node) {
	set.set.remove(node.id)
}

func (set NodeSet) Has(node *Node) bool {
	return set.set.has(node.id)
}

func (set NodeSet) Len() int {
	return set.set.size()
}

func (set NodeSet) Iterator() NodeIterator {
	i := set.set.iterator()
	return nodeIterator{i}
}

func (set NodeSet) Slice() []*Node {
	return NodeSlice(set.Iterator())
}

// EdgeSet keeps an unordered set of edges
type EdgeSet struct {
	set fastSet
}

func NewEdgeSet() *EdgeSet {
	es := &EdgeSet{}
	es.set.init()
	return es
}

func (set *EdgeSet) Add(edge *Edge) {
	set.set.add(edge.id, edge)
}

func (set EdgeSet) Remove(edge *Edge) {
	set.set.remove(edge.id)
}

func (set EdgeSet) Len() int {
	return set.set.size()
}

func (set EdgeSet) Iterator() EdgeIterator {
	i := set.set.iterator()
	return edgeIterator{i}
}

func (set EdgeSet) Slice() *Path {
	return EdgeSlice(set.Iterator())
}
