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

// Cursor is a convenience class to move around a graph
type Cursor struct {
	node *Node
	path *Path
}

// Sets the node cursor is pointing at.
func (c *Cursor) Set(node *Node) *Cursor {
	c.node = node
	return c
}

// StartPath starts a new path at the current node. Panics if the node is not valid.
func (c *Cursor) StartPath() *Cursor {
	c.path = PathFromNode(c.node)
	return c
}

// PushToPath pushes the edge onto the path. The path must be started,
// and the pushed edge must be connected to the current node. This
// also advances the cursor to the target node.
func (c *Cursor) PushToPath(edge *Edge) *Cursor {
	c.path.Append(PathElement{
		Edge: edge,
	})
	c.node = edge.GetTo()
	return c
}

// GetPath returns the recorded path. This is the internal copy of the
// path, so caller must clone if it changes are necessary
func (c *Cursor) GetPath() *Path {
	return c.path
}

// PopFromPath removes the last edge from the path. This also moves the cursor to the previous node
func (c *Cursor) PopFromPath() *Cursor {
	c.path.RemoveLast()
	return c
}

// Edges returns an iterator of edges of the current node. If the
// current node is invalid, returns an empty iterator
func (c *Cursor) Edges(dir EdgeDir) EdgeIterator {
	if c.node == nil {
		return &edgeIterator{emptyIterator{}}
	}
	return c.node.GetEdges(dir)
}

// EdgesWith returns an iterator of edges of the current node with the
// given label. If the current node is invalid, returns an empty
// iterator
func (c *Cursor) EdgesWith(dir EdgeDir, label string) EdgeIterator {
	if c.node == nil {
		return &edgeIterator{emptyIterator{}}
	}
	return c.node.GetEdgesWithLabel(dir, label)
}

// Nodes returns an iterator over the nodes that are reachable by a
// single edge from the current node. If the current node is not
// valid, returns an empty iterator
func (c *Cursor) Nodes(dir EdgeDir) NodeIterator {
	if c.node == nil {
		return &nodeIterator{emptyIterator{}}
	}
	return newNodesFromEdges(c.node.GetEdges(dir), dir)
}

// NodesWith returns an iterator over the nodes that are reachable by
// a single edge with the given label from the current node. If the
// current node is not valid, returns an empty iterator
func (c *Cursor) NodesWith(dir EdgeDir, label string) NodeIterator {
	if c.node == nil {
		return &nodeIterator{emptyIterator{}}
	}
	return newNodesFromEdges(c.node.GetEdgesWithLabel(dir, label), dir)
}

// Forward returns an edge iterator for outgoing edges of the current
// node. If the current node is invalid, returns an empty iterator
func (c *Cursor) Forward() EdgeIterator { return c.Edges(OutgoingEdge) }

// Backward returns an edge iterator for incoming edges of the current
// node. If the current node is invalid, returns an empty iterator
func (c *Cursor) Backward() EdgeIterator { return c.Edges(IncomingEdge) }

// ForwardWith returns an edge iterator for outgoing edges of the current
// node with the given label. If the current node is invalid, returns an empty iterator
func (c *Cursor) ForwardWith(label string) EdgeIterator { return c.EdgesWith(OutgoingEdge, label) }

// BackwardWith returns an edge iterator for incoming edges of the
// current node with the given label. If the current node is invalid,
// returns an empty iterator
func (c *Cursor) BackwardWith(label string) EdgeIterator { return c.EdgesWith(IncomingEdge, label) }

// NextNodes returns a node iterator that can be reached at one step
// from the current node. If current node is invalid, returns an empty iterator.
func (c *Cursor) NextNodes() NodeIterator { return c.Nodes(OutgoingEdge) }

// NextNodesWith returns a node iterator that can be reached at one step
// from the current node with the given label. If current node is invalid, returns an empty iterator.
func (c *Cursor) NextNodesWith(label string) NodeIterator { return c.NodesWith(OutgoingEdge, label) }

// PrevNodes returns a node iterator that can be reached at one step
// backwards from the current node. If current node is invalid, returns an empty iterator.
func (c *Cursor) PrevNodes() NodeIterator { return c.Nodes(IncomingEdge) }

// PrevNodesWith returns a node iterator that can be reached at one step
// backwards from the current node with the given label. If current node is invalid, returns an empty iterator.
func (c *Cursor) PrevNodesWith(label string) NodeIterator { return c.NodesWith(IncomingEdge, label) }
