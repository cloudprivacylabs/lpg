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

type edgeList struct {
	head *Edge
	tail *Edge
	n    int
}

type edgeElement struct {
	next *Edge
	prev *Edge
}

func (list *edgeList) add(edge *Edge, ix int) {
	edge.listElements[ix].prev = list.tail
	if list.tail != nil {
		list.tail.listElements[ix].next = edge
	}
	edge.listElements[ix].next = nil
	list.tail = edge

	if list.head == nil {
		list.head = edge
	}
	list.n++
}

func (list *edgeList) remove(edge *Edge, ix int) {
	el := &edge.listElements[ix]

	if el.prev != nil {
		el.prev.listElements[ix].next = el.next
	} else {
		list.head = el.next
	}
	if el.next != nil {
		el.next.listElements[ix].prev = el.prev
	} else {
		list.tail = el.prev
	}
	el.next = nil
	el.prev = nil
	list.n--
}

type edgeListIterator struct {
	current, next *Edge
	n             int
	ix            int
}

func (e *edgeListIterator) Next() bool {
	e.current = e.next
	if e.next != nil {
		e.next = e.next.listElements[e.ix].next
		return true
	}
	return false
}

func (e *edgeListIterator) Value() interface{} {
	return e.current
}

func (e *edgeListIterator) Edge() *Edge {
	return e.current
}

func (e *edgeListIterator) MaxSize() int {
	return e.n
}

type nodeList struct {
	head *Node
	tail *Node
	n    int
}

type nodeElement struct {
	next *Node
	prev *Node
}

func (list *nodeList) add(node *Node) {
	node.prev = list.tail
	if list.tail != nil {
		list.tail.next = node
	}
	list.tail = node

	if list.head == nil {
		list.head = node
	}
	list.n++
}

func (list *nodeList) remove(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		list.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		list.tail = node
	}
	list.n--
}

type nodeListIterator struct {
	current, next *Node
	n             int
	ix            int
}

func (n *nodeListIterator) Next() bool {
	n.current = n.next
	if n.next != nil {
		n.next = n.next.next
		return true
	}
	return false
}

func (n *nodeListIterator) Value() interface{} {
	return n.current
}

func (n *nodeListIterator) Node() *Node {
	return n.current
}

func (n *nodeListIterator) MaxSize() int {
	return n.n
}
