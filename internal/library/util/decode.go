package util

import (
	"encoding/json"
	"errors"
)

func Decode(dest interface{}, src interface{}) error {
	if !IsPoint(dest) {
		return errors.New("decode: dest expect a point")
	}

	if b, err := json.Marshal(src); err != nil {
		return err
	} else {
		if err := json.Unmarshal(b, dest); err != nil {
			return err
		}
	}
	return nil
}
