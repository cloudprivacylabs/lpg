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
	"reflect"
)

// CollectAllPaths iterates the variable length paths that have the
// edges in firstLeg. For each edge, it calls the edgeFilter
// function. If the edge is accepted, it recursively descends and
// calls accumulator.AddPath for each discovered path until AddPath
// returns false

// 4 - 7 - 7 - 8 - 3 - 3 - 2 - 2 - 1
// 4 - 7 - 7 - 8 - 3 - 4 - 4 - 7 - 8 - 3 - 3 - 2 - 2 - 1
// 4 - 7 - 7 - 8 - 3 - 3 - 2 - 2 - 1 - 3 - 4 - 4 - 7 - 8
// 3 - 3 - 2 - 2 - 1
func CollectAllPaths(graph *Graph, fromNode *Node, firstLeg EdgeIterator, edgeFilter func(*Edge) bool, dir EdgeDir, min, max int, accumulator func([]*Edge, *Node) bool) {
	var recurse func([]*Edge, *Node) bool
	// return if all edges of p2 exist in p1
	prefixPath := func(p1, p2 []*Edge) bool {
		if len(p2) == 0 {
			return false
		}
	path1Start:
		for path1Idx, e1 := range p1 {
			for _, e2 := range p2 {
				if e2 == e1 {
					if len(p1[path1Idx:]) < len(p2) {
						return false
					}
					if !reflect.DeepEqual(p2, p1[path1Idx:path1Idx+len(p2)]) {
						continue path1Start
					}
					return true
				}
			}
		}
		return false
	}
	isLoop := func(nextEdge *Edge, path []*Edge) bool {
		cmpPath := make([]*Edge, 0)
		for _, step := range path {
			if nextEdge.GetTo() == step.GetFrom() {
				cmpPath = append(cmpPath, nextEdge)
			}
		}
		return prefixPath(path, cmpPath)
	}
	// isLoop := func(node *Node, edges []*Edge) bool {
	// 	for _, e := range edges {
	// 		if e.GetFrom() == node {
	// 			return true
	// 		}
	// 	}
	// 	if len(edges) > 0 {
	// 		return edges[len(edges)-1].GetTo() == node
	// 	}
	// 	return false
	// }

	recurse = func(prefix []*Edge, lastNode *Node) bool {
		var endNode *Node
		switch dir {
		case OutgoingEdge:
			endNode = prefix[len(prefix)-1].GetTo()
		case IncomingEdge:
			endNode = prefix[len(prefix)-1].GetFrom()
		case AnyEdge:
			if prefix[len(prefix)-1].GetTo() == lastNode {
				endNode = prefix[len(prefix)-1].GetFrom()
			} else {
				endNode = prefix[len(prefix)-1].GetTo()
			}
		}

		if (min == -1 || len(prefix) >= min) && (max == -1 || len(prefix) <= max) {
			// A valid path
			entry := make([]*Edge, len(prefix))
			copy(entry, prefix)
			if !accumulator(entry, endNode) {
				return false
			}
		}

		if max != -1 && len(prefix) >= max {
			return true
		}

		// if isLoop(endNode, prefix[:len(prefix)-1]) {
		// 	return true
		// }
		itr := edgeIterator{
			&filterIterator{
				itr: endNode.GetEdges(dir),
				filter: func(item interface{}) bool {
					return edgeFilter(item.(*Edge))
				},
			},
		}
		for itr.Next() {
			edge := itr.Edge()
			if isLoop(edge, prefix) {
				return false
			}
			// if isLoop(edge, prefix[:len(prefix)-1]) {
			// 	return false
			// }
			if !recurse(append(prefix, edge), endNode) {
				return false
			}
		}
		return true
	}

	for firstLeg.Next() {
		edge := firstLeg.Edge()
		if !recurse([]*Edge{edge}, fromNode) {
			break
		}
	}
}
