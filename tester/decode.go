package tester

import "encoding/json"

func Decode(source interface{}, dest interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(source); err != nil {
		return
	}

	if err = json.Unmarshal(b, dest); err != nil {
		return
	}
	return
}
