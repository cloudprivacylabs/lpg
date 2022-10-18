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

// A Path can be a node, or node-edge-node...-edge-node sequence.
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

// Clone returns a copy of the path
func (p *Path) Clone() *Path {
	ret := Path{
		only: p.only,
		fwd:  make([]*Edge, len(p.fwd)+len(p.bk)),
	}
	j := 0
	for i := len(p.bk) - 1; i >= 0; i-- {
		ret.fwd[j] = p.bk[i]
		j++
	}
	copy(ret.fwd[j:], p.fwd)
	return &ret
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

// RemoveLast removes the last edge from the path. If the path has one
// edge, the path becomes a single node containing the source. If the
// path has only a node, path becomes empty
func (p *Path) RemoveLast() *Path {
	if p.only != nil {
		p.only = nil
		return p
	}
	if len(p.fwd)+len(p.bk) == 1 {
		if len(p.fwd) == 1 {
			p.only = p.fwd[0].GetFrom()
			p.fwd = p.fwd[:0]
		} else {
			p.only = p.bk[0].GetTo()
			p.bk = p.bk[:0]
		}
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
	if len(p.fwd)+len(p.bk) == 1 {
		if len(p.fwd) == 1 {
			p.only = p.fwd[0].GetTo()
			p.fwd = p.fwd[:0]
		} else {
			p.only = p.bk[0].GetFrom()
			p.bk = p.bk[:0]
		}
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
