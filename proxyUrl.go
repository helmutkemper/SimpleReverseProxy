package marketPlaceProcy

import "time"

type ProxyUrl struct {
  /*
  Url da rota para o proxy
  */
  Url                             string                  `json:"url"`

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
  Habilitada / Desabilitada temporariamente para esperar a rota voltar a responder
  */
  Enabled                         bool                    `json:"enabled"`

  /*
  Quando marcado true, a rotina que reabita ignora a rota e a mantém desabilitada
  */
  Forever                         bool                    `json:"forever"`

  /*
  Total de erros durante a execução da rota do proxy
  */
  ErrorCounter                    int64                   `json:"errorCounter"`

  /*
  Conta quantos erros seguidos houveram para poder decidir se desabilita a roda do proxy
  */
  ErrorConsecutiveCounter         int64                   `json:"errorConsecutiveCounter"`

  /*
  Arquiva o tempo desabilitado para poder reabilitar por time out
  */
  DisabledSince                   time.Time               `json:"-"`

  /*
  Usado pelo código para evitar que uma rota fique em loop infinito
  */
  LastLoopError                   bool                    `json:"lastLoopError"`

  /*
  Usado para indicar que a rota já foi usada
  A ideia é escolher a próximo rota livre em vez de ficar repetindo
  */
  LastLoopOk                      bool                    `json:"lastLoopOk"`
}
