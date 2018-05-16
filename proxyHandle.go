package marketPlaceProcy

import (
  "time"
  "encoding/json"
  "runtime"
  "reflect"
)

type ProxyHandle struct {
  /*
  Nome da rota para manter organizado
  */
  Name                            string                  `json:"name"`

  /*
  Tempo total de execução da rota.
  A soma de todos os tempos de resposta
  */
  TotalTime                       time.Duration           `json:"totalTime"`

  /*
  Quantidades de usos sem erro
  */
  UsedSuccessfully                int64                   `json:"usedSuccessfully"`

  /*
  Função a ser servida
  */
  Handle                          ProxyHandlerFunc        `json:"-"`
  
  
  HandleAsString                  string                  `json:"handle"`
}
func (el *ProxyHandle) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct{
    Name                            string                  `json:"name"`
    TotalTime                       time.Duration           `json:"totalTime"`
    UsedSuccessfully                int64                   `json:"usedSuccessfully"`
    HandleAsString                  string                  `json:"handleAsString"`
  }{
    Name:             el.Name,
    TotalTime:        el.TotalTime,
    UsedSuccessfully: el.UsedSuccessfully,
    HandleAsString:    runtime.FuncForPC( reflect.ValueOf( el.Handle ).Pointer() ).Name(),
  })
}
func (el *ProxyHandle) UnmarshalJSON(data []byte) error {
  type tmpStt struct{
    Name                            string                  `json:"name"`
    TotalTime                       time.Duration           `json:"totalTime"`
    UsedSuccessfully                int64                   `json:"usedSuccessfully"`
    HandleAsString                  string                  `json:"handleAsString"`
  }
  var tmp = tmpStt{}
  err := json.Unmarshal( data, &tmp )
  if err != nil {
    return err
  }

  if tmp.HandleAsString != "" {
    el.Handle = FuncMap[ tmp.HandleAsString ].(ProxyHandlerFunc)
  }

  el.Name               = tmp.Name
  el.TotalTime          = tmp.TotalTime
  el.UsedSuccessfully   = tmp.UsedSuccessfully

  return nil
}
