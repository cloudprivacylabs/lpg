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
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// An EdgeMap stores edges indexed by edge label
type EdgeMap struct {
	m *linkedhashmap.Map
	n int
}

func NewEdgeMap() *EdgeMap {
	return &EdgeMap{
		m: linkedhashmap.New(),
	}
}

func (em *EdgeMap) Clear() {
	em.m = linkedhashmap.New()
	em.n = 0
}

func (em *EdgeMap) Add(edge *Edge) {
	var set *FastSet
	v, found := em.m.Get(edge.label)
	if !found {
		set = NewFastSet()
		em.m.Put(edge.label, set)
	} else {
		set = v.(*FastSet)
	}
	if set.Add(edge.id, edge) {
		em.n++
	}
}

func (em EdgeMap) Remove(edge *Edge) {
	var set *FastSet
	v, found := em.m.Get(edge.label)
	if !found {
		return
	}
	set = v.(*FastSet)
	if set.Remove(edge.id, edge) {
		em.n--
	}
	if set.Size() == 0 {
		em.m.Remove(edge.label)
	}
}

func (em EdgeMap) IsEmpty() bool {
	if em.m.Size() == 0 {
		return true
	}
	return false
}

func (em EdgeMap) Len() int { return em.n }

type edgeMapIterator struct {
	labels  iteratorWithoutSize
	current EdgeIterator
	size    int
}

func (itr *edgeMapIterator) Next() bool {
	if itr.current != nil {
		if itr.current.Next() {
			return true
		}
		itr.current = nil
	}
	if itr.labels == nil {
		return false
	}
	if !itr.labels.Next() {
		return false
	}
	set := itr.labels.Value().(*FastSet)
	setItr := set.Iterator()
	itr.current = &edgeIterator{withSize(setItr, set.Size())}
	itr.current.Next()
	return true
}

func (itr *edgeMapIterator) Value() interface{} {
	return itr.current.Value()
}

func (itr *edgeMapIterator) Edge() *Edge {
	return itr.current.Edge()
}

func (itr *edgeMapIterator) MaxSize() int { return itr.size }

func (em EdgeMap) Iterator() EdgeIterator {
	i := em.m.Iterator()
	return &edgeMapIterator{labels: &i, size: em.Len()}
}

func (em EdgeMap) IteratorLabel(label string) EdgeIterator {
	v, found := em.m.Get(label)
	if !found {
		return &edgeIterator{&emptyIterator{}}
	}
	set := v.(*FastSet)
	i := set.Iterator()
	return &edgeIterator{withSize(i, set.Size())}
}

func (em EdgeMap) IteratorAnyLabel(labels StringSet) EdgeIterator {
	strings := labels.Slice()
	return &edgeIterator{&funcIterator{
		iteratorFunc: func() Iterator {
			for len(strings) != 0 {
				v, found := em.m.Get(strings[0])
				strings = strings[1:]
				if !found {
					continue
				}
				itr := v.(*FastSet).Iterator()
				return withSize(itr, -1)
			}
			return nil
		},
	},
	}
}

// An NodeMap stores nodes indexed by node labels
type NodeMap struct {
	// m[string]*FastSet
	m        *linkedhashmap.Map
	nolabels FastSet
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		m:        linkedhashmap.New(),
		nolabels: *NewFastSet(),
	}
}

func (nm *NodeMap) Replace(node *Node, oldLabels, newLabels StringSet) {
	if oldLabels.Len() == 0 {
		if newLabels.Len() == 0 {
			return
		}
		nm.nolabels.Remove(node.id, node)
	}
	if newLabels.Len() == 0 {
		nm.nolabels.Add(node.id, node)
		return
	}
	var set *FastSet
	// Process removed labels
	for label := range oldLabels.M {
		if !newLabels.Has(label) {
			v, found := nm.m.Get(label)
			if !found {
				continue
			}
			set = v.(*FastSet)
			set.Remove(node.id, node)
			if set.Len() == 0 {
				nm.m.Remove(label)
			}
		}
	}
	// Process added labels
	for label := range newLabels.M {
		if !oldLabels.Has(label) {
			v, found := nm.m.Get(label)
			if !found {
				set = NewFastSet()
				nm.m.Put(label, set)
			} else {
				set = v.(*FastSet)
			}
			set.Add(node.id, node)
		}
	}
}

func (nm *NodeMap) Add(node *Node) {
	if node.labels.Len() == 0 {
		nm.nolabels.Add(node.id, node)
		return
	}

	var set *FastSet
	for label := range node.labels.M {
		v, found := nm.m.Get(label)
		if !found {
			set = NewFastSet()
			nm.m.Put(label, set)
		} else {
			set = v.(*FastSet)
		}
		set.Add(node.id, node)
	}
}

func (nm NodeMap) Remove(node *Node) {
	if node.labels.Len() == 0 {
		nm.nolabels.Remove(node.id, node)
		return
	}
	var set *FastSet
	for label := range node.labels.M {
		v, found := nm.m.Get(label)
		if !found {
			continue
		}
		set = v.(*FastSet)
		set.Remove(node.id, node)
		if set.Len() == 0 {
			nm.m.Remove(label)
		}
	}
}

func (nm NodeMap) IsEmpty() bool {
	if nm.m.Size() == 0 {
		return true
	}
	return false
}

type nodeMapIterator struct {
	labels     *linkedhashmap.Iterator
	seenLabels []string
	current    NodeIterator
}

func (itr *nodeMapIterator) Next() bool {
	if itr.current != nil {
		if itr.current.Next() {
			return true
		}
		itr.current = nil
	}
	if itr.labels == nil {
		return false
	}
	if !itr.labels.Next() {
		return false
	}
	itr.seenLabels = append(itr.seenLabels, itr.labels.Key().(string))
	set := itr.labels.Value().(*FastSet)
	setItr := set.Iterator()
	itr.current = &nodeIterator{withSize(setItr, -1)}
	itr.current.Next()
	return true
}

func (itr *nodeMapIterator) Value() interface{} {
	return itr.current.Value()
}

func (itr *nodeMapIterator) Node() *Node {
	return itr.current.Node()
}

func (nm NodeMap) Iterator() NodeIterator {
	i := nm.m.Iterator()

	nmIterator := &nodeMapIterator{labels: &i}
	return &nodeIterator{
		MultiIterator(
			&filterIterator{
				itr: withSize(nmIterator, -1),
				filter: func(node interface{}) bool {
					onode := node.(*Node)
					nSeen := 0
					for _, l := range nmIterator.seenLabels {
						if onode.labels.Has(l) {
							nSeen++
						}
					}
					return nSeen < 2
				},
			},
			nm.nolabels.Iterator(),
		),
	}
}

func (nm NodeMap) IteratorAllLabels(labels StringSet) NodeIterator {
	// Find the smallest map element, iterate that
	var minSet *FastSet
	for label := range labels.M {
		v, found := nm.m.Get(label)
		if !found {
			return &nodeIterator{&emptyIterator{}}
		}
		mp := v.(*FastSet)
		if minSet == nil || minSet.Len() > mp.Len() {
			minSet = mp
		}
	}
	itr := minSet.Iterator()
	flt := &filterIterator{
		itr: withSize(itr, minSet.Len()),
		filter: func(item interface{}) bool {
			onode := item.(*Node)
			return onode.labels.HasAllSet(labels)
		},
	}
	return &nodeIterator{flt}
}
