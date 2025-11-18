package interpreter

import (
	"fmt"
	"sync"
)

type ArrayValue struct {
	Values []Value
	mu     *sync.RWMutex
}

func (av ArrayValue) Content() string {
	av.mu.RLock()
	defer av.mu.RUnlock()
	str := "["
	for i, val := range av.Values {
		str += val.Content()
		if i < len(av.Values)-1 {
			str += ", "
		}
	}
	str += "]"
	return str
}

func (av ArrayValue) IsTruthy() bool {
	av.mu.RLock()
	defer av.mu.RUnlock()
	return len(av.Values) > 0
}

func (av ArrayValue) Type() ValueType {
	return VAL_ARRAY
}

func (av ArrayValue) Data() any {
	av.mu.RLock()
	defer av.mu.RUnlock()
	return av.Values
}

func (av ArrayValue) TypeName() string {
	return "ARRAY"
}

func (av *ArrayValue) Get(index int) (Value, bool) {
	av.mu.RLock()
	defer av.mu.RUnlock()
	if index < 0 || index >= len(av.Values) {
		return nil, false
	}
	return av.Values[index], true
}

func (av *ArrayValue) Set(index int, val Value) bool {
	av.mu.Lock()
	defer av.mu.Unlock()
	if index < 0 || index >= len(av.Values) {
		return false
	}
	av.Values[index] = val
	return true
}

func (av *ArrayValue) Append(val Value) {
	av.mu.Lock()
	defer av.mu.Unlock()
	av.Values = append(av.Values, val)
}

func (av *ArrayValue) Len() int {
	av.mu.RLock()
	defer av.mu.RUnlock()
	return len(av.Values)
}

func (av *ArrayValue) SetAt(indexes []int, val Value) error {
	if len(indexes) == 0 {
		return fmt.Errorf("no indexes provided")
	}
	av.mu.Lock()
	defer av.mu.Unlock()

	idx := indexes[0]
	if idx < 0 || idx >= len(av.Values) {
		return fmt.Errorf("array index %d out of bounds", idx)
	}
	if len(indexes) == 1 {
		av.Values[idx] = val
		return nil
	}
	subArray, ok := av.Values[idx].(ArrayValue)
	if !ok {
		return fmt.Errorf("expected array at index %d", idx)
	}
	return subArray.SetAt(indexes[1:], val)
}
