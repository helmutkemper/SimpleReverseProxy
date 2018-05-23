package SimpleReverseProxy

import (
  "io"
  "io/ioutil"
  "encoding/json"
  "errors"
)

func Input(r *ProxyRequest, inData interface{}) (error, string) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    return err, ""
  }

  err = r.Body.Close()
  if err != nil {
    return err, ""
  }

  if len(body) != 0 {
    err = json.Unmarshal(body, inData)
    return err, ""
  } else {
    return nil, r.ExpRegMatches["id"]

  }

  return errors.New("raw input data error"), ""
}