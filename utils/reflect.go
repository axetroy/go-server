package utils

import "encoding/json"

func Reflect(source interface{}, desc interface{}) (err error) {
	var jsonb []byte

	if jsonb, err = json.Marshal(source); err != nil {
		return
	}

	if err = json.Unmarshal(jsonb, desc); err != nil {
		return
	}

	return
}
