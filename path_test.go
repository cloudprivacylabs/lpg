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

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestPathBasic(t *testing.T) {
	f, err := os.Open("testdata/g1.json")
	if err != nil {
		t.Error(err)
		return
	}
	target := NewGraph()
	err = JSON{}.Decode(target, json.NewDecoder(f))
	if err != nil {
		t.Error(err)
		return
	}

	find := func(id string) *Node {
		for nodes := target.GetNodes(); nodes.Next(); {
			node := nodes.Node()
			if s, _ := node.GetProperty("id"); s == id {
				return node
			}
		}
		return nil
	}
	cursor := Cursor{}
	cursor.Set(find("0"))
	cursor.StartPath()
	if cursor.GetPath().NumNodes() != 1 {
		t.Errorf("Wrong numnodes: %v", cursor.GetPath())
	}
	if cursor.GetPath().NumEdges() != 0 {
		t.Errorf("Wrong numEdges: %v", cursor.GetPath())
	}
	itr := cursor.ForwardWith("a")
	itr.Next()
	cursor.PushToPath(itr.Edge())
	if cursor.GetPath().NumNodes() != 2 {
		t.Errorf("Wrong numnodes: %v", cursor.GetPath())
	}
	if cursor.GetPath().NumEdges() != 1 {
		t.Errorf("Wrong numEdges: %v", cursor.GetPath())
	}
}

func TestPathSlice(t *testing.T) {
	f, err := os.Open("testdata/g1.json")
	if err != nil {
		t.Error(err)
		return
	}
	target := NewGraph()
	err = JSON{}.Decode(target, json.NewDecoder(f))
	if err != nil {
		t.Error(err)
		return
	}
	path := Path{path: make([]PathElement, 0)}
	for itr := target.GetEdges(); itr.Next(); {
		path.path = append(path.path, PathElement{
			Edge: itr.Edge(),
		})
	}
	ps := []struct {
		start    int
		end      int
		expected Path
	}{
		{
			start:    0,
			end:      -1,
			expected: Path{path: path.path},
		},
		{
			start:    0,
			end:      1,
			expected: Path{path: path.path[0:0]},
		},
	}
	for _, sp := range ps {
		if !reflect.DeepEqual(sp.expected, path.Slice(sp.start, sp.end)) {
			t.Errorf("expected %v, got: %v", sp.expected, path.Slice(sp.start, sp.end))
		}
	}
}

func TestPathHasPrefix(t *testing.T) {

}

func TestAppend(t *testing.T) {
	f, err := os.Open("testdata/g1.json")
	if err != nil {
		t.Error(err)
		return
	}
	target := NewGraph()
	err = JSON{}.Decode(target, json.NewDecoder(f))
	if err != nil {
		t.Error(err)
		return
	}
	path := Path{path: make([]PathElement, 0)}
	for itr := target.GetEdges(); itr.Next(); {
		path.path = append(path.path, PathElement{
			Edge: itr.Edge(),
		})
	}
	cp := make([]PathElement, len(path.path))
	copy(cp, path.path)
	fmt.Println(cp)
	fmt.Println(path.path)
	path.Append(cp...)
	t.Error()
}

// test
// slice
// hasprefix
