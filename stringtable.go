package lpg

import (
// "fmt"
)

// stringTable keeps unique strings in a slice indexed by a map, so
// identical strings are not repeated, and integers can be used
// instead of strings.
type stringTable struct {
	strs      []rcstring
	strmap    map[string]int
	firstFree int
}

type rcstring struct {
	str string
	c   int
}

func (table *stringTable) init() {
	table.strs = make([]rcstring, 0, 16)
	table.strmap = make(map[string]int)
	table.firstFree = -1
}

// allocate will allocate a string and return its index
func (table *stringTable) allocate(s string) int {
	// defer func() {
	// 	fmt.Printf("allocate %p %s %v\n", table, s, table)
	// }()
	// If string exists, increment reference count
	i, exists := table.strmap[s]
	if exists {
		table.strs[i].c++
		return i
	}
	if table.firstFree != -1 {
		ret := table.firstFree
		table.firstFree = table.strs[ret].c
		table.strs[ret].c = 1
		table.strs[ret].str = s
		table.strmap[s] = ret
		return ret
	}
	ret := len(table.strs)
	table.strmap[s] = ret
	table.strs = append(table.strs, rcstring{str: s, c: 1})
	return ret
}

// lookup will lookup a string index without allocating
func (table *stringTable) lookup(s string) (int, bool) {
	i, exists := table.strmap[s]
	if exists {
		return i, true
	}
	return 0, false
}

func (table *stringTable) str(i int) string {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("***panic, table: %p %v i %d\n", table, table, i)
	// 	}
	// }()
	return table.strs[i].str
}

func (table *stringTable) free(i int) {
	// defer func() {
	// 	fmt.Println("tree", i, table)
	// }()
	table.strs[i].c--
	if table.strs[i].c > 0 {
		return
	}
	table.strs[i].c = table.firstFree
	table.firstFree = i
	delete(table.strmap, table.strs[i].str)
}
