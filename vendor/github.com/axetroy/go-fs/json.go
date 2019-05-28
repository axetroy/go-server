package fs

import (
  "io/ioutil"
  "encoding/json"
  "os"
  "path"
)

func ReadJson(filepath string) (map[string]interface{}, error) {
  if b, err := ioutil.ReadFile(filepath); err != nil {
    return nil, err
  } else {
    m := new(map[string]interface{})
    if err = json.Unmarshal(b, m); err != nil {
      return nil, err
    } else {
      return *m, nil
    }
  }
}

func WriteJson(filepath string, data []byte) (error) {
  m := new(map[string]interface{})
  if err := json.Unmarshal(data, m); err != nil {
    return err
  }
  if err := ioutil.WriteFile(filepath, data, os.ModePerm); err != nil {
    return err
  }
  return nil
}

func OuputJson(filepath string, data []byte) error {
  if err := EnsureDir(path.Dir(filepath)); err != nil {
    return err
  }
  return WriteJson(filepath, data)
}
