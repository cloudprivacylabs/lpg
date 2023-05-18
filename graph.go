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
	allNodes nodeList
	allEdges edgeMap
	idBase   int
}

// NewGraph constructs and returns a new graph. The new graph has no
// nodes or edges.
func NewGraph() *Graph {
	return &Graph{
		index:    newGraphIndex(),
		allNodes: nodeList{},
	}
}

// NewNode creates a new node with the given labels and properties
func (g *Graph) NewNode(labels []string, props map[string]interface{}) *Node {
	var p properties
	if len(props) > 0 {
		p = make(properties, len(props))
		for k, v := range props {
			p[k] = v
		}
	}
	return g.FastNewNode(NewStringSet(labels...), p)
}

// FastNewNode creates a new node with the given labels and
// properties. This version does not copy the labels and properties,
// but uses the given label set and map directly
func (g *Graph) FastNewNode(labels StringSet, props map[string]interface{}) *Node {
	node := &Node{
		labels:     labels,
		graph:      g,
		properties: properties(props),
	}
	node.id = g.idBase
	g.idBase++
	g.addNode(node)
	return node
}

// NewEdge creates a new edge between the two nodes of the graph. Both
// nodes must be nodes of this graph, otherwise this call panics
func (g *Graph) NewEdge(from, to *Node, label string, props map[string]any) *Edge {
	var p properties
	if len(props) > 0 {
		p = make(properties, len(props))
		for k, v := range props {
			p[k] = v
		}
	}
	return g.FastNewEdge(from, to, label, p)
}

// FastNewEdge creates a new edge between the two nodes of the
// graph. Both nodes must be nodes of this graph, otherwise this call
// panics. This version uses the given properties map directly without
// copying it.
func (g *Graph) FastNewEdge(from, to *Node, label string, props map[string]any) *Edge {
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
		id:         g.idBase,
		properties: properties(props),
	}
	g.idBase++
	g.allEdges.add(newEdge, 0)
	g.connect(newEdge)
	g.index.addEdgeToIndex(newEdge, g)
	return newEdge
}

// NumNodes returns the number of nodes in the graph
func (g *Graph) NumNodes() int {
	return g.allNodes.n
}

// NumEdges returns the number of edges in the graph
func (g *Graph) NumEdges() int {
	return g.allEdges.size()
}

// GetNodes returns a node iterator that goes through all the nodes of
// the graph. The behavior of the returned iterator is undefined if
// during iteration nodes are updated, new nodes are added, or
// existing nodes are deleted.
func (g *Graph) GetNodes() NodeIterator {
	return &nodeListIterator{next: g.allNodes.head, n: g.allNodes.n}
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
	return g.allEdges.iterator(0)
}

// GetEdgesWithAnyLabel returns an iterator that goes through the
// edges of the graph that are labeled with one of the labels in the
// given set. The behavior of the returned iterator is undefined if
// during iteration edges are updated, new edges are added, or
// existing edges are deleted.
func (g *Graph) GetEdgesWithAnyLabel(set StringSet) EdgeIterator {
	return g.allEdges.iteratorAnyLabel(set, 0)
}

// AddEdgePropertyIndex adds an index for the given edge property
func (g *Graph) AddEdgePropertyIndex(propertyName string, ix IndexType) {
	g.index.EdgePropertyIndex(propertyName, g, ix)
}

// AddNodePropertyIndex adds an index for the given node property
func (g *Graph) AddNodePropertyIndex(propertyName string, ix IndexType) {
	g.index.NodePropertyIndex(propertyName, g, ix)
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
	return nodeIterator{&filterIterator{
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
		return nodeIterator{
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
		return nodeIterator{
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
	nix := g.index.isNodePropertyIndexed(key)
	if node.properties == nil {
		node.properties = make(properties)
	} else {
		oldValue, exists := node.properties[key]
		if exists {
			if nix != nil {
				nix.remove(oldValue, node.id)
			}
		}
	}
	node.properties[key] = value
	if nix != nil {
		nix.add(value, node.id, node)
	}
}

func (g *Graph) cloneNode(sourceGraph *Graph, sourceNode *Node, cloneProperty func(string, interface{}) interface{}) *Node {
	newNode := &Node{
		labels: sourceNode.labels.Clone(),
		graph:  g,
	}
	if sourceNode.properties != nil {
		newNode.properties = sourceNode.properties.clone(sourceGraph, g, cloneProperty)
	}
	newNode.id = g.idBase
	g.idBase++
	g.addNode(newNode)
	return newNode
}

func (g *Graph) addNode(node *Node) {
	g.allNodes.add(node)
	g.index.addNodeToIndex(node, g)
}

func (g *Graph) removeNodeProperty(node *Node, key string) {
	if node.properties == nil {
		return
	}
	value, exists := node.properties[key]
	if !exists {
		return
	}
	nix := g.index.isNodePropertyIndexed(key)
	if nix != nil {
		nix.remove(value, node.id)
	}
	delete(node.properties, key)
}

func (g *Graph) detachRemoveNode(node *Node) {
	g.detachNode(node)
	g.allNodes.remove(node)
	g.index.removeNodeFromIndex(node, g)
}

func (g *Graph) detachNode(node *Node) {
	for _, edge := range EdgeSlice(node.incoming.iterator(2)) {
		g.disconnect(edge)
		g.allEdges.remove(edge, 0)
		g.index.removeEdgeFromIndex(edge, g)
	}
	node.incoming = edgeMap{}
	for _, edge := range EdgeSlice(node.outgoing.iterator(1)) {
		g.disconnect(edge)
		g.allEdges.remove(edge, 0)
		g.index.removeEdgeFromIndex(edge, g)
	}
	node.outgoing = edgeMap{}
}

func (g *Graph) cloneEdge(from, to *Node, sourceEdge *Edge, cloneProperty func(string, interface{}) interface{}) *Edge {
	if from.graph != g {
		panic("from node is not in graph")
	}
	if to.graph != g {
		panic("to node is not in graph")
	}
	newEdge := &Edge{
		from:  from,
		to:    to,
		label: sourceEdge.label,
		id:    g.idBase,
	}
	if sourceEdge.properties != nil {
		newEdge.properties = sourceEdge.properties.clone(to.graph, g, cloneProperty)
	}
	g.idBase++
	g.allEdges.add(newEdge, 0)
	g.connect(newEdge)
	g.index.addEdgeToIndex(newEdge, g)
	return newEdge
}

func (g *Graph) connect(edge *Edge) {
	edge.to.incoming.add(edge, 2)
	edge.from.outgoing.add(edge, 1)
}

func (g *Graph) disconnect(edge *Edge) {
	edge.to.incoming.remove(edge, 2)
	edge.from.outgoing.remove(edge, 1)
}

func (g *Graph) setEdgeLabel(edge *Edge, label string) {
	g.disconnect(edge)
	edge.label = label
	g.connect(edge)
}

func (g *Graph) removeEdge(edge *Edge) {
	g.disconnect(edge)
	g.allEdges.remove(edge, 0)
}

func (g *Graph) setEdgeProperty(edge *Edge, key string, value interface{}) {
	nix := g.index.isEdgePropertyIndexed(key)
	if edge.properties == nil {
		edge.properties = make(properties)
	} else {
		oldValue, exists := edge.properties[key]
		if exists {
			if nix != nil {
				nix.remove(oldValue, edge.id)
			}
		}
	}
	edge.properties[key] = value
	if nix != nil {
		nix.add(value, edge.id, edge)
	}
}

func (g *Graph) removeEdgeProperty(edge *Edge, key string) {
	if edge.properties == nil {
		return
	}
	oldValue, exists := edge.properties[key]
	if !exists {
		return
	}
	nix := g.index.isEdgePropertyIndexed(key)
	if nix != nil {
		nix.remove(oldValue, edge.id)
	}
	delete(edge.properties, key)
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
