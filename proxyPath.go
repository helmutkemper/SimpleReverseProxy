package marketPlaceProcy

type ProxyPath struct {
  /*
  [opcional] Quando omitido, juntamente com ExpReg, faz com que todo o subdomínio seja usado para a rota
  */
  Path                            string                  `json:"path"`

  /*
  [opcional] Método da chamada GET/POST/DELETE...
  */
  Method                          string                  `json:"method"`

  /*
  true faz com que o path seja checado por expressão regular
  */
  ExpReg                          string                  `json:"expReg"`
}

