// Copyright 2021 Cloud Privacy Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lpg

import (
	"context"
)

// Sources finds all the source nodes in the graph
func SourcesItr(graph *Graph) NodeIterator {
	return nodeIterator{
		&filterIterator{
			itr: graph.GetNodes(),
			filter: func(item interface{}) bool {
				node := item.(*Node)
				if edges := node.GetEdges(IncomingEdge); edges.Next() {
					return false
				}
				return true
			},
		},
	}
}

// Sources finds all the source nodes in the graph
func Sources(graph *Graph) []*Node {
	return NodeSlice(SourcesItr(graph))
}

// Sinks finds all the sink nodes in the graph
func SinksItr(graph *Graph) NodeIterator {
	return nodeIterator{
		&filterIterator{
			itr: graph.GetNodes(),
			filter: func(item interface{}) bool {
				node := item.(*Node)
				if edges := node.GetEdges(OutgoingEdge); edges.Next() {
					return false
				}
				return true
			},
		},
	}
}

// Sinks finds all the sink nodes in the graph
func Sinks(graph *Graph) []*Node {
	return NodeSlice(SinksItr(graph))
}

// EdgesBetweenNodes finds all the edges that go from 'from'  to 'to'
func EdgesBetweenNodes(from, to *Node) []*Edge {
	ret := make([]*Edge, 0)
	for edges := from.GetEdges(OutgoingEdge); edges.Next(); {
		edge := edges.Edge()
		if edge.GetTo() == to {
			ret = append(ret, edge)
		}
	}
	return ret
}

// CheckIsomoprhism checks to see if graphs given are equal as defined
// by the edge equivalence and node equivalence functions. The
// nodeEquivalenceFunction will be called for all pairs of nodes. The
// edgeEquivalenceFunction will be called for edges connecting
// equivalent nodes.
//
// This is a potentially long running function. Cancel the context to
// stop. If the function returns because of context cancellation,
// error will be ctx.Err()
func CheckIsomorphism(ctx context.Context, g1, g2 *Graph, nodeEquivalenceFunc func(n1, n2 *Node) bool, edgeEquivalenceFunc func(e1, e2 *Edge) bool) (bool, error) {
	if g1.NumNodes() != g2.NumNodes() || g1.NumEdges() != g2.NumEdges() {
		return false, nil
	}

	// Slice of all nodes of g1
	all1Nodes := NodeSlice(g1.GetNodes())
	// Possible node equivalences. equivalences[i] is a slices of nodes of n2 that are possibly equivalent to all1Nodes[i]
	equivalences := make([][]*Node, len(all1Nodes))

	// Fill possible equivalences
	for i, node1 := range all1Nodes {
		for nodes := g2.GetNodes(); nodes.Next(); {
			node2 := nodes.Node()
			if nodeEquivalenceFunc(node1, node2) {
				equivalences[i] = append(equivalences[i], node2)
			}
		}
		if len(equivalences[i]) == 0 {
			return false, nil
		}
		if err := ctx.Err(); err != nil {
			return false, err
		}
	}

	nodeEquivalences := make([]int, len(all1Nodes))

	// build a node equivalence map based on the current state of nodeEquivalences. nodeEquivalences must be valid
	buildNodeEquivalence := func() map[*Node]*Node {
		eq := make(map[*Node]*Node)
		for i := range nodeEquivalences {
			node1 := all1Nodes[i]
			node2 := equivalences[i][nodeEquivalences[i]]
			eq[node1] = node2
		}
		return eq
	}

	// Increment node equivalences to the next node permutation
	next := func() bool {
		for index := range nodeEquivalences {
			nodeEquivalences[index]++
			if nodeEquivalences[index] < len(equivalences[index]) {
				return true
			}
			nodeEquivalences[index] = 0
		}
		return false
	}

	isIsomorphism := func(nodeMapping map[*Node]*Node) bool {
		for node1, node2 := range nodeMapping {
			// node1 and node2 are equivalent. Now we check if equivalent edges go to equivalent nodes
			edges1 := EdgeSlice(node1.GetEdges(OutgoingEdge))
			edges2 := EdgeSlice(node2.GetEdges(OutgoingEdge))
			// There must be same number of edges
			if len(edges1) != len(edges2) {
				return false
			}

			for _, edge1 := range edges1 {
				found := false
				for _, edge2 := range edges2 {
					// edge1.Edge.GetTo()
					if nodeMapping[edge1.GetTo()] == edge2.GetTo() &&
						edgeEquivalenceFunc(edge1, edge2) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
		return true
	}

	// Iterate possible node equivalences, and check isomorphism
	for {
		nodeMapping := buildNodeEquivalence()
		if isIsomorphism(nodeMapping) {
			return true, nil
		}
		if !next() {
			break
		}
		if err := ctx.Err(); err != nil {
			return false, err
		}
	}
	return false, nil
}

// ForEachNode iterates through all the nodes of g until predicate
// returns false or all nodes are processed.
func ForEachNode(g *Graph, predicate func(*Node) bool) bool {
	for nodes := g.GetNodes(); nodes.Next(); {
		if !predicate(nodes.Node()) {
			return false
		}
	}
	return true
}
