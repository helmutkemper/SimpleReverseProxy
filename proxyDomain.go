package marketPlaceProcy

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
