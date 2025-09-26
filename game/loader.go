package game

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var LoadMap = map[string]any{} // game specific data loaded from yaml files

var ErrNoObj = errors.New("unknown object id")

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

		var obj, ok = LoadMap[id]
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

func LoadChain(fpath string) (err error) {
	var b []byte
	if b, err = os.ReadFile(fpath); err != nil {
		return
	}
	return ReadChain(b)
}
