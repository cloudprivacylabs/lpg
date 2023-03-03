package lpg

import (
	"fmt"
	"testing"

	sm "github.com/bserdar/slicemap"
)

func TestCollectAllPaths(t *testing.T) {
	graph, nodes := GetLineGraphWithSelfLoops(10, true)
	nodes[2].SetProperty("key", "value")
	nodes[2].SetLabels(NewStringSet("node2"))
	nodes[3].SetLabels(NewStringSet("node3"))
	nodes[3].SetProperty("key", "value")
	acc := &DefaultMatchAccumulator{}
	// allEdgePaths := make([][]*Edge, 0)
	set := sm.SliceMap[*Edge, struct{}]{}
	CollectAllPaths(graph, nodes[2], nodes[2].GetEdges(OutgoingEdge), func(e *Edge) bool { return true }, OutgoingEdge, 1, -1, func(e []*Edge, n *Node) bool {
		if _, seen := set.Get(e); seen {
			return false
		}
		set.Put(e, struct{}{})
		// allEdgePaths = append(allEdgePaths, e)
		acc.Paths = append(acc.Paths, e)
		return true
	})
	set.ForEach(func(k []*Edge, s struct{}) bool {
		for _, e := range k {
			fmt.Println(e.GetFrom(), e.GetTo())

		}
		return true
	})
	// filteredPaths := make(map[*[]*Edge]struct{})
	// for idx := 0; idx < len(allEdgePaths); idx++ {
	// 	for ix := 1; ix < len(allEdgePaths); ix++ {
	// 		if allEdgePaths[ix][0] == allEdgePaths[idx][0] {
	// 			// fmt.Println("hererer")
	// 			continue
	// 		}
	// 	}
	// 	filteredPaths[&allEdgePaths[idx]] = struct{}{}
	// }
	// uniquePaths := make([]*Edge, 0)
	// contains := func(e *Edge) bool {
	// 	for _, path := range uniquePaths {
	// 		if e == path {
	// 			return true
	// 		}
	// 	}
	// 	if e.GetFrom() == e.GetTo() {
	// 		// fmt.Println(e.GetFrom(), e.GetTo())
	// 	}
	// 	// return e.GetFrom() == e.GetTo()
	// 	return false
	// }
	// for ed := range filteredPaths {
	// 	fmt.Println((*ed)[0].GetFrom(), (*ed)[0].GetTo())
	// }
	fmt.Println(len(acc.Paths))
	// for _, se := range allEdgePaths {
	// 	for _, e := range se {
	// 		if contains(e) {
	// 			continue
	// 		}
	// 		uniquePaths = append(uniquePaths, e)
	// 	}
	// }
	// fmt.Println(len(uniquePaths))
	// fmt.Println(uniquePaths)
	// for _, path := range uniquePaths {
	// 	fmt.Println(path.GetFrom(), path.GetTo())
	// 	// fmt.Println()
	// }
	// for _, p := range acc.Paths {
	// 	for _, path := range p.([]*Edge) {
	// 		fmt.Println(path.GetFrom(), path.GetTo())
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println(acc.Paths...)
	t.Fail()
}
