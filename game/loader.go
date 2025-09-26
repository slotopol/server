package game

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var DataRouter = map[string]any{} // object router for data loading

var LoadMap = [][]byte{} // storage for loaded data

var ErrNoObj = errors.New("unknown object id")

// ReadChain reads and decodes a sequence of objects from a byte slice.
// Each object is identified by a string id, followed by its YAML representation.
// The objects are looked up in the DataRouter map.
// If an object implements the Clearer interface, its Clear method is called before decoding.
func ReadChain(b []byte) (err error) {
	type Clearer interface {
		Clear()
	}

	var dec = yaml.NewDecoder(bytes.NewReader(b))
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
func MustReadChain(b []byte) {
	if err := ReadChain(b); err != nil {
		panic(err)
	}
}

// LoadChain reads and decodes a sequence of objects from a YAML file.
func LoadChain(fpath string) (err error) {
	var b []byte
	if b, err = os.ReadFile(fpath); err != nil {
		return
	}
	return ReadChain(b)
}
