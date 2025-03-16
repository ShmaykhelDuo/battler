package model

import (
	"reflect"

	tf "github.com/wamuir/graft/tensorflow"
)

func NewTensor(value any) (*tf.Tensor, error) {
	val := reflect.ValueOf(value)
	t := val.Type()

	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		val = wrapSlice(val)
	}

	val = wrapSlice(val)

	return tf.NewTensor(val.Interface())
}

func wrapSlice(value reflect.Value) reflect.Value {
	resType := reflect.SliceOf(value.Type())
	res := reflect.MakeSlice(resType, 0, 1)
	res = reflect.Append(res, value)
	return res
}
