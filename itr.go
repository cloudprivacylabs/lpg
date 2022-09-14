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

// An Iterator iterates the items of a collection
type Iterator interface {
	// Next moves to the next item in the iterator, and returns true if
	// move was successful. If there are no next items remaining, returns false
	Next() bool

	// Value returns the current item in the iterator. This is undefined
	// before the first call to Next, and after Next returns false.
	Value() interface{}

	// MaxSize returns an estimation of maximum number of elements. If unknown, returns -1
	MaxSize() int
}

// NodeIterator iterates nodes of an underlying list
type NodeIterator interface {
	Iterator
	// Returns the current node
	Node() *Node
}

// NodeSlice reads all the remaining items of a node iterator and returns them in a slice
func NodeSlice(in NodeIterator) []*Node {
	ret := make([]*Node, 0)
	for in.Next() {
		ret = append(ret, in.Node())
	}
	return ret
}

// EdgeIterator iterates the edges of an underlying list
type EdgeIterator interface {
	Iterator
	// Returns the current edge
	Edge() *Edge
}

// EdgeSlice reads all the remaining items of an edge iterator and returns them in a slice
func EdgeSlice(in EdgeIterator) []*Edge {
	ret := make([]*Edge, 0)
	for in.Next() {
		ret = append(ret, in.Edge())
	}
	return ret
}

// TargetNodes returns the target nodes of all edges
func TargetNodes(in EdgeIterator) []*Node {
	set := make(map[*Node]struct{})
	for in.Next() {
		set[in.Edge().GetTo()] = struct{}{}
	}
	ret := make([]*Node, 0, len(set))
	for x := range set {
		ret = append(ret, x)
	}
	return ret
}

// SourceNodes returns the source nodes of all edges
func SourceNodes(in EdgeIterator) []*Node {
	set := make(map[*Node]struct{})
	for in.Next() {
		set[in.Edge().GetFrom()] = struct{}{}
	}
	ret := make([]*Node, 0, len(set))
	for x := range set {
		ret = append(ret, x)
	}
	return ret
}

type emptyIterator struct{}

func (emptyIterator) Next() bool         { return false }
func (emptyIterator) Value() interface{} { return nil }
func (emptyIterator) MaxSize() int       { return 0 }

// filterIterator filters the items of the underlying iterator
type filterIterator struct {
	itr     Iterator
	filter  func(interface{}) bool
	current interface{}
}

func (itr *filterIterator) Next() bool {
	for itr.itr.Next() {
		itr.current = itr.itr.Value()
		if itr.filter(itr.current) {
			return true
		}
		itr.current = nil
	}
	return false
}

func (itr *filterIterator) Value() interface{} {
	return itr.current
}

func (itr *filterIterator) MaxSize() int { return itr.itr.MaxSize() }

// ProcIterator calls the function for every element
type procIterator struct {
	itr  Iterator
	proc func(interface{}) interface{}
}

func (itr *procIterator) Next() bool {
	return itr.itr.Next()
}

func (itr *procIterator) Value() interface{} {
	return itr.proc(itr.itr.Value())
}

func (itr *procIterator) MaxSize() int { return itr.itr.MaxSize() }

// makeUniqueIterator returns a filter iterator that will filter out duplicates
func makeUniqueIterator(itr Iterator) Iterator {
	seenItems := make(map[interface{}]struct{})
	return &filterIterator{
		itr: itr,
		filter: func(item interface{}) bool {
			if _, seen := seenItems[item]; seen {
				return false
			}
			seenItems[item] = struct{}{}
			return true
		},
	}
}

// funcIterator iterates through a set of underlying iterators obtained from a function
type funcIterator struct {
	// Returns a new iterator every time it is called. When returns nil, iteration stops
	iteratorFunc func() Iterator
	current      Iterator
}

func (itr *funcIterator) Next() bool {
	for {
		if itr.current != nil {
			if itr.current.Next() {
				return true
			}
			itr.current = nil
		}
		itr.current = itr.iteratorFunc()
		if itr.current == nil {
			return false
		}
	}
}

func (itr *funcIterator) Value() interface{} {
	return itr.current.Value()
}

func (itr *funcIterator) MaxSize() int { return -1 }

// MultiIterator returns an iterator that contatenates all the given iterators
func MultiIterator(iterators ...Iterator) Iterator {
	return &funcIterator{
		iteratorFunc: func() Iterator {
			if len(iterators) == 0 {
				return nil
			}
			ret := iterators[0]
			iterators = iterators[1:]
			return ret
		},
	}
}

// nodeIterator is a type-safe iterator for nodes
type nodeIterator struct {
	Iterator
}

func (n nodeIterator) Node() *Node {
	return n.Value().(*Node)
}

// edgeIterator is a type-safe iterator for edges
type edgeIterator struct {
	Iterator
}

func (n edgeIterator) Edge() *Edge {
	return n.Value().(*Edge)
}

type iteratorWithoutSize interface {
	Next() bool
	Value() interface{}
}

type iteratorWithSize struct {
	itr  iteratorWithoutSize
	size int
}

func (i *iteratorWithSize) Next() bool         { return i.itr.Next() }
func (i *iteratorWithSize) Value() interface{} { return i.itr.Value() }
func (i *iteratorWithSize) MaxSize() int       { return i.size }

func withSize(itr iteratorWithoutSize, size int) Iterator {
	return &iteratorWithSize{
		itr:  itr,
		size: size}
}

type listIterator struct {
	next, current *list.Element
	size          int
}

func (l *listIterator) Next() bool {
	l.current = l.next
	if l.next != nil {
		l.next = l.next.Next()
	}
	return l.current != nil
}

func (l *listIterator) Value() interface{} {
	return l.current.Value
}

func (l *listIterator) MaxSize() int { return l.size }
