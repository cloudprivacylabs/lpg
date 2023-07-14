[![GoDoc](https://godoc.org/github.com/cloudprivacylabs/lpg?status.svg)](https://godoc.org/github.com/cloudprivacylabs/lpg/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudprivacylabs/lpg)](https://goreportcard.com/report/github.com/cloudprivacylabs/lpg/v2)
[![Build Status](https://github.com/cloudprivacylabs/lpg/actions/workflows/CI.yml/badge.svg?branch=main)](https://github.com/cloudprivacylabs/lpg/actions/workflows/CI.yml)
# Labeled property graphs

This labeled property graph package implements the openCypher model of
labeled property graphs. A labeled property graph (LPG) contains nodes
and directed edges between those nodes. Every node contains:

  * Labels: Set of string tokens that usually identify the type of the
    node,
  * Properties: Key-value pairs.
  
Every edge contains:
  * A label: String token that identifies a relationship, and
  * Properties: Key-value pairs.

A `Graph` objects keeps an index of the nodes and edges included in
it. Create a graph using `NewGraph` function:

```
g := lpg.NewGraph()
// Create two nodes
n1 := g.NewNode([]string{"label1"},map[string]interface{}{"prop": "value1" })
n2 := g.NewNode([]string{"label2"},map[string]interface{}{"prop": "value2" })
// Connect the two nodes with an edge
edge:=g.NewEdge(n1,n2,"relatedTo",nil)
```

The LPG library uses iterators to address nodes and edges.

``` go
for nodes:=graph.GetNodes(); nodes.Next(); {
  node:=nodes.Node()
}
for edges:=graph.GetEdges(); edges.Next(); {
  edge:edges.Edge()
}
```

Every node knows its adjacent edges. 

```go
// Get outgoing edges
for edges:=node1.GetEdges(lpg.OutgoingEdge); edges.Next(); {
  edge:=edges.Edge
}

// Get all edges
for edges:=node1.GetEdges(lpg.AnyEdge); edges.Next(); {
  edge:=edges.Edge
}
```

The graph indexes nodes by label, so access to nodes using labels is
fast. You can add additional indexes on properties:

```
g := lpg.NewGraph()
// Index all nodes with property 'prop'
g.AddNodePropertyIndex("prop")

// This access should be fast
nodes := g.GetNodesWithProperty("prop")

// This will go through all nodes
slowNodes:= g.GetNodesWithProperty("propWithoutIndex")
```

## Pattern Searches

Graph library supports searching patterns within a graph. The
following example searches for the pattern that match

```
(:label1) -[]->({prop:value})`
```

and returns the head nodes for every matching path:

```
pattern := lpg.Pattern{ 
 // Node containing label 'label1'
 {
   Labels: lpg.NewStringSet("label1"),
 },
 // Edge of length 1
 {
   Min: 1, 
   Max: 1,
 },
 // Node with property prop=value
 {
   Properties: map[string]interface{} {"prop":"value"},
 }}
nodes, err:=pattern.FindNodes(g,nil)
```

Variable length paths are supported:

``` go
pattern := lpg.Pattern{ 
 // Node containing label 'label1'
 {
   Labels: lpg.NewStringSet("label1"),
 },
 // Minimum paths of length 2, no maximum length
 {
   Min: 2, 
   Max: -1,
 },
 // Node with property prop=value
 {
   Properties: map[string]interface{} {"prop":"value"},
 }}

```

## JSON Encoding

This graph library uses the following JSON representation:

```
{
  "nodes": [
     {
       "n": 0,
       "labels": [ "l1", "l2",... ],
       "properties": {
          "key1": value,
          "key2": value,
          ...
        },
        "edges": [
           {
             "to": "1",
             "label": "edgeLabel",
             "properties": {
               "key1": value,
               "key2": value,
               ...
             }
           },
           ...
        ]
     },
      ...
  ],
  "edges": [
     {
        "from": 0,
        "to": 1,
        "label": "edgeLabel",
        "properties": {
           "key1": value1,
           "key2": value2,
           ...
        }
     },
     ...
  ]
}
```

All graph nodes are under the `nodes` key as an array. The `n` key
identifies the node using a unique index. All node references in edges
use these indexes. A node may include all outgoing edges embedded in
it, or edges may be included under a separate top-level array
`edges`. If the edge is included in the node, the edge only has a `to`
field that gives the target node index as the node containing the edge
is assumed to be the source node. Edges under the top-level `edges`
array include both a `from` and a `to` index.

Standard library JSON marshaler/unmarshaler does not work with graphs,
because the edge and node property values are of type
`interface{}`. The `JSON` struct can be used to marshal and unmarshal
graphs with custom property marshaler and unmarshalers.

This Go module is part of the [Layered Schema
Architecture](https://layeredschemas.org).

