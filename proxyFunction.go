package SimpleReverseProxy

import (
	log "github.com/helmutkemper/seelog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func ProxyFunc(w http.ResponseWriter, r *http.Request) {

	// Espera uma nova chamada para que a nova rota tenha efeito
	if ProxyNewRootConfig.Len() > 0 {
		ProxyRootConfig.Routes = ProxyNewRootConfig.Get()
		ProxyNewRootConfig.Clear()
	}

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

			loopCounter := 0

			for {
				passEnabled := false
				keyUrlToUse := 0
				externalServerUrl := ""
				passNextRoute := false
				// Procura pela próxima rota para uso que esteja habilitada
				for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
					tmpFomData := data.(ProxyRoute).ProxyServers.GetKey(urlKey)
					if tmpFomData.LastLoopOk == false && tmpFomData.Enabled == true && tmpFomData.LastLoopError == false {
						passNextRoute = true
						passEnabled = true
						tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(urlKey)
						tmpFromRootConfig.LastLoopOk = true
						ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFromRootConfig)
						keyUrlToUse = urlKey
						break
					}
				}

				// A próxima rota não foi encontrada
				if passNextRoute == false {
					// Limpa todas as indicações de próxima rota
					for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
						tmpFomData := data.(ProxyRoute).ProxyServers.GetKey(urlKey)
						tmpFomData.LastLoopOk = false
						ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFomData)
					}

					// Procura por uma rota habilitada e que não houve um erro na tentativa anterior
					for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
						tmpFomData := data.(ProxyRoute).ProxyServers.GetKey(urlKey)
						if tmpFomData.Enabled == true && tmpFomData.LastLoopError == false {
							tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(urlKey)
							tmpFromRootConfig.LastLoopOk = true
							ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFromRootConfig)
							passEnabled = true
							keyUrlToUse = urlKey
							break
						}
					}

					// Todas as rotas estão desabilitadas ou houveram erros na tentativa anterior
					if passEnabled == false {

						// Todas as rotas estão desabilitadas ou houveram erros na tentativa anterior
						log.Warnf("All routes reported error on previous attempt or are disabled. Host: %v", r.Host)

						// Desabilita a indicação de erro na etapa anterior
						for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
							tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(urlKey)
							tmpFromRootConfig.LastLoopError = false
							ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFromRootConfig)
						}

						// Procura por uma rota habilitada mesmo que tenha tido erro na etapa anterior
						// Uma rota desabilitada teve vários erros consecutivos, por isto, foi desabilitada temporariamente
						for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
							tmpFomData := data.(ProxyRoute).ProxyServers.GetKey(urlKey)
							if tmpFomData.Enabled == true {
								tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(urlKey)
								tmpFromRootConfig.LastLoopOk = true
								ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFromRootConfig)
								passEnabled = true
								keyUrlToUse = urlKey
								break
							}
						}
					}
				}

				// Todas as rotas estão desabilitada por erro
				// Habilita todas as rotas e tenta novamente
				if passEnabled == false {
					for urlKey := range data.(ProxyRoute).ProxyServers.Get() {
						tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(urlKey)
						tmpFromRootConfig.Enabled = true
						ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(urlKey, tmpFromRootConfig)
					}

					//aconteceu um erro grave, todas as rotas falharam com erros consecutivos e foram habilitadas a força para tentar de qualquer modo
					log.Warnf("All %v domain routes are disabled by error and the system is trying all routes anyway.", r.Host)

					loopCounter += 1
					continue
				}

				externalServerUrl = data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Url

				containerUrl, err := url.Parse(externalServerUrl)
				if err != nil {
					// Avisar que houve erro no parser
					log.Criticalf("The route '%v - %v' of the domain '%v' is wrong. Error: %v", data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Name, externalServerUrl, r.Host, err.Error())
					loopCounter += 1

					tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrlToUse)

					tmpFromRootConfig.ErrorCounter += 1
					tmpFromRootConfig.ErrorConsecutiveCounter += 1
					tmpFromRootConfig.LastLoopError = true
					ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrlToUse, tmpFromRootConfig)

					if data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {

						// avisar que rota foi removida
						log.Criticalf("The route '%v - %v' of the domain '%v' is wrong and has been disabled indefinitely until it is corrected by the admin.", data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Name, externalServerUrl, r.Host)

						tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrlToUse)

						tmpFromRootConfig.Enabled = false
						tmpFromRootConfig.Forever = true
						tmpFromRootConfig.DisabledSince = NetworkTime.Now()

						ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrlToUse, tmpFromRootConfig)
					}

					// Houveram erros excessivos e o processo foi abortado
					if loopCounter >= ProxyRootConfig.MaxLoopTry {

						// Página de erro específica do domínio
						if data.(ProxyRoute).Domain.ErrorHandle != nil {
							data.(ProxyRoute).Domain.ErrorHandle(responseWriter, request)

							// Página de erro do sistema
						} else if ProxyRootConfig.ErrorHandle != nil {
							ProxyRootConfig.ErrorHandle(responseWriter, request)
						}

						timeMeasure(start, handleName)
						return
					}

					continue
				}

				transport := &transport{
					RoundTripper: http.DefaultTransport,
					Error:        nil,
				}
				proxy := NewSingleHostReverseProxy(containerUrl)
				proxy.Transport = transport
				proxy.ServeHTTP(w, r)

				if transport.Error != nil {
					// avisar que houve erro na leitura da rota
					log.Warnf("The route '%v - %v' of the domain '%v' returned an error. Error: %v", data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Name, externalServerUrl, r.Host, transport.Error.Error())
					loopCounter += 1

					tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrlToUse)

					tmpFromRootConfig.ErrorCounter += 1
					tmpFromRootConfig.ErrorConsecutiveCounter += 1
					tmpFromRootConfig.LastLoopError = true

					ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrlToUse, tmpFromRootConfig)

					if data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {
						// avisar que rota foi removida
						log.Warnf("The route '%v - %v' of the domain '%v' returned many consecutive errors and was temporarily disabled.", data.(ProxyRoute).ProxyServers.GetKey(keyUrlToUse).Name, externalServerUrl, r.Host)

						tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrlToUse)

						tmpFromRootConfig.Enabled = false
						tmpFromRootConfig.DisabledSince = NetworkTime.Now()

						ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrlToUse, tmpFromRootConfig)
					}

					// Houveram erros excessivos e o processo foi abortado
					if loopCounter >= ProxyRootConfig.MaxLoopTry {

						log.Criticalf("The '%v' domain returned more %v consecutive errors and the error page was displayed to the user.", r.Host, ProxyRootConfig.MaxLoopTry)

						// Página de erro específica do domínio
						if data.(ProxyRoute).Domain.ErrorHandle != nil {
							data.(ProxyRoute).Domain.ErrorHandle(responseWriter, request)

							// Página de erro do sistema
						} else if ProxyRootConfig.ErrorHandle != nil {
							ProxyRootConfig.ErrorHandle(responseWriter, request)
						}

						timeMeasure(start, handleName)
						return
					}

					continue
				}

				// rodou sem erro

				tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrlToUse)

				tmpFromRootConfig.ErrorConsecutiveCounter = 0
				tmpFromRootConfig.UsedSuccessfully += 1
				tmpFromRootConfig.TotalTime += NetworkTime.Since(start) * time.Nanosecond

				ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrlToUse, tmpFromRootConfig)

				// LastLoopError evita um loop infinito em rotas com erro de resposta
				for keyUrl := range data.(ProxyRoute).ProxyServers.Get() {
					tmpFromRootConfig := ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.GetKey(keyUrl)
					tmpFromRootConfig.LastLoopError = false
					ProxyRootConfig.Routes[data.(ProxyRoute).Index].ProxyServers.SetKey(keyUrl, tmpFromRootConfig)
				}

				timeMeasure(start, handleName)
				return
			}

		} else {

			if data.(ProxyRoute).Handle.Handle != nil {

				if data.(ProxyRoute).Path.ExpReg != "" {
					matched, err := regexp.MatchString(data.(ProxyRoute).Path.ExpReg, r.URL.Path)
					if err != nil {
						// fixme: colocar um erro aqui
					}
					if matched == true {
						re := regexp.MustCompile(data.(ProxyRoute).Path.ExpReg)
						for k, v := range re.SubexpNames() {
							if k == 0 || v == "" {
								continue
							}

							request.ExpRegMatches[v] = re.ReplaceAllString(r.URL.Path, `${`+v+`}`)
						}
					} else {
						// fixme: page not found
						return
					}
				}

				data.(ProxyRoute).Handle.Handle(responseWriter, request)
				ProxyRootConfig.Routes[data.(ProxyRoute).Index].Handle.TotalTime += NetworkTime.Since(start) * time.Nanosecond
				ProxyRootConfig.Routes[data.(ProxyRoute).Index].Handle.UsedSuccessfully += 1
				timeMeasure(start, handleName)

				return
			}
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
