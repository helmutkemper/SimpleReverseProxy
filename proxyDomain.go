package marketPlaceProcy

import (
  "encoding/json"
  "runtime"
  "reflect"
)

type ProxyDomain struct {

  /*
  Função de erro padrão para o domínio.
  */
  ErrorHandle                     ProxyHandlerFunc        `json:"-"`

  /*
  Função de page not found padrão para o domínio
  */
  NotFoundHandle                  ProxyHandlerFunc        `json:"-"`

  Host                            string                  `json:"host"`

  QueryStringEnable               bool                    `json:"queryStringEnable"`
}
func (el *ProxyDomain) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    Host                            string                  `json:"host"`
    QueryStringEnable               bool                    `json:"queryStringEnable"`
  }{
    ErrorHandleAsString:    runtime.FuncForPC( reflect.ValueOf( el.ErrorHandle ).Pointer() ).Name(),
    NotFoundHandleAsString: runtime.FuncForPC( reflect.ValueOf( el.NotFoundHandle ).Pointer() ).Name(),
    Host:                   el.Host,
    QueryStringEnable:      el.QueryStringEnable,
  })
}
func (el *ProxyDomain) UnmarshalJSON(data []byte) error {
  type tmpStt struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    Host                            string                  `json:"host"`
    QueryStringEnable               bool                    `json:"queryStringEnable"`
  }
  var tmp = tmpStt{}
  err := json.Unmarshal( data, &tmp )
  if err != nil {
    return err
  }

  if tmp.ErrorHandleAsString != "" {
    el.ErrorHandle = FuncMap[ tmp.ErrorHandleAsString ].( ProxyHandlerFunc )
  }

  if tmp.NotFoundHandleAsString != "" {
    el.NotFoundHandle = FuncMap[ tmp.NotFoundHandleAsString ].(ProxyHandlerFunc)
  }

  el.Host              = tmp.Host
  el.QueryStringEnable = tmp.QueryStringEnable

  return nil
}
