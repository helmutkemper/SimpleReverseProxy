package SimpleReverseProxy

import "net/http"

type ProxyResponseWriter struct {
  http.ResponseWriter
}