package lpg

import (
	"fmt"
	"testing"
)

func TestCollectAllPaths(t *testing.T) {
	graph, nodes := GetLineGraphWithSelfLoops(2, true)
	nodes[0].SetProperty("key", "value")
	nodes[0].SetLabels(NewStringSet("node2"))
	nodes[1].SetLabels(NewStringSet("node3"))
	nodes[1].SetProperty("key", "value")
	acc := &DefaultMatchAccumulator{}
	CollectAllPaths(graph, nodes[1], nodes[1].GetEdges(OutgoingEdge), func(e *Edge) bool { return true }, OutgoingEdge, 1, -1, func(e *Path, n *Node) bool {
		acc.Paths = append(acc.Paths, e)
		return true
	})
	fmt.Println(len(acc.Paths))
	// for _, p := range acc.Paths {
	// 	for _, path := range p.(*Path).path {
	// 		fmt.Println(path)
	// 		fmt.Printf("%p %p", path.GetSourceNode(), path.GetTargetNode())
	// 		fmt.Println()
	// 	}
	// 	fmt.Println()
	// }
	t.Fail()
}
