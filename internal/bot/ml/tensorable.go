package ml

import "encoding/json"

type Tensorable interface {
	Value() any
}

type TensorableScalar interface {
	~float32 | ~float64 | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string | ~bool
}

type TensorableValue[T TensorableScalar] struct {
	Item T
}

func (t TensorableValue[T]) Value() any {
	return [1][1]T{{t.Item}}
}

func (t TensorableValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Item)
}

type TensorableSlice[T TensorableScalar] struct {
	Items []T
}

func (t TensorableSlice[T]) Value() any {
	return [1][]T{t.Items}
}

func (t TensorableSlice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Items)
}
