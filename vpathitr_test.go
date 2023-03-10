package lpg

import (
	"fmt"
	"testing"
)

// match (:b)-[e*]-() return e
/*
27 = a
28 = b
29 = c

(b->c)
(b->c) (c->c)
(a->b)
(a->b) (a->a)
(b->b)
(b->b) (b->c)
(b->b) (b->c) (c->c)
(b->b) (a->b)
(b->b) (a->b) (a->a)
*/

/*
(:a)->(:b) Reverse
(:a)->(:b) Reverse (:a)->(:a)
(:b)->(:b)
(:b)->(:b) (:a)->(:b) Reverse
(:b)->(:b) (:a)->(:b) Reverse (:a)->(:a)
(:b)->(:a {})
*/

func TestCollectAllPaths(t *testing.T) {
	graph, nodes := GetLineGraphWithSelfLoops(3, true)
	// x := JSON{}
	// buf := &bytes.Buffer{}
	// if err := x.Encode(graph, buf); err != nil {
	// 	t.Error(err)
	// }
	// nodes[0].SetProperty("key", "value")
	nodes[0].SetLabels(NewStringSet("node1"))
	nodes[1].SetLabels(NewStringSet("node2"))
	// nodes[1].SetProperty("key", "value")
	acc := &DefaultMatchAccumulator{}
	CollectAllPaths(graph, nodes[1], nodes[1].GetEdges(AnyEdge), func(e *Edge) bool { return true }, AnyEdge, 1, -1, func(e *Path) bool {
		// fmt.Println(e.String())
		// fmt.Println(e)
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
