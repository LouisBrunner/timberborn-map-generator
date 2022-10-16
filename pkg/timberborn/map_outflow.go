package timberborn

import (
	"fmt"
	"strings"
)

type MapOutflow struct {
	// TODO: no idea what does values mean
	A int
	B int
	C int
	D int
}

func (me MapOutflow) String() string {
	// FIXME: just need a generic ForEach
	return strings.Join([]string{
		fmt.Sprint(me.A),
		fmt.Sprint(me.B),
		fmt.Sprint(me.C),
		fmt.Sprint(me.D),
	}, ":")
}
