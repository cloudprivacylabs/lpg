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
