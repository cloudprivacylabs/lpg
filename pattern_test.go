package lpg

import (
	"fmt"
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

// (n1)->(n2)->(n3)...
func GetLineGraph(n int, withIndex bool) (*Graph, []*Node) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < n; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < n-1; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	return graph, nodes
}

// / (n1)->(n1)->(n2)->(n2)->(n3)->(n3)...
func GetLineGraphWithSelfLoops(n int, withIndex bool) (*Graph, []*Node) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < n; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < n-1; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
		graph.NewEdge(nodes[i], nodes[i], "label", nil)
	}
	return graph, nodes
}

// (n1)->(n2)->(n3)->(n1)
func GetCircleGraph(n int, withIndex bool) (*Graph, []*Node) {
	graph := NewGraph()
	if withIndex {
		graph.index.NodePropertyIndex("key", graph, BtreeIndex)
	}
	nodes := make([]*Node, 0)
	for i := 0; i < n; i++ {
		nodes = append(nodes, graph.NewNode([]string{"a"}, nil))
	}
	for i := 0; i < n-1; i++ {
		graph.NewEdge(nodes[i], nodes[i+1], "label", nil)
	}
	graph.NewEdge(nodes[n-1], nodes[0], "label", nil)
	return graph, nodes
}

func GetTreeGraph(n, width int) (*Graph, []*Node) {
	graph := NewGraph()
	nodes := make([]*Node, 0)
	return graph, nodes
}

func TestSimpleDirectedPathPatternWithIndex(t *testing.T) {
	testSimpleDirectedPathPattern(t, true)
}
func TestSimpleDirectedPathPatternWithoutIndex(t *testing.T) {
	testSimpleDirectedPathPattern(t, false)
}

// (n:{prop:val}) -[*]-> (m:{prop:val})
func testSimpleDirectedPathPattern(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
	nodes[5].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[0].([]*Edge)[0].GetFrom() != nodes[5] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[0].([]*Edge)[0].GetFrom(), nodes[5])
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[6] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[6])
	}
}

func TestSimpleDirectedPathPatternSelfLoopsWithIndex(t *testing.T) {
	testSimpleDirectedPathPatternWithSelfLoops(t, true)
}

func TestSimpleDirectedPathPatternSelfLoopsWithoutIndex(t *testing.T) {
	testSimpleDirectedPathPatternWithSelfLoops(t, false)
}

// fail
func testSimpleDirectedPathPatternWithSelfLoops(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
	nodes[5].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if len(acc.Paths) != 4 {
		t.Errorf("Expected number of paths to be 4, got %d", len(acc.Paths))
	}
	n5 := 0
	n6 := 0
	for i := range acc.Paths {
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[5] {
			n5++
		}
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[6] {
			n6++
		}
	}
	if n5 != 2 {
		t.Errorf("Expected number of paths through n5 to be 2, got %d", n5)
	}
	if n6 != 2 {
		t.Errorf("Expected number of paths through n6 to be 2, got %d", n6)
	}
}

func TestSimpleDirectedPathPatternCircleGraphWithIndex(t *testing.T) {
	testSimpleDirectedPathPatternCircleGraph(t, true)
}

func TestSimpleDirectedPathPatternCircleGraphWithoutIndex(t *testing.T) {
	testSimpleDirectedPathPatternCircleGraph(t, false)
}

func testSimpleDirectedPathPatternCircleGraph(t *testing.T, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
	nodes[5].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[0].([]*Edge)[0].GetFrom() != nodes[5] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[0].([]*Edge)[0].GetFrom(), nodes[5])
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[6] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[6])
	}
}

func TestSimplePathPatternPatternWithIndex(t *testing.T) {
	testSimplePathPattern(t, true)
}

func TestSimplePathPatternWithoutIndex(t *testing.T) {
	testSimplePathPattern(t, false)
}

// (n:{prop:val})-[]-(m:{prop:val})
func testSimplePathPattern(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: 1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[2] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[2])
	}
	if acc.Paths[3].([]*Edge)[0].GetFrom() != nodes[8] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[3].([]*Edge)[0].GetFrom(), nodes[8])
	}
}

func TestSimplePathPatternWithSelfLoopsWithIndex(t *testing.T) {
	testSimplePathPatternWithSelfLoops(t, true)
}
func TestSimplePathPatternWithSelfLoopsWithoutIndex(t *testing.T) {
	testSimplePathPatternWithSelfLoops(t, false)
}

func testSimplePathPatternWithSelfLoops(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
	nodes[2].SetProperty("key", "valor")
	nodes[3].SetProperty("key", "valor")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "valor"}},
		{Min: -1, Max: 1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	n2 := 0
	n3 := 0
	for i := range acc.Paths {
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[2] {
			n2++
		}
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[3] {
			n3++
		}
	}
	if len(acc.Paths) != 5 {
		t.Errorf("expected length of path accumulator to be 5, got %d", len(acc.Paths))
	}
	if n2 != 2 {
		t.Errorf("Expected number of paths through n2 to be 2, got %d", n2)
	}
	if n3 != 2 {
		t.Errorf("Expected number of paths through n3 to be 2, got %d", n3)
	}
}

func TestSimplePathPatternCircleGraphWithIndex(t *testing.T) {
	testSimplePathPatternCircleGraph(t, true)
}
func TestSimplePathPatternCircleGraphWithoutIndex(t *testing.T) {
	testSimplePathPatternCircleGraph(t, false)
}

func testSimplePathPatternCircleGraph(t *testing.T, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: 1, Undirected: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[2] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[2])
	}
	if acc.Paths[3].([]*Edge)[0].GetFrom() != nodes[8] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[3].([]*Edge)[0].GetFrom(), nodes[8])
	}
}

func TestVariablePathPatternWithIndex(t *testing.T) {
	testVariablePathPattern(t, true)
}

func TestVariablePathPatternWithoutIndex(t *testing.T) {
	testVariablePathPattern(t, false)
}

// (n:{prop:val})-[*]-(m:{prop:val})
func testVariablePathPattern(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: -1, Undirected: true},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[2] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[2])
	}
	if acc.Paths[7].([]*Edge)[0].GetFrom() != nodes[8] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[7].([]*Edge)[0].GetFrom(), nodes[8])
	}
}

func TestVariablePathPatternSelfLoopsWithIndex(t *testing.T) {
	testVariablePathPatternWithSelfLoops(t, true)
}

func TestVariablePathPatternSelfLoopsWithoutIndex(t *testing.T) {
	testVariablePathPatternWithSelfLoops(t, false)
}

// fail
func testVariablePathPatternWithSelfLoops(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[2].SetLabels(NewStringSet("node2"))
	nodes[3].SetLabels(NewStringSet("node3"))
	nodes[3].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: -1},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(acc.Paths))
	for _, p := range acc.Paths {
		for _, path := range p.([]*Edge) {
			fmt.Println(path.GetFrom(), path.GetTo())
		}
		fmt.Println()
		// fmt.Println(acc.Paths[i].([]*Edge)[0].GetFrom())
		// fmt.Println(acc.Paths[i].([]*Edge)[0].GetTo())
	}
	// if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[2] {
	// 	t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[2])
	// }
	// if acc.Paths[7].([]*Edge)[0].GetFrom() != nodes[8] {
	// 	t.Errorf("Expecting %v, got: %v", acc.Paths[7].([]*Edge)[0].GetFrom(), nodes[8])
	// }
	if len(acc.Paths) != 6 {
		t.Errorf("Expected length of paths")
	}
	n2, n3 := 0, 0
	for i := range acc.Paths {
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[2] {
			n2++
		}
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[3] {
			n3++
		}
	}
}

func TestVariablePathPatternCircleGraphWithIndex(t *testing.T) {
	testVariablePathPatternCircleGraph(t, true)
}

func TestVariablePathPatternCircleGraphWithoutIndex(t *testing.T) {
	testVariablePathPatternCircleGraph(t, false)
}

// fail
func testVariablePathPatternCircleGraph(t *testing.T, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[8].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: -1, Undirected: true},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	// 16 paths
	// 4, 5, 6, 7
	if acc.Paths[2].([]*Edge)[0].GetFrom() != nodes[2] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[1].([]*Edge)[0].GetFrom(), nodes[2])
	}
	// 12, 13, 14, 15
	if acc.Paths[7].([]*Edge)[0].GetFrom() != nodes[8] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[7].([]*Edge)[0].GetFrom(), nodes[8])
	}
}

func TestPathLengthTwoPatternWithIndex(t *testing.T) {
	testPathLengthTwoPattern(t, true)
}

func TestPathLengthTwoPatternWithoutIndex(t *testing.T) {
	testPathLengthTwoPattern(t, false)
}

func testPathLengthTwoPattern(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
	nodes[4].SetProperty("key", "value")
	nodes[7].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 2, Max: 2, Undirected: true},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[4] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[6].([]*Edge)[0].GetFrom(), nodes[7])
	}
	if acc.Paths[3].([]*Edge)[0].GetFrom() != nodes[7] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[6].([]*Edge)[0].GetFrom(), nodes[7])
	}
}

func TestPathLengthTwoPatternWithSelfLoopsWithIndex(t *testing.T) {
	testPathLengthTwoPatternWithSelfLoops(t, true)
}
func TestPathLengthTwoPatternWithSelfLoopsWithoutIndex(t *testing.T) {
	testPathLengthTwoPatternWithSelfLoops(t, false)
}

// fail
func testPathLengthTwoPatternWithSelfLoops(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
	nodes[4].SetProperty("key", "value")
	nodes[7].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 2, Max: 2, Undirected: true},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	// 12 paths
	n4, n7 := 0, 0
	for i := range acc.Paths {
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[4] {
			n4++
		}
		if acc.Paths[i].([]*Edge)[0].GetFrom() == nodes[4] {
			n7++
		}
	}
	// 5
	if n4 != 0 {
		t.Errorf("Expected number of paths through n4 to be 3, got %d", n4)
	}
	// 5
	if n7 != 0 {
		t.Errorf("Expected number of paths through n4 to be 3, got %d", n7)
	}
}

func TestPathLengthTwoPatternCircleGraphWithIndex(t *testing.T) {
	testPathLengthTwoPatternCircleGraph(t, true)
}
func TestPathLengthTwoPatternCircleGraphWithoutIndex(t *testing.T) {
	testPathLengthTwoPatternCircleGraph(t, false)
}

func testPathLengthTwoPatternCircleGraph(t *testing.T, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
	nodes[4].SetProperty("key", "value")
	nodes[7].SetProperty("key", "value")
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: 2, Max: 2, Undirected: true},
		{Name: "n2", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if acc.Paths[1].([]*Edge)[0].GetFrom() != nodes[4] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[6].([]*Edge)[0].GetFrom(), nodes[7])
	}
	if acc.Paths[3].([]*Edge)[0].GetFrom() != nodes[7] {
		t.Errorf("Expecting %v, got: %v", acc.Paths[6].([]*Edge)[0].GetFrom(), nodes[7])
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
