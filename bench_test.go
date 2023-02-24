package lpg

import (
	"fmt"
	"testing"
)

func BenchmarkPropNonExistsGraph(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h"}))
	}
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			nodes.Node().GetProperty("z")
		}
	}
}

func BenchmarkPropExistsGraph(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		if i < 500 {
			nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h"}))
		} else {
			nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h", "z": "zz"}))
		}
	}
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			nodes.Node().GetProperty("z")
		}
	}
}

func BenchmarkDeleteEdge(b *testing.B) {
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
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			for {
				edgeRemoved := false
				for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
					edges.Edge().Remove()
					edgeRemoved = true
					break
				}
				if !edgeRemoved {
					break
				}
			}
		}
	}
}

func BenchmarkFindEdgeLabel(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for i := 0; i < len(nodes)-1; i++ {
		for _, label := range labels {
			g.NewEdge(nodes[i], nodes[i+1], label, nil)
		}
	}
	edgeHasLabel := func(edge *Edge, str string) bool {
		return edge.GetLabel() == str
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
				edgeHasLabel(edges.Edge(), "h")
			}
		}
	}
}

func BenchmarkFindEdgeProp(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for i := 0; i < len(nodes)-1; i++ {
		for _, label := range labels {
			if i < len(nodes)/2 {
				g.NewEdge(nodes[i], nodes[i+1], label, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h", "z": "zz"})
			} else {
				g.NewEdge(nodes[i], nodes[i+1], label, nil)
			}
		}
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
				edges.Edge().GetProperty("z")
			}
		}
	}
}
