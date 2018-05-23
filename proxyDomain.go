package SimpleReverseProxy

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

  Path                            string                  `json:"path"`
  PathExpReg                      string                  `json:"pathExpReg"`

  QueryStringEnable               bool                    `json:"queryStringEnable"`
}
func (el *ProxyDomain) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    Host                            string                  `json:"host"`
    Path                            string                  `json:"path"`
    PathExpReg                      string                  `json:"pathExpReg"`
  }{
    ErrorHandleAsString:    runtime.FuncForPC( reflect.ValueOf( el.ErrorHandle ).Pointer() ).Name(),
    NotFoundHandleAsString: runtime.FuncForPC( reflect.ValueOf( el.NotFoundHandle ).Pointer() ).Name(),
    Host:                   el.Host,
    Path:                   el.Path,
    PathExpReg:             el.PathExpReg,
  })
}
func (el *ProxyDomain) UnmarshalJSON(data []byte) error {
  type tmpStt struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    Host                            string                  `json:"host"`
    Path                            string                  `json:"path"`
    PathExpReg                      string                  `json:"pathExpReg"`
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

  el.Host       = tmp.Host
  el.Path       = tmp.Path
  el.PathExpReg = tmp.PathExpReg

  return nil
}
