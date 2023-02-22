package lpg

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkPropNonExistsGraph(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 1000; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
	}
	hasProp := func(node *Node) bool {
		return node.ForEachProperty(func(s string, i interface{}) bool {
			_, ok := node.GetProperty(s)
			return ok
		})
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			hasProp(nodes.Node())
		}
	}
}

func BenchmarkPropExistsGraph(b *testing.B) {
	g := NewGraph()
	nodes := make([]*Node, 0)
	for i := 0; i < 999; i++ {
		nodes = append(nodes, g.NewNode([]string{fmt.Sprint(i)}, nil))
	}
	nodes = append(nodes, g.NewNode([]string{fmt.Sprint(999)}, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h"}))
	s := nodes[len(nodes)-1]
	nodes[len(nodes)-1] = nodes[len(nodes)/2]
	nodes[len(nodes)/2] = s
	labels := []string{"a", "b", "c", "d"}
	for i := 0; i < len(nodes)-1; i++ {
		g.NewEdge(nodes[i], nodes[i+1], labels[i%4], nil)
	}
	hasProp := func(node *Node) bool {
		return node.ForEachProperty(func(s string, i interface{}) bool {
			_, ok := node.GetProperty(s)
			return ok
		})
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			hasProp(nodes.Node())
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
		nc := 0
		rn := rand.Intn(100)
		for nodes := g.GetNodes(); nodes.Next(); {
			nc++
			node := nodes.Node()
			for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
				if nc == rn {
					g.disconnect(edges.Edge())
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
				edgeHasLabel(edges.Edge(), labels[rand.Intn(len(labels))])
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
			if i == rand.Intn(len(nodes)) {
				g.NewEdge(nodes[i], nodes[i+1], label, map[string]interface{}{"a": "b", "c": "d", "e": "f", "g": "h"})
			} else {
				g.NewEdge(nodes[i], nodes[i+1], label, nil)
			}
		}
	}
	edgeHasProp := func(edge *Edge) bool {
		return edge.ForEachProperty(func(s string, i interface{}) bool {
			_, ok := edge.GetProperty(s)
			return ok
		})
	}
	for n := 0; n < b.N; n++ {
		for nodes := g.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			for edges := node.GetEdges(OutgoingEdge); edges.Next(); {
				edgeHasProp(edges.Edge())
			}
		}
	}
}
