package lpg

import (
	"fmt"
	"testing"
)

func TestStringTable(t *testing.T) {
	table := stringTable{}
	table.init()

	m := make(map[string]int)

	for i := 0; i < 1000; i++ {
		str := fmt.Sprint(i)
		m[str] = table.allocate(fmt.Sprint(i))
	}
	if len(table.strs) != len(m) {
		t.Errorf("Wrong table len, table: %d m: %d", len(table.strs), len(m))
	}
	if len(table.strmap) != len(m) {
		t.Errorf("Wrong table len")
	}
	if table.firstFree != -1 {
		t.Errorf("There are free elements")
	}
	for k, v := range m {
		if table.str(v) != k {
			t.Errorf("Wrong string at %d", v)
		}
		if l, ok := table.lookup(k); l != v || !ok {
			t.Errorf("Wrong index %s", k)
		}
	}
	for i := 0; i < 1000; i += 2 {
		table.free(i)
		delete(m, fmt.Sprint(i))
	}
	if len(table.strmap) != len(m) {
		t.Errorf("Wrong table len")
	}
	if table.firstFree == -1 {
		t.Errorf("No free elements")
	}
	for i := 2000; i < 3000; i++ {
		str := fmt.Sprint(i)
		m[str] = table.allocate(str)
	}
	if len(table.strs) != len(m) {
		t.Errorf("Wrong table len, table: %d m: %d", len(table.strs), len(m))
	}
	if len(table.strmap) != len(m) {
		t.Errorf("Wrong table len: %d %d", len(table.strmap), len(m))
	}
	if table.firstFree != -1 {
		t.Errorf("There are free elements")
	}
}
