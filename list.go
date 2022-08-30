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
	list.tail = edge

	edge.listElements[ix].next = nil
	if list.head == nil {
		list.head = edge
	}
	list.n++
}

func (list *edgeList) remove(edge *Edge, ix int) {
	if edge.listElements[ix].prev != nil {
		edge.listElements[ix].prev.listElements[ix].next = edge.listElements[ix].next
	} else {
		list.head = edge.listElements[ix].next
	}
	if edge.listElements[ix].next != nil {
		edge.listElements[ix].next.listElements[ix].prev = edge.listElements[ix].prev
	} else {
		list.tail = edge
	}
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
	}
	return e.current != nil

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
