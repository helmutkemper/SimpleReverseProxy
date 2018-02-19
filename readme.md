# Market Place Proxy

Exemplo de uso
```
package main

import (
  mkp "github.com/helmutkemper/marketPlaceProxy"
  "net/http"
)

func hello(w mkp.ProxyResponseWriter, r *mkp.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  w.Write( []byte( "controller: " ) )
  w.Write( []byte( r.ExpRegMatches[ "controller" ] ) )
  w.Write( []byte( "<br>" ) )

  w.Write( []byte( "module: " ) )
  w.Write( []byte( r.ExpRegMatches[ "module" ] ) )
  w.Write( []byte( "<br>" ) )

  w.Write( []byte( "site: " ) )
  w.Write( []byte( r.ExpRegMatches[ "site" ] ) )
  w.Write( []byte( "<br>" ) )
}

func main() {
  mkp.ProxyRootConfig = mkp.ProxyConfig{
    ListenAndServe: ":8888",
    Routes: []mkp.ProxyRoute{
      {
        Name: "blog",
        Domain: mkp.ProxyDomain{
          SubDomain: "blog",
          Domain: "localhost",
          Port: "8888",
        },
        ProxyEnable: true,
        ProxyServers: []mkp.ProxyUrl{
          {
            Name: "docker 1 - ok",
            Url: "http://localhost:2368",
          },
          {
            Name: "docker 2 - error",
            Url: "http://localhost:2367",
          },
          {
            Name: "docker 3 - error",
            Url: "http://localhost:2367",
          },
        },
      },
      {
        Name: "hello",
        Domain: mkp.ProxyDomain{
          NotFoundHandle: mkp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mkp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mkp.ProxyPath{
          ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mkp.ProxyHandle{
          Handle: hello,
        },
      },
      {
        Name: "addTest",
        Domain: mkp.ProxyDomain{
          NotFoundHandle: mkp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mkp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mkp.ProxyPath{
          Path : "/add",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mkp.ProxyHandle{
          Handle: mkp.ProxyRootConfig.RouteAdd,
        },
      },
      {
        Name: "removeTest",
        Domain: mkp.ProxyDomain{
          NotFoundHandle: mkp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mkp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mkp.ProxyPath{
          Path : "/remove",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mkp.ProxyHandle{
          Handle: mkp.ProxyRootConfig.RouteDelete,
        },
      },
      {
        Name: "panel",
        Domain: mkp.ProxyDomain{
          NotFoundHandle: mkp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mkp.ProxyRootConfig.ProxyError,
          SubDomain: "root",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mkp.ProxyPath{
          Path: "/statistics",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mkp.ProxyHandle{
          Handle: mkp.ProxyRootConfig.ProxyStatistics,
        },
      },
    },
  }
  mkp.ProxyRootConfig.Prepare()
  go mkp.ProxyRootConfig.VerifyDisabled()

  http.HandleFunc("/", mkp.ProxyFunc)
  http.ListenAndServe(mkp.ProxyRootConfig.ListenAndServe, nil)
}
```