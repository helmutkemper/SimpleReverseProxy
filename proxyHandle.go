package marketPlaceProcy

import "time"

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
}

