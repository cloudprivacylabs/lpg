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

import "fmt"

// CollectAllPaths iterates the variable length paths that have the
// edges in firstLeg. For each edge, it calls the edgeFilter
// function. If the edge is accepted, it recursively descends and
// calls accumulator.AddPath for each discovered path until AddPath
// returns false
func CollectAllPaths(graph *Graph, fromNode *Node, firstLeg EdgeIterator, edgeFilter func(*Edge) bool, dir EdgeDir, min, max int, accumulator func(*Path) bool) {
	var recurse func(*Path) bool
	isLoop := func(path *Path, nextPath PathElement) bool {
		for _, p := range path.path {
			if nextPath.Edge == p.Edge {
				fmt.Println(path)
				return true
			}
		}
		// var lastOccurrenceIdx int
		// var loopCount int
		// for i := 0; i < path.NumNodes(); i++ {
		// 	node := path.GetNode(i)
		// 	if nextPath.GetTargetNode() == node {
		// 		lastOccurrenceIdx = i
		// 		loopCount++
		// 	}
		// }
		// if loopCount < 2 {
		// 	return false
		// }
		// shortPath := &Path{path: make([]PathElement, 0)}
		// for i := 0; i < path.NumNodes(); i++ {
		// 	node := path.GetNode(i)
		// 	// pontentially a loop
		// 	if nextPath.GetTargetNode() == node {
		// 		if i == lastOccurrenceIdx {
		// 			return false
		// 		}
		// 		copy(shortPath.path, path.path[lastOccurrenceIdx:])
		// 		return path.Slice(i, -1).HasPrefixPath(shortPath.Append(nextPath))
		// 	}
		// }
		return false
	}
	recurse = func(prefix *Path) bool {
		if (min == -1 || prefix.NumEdges() >= min) && (max == -1 || prefix.NumEdges() <= max) {
			if !accumulator(prefix.Clone()) {
				return false
			}
		}

		if max != -1 && prefix.NumEdges() >= max {
			return true
		}

		itr := edgeIterator{
			&filterIterator{
				itr: prefix.Last().GetEdges(dir),
				filter: func(item interface{}) bool {
					return edgeFilter(item.(*Edge))
				},
			},
		}

		for itr.Next() {
			edge := itr.Edge()
			pe := PathElement{Edge: edge}
			if edge.GetFrom() != edge.GetTo() {
				if edge.GetTo() == prefix.Last() {
					pe.Reverse = true
				}
			}
			if isLoop(prefix, pe) {
				return true
			}
			if !recurse(prefix.Clone().Append(pe)) {
				return false
			}
		}
		return true
	}

	edgeIr := makeUniqueIterator(firstLeg)
	for edgeIr.Next() {
		edge := firstLeg.Edge()
		pe := PathElement{Edge: edge}
		if edge.GetFrom() != edge.GetTo() {
			if edge.GetTo() == fromNode {
				pe.Reverse = true
			}
		}
		fmt.Println(edge)
		if !recurse(&Path{path: []PathElement{pe}}) {
			break
		}
	}
}
