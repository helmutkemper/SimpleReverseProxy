package SimpleReverseProxy

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

  Index                           int                     `json:"index"`
}
func (el *ProxyRoute) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct{
    Name                            string                  `json:"name"`
    Domain                          ProxyDomain             `json:"domain"`
    Path                            ProxyPath                `json:"path"`
    Handle                          ProxyHandle             `json:"handle"`
    ProxyEnable                     bool                    `json:"proxyEnable"`
    ProxyServers                    []ProxyUrl              `json:"proxyServers"`
  }{
    Name:         el.Name,
    Domain:       el.Domain,
    Path:         el.Path,
    Handle:       el.Handle,
    ProxyEnable:  el.ProxyEnable,
    ProxyServers: el.ProxyServers,
  })
}