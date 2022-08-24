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

// A Graph is a labeled property graph containing nodes, and directed
// edges combining those nodes.
//
// Every node of the graph contains a set of labels, and a map of
// properties. Each node provides access to the edges adjacent to
// it. Every edge contains a from and to node, a label, and a map of
// properties.
//
// Nodes and edges are "owned" by a graph object. Given a node or
// edge, you can always get its graph using `node.GetGraph()` or
// `edge.GetGraph()` method. You cannot use an edge to connect nodes
// of different graphs.
//
// Zero value for a Graph is not usable. Use `NewGraph` to construct a
// new graph.
type Graph struct {
	index    graphIndex
	allNodes NodeSet
	allEdges EdgeMap
	idBase   int
}

// NewGraph constructs and returns a new graph. The new graph has no
// nodes or edges.
func NewGraph() *Graph {
	return &Graph{
		index: newGraphIndex(),
	}
}

// NewNode creates a new node with the given labels and properties
func (g *Graph) NewNode(labels []string, properties map[string]interface{}) *Node {
	node := &Node{labels: NewStringSet(labels...), Properties: Properties(properties), graph: g}
	node.id = g.idBase
	g.idBase++
	g.addNode(node)
	return node
}

// NewEdge creates a new edge between the two nodes of the graph. Both
// nodes must be nodes of this graph, otherwise this call panics
func (g *Graph) NewEdge(from, to *Node, label string, properties map[string]interface{}) *Edge {
	if from.graph != g {
		panic("from node is not in graph")
	}
	if to.graph != g {
		panic("to node is not in graph")
	}
	newEdge := &Edge{
		from:       from,
		to:         to,
		label:      label,
		Properties: Properties(properties),
		id:         g.idBase,
	}
	g.idBase++
	g.allEdges.Add(newEdge)
	g.connect(newEdge)
	g.index.addEdgeToIndex(newEdge)
	return newEdge
}

// NumNodes returns the number of nodes in the graph
func (g *Graph) NumNodes() int {
	return g.allNodes.Len()
}

// NumEdges returns the number of edges in the graph
func (g *Graph) NumEdges() int {
	return g.allEdges.Len()
}

// GetNodes returns a node iterator that goes through all the nodes of
// the graph. The behavior of the returned iterator is undefined if
// during iteration nodes are updated, new nodes are added, or
// existing nodes are deleted.
func (g *Graph) GetNodes() NodeIterator {
	return g.allNodes.Iterator()
}

// GetNodesWithAllLabels returns an iterator that goes through the
// nodes that have all the nodes given in the set. The nodes may have
// more nodes than given in the set. The behavior of the returned
// iterator is undefined if during iteration nodes are updated, new
// nodes are added, or existing nodes are deleted.
func (g *Graph) GetNodesWithAllLabels(labels StringSet) NodeIterator {
	return g.index.nodesByLabel.IteratorAllLabels(labels)
}

// GetEdges returns an edge iterator that goes through all the edges
// of the graph. The behavior of the returned iterator is undefined if
// during iteration edges are updated, new edges are added, or
// existing edges are deleted.
func (g *Graph) GetEdges() EdgeIterator {
	return g.allEdges.Iterator()
}

// GetEdgesWithAnyLabel returns an iterator that goes through the
// edges of the graph that are labeled with one of the labels in the
// given set. The behavior of the returned iterator is undefined if
// during iteration edges are updated, new edges are added, or
// existing edges are deleted.
func (g *Graph) GetEdgesWithAnyLabel(set StringSet) EdgeIterator {
	return g.allEdges.IteratorAnyLabel(set)
}

// AddEdgePropertyIndex adds an index for the given edge property
func (g *Graph) AddEdgePropertyIndex(propertyName string) {
	g.index.EdgePropertyIndex(propertyName, g)
}

// AddNodePropertyIndex adds an index for the given node property
func (g *Graph) AddNodePropertyIndex(propertyName string) {
	g.index.NodePropertyIndex(propertyName, g)
}

// GetNodesWithProperty returns an iterator for the nodes that has the
// property. If there is an index for the node property, and iterator
// over that index is returned. Otherwise, an iterator that goes
// through all nodes while filtering them is returned. The behavior of
// the returned iterator is undefined if during iteration nodes are
// updated, new nodes are added, or existing nodes are deleted.
func (g *Graph) GetNodesWithProperty(property string) NodeIterator {
	itr := g.index.NodesWithProperty(property)
	if itr != nil {
		return itr
	}
	return &nodeIterator{&filterIterator{
		itr: g.GetNodes(),
		filter: func(v interface{}) bool {
			wp, ok := v.(*Node)
			if !ok {
				return false
			}
			_, exists := wp.GetProperty(property)
			return exists
		},
	}}
}

// GetEdgesWithProperty returns an iterator for the edges that has the
// property. If there is an index for the property, and iterator over
// that index is returned. Otherwise, an iterator that goes through
// all edges while filtering them is returned. The behavior of the
// returned iterator is undefined if during iteration edges are
// updated, new edges are added, or existing edges are deleted.
func (g *Graph) GetEdgesWithProperty(property string) EdgeIterator {
	itr := g.index.EdgesWithProperty(property)
	if itr != nil {
		return itr
	}
	return &edgeIterator{&filterIterator{
		itr: g.GetEdges(),
		filter: func(v interface{}) bool {
			wp, ok := v.(*Edge)
			if !ok {
				return false
			}
			_, exists := wp.GetProperty(property)
			return exists
		},
	}}

}

// FindNodes returns an iterator that will iterate through all the
// nodes that have all of the given labels and properties. If
// allLabels is nil or empty, it does not look at the labels. If
// properties is nil or empty, it does not look at the properties
func (g *Graph) FindNodes(allLabels StringSet, properties map[string]interface{}) NodeIterator {
	if allLabels.Len() == 0 && len(properties) == 0 {
		// Return all nodes
		return g.GetNodes()
	}

	var nodesByLabelItr NodeIterator
	if allLabels.Len() > 0 {
		nodesByLabelItr = g.index.nodesByLabel.IteratorAllLabels(allLabels)
	}
	// Select the iterator with minimum max size
	nodesByLabelSize := nodesByLabelItr.MaxSize()
	propertyIterators := make(map[string]NodeIterator)
	if len(properties) > 0 {
		for k, v := range properties {
			itr := g.index.GetIteratorForNodeProperty(k, v)
			if itr == nil {
				continue
			}
			propertyIterators[k] = itr
		}
	}
	var minimumPropertyItrKey string
	minPropertySize := -1
	for k, itr := range propertyIterators {
		maxSize := itr.MaxSize()
		if maxSize == -1 {
			continue
		}
		if minPropertySize == -1 || minPropertySize > maxSize {
			minPropertySize = maxSize
			minimumPropertyItrKey = k
		}
	}

	nodeFilterFunc := GetNodeFilterFunc(allLabels, properties)
	// Iterate the minimum iterator, with a filter
	if nodesByLabelSize != -1 && (minPropertySize == -1 || minPropertySize > nodesByLabelSize) {
		// Iterate by node label
		// build a filter from properties
		return &nodeIterator{
			&filterIterator{
				itr: nodesByLabelItr,
				filter: func(item interface{}) bool {
					return nodeFilterFunc(item.(*Node))
				},
			},
		}
	}
	if minPropertySize != -1 {
		// Iterate by property
		return &nodeIterator{
			&filterIterator{
				itr: propertyIterators[minimumPropertyItrKey],
				filter: func(item interface{}) bool {
					return nodeFilterFunc(item.(*Node))
				},
			},
		}
	}
	// Iterate all
	return g.GetNodes()
}

// FindEdges returns an iterator that will iterate through all the
// edges whose label is in the given labels and have all the
// properties. If labels is nil or empty, it does not look at the
// labels. If properties is nil or empty, it does not look at the
// properties
func (g *Graph) FindEdges(labels StringSet, properties map[string]interface{}) EdgeIterator {
	if labels.Len() == 0 && len(properties) == 0 {
		// Return all edges
		return g.GetEdges()
	}

	var edgesByLabelItr EdgeIterator
	if labels.Len() > 0 {
		edgesByLabelItr = g.GetEdgesWithAnyLabel(labels)
	}
	// Select the iterator with minimum max size
	edgesByLabelSize := edgesByLabelItr.MaxSize()
	propertyIterators := make(map[string]EdgeIterator)
	if len(properties) > 0 {
		for k, v := range properties {
			itr := g.index.GetIteratorForEdgeProperty(k, v)
			if itr == nil {
				continue
			}
			propertyIterators[k] = itr
		}
	}
	var minimumPropertyItrKey string
	minPropertySize := -1
	for k, itr := range propertyIterators {
		maxSize := itr.MaxSize()
		if maxSize == -1 {
			continue
		}
		if minPropertySize == -1 || minPropertySize > maxSize {
			minPropertySize = maxSize
			minimumPropertyItrKey = k
		}
	}

	edgeFilterFunc := GetEdgeFilterFunc(labels, properties)
	// Iterate the minimum iterator, with a filter
	if edgesByLabelSize != -1 && (minPropertySize == -1 || minPropertySize > edgesByLabelSize) {
		// Iterate by edge label
		// build a filter from properties
		return &edgeIterator{
			&filterIterator{
				itr: edgesByLabelItr,
				filter: func(item interface{}) bool {
					return edgeFilterFunc(item.(*Edge))
				},
			},
		}
	}
	if minPropertySize != -1 {
		// Iterate by property
		return &edgeIterator{
			&filterIterator{
				itr: propertyIterators[minimumPropertyItrKey],
				filter: func(item interface{}) bool {
					return edgeFilterFunc(item.(*Edge))
				},
			},
		}
	}
	// Iterate all
	return g.GetEdges()
}

// GetNodeFilterFunc returns a filter function that can be used to select
// nodes that have all the specified labels, with correct property
// values
func GetNodeFilterFunc(labels StringSet, properties map[string]interface{}) func(*Node) bool {
	return func(node *Node) (cmp bool) {
		if labels.Len() > 0 {
			if !node.labels.HasAllSet(labels) {
				return false
			}
		}
		defer func() {
			if r := recover(); r != nil {
				cmp = false
			}
		}()
		for k, v := range properties {
			nodeValue, exists := node.GetProperty(k)
			if !exists {
				if v != nil {
					return false
				}
			}
			if ComparePropertyValue(v, nodeValue) != 0 {
				return false
			}
		}
		return true
	}
}

// GetEdgeFilterFunc returns a function that can be used to select edges
// that have at least one of the specified labels, with correct
// property values
func GetEdgeFilterFunc(labels StringSet, properties map[string]interface{}) func(*Edge) bool {
	return func(edge *Edge) (cmp bool) {
		if labels.Len() > 0 {
			if !labels.Has(edge.label) {
				return false
			}
		}
		defer func() {
			if r := recover(); r != nil {
				cmp = false
			}
		}()
		for k, v := range properties {
			edgeValue, exists := edge.GetProperty(k)
			if !exists {
				if v != nil {
					return false
				}
			}
			if ComparePropertyValue(v, edgeValue) != 0 {
				return false
			}
		}
		return true
	}
}

func (g *Graph) setNodeLabels(node *Node, labels StringSet) {
	g.index.nodesByLabel.Replace(node, node.GetLabels(), labels)
	node.labels = labels.Clone()
}

func (g *Graph) setNodeProperty(node *Node, key string, value interface{}) {
	if node.Properties == nil {
		node.Properties = make(Properties)
	}
	oldValue, exists := node.Properties[key]
	nix := g.index.isNodePropertyIndexed(key)
	if nix != nil && exists {
		nix.remove(oldValue, node.id, node)
	}
	node.Properties[key] = value
	if nix != nil {
		nix.add(value, node.id, node)
	}
}

func (g *Graph) cloneNode(node *Node, cloneProperty func(string, interface{}) interface{}) *Node {
	newNode := &Node{
		labels:     node.labels.Clone(),
		Properties: node.Properties.clone(cloneProperty),
		graph:      g,
	}
	newNode.id = g.idBase
	g.idBase++
	g.addNode(newNode)
	return newNode
}

func (g *Graph) addNode(node *Node) {
	g.allNodes.Add(node)
	g.index.addNodeToIndex(node)
}

func (g *Graph) removeNodeProperty(node *Node, key string) {
	if node.Properties == nil {
		return
	}
	value, exists := node.Properties[key]
	if !exists {
		return
	}
	nix := g.index.isNodePropertyIndexed(key)
	if nix != nil {
		nix.remove(value, node.id, node)
	}
	delete(node.Properties, key)
}

func (g *Graph) detachRemoveNode(node *Node) {
	g.detachNode(node)
	g.allNodes.Remove(node)
	g.index.removeNodeFromIndex(node)
}

func (g *Graph) detachNode(node *Node) {
	for _, edge := range EdgeSlice(node.incoming.Iterator()) {
		g.disconnect(edge)
		g.allEdges.Remove(edge)
		g.index.removeEdgeFromIndex(edge)
	}
	node.incoming.Clear()
	for _, edge := range EdgeSlice(node.outgoing.Iterator()) {
		g.disconnect(edge)
		g.allEdges.Remove(edge)
		g.index.removeEdgeFromIndex(edge)
	}
	node.outgoing.Clear()
}

func (g *Graph) cloneEdge(from, to *Node, edge *Edge, cloneProperty func(string, interface{}) interface{}) *Edge {
	if from.graph != g {
		panic("from node is not in graph")
	}
	if to.graph != g {
		panic("to node is not in graph")
	}
	newEdge := &Edge{
		from:       from,
		to:         to,
		label:      edge.label,
		Properties: edge.Properties.clone(cloneProperty),
		id:         g.idBase,
	}
	g.idBase++
	g.allEdges.Add(newEdge)
	g.connect(newEdge)
	g.index.addEdgeToIndex(newEdge)
	return newEdge
}

func (g *Graph) connect(edge *Edge) {
	edge.to.incoming.Add(edge)
	edge.from.outgoing.Add(edge)
}

func (g *Graph) disconnect(edge *Edge) {
	edge.to.incoming.Remove(edge)
	edge.from.outgoing.Remove(edge)
}

func (g *Graph) setEdgeLabel(edge *Edge, label string) {
	g.disconnect(edge)
	edge.label = label
	g.connect(edge)
}

func (g *Graph) removeEdge(edge *Edge) {
	g.disconnect(edge)
	g.allEdges.Remove(edge)
}

func (g *Graph) setEdgeProperty(edge *Edge, key string, value interface{}) {
	if edge.Properties == nil {
		edge.Properties = make(Properties)
	}
	oldValue, exists := edge.Properties[key]
	nix := g.index.isEdgePropertyIndexed(key)
	if nix != nil && exists {
		nix.remove(oldValue, edge.id, edge)
	}
	edge.Properties[key] = value
	if nix != nil {
		nix.add(value, edge.id, edge)
	}
}

func (g *Graph) removeEdgeProperty(edge *Edge, key string) {
	if edge.Properties == nil {
		return
	}
	oldValue, exists := edge.Properties[key]
	if !exists {
		return
	}
	nix := g.index.isEdgePropertyIndexed(key)
	if nix != nil {
		nix.remove(oldValue, edge.id, edge)
	}
	delete(edge.Properties, key)
}

type WithProperties interface {
	GetProperty(key string) (interface{}, bool)
}

func buildPropertyFilterFunc(key string, value interface{}) func(WithProperties) bool {
	return func(properties WithProperties) (cmp bool) {
		pvalue, exists := properties.GetProperty(key)
		if !exists {
			return value == nil
		}

		defer func() {
			if r := recover(); r != nil {
				cmp = false
			}
		}()
		return ComparePropertyValue(value, pvalue) == 0
	}
}
