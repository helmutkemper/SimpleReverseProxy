package marketPlaceProcy

import (
  "encoding/json"
)

type ProxyRoute struct {
  /*
  Nome para o log e outras funções, deve ser único e começar com letra ou '_'
  */
  Name                            string                  `json:"name"`

  /*
  Dados do domínio
  */
  Domain                          ProxyDomain             `json:"domain"`

  /*
  [opcional] Dados do caminho dentro do domínio
  */
  Path                            ProxyPath                `json:"path"`

  /*
  [opcional] Dado da aplicação local
  */
  Handle                          ProxyHandle             `json:"handle"`

  /*
  Habilita a funcionalidade do proxy, caso contrário, será chamada a função handle
  */
  ProxyEnable                     bool                    `json:"proxyEnable"`

  /*
  Lista de todas as URLs para os containers com a aplicação
  */
  ProxyServers                    []ProxyUrl              `json:"proxyServers"`
}
func (el *ProxyRoute) MarshalJSON() ([]byte, error) {
  return json.Marshal(&ProxyRoute{
    Name: el.Name,
    Domain: el.Domain,
    Path: el.Path,
    Handle: el.Handle,
    ProxyEnable: el.ProxyEnable,
    ProxyServers: el.ProxyServers,
  })
}
