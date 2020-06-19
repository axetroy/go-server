package proto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
)

type Name string

func (n Name) String() string {
	return string(n)
}

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

func (p Proto) Data() ([]byte, error) {
	b, err := json.Marshal(p.params)

	return b, err
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

func Parse(link string) (*Proto, error) {
	reg := regexp.MustCompile(`^([\w_]+)://(.*)$`)

	if !reg.MatchString(link) {
		return nil, fmt.Errorf("invalid link '%s'", link)
	}

	matcher := reg.FindAllStringSubmatch(link, 1)

	if len(matcher) == 0 {
		return nil, fmt.Errorf("invalid link '%s'", link)
	}

	firstMatch := matcher[0]
	schema := firstMatch[1]
	data := firstMatch[2]

	var params map[string]interface{}

	b, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &params); err != nil {
		return nil, err
	}

	p := NewProto(Name(schema), params)

	return &p, nil
}
