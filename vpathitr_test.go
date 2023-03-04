package lpg

import (
	"fmt"
	"testing"
)

func TestCollectAllPaths(t *testing.T) {
	graph, nodes := GetLineGraphWithSelfLoops(10, true)
	nodes[2].SetProperty("key", "value")
	nodes[2].SetLabels(NewStringSet("node2"))
	nodes[3].SetLabels(NewStringSet("node3"))
	nodes[3].SetProperty("key", "value")
	acc := &DefaultMatchAccumulator{}
	CollectAllPaths(graph, nodes[2], nodes[2].GetEdges(OutgoingEdge), func(e *Edge) bool { return true }, OutgoingEdge, 1, -1, func(e []*Edge, n *Node) bool {
		if len(e) == 2 {
			if e[0].GetFrom() == e[0].GetTo() {
				if e[1].GetFrom() == e[0].GetTo() {
					if e[1].GetFrom() == e[1].GetTo() {
						fmt.Printf("%p ", (e))
						fmt.Println(e)
					}
				}
			}
		}
		acc.Paths = append(acc.Paths, e)
		return true
	})
	fmt.Println(len(acc.Paths))
	for _, p := range acc.Paths {
		for _, path := range p.([]*Edge) {
			// fmt.Println(path.GetFrom(), path.GetTo())
			fmt.Printf("%p ", (path))
		}
		fmt.Println()
	}
	t.Fail()
}
