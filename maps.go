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

// An NodeMap stores nodes indexed by node labels
type NodeMap struct {
	// m[string]*fastSet
	m        *linkedhashmap.Map
	nolabels fastSet
}

func NewNodeMap() *NodeMap {
	nm := &NodeMap{
		m: linkedhashmap.New(),
	}
	nm.nolabels.init()
	return nm
}

func (nm *NodeMap) Replace(node *Node, oldLabels, newLabels StringSet) {
	if oldLabels.Len() == 0 {
		if newLabels.Len() == 0 {
			return
		}
		nm.nolabels.remove(node.id, node)
	}
	if newLabels.Len() == 0 {
		nm.nolabels.add(node.id, node)
		return
	}
	var set *fastSet
	// Process removed labels
	for label := range oldLabels.M {
		if !newLabels.Has(label) {
			v, found := nm.m.Get(label)
			if !found {
				continue
			}
			set = v.(*fastSet)
			set.remove(node.id, node)
			if set.size() == 0 {
				nm.m.Remove(label)
			}
		}
	}
	// Process added labels
	for label := range newLabels.M {
		if !oldLabels.Has(label) {
			v, found := nm.m.Get(label)
			if !found {
				set = newFastSet()
				nm.m.Put(label, set)
			} else {
				set = v.(*fastSet)
			}
			set.add(node.id, node)
		}
	}
}

func (nm *NodeMap) Add(node *Node) {
	if node.labels.Len() == 0 {
		nm.nolabels.add(node.id, node)
		return
	}

	var set *fastSet
	for label := range node.labels.M {
		v, found := nm.m.Get(label)
		if !found {
			set = newFastSet()
			nm.m.Put(label, set)
		} else {
			set = v.(*fastSet)
		}
		set.add(node.id, node)
	}
}

func (nm NodeMap) Remove(node *Node) {
	if node.labels.Len() == 0 {
		nm.nolabels.remove(node.id, node)
		return
	}
	var set *fastSet
	for label := range node.labels.M {
		v, found := nm.m.Get(label)
		if !found {
			continue
		}
		set = v.(*fastSet)
		set.remove(node.id, node)
		if set.size() == 0 {
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
	set := itr.labels.Value().(*fastSet)
	setItr := set.iterator()
	itr.current = nodeIterator{withSize(setItr, -1)}
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
	return nodeIterator{
		MultiIterator(
			&filterIterator{
				itr: withSize(nmIterator, -1),
				filter: func(node interface{}) bool {
					onode := node.(*Node)
					nSeen := 0
					for _, l := range nmIterator.seenLabels {
						if onode.labels.Has(l) {
							nSeen++
							if nSeen > 1 {
								return false
							}
						}
					}
					return true
				},
			},
			nm.nolabels.iterator(),
		),
	}
}

func (nm NodeMap) IteratorAllLabels(labels StringSet) NodeIterator {
	// Find the smallest map element, iterate that
	var minSet *fastSet
	for label := range labels.M {
		v, found := nm.m.Get(label)
		if !found {
			return nodeIterator{&emptyIterator{}}
		}
		mp := v.(*fastSet)
		if minSet == nil || minSet.size() > mp.size() {
			minSet = mp
		}
	}
	itr := minSet.iterator()
	flt := &filterIterator{
		itr: withSize(itr, minSet.size()),
		filter: func(item interface{}) bool {
			onode := item.(*Node)
			return onode.labels.HasAllSet(labels)
		},
	}
	return nodeIterator{flt}
}
