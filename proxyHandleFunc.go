package SimpleReverseProxy

type ProxyHandlerFunc func(ProxyResponseWriter, *ProxyRequest)
