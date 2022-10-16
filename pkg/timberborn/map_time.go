package timberborn

import (
	"bytes"
	"time"
)

type MapTime struct {
	Time time.Time
}

func (me MapTime) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(`"`)
	buf.WriteString(me.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(`"`)

	return buf.Bytes(), nil
}
