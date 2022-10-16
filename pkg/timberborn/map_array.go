package timberborn

import (
	"bytes"
	"fmt"
)

type MapArray[T any] struct {
	width   int
	content []T
}

func NewMapArray[T any](width, height int) MapArray[T] {
	return MapArray[T]{
		width:   width,
		content: make([]T, width*height),
	}
}

func (me *MapArray[T]) Set(x, y int, value T) error {
	position := y*me.width + x
	if position >= len(me.content) {
		return fmt.Errorf("could not set %v,%v as it is out-of-range", x, y)
	}
	me.content[position] = value
	return nil
}

func (me *MapArray[T]) Get(x, y int) (T, error) {
	position := y*me.width + x
	if position >= len(me.content) {
		return me.content[0], fmt.Errorf("could not set %v,%v as it is out-of-range", x, y)
	}
	return me.content[position], nil
}

func (me MapArray[T]) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(`{"Array":"`)
	for i, v := range me.content {
		raw := fmt.Sprint(v)
		buf.WriteString(raw)
		if i+1 < len(me.content) {
			buf.WriteString(` `)
		}
	}
	buf.WriteString(`"}`)

	return buf.Bytes(), nil
}
