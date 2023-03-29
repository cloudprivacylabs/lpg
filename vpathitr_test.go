package lpg

import (
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
	graph, nodes := GetLineGraphWithSelfLoops(4, true)
	// nodes[2].SetProperty("key", "value")
	nodes[2].SetLabels(NewStringSet("node2"))
	nodes[3].SetLabels(NewStringSet("node3"))
	nodes[0].SetLabels(NewStringSet("node0"))
	nodes[1].SetLabels(NewStringSet("node1"))
	// nodes[3].SetProperty("key", "value")
	acc := &DefaultMatchAccumulator{}
	CollectAllPaths(graph, nodes[1], nodes[1].GetEdges(AnyEdge), func(e *Edge) bool { return true }, AnyEdge, 1, -1, func(e *Path) bool {
		acc.Paths = append(acc.Paths, e)
		return true
	})
}
