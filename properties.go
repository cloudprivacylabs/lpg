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
	"fmt"
	"strings"
)

type properties map[int]any

// GetProperty returns the value for the key, and whether or not key
// exists. p can be nil
func (p *properties) getProperty(strTable *stringTable, key string) (interface{}, bool) {
	if *p == nil {
		return nil, false
	}
	lookupIdx, ok := strTable.lookup(key)
	if !ok {
		return nil, false
	}
	x, ok := (*p)[lookupIdx]
	return x, ok
}

// ForEachProperty calls f for each property in p until f returns
// false. Returns false if f returned false. p can be nil
func (p *properties) forEachProperty(strTable *stringTable, f func(string, interface{}) bool) bool {
	if *p == nil {
		return true
	}
	for k, v := range *p {
		if !f(strTable.str(k), v) {
			return false
		}
	}
	return true
}

// WithNativeValue is used to return a native value for property
// values. If the property value implements this interface, the
// underlying native value for indexing and comparison is obtained
// from this interface.
type WithNativeValue interface {
	GetNativeValue() interface{}
}

// ComparePropertyValue compares a and b. They must be
// comparable. Supported types are
//
//	int
//	string
//	[]int
//	[]string
//	[]interface
//
// The []interface must have one of the supported types as its elements
//
// If one of the values implement GetNativeValue() method, then it is
// called to get the underlying value
func ComparePropertyValue(a, b interface{}) int {

	if n, ok := a.(WithNativeValue); ok {
		return ComparePropertyValue(n.GetNativeValue(), b)
	}
	if n, ok := b.(WithNativeValue); ok {
		return ComparePropertyValue(a, n.GetNativeValue())
	}

	switch v1 := a.(type) {
	case string:
		if v2, ok := b.(string); ok {
			if v1 == v2 {
				return 0
			}
			if v1 < v2 {
				return -1
			}
			return 1
		}

	case int:
		if v2, ok := b.(int); ok {
			if v1 == v2 {
				return 0
			}
			if v1 < v2 {
				return -1
			}
			return 1
		}

	case []string:
		if v2, ok := b.([]string); ok {
			l1 := len(v1)
			l2 := len(v2)
			for i := 0; i < l1 && i < l2; i++ {
				if v1[i] < v2[i] {
					return -1
				}
				if v1[i] > v2[i] {
					return 1
				}
			}
			if l1 < l2 {
				return -1
			}
			if l1 > l2 {
				return 1
			}
			return 0
		}
		if v2, ok := b.([]interface{}); ok {
			l1 := len(v1)
			l2 := len(v2)
			for i := 0; i < l1 && i < l2; i++ {
				switch ComparePropertyValue(v1[i], v2[i]) {
				case -1:
					return -1
				case 1:
					return 1
				}
			}
			if l1 < l2 {
				return -1
			}
			if l1 > l2 {
				return 1
			}
			return 0
		}

	case []int:
		if v2, ok := b.([]int); ok {
			l1 := len(v1)
			l2 := len(v2)
			for i := 0; i < l1 && i < l2; i++ {
				if v1[i] < v2[i] {
					return -1
				}
				if v1[i] > v2[i] {
					return 1
				}
			}
			if l1 < l2 {
				return -1
			}
			if l1 > l2 {
				return 1
			}
			return 0
		}
		if v2, ok := b.([]interface{}); ok {
			l1 := len(v1)
			l2 := len(v2)
			for i := 0; i < l1 && i < l2; i++ {
				switch ComparePropertyValue(v1[i], v2[i]) {
				case -1:
					return -1
				case 1:
					return 1
				}
			}
			if l1 < l2 {
				return -1
			}
			if l1 > l2 {
				return 1
			}
			return 0
		}

	case []interface{}:
		if v2, ok := b.([]interface{}); ok {
			l1 := len(v1)
			l2 := len(v2)
			for i := 0; i < l1 && i < l2; i++ {
				switch ComparePropertyValue(v1[i], v2[i]) {
				case -1:
					return -1
				case 1:
					return 1
				}
			}
			if l1 < l2 {
				return -1
			}
			if l1 > l2 {
				return 1
			}
			return 0
		}
		if v2, ok := b.([]string); ok {
			return -ComparePropertyValue(v2, v1)
		}
		if v2, ok := b.([]int); ok {
			return -ComparePropertyValue(v2, v1)
		}
	}
	panic(fmt.Sprintf("Incomparable values: %v (%T) vs %v (%T)", a, a, b, b))
}

func (p properties) String() string {
	elements := make([]string, 0, len(p))
	for k, v := range p {
		if _, node := v.(*Node); node {
			continue
		}
		if _, edge := v.(*Edge); edge {
			continue
		}
		elements = append(elements, fmt.Sprintf("%d:%v", k, v))
	}
	return "{" + strings.Join(elements, " ") + "}"
}

// lookup proprs from source. allocate to target
func (p properties) clone(sourceGraph, targetGraph *Graph, cloneProperty func(string, interface{}) interface{}) properties {
	ret := make(properties, len(p))
	for k, v := range p {
		key := sourceGraph.stringTable.str(k)
		ix := targetGraph.stringTable.allocate(key)
		ret[ix] = cloneProperty(key, v)
	}
	return ret
}
