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
	"container/list"
)

// A hashIndex is a hash table index
type hashIndex struct {
	values   map[interface{}]*fastSet
	elements list.List
}

func (ix *hashIndex) add(value interface{}, id int, item interface{}) {
	if ix.values == nil {
		ix.values = make(map[interface{}]*fastSet)
	}

	el := ix.elements.PushBack(item)
	fs, ok := ix.values[value]
	if !ok {
		fs = newFastSet()
		ix.values[value] = fs
	}
	fs.add(id, el)
}

func (ix *hashIndex) remove(value interface{}, id int) {
	if ix.values == nil {
		return
	}
	fs, ok := ix.values[value]
	if !ok {
		return
	}
	el, ok := fs.get(id)
	if !ok {
		return
	}
	fs.remove(id)
	ix.elements.Remove(el.(*list.Element))
}

// find returns the iterator and expected size.
func (ix *hashIndex) find(value interface{}) Iterator {
	if ix.values == nil {
		return emptyIterator{}
	}
	v, found := ix.values[value]
	if !found {
		return emptyIterator{}
	}
	itr := &procIterator{itr: v.iterator(), proc: func(in interface{}) interface{} { return in.(*list.Element).Value }}
	return withSize(itr, v.size())
}

func (ix *hashIndex) valueItr() Iterator {
	if ix.values == nil {
		return emptyIterator{}
	}
	return &listIterator{next: ix.elements.Front(), size: ix.elements.Len()}
}
