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
	path []PathElement
}

type PathElement struct {
	Edge      *Edge
	isReverse bool
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
		path: make([]PathElement, len(p.path)),
	}
	copy(ret.path, p.path)
	return &ret
}

// SetOnlyNode sets the path to a single node
func (p *Path) SetOnlyNode(node *Node) *Path {
	p.only = node
	return p
}

// Clear the path
func (p *Path) Clear() *Path {
	p.only = nil
	p.path = make([]PathElement, 0)
	return p
}

// Last returns the last node of the path
func (p *Path) Last() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.path) != 0 {
		return p.path[len(p.path)-1].Edge.GetTo()
	}
	return nil
}

// First returns the first node of the path
func (p *Path) First() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.path) != 0 {
		return p.path[0].Edge.GetFrom()
	}
	return nil
}

func (p *Path) GetSourceNode() *Node {
	if p.path[0].isReverse {
		return p.path[0].Edge.GetTo()
	}
	return p.path[0].Edge.GetFrom()
}

func (p *Path) GetTargetNode() *Node {
	if p.path[1].isReverse {
		return p.path[1].Edge.GetFrom()
	}
	return p.path[1].Edge.GetTo()
}

// Append an edge to the end of the path. The edge must be outgoing from the last node of the path
func (p *Path) Append(path []PathElement) *Path {
	last := p.Last()
	if last != nil && last != path[0].Edge.GetFrom() {
		panic("Appended edge is disconnected from path")
	}
	if p.only != nil {
		p.only = nil
	}
	for i := range path {
		if i == 0 {
			continue
		}
		if path[i].Edge.GetFrom() != path[i-1].Edge.GetTo() {
			panic("Appended edges are disconnected")
		}
	}
	p.path = append(p.path, path...)
	return p
}

// GetEdge returns the nth edge
func (p *Path) GetEdge(n int) *Edge {
	if p.only != nil {
		return nil
	}
	if n < len(p.path) {
		return p.path[n].Edge
	}
	return nil
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

// return if all edges of p2 exist in p1
func (p *Path) IsPrefix(p1, p2 []PathElement) bool {
	if len(p2) == 0 {
		return false
	}
	if len(p2) > len(p1) {
		return false
	}
	for path2Idx, e2 := range p2 {
		if e2 != p1[path2Idx] {
			return false
		}
	}
	return true
}

func (p *Path) Slice(startNodeIndex, endNodeIndex int) []PathElement {
	if startNodeIndex > 0 && startNodeIndex < len(p.path) && endNodeIndex > startNodeIndex && endNodeIndex <= len(p.path) {
		return p.path[startNodeIndex : endNodeIndex+1]
	}
	return nil
}

// RemoveLast removes the last edge from the path. If the path has one
// edge, the path becomes a single node containing the source. If the
// path has only a node, path becomes empty
func (p *Path) RemoveLast() *Path {
	if p.only != nil {
		p.only = nil
		return p
	}
	if len(p.path) > 0 {
		p.path[len(p.path)-1].Edge.Remove()
		p.path = p.path[:len(p.path)-1]
	}
	// if len(p.fwd)+len(p.bk) == 1 {
	// 	if len(p.fwd) == 1 {
	// 		p.only = p.fwd[0].GetFrom()
	// 		p.fwd = p.fwd[:0]
	// 	} else {
	// 		p.only = p.bk[0].GetTo()
	// 		p.bk = p.bk[:0]
	// 	}
	// 	return p
	// }
	// if len(p.fwd) != 0 {
	// 	p.fwd = p.fwd[:len(p.fwd)-1]
	// 	return p
	// }
	// if len(p.bk) != 0 {
	// 	p.fwd = make([]*Edge, len(p.bk))
	// 	for i := range p.bk {
	// 		p.fwd[len(p.fwd)-1-i] = p.bk[i]
	// 	}
	// }
	return p
}

// RemoveFirst removes the first edge from the graph. If the path has
// only a node, path becomes empty
func (p *Path) RemoveFirst() *Path {
	if p.only != nil {
		p.only = nil
		return p
	}
	if len(p.path) > 0 {
		p.path[0].Edge.Remove()
		p.path = p.path[1:]
	}
	// if len(p.fwd)+len(p.bk) == 1 {
	// 	if len(p.fwd) == 1 {
	// 		p.only = p.fwd[0].GetTo()
	// 		p.fwd = p.fwd[:0]
	// 	} else {
	// 		p.only = p.bk[0].GetFrom()
	// 		p.bk = p.bk[:0]
	// 	}
	// 	return p
	// }
	// if len(p.bk) != 0 {
	// 	p.bk = p.bk[:len(p.bk)-1]
	// 	return p
	// }
	// if len(p.fwd) != 0 {
	// 	p.fwd = p.fwd[1:]
	// }
	return p
}

func (p *Path) NumEdges() int {
	return len(p.path)
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
	return p.only == nil && len(p.path) == 0
}
