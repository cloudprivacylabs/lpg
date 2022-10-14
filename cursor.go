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

type Path struct {
	only *Node
	// forward path, fwd[i] -> fwd[i+1]
	fwd []*Edge
	// backward path, bk[i+1]->bk[i]
	bk []*Edge
}

// PathFromNode creates a path containing a single node
func PathFromNode(node *Node) *Path {
	if node == nil {
		panic("Nil node in path")
	}
	return &Path{
		only: node,
	}
}

// SetOnlyNode sets the path to a single node
func (p *Path) SetOnlyNode(node *Node) *Path {
	p.only = node
	p.fwd = nil
	p.bk = nil
	return p
}

// Clear the path
func (p *Path) Clear() *Path {
	p.only = nil
	p.fwd = nil
	p.bk = nil
	return p
}

// Last returns the last node of the path
func (p *Path) Last() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.fwd) != 0 {
		return p.fwd[len(p.fwd)-1].GetTo()
	}
	if len(p.bk) != 0 {
		return p.bk[0].GetTo()
	}
	return nil
}

// First returns the first node of the path
func (p *Path) First() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.bk) != 0 {
		return p.bk[len(p.bk)-1].GetFrom()
	}
	if len(p.fwd) != 0 {
		return p.fwd[0].GetFrom()
	}
	return nil
}

// Append an edge to the end of the path. The edge must be outgoing from the last node of the path
func (p *Path) Append(edge ...*Edge) *Path {
	last := p.Last()
	if last != nil && last != edge[0].GetFrom() {
		panic("Appended edge is disconnected from path")
	}
	if p.only != nil {
		p.only = nil
	}
	for i := range edge {
		if i == 0 {
			continue
		}
		if edge[i].GetFrom() != edge[i-1].GetTo() {
			panic("Appended edges are disconnected")
		}
	}
	p.fwd = append(p.fwd, edge...)
	return p
}

// Prepend an edge to tbe beginning of the path. The edge must be
// incoming to the first node of the path
func (p *Path) Prepend(edge ...*Edge) *Path {
	first := p.First()
	if first != nil && first != edge[0].GetTo() {
		panic("Prepended edge is disconnected from path")
	}
	if p.only != nil {
		p.only = nil
	}
	for i := range edge {
		if i == 0 {
			continue
		}
		if edge[i].GetTo() != edge[i-1].GetFrom() {
			panic("Prepended edges are disconnected")
		}
	}
	p.bk = append(p.bk, edge...)
	return p
}

// GetEdge returns the nth edge
func (p *Path) GetEdge(n int) *Edge {
	if p.only != nil {
		return nil
	}
	if n < len(p.bk) {
		return p.bk[len(p.bk)-1-n]
	}
	n -= len(p.bk)
	return p.fwd[n]
}

// GetNode returns the nth node
func (p *Path) GetNode(n int) *Node {
	if p.only != nil && n == 0 {
		return p.only
	}
	if n == 0 {
		e := p.First()
		if e != nil {
			return e
		}
		panic("Invalid node index")
	}
	e := p.GetEdge(n - 1)
	return e.GetTo()
}

// RemoveLast removes the last edge from the path. If the path has
// only a node, path becomes empty
func (p *Path) RemoveLast() *Path {
	if p.only != nil {
		p.only = nil
		return p
	}
	if len(p.fwd) != 0 {
		p.fwd = p.fwd[:len(p.fwd)-1]
		return p
	}
	if len(p.bk) != 0 {
		p.fwd = make([]*Edge, len(p.bk))
		for i := range p.bk {
			p.fwd[len(p.fwd)-1-i] = p.bk[i]
		}
	}
	return p
}

// RemoveFirst removes the first edge from the graph. If the path has
// only a node, path becomes empty
func (p *Path) RemoveFirst() *Path {
	if p.only != nil {
		p.only = nil
		return p
	}
	if len(p.bk) != 0 {
		p.bk = p.bk[:len(p.bk)-1]
		return p
	}
	if len(p.fwd) != 0 {
		p.fwd = p.fwd[1:]
	}
	return p
}

func (p *Path) NumEdges() int {
	return len(p.bk) + len(p.fwd)
}

func (p *Path) NumNodes() int {
	if p.only != nil {
		return 1
	}
	n := p.NumEdges()
	if n == 0 {
		return 0
	}
	return n + 1
}

func (p *Path) IsEmpty() bool {
	return p.only == nil && len(p.bk) == 0 && len(p.fwd) == 0
}

// Cursor is a convenience class to move around a graph
type Cursor struct {
	node *Node
	path *Path
}

// Sets the node cursor is pointing at.
func (c *Cursor) Set(node *Node) *Cursor {
	c.node = node
	return c
}

// StartPath starts a new path at the current node. Panics if the node is not valid.
func (c *Cursor) StartPath() *Cursor {
	c.path = PathFromNode(c.node)
	return c
}

func (c *Cursor) PushToPath() *Cursor {
	return c
}

func (c *Cursor) PopFromPath() *Cursor {
	return c
}

// Edges returns an iterator of edges of the current node. If the
// current node is invalid, returns an empty iterator
func (c *Cursor) Edges(dir EdgeDir) EdgeIterator {
	if c.node == nil {
		return &edgeIterator{emptyIterator{}}
	}
	return c.node.GetEdges(dir)
}

// EdgesWith returns an iterator of edges of the current node with the
// given label. If the current node is invalid, returns an empty
// iterator
func (c *Cursor) EdgesWith(dir EdgeDir, label string) EdgeIterator {
	if c.node == nil {
		return &edgeIterator{emptyIterator{}}
	}
	return c.node.GetEdgesWithLabel(dir, label)
}

// Nodes returns an iterator over the nodes that are reachable by a
// single edge from the current node. If the current node is not
// valid, returns an empty iterator
func (c *Cursor) Nodes(dir EdgeDir) NodeIterator {
	if c.node == nil {
		return &nodeIterator{emptyIterator{}}
	}
	return newNodesFromEdges(c.node.GetEdges(dir), dir)
}

// NodesWith returns an iterator over the nodes that are reachable by
// a single edge with the given label from the current node. If the
// current node is not valid, returns an empty iterator
func (c *Cursor) NodesWith(dir EdgeDir, label string) NodeIterator {
	if c.node == nil {
		return &nodeIterator{emptyIterator{}}
	}
	return newNodesFromEdges(c.node.GetEdgesWithLabel(dir, label), dir)
}

// Forward returns an edge iterator for outgoing edges of the current
// node. If the current node is invalid, returns an empty iterator
func (c *Cursor) Forward() EdgeIterator { return c.Edges(OutgoingEdge) }

// Backward returns an edge iterator for incoming edges of the current
// node. If the current node is invalid, returns an empty iterator
func (c *Cursor) Backward() EdgeIterator { return c.Edges(IncomingEdge) }

// ForwardWith returns an edge iterator for outgoing edges of the current
// node with the given label. If the current node is invalid, returns an empty iterator
func (c *Cursor) ForwardWith(label string) EdgeIterator { return c.EdgesWith(OutgoingEdge, label) }

// BackwardWith returns an edge iterator for incoming edges of the
// current node with the given label. If the current node is invalid,
// returns an empty iterator
func (c *Cursor) BackwardWith(label string) EdgeIterator { return c.EdgesWith(IncomingEdge, label) }
