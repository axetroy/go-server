package proto

import (
	"encoding/base64"
	"encoding/json"
)

type Name string

var (
	Auth Name = "auth"
)

type Proto struct {
	name   Name
	params interface{}
}

func NewProto(name Name, params interface{}) Proto {
	return Proto{
		name:   name,
		params: params,
	}
}

func (p Proto) String() (string, error) {
	b, err := json.Marshal(p.params)

	if err != nil {
		return "", err
	}

	text := base64.StdEncoding.EncodeToString(b)

	path := string(p.name) + "://" + text

	return path, nil
}
