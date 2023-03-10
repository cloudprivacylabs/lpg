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
	"strings"
)

// A Path can be a node, or node-edge-node...-edge-node sequence.
type Path struct {
	only *Node
	path []PathElement
}

type PathElement struct {
	Edge    *Edge
	Reverse bool
}

func NewPathFromElements(elements ...PathElement) *Path {
	return &Path{path: elements}
}

func NewPathElementsFromEdges(edges []*Edge) []PathElement {
	pe := make([]PathElement, len(edges))
	for idx, e := range edges {
		pe[idx].Edge = e
	}
	return pe
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
	p.path = nil
	return p
}

// Clear the path
func (p *Path) Clear() *Path {
	p.only = nil
	p.path = nil
	return p
}

// Last returns the last node of the path
func (p *Path) Last() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.path) != 0 {
		return p.path[len(p.path)-1].GetTargetNode()
	}
	return nil
}

// First returns the first node of the path
func (p *Path) First() *Node {
	if p.only != nil {
		return p.only
	}
	if len(p.path) != 0 {
		return p.path[0].GetSourceNode()
	}
	return nil
}

func (p PathElement) GetSourceNode() *Node {
	if p.Reverse {
		return p.Edge.GetTo()
	}
	return p.Edge.GetFrom()
}

func (p PathElement) GetTargetNode() *Node {
	if p.Reverse {
		return p.Edge.GetFrom()
	}
	return p.Edge.GetTo()
}

// Append an edge to the end of the path. The edge must be outgoing from the last node of the path
func (p *Path) Append(path ...PathElement) *Path {
	if len(path) == 0 {
		return p
	}
	if p.NumNodes() == 0 {
		p.path = make([]PathElement, len(path))
		copy(p.path, path)
		return p
	}
	last := p.Last()
	if last != nil && last != path[0].GetSourceNode() {
		// fmt.Println(last, path[0].GetSourceNode())
		panic("Appended edge is disconnected from path")
	}
	if p.only != nil {
		p.only = nil
	}
	for _, pe := range path {
		if pe.GetSourceNode() != last {
			panic("Appended edges are disconnected")
		}
		last = pe.GetTargetNode()
	}
	p.path = append(p.path, path...)
	return p
}

// can append a single node
func (p *Path) AppendPath(path Path) *Path {
	switch p.NumNodes() {
	case 0:
		copy(p.path, path.path)
		return p
	case 1:
		if p.only != nil {
			p.only = nil
		}
	default:
		switch path.NumNodes() {
		case 0:
			return p
		case 1:
			if path.only != nil {
				path.only = nil
			}
		}
	}
	return p.Append(path.path...)
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
	return p.path[n-1].GetTargetNode()
}

// String returns the path as a string
func (p *Path) String() string {
	sb := strings.Builder{}
	for _, p := range p.path {
		sb.WriteString(p.Edge.GetFrom().String() + "->" + p.Edge.GetTo().String() + " ")
		if p.Reverse {
			sb.WriteString("Reverse ")
		}
		// p.Edge.GetFrom().String()
		// if p.Reverse {
		// 	sb.WriteString(p.GetSourceNode().String() + "<-" + p.GetTargetNode().String())
		// } else {
		// 	sb.WriteString(p.GetSourceNode().String() + "->" + p.GetTargetNode().String())
		// }
	}
	return sb.String()
}

// HasPrefix return if all edges of p1 exist in path
func (p *Path) HasPrefix(p1 []PathElement) bool {
	if p.NumNodes() < 2 {
		return false
	}
	if len(p1) == 0 {
		return true
	}
	if len(p1) > p.NumEdges() {
		return false
	}
	for path1Idx, e1 := range p1 {
		// fmt.Printf("%p %p", e1.Edge, p.path[path1Idx].Edge)
		// fmt.Println()
		// fmt.Println(e1.Edge == p.path[path1Idx].Edge)
		if e1 != p.path[path1Idx] {
			return false
		}
	}
	return true
}

// HasPrefixPaths returns if all paths elements of p1 exist in path
func (p *Path) HasPrefixPath(p1 *Path) bool {
	switch p.NumNodes() {
	case 0:
		return p1.NumNodes() == 0
	case 1:
		switch p1.NumNodes() {
		case 0:
			return true
		case 1:
			return p.only == p1.only
		default:
			return false
		}
	default:
		switch p1.NumNodes() {
		case 0:
			return true
		case 1:
			return p.First() == p1.First()
		default:
			return p.HasPrefix(p1.path)
		}
	}
}

// Slice returns a copy of p partitioned by the args start and end index
func (p *Path) Slice(startNodeIndex, endNodeIndex int) *Path {
	if p.NumNodes() == 0 {
		panic("index error")
	}
	if startNodeIndex == 0 && endNodeIndex == 0 {
		return &Path{
			only: p.First(),
		}
	}
	if startNodeIndex > p.NumNodes() {
		panic("start index greater than length of num nodes")
	}
	if endNodeIndex > p.NumNodes() || endNodeIndex == -1 {
		endNodeIndex = p.NumNodes()
	}
	if endNodeIndex < startNodeIndex {
		panic("")
	}
	if startNodeIndex == endNodeIndex {
		return &Path{
			only: p.path[startNodeIndex].GetSourceNode(),
		}
	}
	if startNodeIndex == p.NumNodes()-1 {
		return &Path{
			only: p.Last(),
		}
	}
	pth := make([]PathElement, endNodeIndex-startNodeIndex-1)
	copy(pth, p.path[startNodeIndex:endNodeIndex-1])
	return &Path{path: pth}
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
