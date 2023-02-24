package lpg

import (
	"testing"
)

func TestPattern(t *testing.T) {
	graph := NewGraph()
	graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[5].SetProperty("key", "value")
	symbols := make(map[string]*PatternSymbol)
	pat := Pattern{
		{},
		{Min: 1, Max: 1},
		{Name: "nodes", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	if _, i := pat.getFastestElement(graph, map[string]*PatternSymbol{}); i != 2 {
		t.Errorf("Expecting 2, got %d", i)
	}
	plan, err := pat.GetPlan(graph, symbols)
	if err != nil {
		t.Error(err)
		return
	}
	acc := &DefaultMatchAccumulator{}
	plan.Run(graph, symbols, acc)
	if _, ok := acc.Symbols[0]["nodes"].(*Node); !ok {
		t.Errorf("Expecting one node, got: %v", acc)
	}

	pat = Pattern{
		{Labels: NewStringSet("bogus")},
		{Min: 1, Max: 1},
		{Name: "nodes", Properties: map[string]interface{}{"key": "value"}},
	}

	symbols = make(map[string]*PatternSymbol)
	plan, err = pat.GetPlan(graph, symbols)
	if err != nil {
		t.Error(err)
		return
	}
	acc = &DefaultMatchAccumulator{}
	plan.Run(graph, symbols, acc)
	if len(acc.Paths) != 0 {
		t.Errorf("Expecting 0 node, got: %+v", acc)
	}

	pat = Pattern{
		{},
		{},
		{Properties: map[string]interface{}{"key": "value2"}},
	}
	if _, i := pat.getFastestElement(graph, map[string]*PatternSymbol{}); i != 2 {
		t.Errorf("Expecting 2, got %d", i)
	}
	pat = Pattern{
		{Properties: map[string]interface{}{"key": "value2"}},
		{},
		{},
	}
	if _, i := pat.getFastestElement(graph, map[string]*PatternSymbol{}); i != 0 {
		t.Errorf("Expecting 0, got %d", i)
	}

}

func TestLoopPattern(t *testing.T) {
	graph := NewGraph()
	graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	symbols := make(map[string]*PatternSymbol)
	symbols["n"] = &PatternSymbol{}
	symbols["n"].Add(nodes[0])
	pat := Pattern{
		{Name: "n"},
		{Min: 1, Max: 1},
		{Name: "n"},
	}
	out := DefaultMatchAccumulator{}
	err := pat.Run(graph, symbols, &out)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out.Symbols) > 0 {
		t.Errorf("Expecting 0 node, got: %+v", out)
	}

	// Create a loop
	graph.NewEdge(nodes[0], nodes[0], "label", nil)
	out = DefaultMatchAccumulator{}
	err = pat.Run(graph, symbols, &out)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out.Symbols) != 1 {
		t.Errorf("Expecting 1 node, got: %v", out)
	}
}

func TestVariableLengthPath(t *testing.T) {
	graph := NewGraph()
	graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}

	nodes[1].SetProperty("property", "value")
	nodes[4].SetProperty("property", "value")

	symbols := make(map[string]*PatternSymbol)
	pat := Pattern{
		{Name: "n", Properties: map[string]interface{}{"property": "value"}},
		{Min: 1, Max: 1},
		{Properties: map[string]interface{}{"property": "value"}},
	}
	out := DefaultMatchAccumulator{}
	err := pat.Run(graph, symbols, &out)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out.Paths) != 0 {
		t.Errorf("Expecting 0 nodes")
	}
	pat = Pattern{
		{Name: "n", Properties: map[string]interface{}{"property": "value"}},
		{Min: 1, Max: 4},
		{Properties: map[string]interface{}{"property": "value"}},
	}
	out = DefaultMatchAccumulator{}
	err = pat.Run(graph, symbols, &out)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out.Paths[0].([]*Edge)) != 3 {
		t.Errorf("Expecting 3 nodes: %+v", out)
	}

}

func TestSimpleDirectedPathPatternWithIndex(t *testing.T) {
	testSimpleDirectedPathPattern(t, true)
}
func TestSimpleDirectedPathPatternWithoutIndex(t *testing.T) {
	testSimpleDirectedPathPattern(t, false)
}

// (n:label) -[*]-> (m:label)
func testSimpleDirectedPathPattern(t *testing.T, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[5].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: 1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[5].([]*Edge)[0].GetFrom() != nodes[5] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[5].([]*Edge)[0].GetFrom(), nodes[5])
	}
}

func TestSimplePathPatternPatternWithIndex(t *testing.T) {
	testSimplePathPattern(t, true)
}

func TestSimplePathPatternWithoutIndex(t *testing.T) {
	testSimplePathPattern(t, false)
}

// (n:label)-[]-(m:label)
func testSimplePathPattern(t *testing.T, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: -1, Max: 1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[4].([]*Edge)[0].GetFrom() != nodes[2] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[9].([]*Edge)[0].GetFrom(), nodes[2])
	}
}

func TestVariablePathPatternWithIndex(t *testing.T) {
	testVariablePathPattern(t, true)
}

func TestVariablePathPatternWithoutIndex(t *testing.T) {
	testVariablePathPattern(t, false)
}

// (n:label)-[*]-(m:label)
func testVariablePathPattern(t *testing.T, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: -1, Max: -1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[179].([]*Edge)[0].GetFrom() != nodes[8] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[9].([]*Edge)[0].GetFrom(), nodes[2])
	}
}

func TestPathLengthTwoPatternWithIndex(t *testing.T) {
	testPathLengthTwoPattern(t, true)
}

func TestPathLengthTwoPatternWithoutIndex(t *testing.T) {
	testPathLengthTwoPattern(t, false)
}

func testPathLengthTwoPattern(t *testing.T, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[4].SetProperty("key", "value")
	nodes[7].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: 2, Max: 2, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[27].([]*Edge)[0].GetFrom() != nodes[7] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[15].([]*Edge)[0].GetFrom(), nodes[7])
	}
}

func benchmarkSimpleDirectedPathPattern(b *testing.B, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[5].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: 1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}

func BenchmarkSimpleDirectedPathPatternWithIndex(b *testing.B) {
	benchmarkSimpleDirectedPathPattern(b, true)
}
func BenchmarkSimpleDirectedPathPatternWithoutIndex(b *testing.B) {
	benchmarkSimpleDirectedPathPattern(b, false)
}

func benchmarkSimplePathPattern(b *testing.B, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: -1, Max: 1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}

func BenchmarkSimplePathPatternWithIndex(b *testing.B) {
	benchmarkSimplePathPattern(b, true)
}
func BenchmarkSimplePathPatternWithoutIndex(b *testing.B) {
	benchmarkSimplePathPattern(b, false)
}

func benchmarkVariablePathPattern(b *testing.B, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: -1, Max: -1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}

func BenchmarkVariablePathPatternWithIndex(b *testing.B) {
	benchmarkVariablePathPattern(b, true)
}
func BenchmarkVariablePathPatternWithoutIndex(b *testing.B) {
	benchmarkVariablePathPattern(b, false)
}

func benchmarkPathLengthTwoPattern(b *testing.B, withIndex bool) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < 9; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	nodes[4].SetProperty("key", "value")
	nodes[7].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}},
		{Min: 2, Max: 2, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}
func BenchmarkPathLengthTwoPatternWithIndex(b *testing.B) {
	benchmarkPathLengthTwoPattern(b, true)
}
func BenchmarkPathLengthTwoPatternWithoutIndex(b *testing.B) {
	benchmarkPathLengthTwoPattern(b, false)
}
