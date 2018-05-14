package marketPlaceProcy

import (
  "time"
  "runtime"
  "reflect"
  "encoding/json"
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
  return json.Marshal(&ProxyHandle{
    Name: el.Name,
    TotalTime: el.TotalTime,
    UsedSuccessfully: el.UsedSuccessfully,
    HandleAsString: runtime.FuncForPC( reflect.ValueOf( el.Handle ).Pointer() ).Name(),
  })
}

