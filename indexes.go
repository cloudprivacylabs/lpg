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

type index interface {
	add(value interface{}, id int, item interface{})
	remove(value interface{}, id int)
	find(value interface{}) Iterator
	valueItr() Iterator
}

type IndexType int

const (
	BtreeIndex IndexType = 0
	HashIndex  IndexType = 1
)

type graphIndex struct {
	nodesByLabel NodeMap

	nodeProperties map[string]index
	edgeProperties map[string]index
}

func newGraphIndex() graphIndex {
	return graphIndex{
		nodesByLabel:   *NewNodeMap(),
		nodeProperties: make(map[string]index),
		edgeProperties: make(map[string]index),
	}
}

// NodePropertyIndex sets up an index for the given node property
func (g *graphIndex) NodePropertyIndex(propertyName string, graph *Graph, it IndexType) {
	_, exists := g.nodeProperties[propertyName]
	if exists {
		return
	}
	var ix index
	if it == BtreeIndex {
		ix = &setTree{}
	} else {
		ix = &hashIndex{}
	}
	g.nodeProperties[propertyName] = ix
	// Reindex
	for nodes := graph.GetNodes(); nodes.Next(); {
		node := nodes.Node()
		value, ok := node.properties[propertyName]
		if ok {
			ix.add(value, node.id, node)
		}
	}
}

func (g *graphIndex) isNodePropertyIndexed(propertyName string) index {
	return g.nodeProperties[propertyName]
}

func (g *graphIndex) isEdgePropertyIndexed(propertyName string) index {
	return g.edgeProperties[propertyName]
}

// GetIteratorForNodeProperty returns an iterator for the given
// key/value, and the max size of the resultset. If no index found,
// returns nil,-1
func (g *graphIndex) GetIteratorForNodeProperty(key string, value interface{}) NodeIterator {
	index, found := g.nodeProperties[key]
	if !found {
		return nil
	}
	itr := index.find(value)
	return nodeIterator{itr}
}

// NodesWithProperty returns an iterator that will go through the
// nodes that has the property
func (g *graphIndex) NodesWithProperty(key string) NodeIterator {
	index, found := g.nodeProperties[key]
	if !found {
		return nil
	}
	return nodeIterator{index.valueItr()}
}

// EdgesWithProperty returns an iterator that will go through the
// edges that has the property
func (g *graphIndex) EdgesWithProperty(key string) EdgeIterator {
	index, found := g.edgeProperties[key]
	if !found {
		return nil
	}
	return edgeIterator{index.valueItr()}
}

func (g *graphIndex) addNodeToIndex(node *Node, graph *Graph) {
	g.nodesByLabel.Add(node)

	for k, v := range node.properties {
		index, found := g.nodeProperties[k]
		if !found {
			continue
		}
		index.add(v, node.id, node)
	}
}

func (g *graphIndex) removeNodeFromIndex(node *Node, graph *Graph) {
	g.nodesByLabel.Remove(node)

	for k, v := range node.properties {
		index, found := g.nodeProperties[k]
		if !found {
			continue
		}
		index.remove(v, node.id)
	}
}

// EdgePropertyIndex sets up an index for the given edge property
func (g *graphIndex) EdgePropertyIndex(propertyName string, graph *Graph, it IndexType) {
	_, exists := g.edgeProperties[propertyName]
	if exists {
		return
	}
	var ix index
	if it == BtreeIndex {
		ix = &setTree{}
	} else {
		ix = &hashIndex{}
	}
	g.edgeProperties[propertyName] = ix
	// Reindex
	for edges := graph.GetEdges(); edges.Next(); {
		edge := edges.Edge()
		value, ok := edge.properties[propertyName]
		if ok {
			ix.add(value, edge.id, edge)
		}
	}
}

func (g *graphIndex) addEdgeToIndex(edge *Edge, graph *Graph) {
	for k, v := range edge.properties {
		index, found := g.edgeProperties[k]
		if !found {
			continue
		}
		index.add(v, edge.id, edge)
	}
}

func (g *graphIndex) removeEdgeFromIndex(edge *Edge, graph *Graph) {
	for k, v := range edge.properties {
		index, found := g.edgeProperties[k]
		if !found {
			continue
		}
		index.remove(v, edge.id)
	}
}

// GetIteratorForEdgeProperty returns an iterator for the given
// key/value, and the max size of the resultset. If no index found,
// returns nil,-1
func (g *graphIndex) GetIteratorForEdgeProperty(key string, value interface{}) EdgeIterator {
	index, found := g.edgeProperties[key]
	if !found {
		return nil
	}
	itr := index.find(value)
	return edgeIterator{itr}
}
