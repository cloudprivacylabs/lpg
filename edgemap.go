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
	edges *list.List
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
	return &edgeMap{
		edgeLabelLists: list.New(),
		labelMap:       make(map[string]*list.Element),
	}
}

func (em *edgeMap) add(edge *Edge) {
	var ell *edgeLabelList

	el := em.labelMap[edge.label]
	if el == nil {
		ell = &edgeLabelList{
			label: edge.label,
			edges: list.New(),
		}
		ell.el = em.edgeLabelLists.PushBack(ell)
		em.labelMap[edge.label] = ell.el
	} else {
		ell = el.Value.(*edgeLabelList)
	}
	edge.el = ell.edges.PushBack(edge)
	em.n++
}

func (em *edgeMap) remove(edge *Edge) {
	el := em.labelMap[edge.label]
	if el == nil {
		return
	}
	ell := el.Value.(*edgeLabelList)
	ell.edges.Remove(edge.el)
	if ell.edges.Len() == 0 {
		em.edgeLabelLists.Remove(ell.el)
		delete(em.labelMap, edge.label)
	}
	em.n--
}

func (em *edgeMap) isEmpty() bool { return em.n == 0 }

func (em *edgeMap) size() int { return em.n }

func (em edgeMap) iterator() EdgeIterator {
	ret := &allEdgesItr{
		labelListCurrent: em.edgeLabelLists.Front(),
	}
	if ret.labelListCurrent != nil {
		ret.labelListNext = ret.labelListCurrent.Next()
		ret.next = ret.labelListCurrent.Value.(*edgeLabelList).edges.Front()
	}
	ret.size = em.n
	return ret
}

func (em edgeMap) iteratorLabel(label string) EdgeIterator {
	l := em.labelMap[label]
	if l == nil {
		return &edgeIterator{&emptyIterator{}}
	}
	ell := l.Value.(*edgeLabelList)
	return &edgeIterator{&listIterator{next: ell.edges.Front(), size: ell.edges.Len()}}
}

func (em edgeMap) iteratorAnyLabel(labels StringSet) EdgeIterator {
	strings := labels.Slice()
	return &edgeIterator{&funcIterator{
		iteratorFunc: func() Iterator {
			for len(strings) != 0 {
				if _, found := em.labelMap[strings[0]]; !found {
					strings = strings[1:]
					continue
				}
				itr := em.iteratorLabel(strings[0])
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
	next, current                   *list.Element
	size                            int
}

func (itr *allEdgesItr) Next() bool {
top:
	itr.current = itr.next
	if itr.next != nil {
		itr.next = itr.next.Next()
	} else {
		itr.labelListCurrent = itr.labelListNext
		if itr.labelListNext != nil {
			itr.labelListNext = itr.labelListNext.Next()
		}
		if itr.labelListCurrent != nil {
			itr.next = itr.labelListCurrent.Value.(*edgeLabelList).edges.Front()
			goto top
		}
		return false
	}
	return itr.current != nil
}

func (itr *allEdgesItr) Value() interface{} {
	return itr.current.Value
}

func (itr *allEdgesItr) Edge() *Edge {
	return itr.current.Value.(*Edge)
}

func (itr *allEdgesItr) MaxSize() int { return itr.size }
