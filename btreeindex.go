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
	"github.com/emirpasic/gods/trees/btree"
)

// A setTree is a B-Tree of linkedhashsets
type setTree struct {
	tree *btree.Tree
}

func (s *setTree) add(key interface{}, id int, item interface{}) {
	if s.tree == nil {
		s.tree = btree.NewWith(16, ComparePropertyValue)
	}
	v, found := s.tree.Get(key)
	if !found {
		v = newFastSet()
		s.tree.Put(key, v)
	}
	set := v.(*fastSet)
	set.add(id, item)
}

func (s setTree) remove(key interface{}, id int) {
	if s.tree == nil {
		return
	}
	v, found := s.tree.Get(key)
	if !found {
		return
	}
	set := v.(*fastSet)
	set.remove(id)
	if set.size() == 0 {
		s.tree.Remove(key)
	}
}

// find returns the iterator and expected size.
func (s setTree) find(key interface{}) Iterator {
	if s.tree == nil {
		return emptyIterator{}
	}
	v, found := s.tree.Get(key)
	if !found {
		return emptyIterator{}
	}
	set := v.(*fastSet)
	itr := set.iterator()
	return withSize(itr, set.size())
}

func (s setTree) valueItr() Iterator {
	if s.tree == nil {
		return emptyIterator{}
	}
	treeItr := s.tree.Iterator()
	return &funcIterator{
		iteratorFunc: func() Iterator {
			if !treeItr.Next() {
				return nil
			}
			set := treeItr.Value().(*fastSet)
			itr := set.iterator()
			return withSize(itr, set.size())
		},
	}
}
