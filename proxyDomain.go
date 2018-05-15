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

  /*
  [opcional] sub domínio sem ponto final. Ex.: blog.domínio.com fica apenas blog
  */
  SubDomain                       string                  `json:"subDomain"`

  /*
  Domínio onde o sistema roda. Foi imaginado para ser textual, por isto, evite ip address
  */
  Domain                          string                  `json:"domain"`

  /*
  [opcional] Coloque apenas o número da porta, sem os ':'. Ex. :8080, fica apenas 8080
  */
  Port                            string                  `json:"port"`
}
func (el *ProxyDomain) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    SubDomain                       string                  `json:"subDomain"`
    Domain                          string                  `json:"domain"`
    Port                            string                  `json:"port"`
  }{
    ErrorHandleAsString:    runtime.FuncForPC( reflect.ValueOf( el.ErrorHandle ).Pointer() ).Name(),
    NotFoundHandleAsString: runtime.FuncForPC( reflect.ValueOf( el.NotFoundHandle ).Pointer() ).Name(),
    SubDomain:              el.SubDomain,
    Domain:                 el.Domain,
    Port:                   el.Port,
  })
}
func (el *ProxyDomain) UnmarshalJSON(data []byte) error {
  type tmpStt struct{
    ErrorHandleAsString             string                  `json:"errorHandleAsString"`
    NotFoundHandleAsString          string                  `json:"notFoundHandleAsString"`
    SubDomain                       string                  `json:"subDomain"`
    Domain                          string                  `json:"domain"`
    Port                            string                  `json:"port"`
  }
  var tmp = tmpStt{}
  err := json.Unmarshal( data, &tmp )
  if err != nil {
    return err
  }

  el.ErrorHandle                 = FuncMap[ tmp.ErrorHandleAsString ].( ProxyHandlerFunc )
  el.NotFoundHandle              = FuncMap[ tmp.NotFoundHandleAsString ].( ProxyHandlerFunc )
  el.SubDomain                   = tmp.SubDomain
  el.Domain                      = tmp.Domain
  el.Port                        = tmp.Port

  return nil
}
