package status

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Step uint8

const (
	Init Step = iota
	Backuped
	Packed
	RemoveExpired
	Finished
)

func (s Step) String() string {
	switch s {
	case Init:
		return `init`
	case Backuped:
		return `backuped`
	case Packed:
		return `packed`
	case RemoveExpired:
		return `remove expired`
	case Finished:
		return `finished`
	}
	return `unknow`
}

type exportElement struct {
	ID       uint64
	Step     string
	Created  time.Time
	Backuped time.Time
}
type Element struct {
	ID       uint64
	Step     Step
	Created  time.Time
	Backuped time.Time
}

func Marshal(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	e := enc.Encode(obj)
	if e != nil {
		return nil, e
	}
	return buffer.Bytes(), nil
}
func Unmarshal(b []byte, obj interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.Decode(obj)
}
