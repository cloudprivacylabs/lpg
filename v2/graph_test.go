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
	"fmt"
	"testing"
)

func TestGraphCRUD(t *testing.T) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], "e", nil)
	}

	if len(NodeSlice(g.GetNodes())) != len(nodes) {
		t.Errorf("Wrong node count")
	}
	if g.NumNodes() != len(nodes) {
		t.Errorf("Wrong numNodes")
	}
	nodes[2].DetachAndRemove()
	if len(NodeSlice(g.GetNodes())) != len(nodes)-1 {
		t.Errorf("Wrong node count")
	}
	if g.NumNodes() != len(nodes)-1 {
		t.Errorf("Wrong numNodes")
	}
}

func BenchmarkAddNode(b *testing.B) {
	g := NewGraph()
	for n := 0; n < b.N; n++ {
		g.NewNode([]string{"a", "b", "c"}, map[string]interface{}{"a": "b", "c": "d"})
	}
}

func benchmarkItrNodes(numNodes int, b *testing.B) {
	g := NewGraph()
	var x *Node
	for i := 0; i < numNodes; i++ {
		g.NewNode([]string{"a", "b", "c"}, map[string]interface{}{"a": "b", "c": "d"})
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			x = nodes.Node()
		}
	}
	_ = x
}

func BenchmarkItrNodes1000(b *testing.B)  { benchmarkItrNodes(1000, b) }
func BenchmarkItrNodes10000(b *testing.B) { benchmarkItrNodes(10000, b) }

func benchmarkItrNodesViaIndex(numNodes int, b *testing.B) {
	g := NewGraph()
	var x *Node
	for i := 0; i < numNodes; i++ {
		g.NewNode([]string{"a", "b", "c"}, map[string]interface{}{"a": "b", "c": "d"})
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.index.nodesByLabel.Iterator(); nodes.Next(); {
			x = nodes.Node()
		}
	}
	_ = x
}

func BenchmarkItrNodesViaIndex1000(b *testing.B)  { benchmarkItrNodesViaIndex(1000, b) }
func BenchmarkItrNodesViaIndex10000(b *testing.B) { benchmarkItrNodesViaIndex(10000, b) }

func BenchmarkCreateEdge(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d"}

	for n := 0; n < b.N; n++ {
		for i := 0; i < len(nodes)-1; i++ {
			g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
		}
	}
}

func BenchmarkItrAllEdge(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
	}
	var edge *Edge

	for n := 0; n < b.N; n++ {
		for edges := g.GetEdges(); edges.Next(); {
			edge = edges.Edge()
		}
	}
	_ = edge
}

func BenchmarkItrNodeEdges(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		for _, label := range labels {
			g.NewEdge(nodes[i], nodes[i+1], label, nil)
		}
	}
	var edge *Edge

	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
				edge = edges.Edge()
			}
		}
	}
	_ = edge
}
