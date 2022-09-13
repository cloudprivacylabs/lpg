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

type edgeLabelList struct {
	label string
	edges edgeList
	el    *list.Element
}

// An edgeMap stores edges indexed by edge label
type edgeMap struct {
	// list of labelstructs
	edgeLabelLists *list.List
	// Map of labels -> *edgeLabelList
	labelMap map[string]*list.Element
	n        int
}

func newEdgeMap() *edgeMap {
	em := &edgeMap{}
	em.init()
	return em
}

func (em *edgeMap) init() {
	em.edgeLabelLists = list.New()
	em.labelMap = make(map[string]*list.Element)
}

func (em *edgeMap) add(edge *Edge, listIndex int) {
	var ell *edgeLabelList

	el := em.labelMap[edge.label]
	if el == nil {
		ell = &edgeLabelList{
			label: edge.label,
		}
		ell.el = em.edgeLabelLists.PushBack(ell)
		em.labelMap[edge.label] = ell.el
	} else {
		ell = el.Value.(*edgeLabelList)
	}
	em.n++
	ell.edges.add(edge, listIndex)
}

func (em *edgeMap) remove(edge *Edge, listIndex int) {
	el := em.labelMap[edge.label]
	if el == nil {
		return
	}
	ell := el.Value.(*edgeLabelList)
	ell.edges.remove(edge, listIndex)
	if ell.edges.n == 0 {
		em.edgeLabelLists.Remove(ell.el)
		delete(em.labelMap, edge.label)
	}
	em.n--
}

func (em *edgeMap) isEmpty() bool { return em.n == 0 }

func (em *edgeMap) size() int { return em.n }

func (em edgeMap) iterator(listIndex int) EdgeIterator {
	ret := &allEdgesItr{
		labelListCurrent: em.edgeLabelLists.Front(),
		ix:               listIndex,
	}
	if ret.labelListCurrent != nil {
		ret.labelListNext = ret.labelListCurrent.Next()
		ret.next = ret.labelListCurrent.Value.(*edgeLabelList).edges.head
	}
	ret.size = em.n
	return ret
}

func (em edgeMap) iteratorLabel(label string, listIndex int) EdgeIterator {
	l := em.labelMap[label]
	if l == nil {
		return edgeIterator{&emptyIterator{}}
	}
	ell := l.Value.(*edgeLabelList)
	return edgeIterator{&edgeListIterator{next: ell.edges.head, n: ell.edges.n, ix: listIndex}}
}

func (em edgeMap) iteratorAnyLabel(labels StringSet, listIndex int) EdgeIterator {
	strings := labels.Slice()
	return edgeIterator{&funcIterator{
		iteratorFunc: func() Iterator {
			for len(strings) != 0 {
				if _, found := em.labelMap[strings[0]]; !found {
					strings = strings[1:]
					continue
				}
				itr := em.iteratorLabel(strings[0], listIndex)
				strings = strings[1:]
				return withSize(itr, -1)
			}
			return nil
		},
	},
	}
}

type allEdgesItr struct {
	labelListNext, labelListCurrent *list.Element
	next, current                   *Edge
	size                            int
	ix                              int
}

func (itr *allEdgesItr) Next() bool {
top:
	itr.current = itr.next
	if itr.next != nil {
		itr.next = itr.next.listElements[itr.ix].next
		return true
	}
	itr.labelListCurrent = itr.labelListNext
	if itr.labelListNext != nil {
		itr.labelListNext = itr.labelListNext.Next()
	}
	if itr.labelListCurrent != nil {
		itr.next = itr.labelListCurrent.Value.(*edgeLabelList).edges.head
		goto top
	}
	return false
}

func (itr *allEdgesItr) Value() interface{} {
	return itr.current
}

func (itr *allEdgesItr) Edge() *Edge {
	return itr.current
}

func (itr *allEdgesItr) MaxSize() int { return itr.size }
