package game

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var DataRouter = map[string]any{} // object router for data loading

var LoadMap = [][]byte{} // storage for loaded data

var (
	ErrNoObj  = errors.New("unknown object id")
	ErrAidHas = errors.New("algorithm ID already declared")
)

// ReadChain reads and decodes a sequence of objects from a stream.
// Each object is identified by a string id, followed by its YAML representation.
// The objects are looked up in the DataRouter map.
// If an object implements the Clearer interface, its Clear method is called before decoding.
func ReadChain(r io.Reader) (err error) {
	type Clearer interface {
		Clear()
	}

	var dec = yaml.NewDecoder(r)
	for {
		var id string
		if err = dec.Decode(&id); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}

		var obj, ok = DataRouter[id]
		if !ok {
			err = fmt.Errorf("%w: %s", ErrNoObj, id)
			return
		}
		if v, ok := obj.(Clearer); ok {
			v.Clear()
		}
		if err = dec.Decode(obj); err != nil {
			return
		}
	}
}

// MustReadChain is like ReadChain but panics if an error occurs.
func MustReadChain(r io.Reader) {
	if err := ReadChain(r); err != nil {
		panic(err)
	}
}
