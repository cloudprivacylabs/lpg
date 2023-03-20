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

func TestReverseSimplePath(t *testing.T) {
	graph, nodes := GetLineGraph(2, true)
	nodes[0].SetProperty("key", "value")
	nodes[1].SetProperty("key", "value")
	nodes[0].SetLabels(NewStringSet("a"))
	nodes[1].SetLabels(NewStringSet("b"))
	pat := Pattern{
		{Name: "n", Labels: StringSet{}, Properties: map[string]interface{}{"key": "value"}},
		{Min: -1, Max: -1, ToLeft: true},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(graph, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if len(acc.Paths) != 1 {
		t.Errorf("Expected length of path to be: %d got: %d", 1, len(acc.Paths))
	}
	if !acc.Paths[0].path[0].Reverse {
		t.Errorf("Expected path to be in reverse")
	}
}

func TestOCGetPattern(t *testing.T) {
	g := NewGraph()
	n1 := g.NewNode([]string{"root"}, nil)
	n2 := g.NewNode([]string{"c1"}, nil)
	n3 := g.NewNode([]string{"c2"}, nil)
	n4 := g.NewNode([]string{"c3"}, nil)

	g.NewEdge(n1, n2, "n1n2", nil)
	g.NewEdge(n1, n3, "n1n3", nil)
	g.NewEdge(n3, n4, "n3n4", nil)
	pat := Pattern{
		{Name: "this", Labels: NewStringSet("c1"), Properties: map[string]interface{}{}},
		{Min: 1, Max: 1, ToLeft: true},
		{Name: "", Labels: StringSet{}, Properties: map[string]interface{}{}},
		{Min: 1, Max: 1},
		{Name: "target", Labels: NewStringSet("c2"), Properties: map[string]interface{}{}},
	}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(g, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if len(acc.Paths) != 1 {
		t.Errorf("Length of accumulated paths should be %d: got %d", 1, len(acc.Paths))
	}
}

func TestOCGetPattern2(t *testing.T) {
	g := NewGraph()
	n1 := g.NewNode([]string{"root"}, nil)
	n2 := g.NewNode([]string{"c1"}, nil)
	n3 := g.NewNode([]string{"c2"}, nil)
	n4 := g.NewNode([]string{"c3"}, nil)

	g.NewEdge(n1, n2, "n1n2", nil)
	g.NewEdge(n1, n3, "n1n3", nil)
	g.NewEdge(n3, n4, "n3n4", nil)
	pat := Pattern{
		{Name: "this", Labels: NewStringSet("c1"), Properties: map[string]interface{}{}},
		{Min: 1, Max: 1, ToLeft: true},
		{Name: "", Labels: StringSet{}, Properties: map[string]interface{}{}},
		{Min: 1, Max: 1},
		{Name: "target", Labels: NewStringSet("c2"), Properties: map[string]interface{}{}},
	}
	symbols := make(map[string]*PatternSymbol)
	symbols["this"] = &PatternSymbol{}
	symbols["this"].Add(n2)
	acc := &DefaultMatchAccumulator{}
	if err := pat.Run(g, symbols, acc); err != nil {
		t.Error(err)
		return
	}
	if len(acc.Paths) != 1 {
		t.Errorf("Length of accumulated paths should be %d: got %d", 1, len(acc.Paths))
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
	if out.Paths[0].NumNodes() != 4 {
		t.Errorf("Expecting 4 nodes: %+v, got num nodes: %d", out, out.Paths[0].NumNodes())
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
	graph.NewEdge(nodes[n-1], nodes[n-1], "label", nil)
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
	n5, n6 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[5] {
			n5++
		}
		if p.GetEdge(0).GetFrom() == nodes[6] {
			n6++
		}
	}
	if n5 != 1 && n6 != 1 {
		t.Errorf("Expected number of paths to be 3, got %d", n6)
	}
}

func TestSimpleDirectedPathPatternSelfLoopsWithIndex(t *testing.T) {
	testSimpleDirectedPathPatternWithSelfLoops(t, true)
}

func TestSimpleDirectedPathPatternSelfLoopsWithoutIndex(t *testing.T) {
	testSimpleDirectedPathPatternWithSelfLoops(t, false)
}

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
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[5] {
			n5++
		}
		if p.GetEdge(0).GetFrom() == nodes[6] {
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
	nodes[5].SetLabels(NewStringSet("n5"))
	nodes[6].SetLabels(NewStringSet("n6"))
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
	n5, n6 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[5] {
			n5++
		}
		if p.GetEdge(0).GetFrom() == nodes[6] {
			n6++
		}
	}
	if n5 != 1 && n6 != 1 {
		t.Errorf("Expected number of paths to be 3, got %d", n6)
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
	n2, n8 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetFrom() == nodes[8] {
			n8++
		}
	}
	if n2 != 1 && n8 != 1 {
		t.Errorf("Expected number of paths to be 3, got %d", n8)
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
	nodes[2].SetProperty("key", "value")
	nodes[3].SetProperty("key", "value")
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
	if len(acc.Paths) != 6 {
		t.Errorf("expected length of path accumulator to be 6, got %d", len(acc.Paths))
	}
	n2, n3 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[3] {
			n3++
		}
	}
	if n2 != 3 {
		t.Errorf("Expected number of paths through n2 to be 2, got %d", n2)
	}
	if n3 != 3 {
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
	n2, n8 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetFrom() == nodes[8] {
			n8++
		}
	}
	if n2 != 1 {
		t.Errorf("Expected number of paths to be 1, got %d", n2)
	}
	if n8 != 1 {
		t.Errorf("Expected number of paths to be 1, got %d", n8)
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
	nodes[3].SetProperty("key", "value")
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
	n2, n3 := 0, 0
	if len(acc.Paths) != 2 {
		t.Errorf("Expected length of paths to be 2 got %d", len(acc.Paths))
	}
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[3] {
			n3++
		}
	}
	if n2 != 2 {
		t.Errorf("Expected number of paths to be 2, got %d", n2)
	}
	if n3 != 2 {
		t.Errorf("Expected number of paths to be 2, got %d", n3)
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
	graph, nodes := GetLineGraphWithSelfLoops(4, withIndex)
	nodes[1].SetProperty("key", "value")
	nodes[1].SetLabels(NewStringSet("b"))
	nodes[2].SetLabels(NewStringSet("c"))
	nodes[2].SetProperty("key", "value")
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
	if len(acc.Paths) != 6 {
		t.Errorf("Expected length of paths to be 6 got %d", len(acc.Paths))
	}
	n2, n3 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[1] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[2] {
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

func testVariablePathPatternCircleGraph(t *testing.T, withIndex bool) {
	graph, nodes := GetCircleGraph(4, withIndex)
	nodes[2].SetProperty("key", "value")
	nodes[3].SetProperty("key", "value")
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
	if len(acc.Paths) != 8 {
		t.Errorf("Expected length of paths to be 8 got %d", len(acc.Paths))
	}
	n2, n3 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[3] {
			n3++
		}
	}
	if n2 != 4 {
		t.Errorf("Expected number of paths to be 4, got %d", n2)
	}
	if n3 != 4 {
		t.Errorf("Expected number of paths to be 4, got %d", n3)
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
	nodes[2].SetProperty("key", "value")
	nodes[4].SetProperty("key", "value")
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
	n2, n4 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[3] {
			n4++
		}
	}
	if n2 != 1 && n4 != 1 {
		t.Errorf("Expected number of paths to be 1, got %d", n4)
	}
}

func TestPathLengthTwoPatternWithSelfLoopsWithIndex(t *testing.T) {
	testPathLengthTwoPatternWithSelfLoops(t, true)
}
func TestPathLengthTwoPatternWithSelfLoopsWithoutIndex(t *testing.T) {
	testPathLengthTwoPatternWithSelfLoops(t, false)
}

func testPathLengthTwoPatternWithSelfLoops(t *testing.T, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
	nodes[4].SetProperty("key", "value")
	nodes[6].SetProperty("key", "value")
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
	n4, n6 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[4] {
			n4++
		}
		if p.GetEdge(0).GetTo() == nodes[6] {
			n6++
		}
	}
	if n4 != 1 {
		t.Errorf("Expected number of paths through n4 to be 1, got %d", n4)
	}
	if n6 != 1 {
		t.Errorf("Expected number of paths through n6 to be 1, got %d", n6)
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
	nodes[2].SetProperty("key", "value")
	nodes[2].SetLabels(NewStringSet("n4"))
	nodes[4].SetLabels(NewStringSet("n7"))
	nodes[4].SetProperty("key", "value")
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
	if len(acc.Paths) != 2 {
		t.Errorf("expected length of path accumulator to be 2, got %d", len(acc.Paths))
	}
	n2, n4 := 0, 0
	for _, p := range acc.Paths {
		if p.GetEdge(0).GetFrom() == nodes[2] {
			n2++
		}
		if p.GetEdge(0).GetTo() == nodes[3] {
			n4++
		}
	}
	if n2 != 1 && n4 != 1 {
		t.Errorf("Expected number of paths to be 1, got %d", n4)
	}
}

func BenchmarkSimpleDirectedPathPatternWithIndex(b *testing.B) {
	benchmarkSimpleDirectedPathPattern(b, true)
	benchmarkSimpleDirectedPathPatternWithSelfLoops(b, true)
	benchmarkSimpleDirectedPathPatternCircleGraph(b, true)
}
func BenchmarkSimpleDirectedPathPatternWithoutIndex(b *testing.B) {
	benchmarkSimpleDirectedPathPattern(b, false)
	benchmarkSimpleDirectedPathPatternWithSelfLoops(b, false)
	benchmarkSimpleDirectedPathPatternCircleGraph(b, false)
}

func benchmarkSimpleDirectedPathPattern(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
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

func benchmarkSimpleDirectedPathPatternWithSelfLoops(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
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
		{Min: -1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}
func benchmarkSimpleDirectedPathPatternCircleGraph(b *testing.B, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
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
		{Min: -1, Max: 1},
		{Name: "n2", Labels: StringSet{}}}
	symbols := make(map[string]*PatternSymbol)
	acc := &DefaultMatchAccumulator{}
	for n := 0; n < b.N; n++ {
		pat.Run(graph, symbols, acc)
	}
}

func BenchmarkSimplePathPatternWithIndex(b *testing.B) {
	benchmarkSimplePathPattern(b, true)
	benchmarkSimplePathPatternWithSelfLoops(b, true)
	benchmarkSimplePathPatternCircleGraph(b, true)
}
func BenchmarkSimplePathPatternWithoutIndex(b *testing.B) {
	benchmarkSimplePathPattern(b, false)
	benchmarkSimplePathPatternWithSelfLoops(b, false)
	benchmarkSimplePathPatternCircleGraph(b, false)
}

func benchmarkSimplePathPattern(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
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

func benchmarkSimplePathPatternWithSelfLoops(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
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
func benchmarkSimplePathPatternCircleGraph(b *testing.B, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
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

func BenchmarkVariablePathPatternWithIndex(b *testing.B) {
	benchmarkVariablePathPattern(b, true)
	benchmarkVariablePathPatternCircleGraph(b, true)
	benchmarkVariablePathPatternWithSelfLoops(b, true)
}
func BenchmarkVariablePathPatternWithoutIndex(b *testing.B) {
	benchmarkVariablePathPattern(b, false)
	benchmarkVariablePathPatternCircleGraph(b, false)
	benchmarkVariablePathPatternWithSelfLoops(b, false)
}

func benchmarkVariablePathPattern(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
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

func benchmarkVariablePathPatternCircleGraph(b *testing.B, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
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

func benchmarkVariablePathPatternWithSelfLoops(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
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

func BenchmarkPathLengthTwoPatternWithIndex(b *testing.B) {
	benchmarkPathLengthTwoPattern(b, true)
	benchmarkPathLengthTwoPatternWithSelfLoops(b, true)
	benchmarkPathLengthTwoPatternCircleGraph(b, true)
}
func BenchmarkPathLengthTwoPatternWithoutIndex(b *testing.B) {
	benchmarkPathLengthTwoPattern(b, false)
	benchmarkPathLengthTwoPatternWithSelfLoops(b, false)
	benchmarkPathLengthTwoPatternCircleGraph(b, false)
}

func benchmarkPathLengthTwoPattern(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraph(10, withIndex)
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
func benchmarkPathLengthTwoPatternWithSelfLoops(b *testing.B, withIndex bool) {
	graph, nodes := GetLineGraphWithSelfLoops(10, withIndex)
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
func benchmarkPathLengthTwoPatternCircleGraph(b *testing.B, withIndex bool) {
	graph, nodes := GetCircleGraph(10, withIndex)
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
