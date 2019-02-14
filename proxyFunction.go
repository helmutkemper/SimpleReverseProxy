package SimpleReverseProxy

import (
	log "github.com/helmutkemper/seelog"
	"net/http"
	"net/url"
	"strings"
)

func ProxyFunc(w http.ResponseWriter, r *http.Request) {

	// Espera uma nova chamada para que a nova rota tenha efeito
	/*if ProxyNewRootConfig.Len() > 0 {
		ProxyRootConfig.Routes = ProxyNewRootConfig.Get()
		ProxyNewRootConfig.Clear()
	}*/

	var responseWriter = ProxyResponseWriter{
		ResponseWriter: w,
	}

	var request = &ProxyRequest{
		R: r,
		//Request: r,
	}

	start := NetworkTime.Now()

	var handleName string

	request.ExpRegMatches = make(map[string]string)
	queryString := make(map[string][]string)

	var err error
	queryString, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		// há um erro na query string
		log.Infof("The query string passed by the user does not appear to be in the correct format. Query String: %v Host: %v%v", r.URL.RawQuery, r.Host, r.URL.Path)
	}

	request.QueryString = queryString

	// Trata todas as rotas
	var method = r.Method
	if method == "" {
		method = "ALL"
	}

	var data, match = ProxyRadix.Get(r.Host)
	if match == true {
		match = data.(ProxyRoute).ProxyEnable
	} else {
		data, match = ProxyRadix.Get("")
		if match == true {
			match = data.(ProxyRoute).ProxyEnable
		}
	}

	if match == false {
		var pathTmp = strings.Split(r.URL.Path, "/")
		var path = "/" + pathTmp[1]
		data, match = ProxyRadix.Get(r.Host + "/" + method + path)

		if match == false {
			data, match = ProxyRadix.Get("/" + method + path)
		}
	}

	if match == true {
		// a rota foi encontrada

		if data.(ProxyRoute).ProxyEnable == true {

			//fixme: string to iota
			if data.(ProxyRoute).LoadBalanceMode == LOAD_BALANCE_ROUND_ROBIN || true {

				for {
					keyUrlToUse := 0
					externalServerUrl := data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Url
					containerUrl, err := url.Parse(externalServerUrl)
					if err != nil {
						log.Critical(err.Error())
					}
					transport := &transport{
						RoundTripper: http.DefaultTransport,
						Error:        nil,
					}
					proxy := NewSingleHostReverseProxy(containerUrl)
					proxy.Transport = transport
					proxy.ServeHTTP(w, r)
					break
				}
			}

		} else {

		}
	}

	// nenhum domínio bateu e está é uma página 404 genérica?
	if ProxyRootConfig.NotFoundHandle != nil {
		ProxyRootConfig.NotFoundHandle(responseWriter, request)
	}
	timeMeasure(start, handleName)
	return

	/*cookie, _ := r.Cookie(sessionName)
	  if cookie == nil {
	    expiration := NetworkTime.Now().Add(365 * 24 * time.Hour)
	    cookie := http.Cookie{Name: sessionName, Value: sessionId(), Expires: expiration}
	    http.SetCookie(w, &cookie)
	  }

	  cookie, _ = r.Cookie(sessionName)
	  fmt.Printf("cookie: %q\n", cookie)*/
}
