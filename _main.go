package main

import (
  log "github.com/helmutkemper/seelog"
  mkpConf "github.com/helmutkemper/marketPlaceProxy/config"
  "net/http"
  "net/url"
  "regexp"
  "time"
  "net"
  "sync"
  "strings"
  "context"
  "io"
  "encoding/json"
  "bytes"
  "strconv"
  "io/ioutil"
  "fmt"
  "os"
  "errors"
)


func hello(w ProxyResponseWriter, r *ProxyRequest) {
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
  ProxyRootConfig = mkpConf.ProxyConfig{
    ListenAndServe: ":8888",
    Routes: []ProxyRoute{
      {
        Name: "blog",
        Domain: ProxyDomain{
          SubDomain: "blog",
          Domain: "localhost",
          Port: "8888",
        },
        ProxyEnable: true,
        ProxyServers: []ProxyUrl{
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
        Domain: ProxyDomain{
          NotFoundHandle: ProxyRootConfig.ProxyNotFound,
          ErrorHandle: ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: ProxyPath{
          ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: ProxyHandle{
          Handle: hello,
        },
      },
      {
        Name: "addTest",
        Domain: ProxyDomain{
          NotFoundHandle: ProxyRootConfig.ProxyNotFound,
          ErrorHandle: ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: ProxyPath{
          Path : "/add",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: ProxyHandle{
          Handle: ProxyRootConfig.RouteAdd,
        },
      },
      {
        Name: "removeTest",
        Domain: ProxyDomain{
          NotFoundHandle: ProxyRootConfig.ProxyNotFound,
          ErrorHandle: ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: ProxyPath{
          Path : "/remove",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: ProxyHandle{
          Handle: ProxyRootConfig.RouteDelete,
        },
      },
      {
        Name: "panel",
        Domain: ProxyDomain{
          NotFoundHandle: ProxyRootConfig.ProxyNotFound,
          ErrorHandle: ProxyRootConfig.ProxyError,
          SubDomain: "root",
          Domain: "localhost",
          Port: "8888",
        },
        Path: ProxyPath{
          Path: "/statistics",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: ProxyHandle{
          Handle: ProxyRootConfig.ProxyStatistics,
        },
      },
    },
  }
  ProxyRootConfig.Prepare()
  go ProxyRootConfig.VerifyDisabled()

  http.HandleFunc("/", ProxyFunc)
  http.ListenAndServe(ProxyRootConfig.ListenAndServe, nil)
}







