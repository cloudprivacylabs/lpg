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
	"fmt"
	"strings"
)

// A Node represents a graph node.
type Node struct {
	labels StringSet
	Properties
	graph    *Graph
	incoming EdgeMap
	outgoing EdgeMap
	id       int
}

// GetGraph returns the graph owning the node
func (node *Node) GetGraph() *Graph { return node.graph }

// GetLabels returns a copy of the node labels
func (node *Node) GetLabels() StringSet { return node.labels.Clone() }

// HasLabel returns true if the node has the given label
func (node *Node) HasLabel(s string) bool { return node.labels.Has(s) }

// GetID returns the unique node ID. The ID is meaningless if the node
// is removed from the graph
func (node *Node) GetID() int { return node.id }

// Returns an edge iterator for incoming or outgoing edges
func (node *Node) GetEdges(dir EdgeDir) EdgeIterator {
	switch dir {
	case IncomingEdge:
		return node.incoming.Iterator()
	case OutgoingEdge:
		return node.outgoing.Iterator()
	}
	i1 := node.incoming.Iterator()
	i2 := node.outgoing.Iterator()
	return &edgeIterator{withSize(MultiIterator(i1, i2), i1.MaxSize()+i2.MaxSize())}
}

// Returns an edge iterator for incoming or outgoing edges with the given label
func (node *Node) GetEdgesWithLabel(dir EdgeDir, label string) EdgeIterator {
	switch dir {
	case IncomingEdge:
		return node.incoming.IteratorLabel(label)
	case OutgoingEdge:
		return node.outgoing.IteratorLabel(label)
	}
	i1 := node.incoming.IteratorLabel(label)
	i2 := node.outgoing.IteratorLabel(label)
	return &edgeIterator{withSize(MultiIterator(i1, i2), i1.MaxSize()+i2.MaxSize())}
}

// Returns an edge iterator for incoming or outgoingn edges that has the given labels
func (node *Node) GetEdgesWithAnyLabel(dir EdgeDir, labels StringSet) EdgeIterator {
	switch dir {
	case IncomingEdge:
		if labels.Len() == 0 {
			return node.incoming.Iterator()
		}
		return node.incoming.IteratorAnyLabel(labels)
	case OutgoingEdge:
		if labels.Len() == 0 {
			return node.outgoing.Iterator()
		}
		return node.outgoing.IteratorAnyLabel(labels)
	}
	i1 := node.GetEdgesWithAnyLabel(IncomingEdge, labels)
	i2 := node.GetEdgesWithAnyLabel(OutgoingEdge, labels)
	return &edgeIterator{withSize(MultiIterator(i1, i2), i1.MaxSize()+i2.MaxSize())}
}

// SetLabels sets the node labels
func (node *Node) SetLabels(labels StringSet) {
	node.graph.setNodeLabels(node, labels)
}

// SetProperty sets a node property
func (node *Node) SetProperty(key string, value interface{}) {
	node.graph.setNodeProperty(node, key, value)
}

// RemoveProperty removes a node property
func (node *Node) RemoveProperty(key string) {
	node.graph.removeNodeProperty(node, key)
}

// Remove all connected edges, and remove the node
func (node *Node) DetachAndRemove() {
	node.graph.detachRemoveNode(node)
}

// Remove all connected edges
func (node *Node) Detach() {
	node.graph.detachNode(node)
}

// String returns the string representation of the node
func (node *Node) String() string {
	labels := strings.Join(node.labels.Slice(), ":")
	if node.labels.Len() > 0 {
		labels = ":" + labels
	}
	return fmt.Sprintf("(%s %s)", labels, node.Properties)
}

// NextNodesWith returns the nodes reachable from source with the given label at one step
func NextNodesWith(source *Node, label string) []*Node {
	return TargetNodes(source.GetEdgesWithLabel(OutgoingEdge, label))
}

// PrevNodesWith returns the nodes reachable from source with the given label at one step
func PrevNodesWith(source *Node, label string) []*Node {
	return SourceNodes(source.GetEdgesWithLabel(IncomingEdge, label))
}
