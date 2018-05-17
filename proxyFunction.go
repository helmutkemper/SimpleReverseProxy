package marketPlaceProcy

import (
  "net/http"
  "time"
  "net/url"
  log "github.com/helmutkemper/seelog"
)

func ProxyFunc(w http.ResponseWriter, r *http.Request) {

  // Espera uma nova chamada para que a nova rota tenha efeito
  if len( ProxyNewRootConfig ) > 0 {
    ProxyRootConfig.Routes = ProxyNewRootConfig
    ProxyNewRootConfig = make([]ProxyRoute, 0)
  }

  var responseWriter = ProxyResponseWriter{
    ResponseWriter: w,
  }

  var request = &ProxyRequest{
    Request: r,
  }

  start := NetworkTime.Now()

  var handleName string

  /*
  fixme: query string habilitada por host para maior desempenho
  request.ExpRegMatches = make( map[string]string )
  queryString := make( map[string][]string )


  queryString, err = url.ParseQuery(r.URL.RawQuery)
  if err != nil {
    // há um erro na query string
    log.Infof( "The query string passed by the user does not appear to be in the correct format. Query String: %v Host: %v%v", r.URL.RawQuery, r.Host, r.URL.Path )
  }

  request.QueryString = queryString
  */

  // Trata todas as rotas

  var method = r.Method
  if method == "" {
    method = "ALL"
  }

  var data, match = ProxyRadix.Get( r.Host + "/" + method + r.URL.Path )
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
        for urlKey := range data.(ProxyRoute).ProxyServers{
          if data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopOk == false && data.(ProxyRoute).ProxyServers[ urlKey ].Enabled == true && data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopError == false {
            passNextRoute = true
            passEnabled = true
            data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopOk = true
            keyUrlToUse = urlKey
            break
          }
        }

        // A próxima rota não foi encontrada
        if passNextRoute == false {
          // Limpa todas as indicações de próxima rota
          for urlKey := range data.(ProxyRoute).ProxyServers{
            data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopOk = false
          }

          // Procura por uma rota habilitada e que não houve um erro na tentativa anterior
          for urlKey := range data.(ProxyRoute).ProxyServers {
            if data.(ProxyRoute).ProxyServers[ urlKey ].Enabled == true && data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopError == false {
              data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopOk = true
              passEnabled = true
              keyUrlToUse = urlKey
              break
            }
          }

          // Todas as rotas estão desabilitadas ou houveram erros na tentativa anterior
          if passEnabled == false {

            // Todas as rotas estão desabilitadas ou houveram erros na tentativa anterior
            log.Warnf( "All routes reported error on previous attempt or are disabled. Host: %v", r.Host )

            // Desabilita a indicação de erro na etapa anterior
            for urlKey := range data.(ProxyRoute).ProxyServers {
              data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopError = false
            }

            // Procura por uma rota habilitada mesmo que tenha tido erro na etapa anterior
            // Uma rota desabilitada teve vários erros consecutivos, por isto, foi desabilitada temporariamente
            for urlKey := range data.(ProxyRoute).ProxyServers {
              if data.(ProxyRoute).ProxyServers[ urlKey ].Enabled == true {
                data.(ProxyRoute).ProxyServers[ urlKey ].LastLoopOk = true
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
          for urlKey := range data.(ProxyRoute).ProxyServers{
            data.(ProxyRoute).ProxyServers[ urlKey ].Enabled = true
          }

          //aconteceu um erro grave, todas as rotas falharam com erros consecutivos e foram habilitadas a força para tentar de qualquer modo
          log.Warnf( "All %v domain routes are disabled by error and the system is trying all routes anyway.", r.Host )

          loopCounter += 1
          continue
        }

        externalServerUrl  = data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Url

        containerUrl, err := url.Parse(externalServerUrl)
        if err != nil {
          // Avisar que houve erro no parser
          log.Criticalf( "The route '%v - %v' of the domain '%v' is wrong. Error: %v", data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host, err.Error() )
          loopCounter += 1

          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorCounter += 1
          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter += 1
          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].LastLoopError = true

          if data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {

            // avisar que rota foi removida
            log.Criticalf( "The route '%v - %v' of the domain '%v' is wrong and has been disabled indefinitely until it is corrected by the admin.", data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host )

            data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Enabled = false
            data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Forever = true
            data.(ProxyRoute).ProxyServers[ keyUrlToUse ].DisabledSince = NetworkTime.Now()
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

            timeMeasure( start, handleName )
            return
          }

          continue
        }

        transport := &transport{
          RoundTripper: http.DefaultTransport,
          Error: nil,
        }
        proxy := NewSingleHostReverseProxy(containerUrl)
        proxy.Transport = transport
        proxy.ServeHTTP(w, r)

        if transport.Error != nil {
          // avisar que houve erro na leitura da rota
          log.Warnf( "The route '%v - %v' of the domain '%v' returned an error. Error: %v", data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host, transport.Error.Error() )
          loopCounter += 1

          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorCounter += 1
          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter += 1
          data.(ProxyRoute).ProxyServers[ keyUrlToUse ].LastLoopError = true

          if data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {
            // avisar que rota foi removida
            log.Warnf( "The route '%v - %v' of the domain '%v' returned many consecutive errors and was temporarily disabled.", data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host )

            data.(ProxyRoute).ProxyServers[ keyUrlToUse ].Enabled = false
            data.(ProxyRoute).ProxyServers[ keyUrlToUse ].DisabledSince = NetworkTime.Now()
          }

          // Houveram erros excessivos e o processo foi abortado
          if loopCounter >= ProxyRootConfig.MaxLoopTry {

            log.Criticalf( "The '%v' domain returned more %v consecutive errors and the error page was displayed to the user.", r.Host, ProxyRootConfig.MaxLoopTry )

            // Página de erro específica do domínio
            if data.(ProxyRoute).Domain.ErrorHandle != nil {
              data.(ProxyRoute).Domain.ErrorHandle(responseWriter, request)

              // Página de erro do sistema
            } else if ProxyRootConfig.ErrorHandle != nil {
              ProxyRootConfig.ErrorHandle(responseWriter, request)
            }

            timeMeasure( start, handleName )
            return
          }

          continue
        }

        // rodou sem erro

        data.(ProxyRoute).ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter = 0
        data.(ProxyRoute).ProxyServers[ keyUrlToUse ].UsedSuccessfully += 1
        data.(ProxyRoute).ProxyServers[ keyUrlToUse ].TotalTime += NetworkTime.Since( start ) * time.Nanosecond

        // LastLoopError evita um loop infinito em rotas com erro de resposta
        for keyUrl := range data.(ProxyRoute).ProxyServers{
          data.(ProxyRoute).ProxyServers[ keyUrl ].LastLoopError = false
        }

        timeMeasure( start, handleName )
        return
      }







    } else {

      if data.(ProxyRoute).Handle.Handle != nil {
        data.(ProxyRoute).Handle.Handle(responseWriter, request)
        //data.(ProxyRoute).Handle.TotalTime += NetworkTime.Since(start) * time.Nanosecond
        //data.(ProxyRoute).Handle.UsedSuccessfully += 1
        timeMeasure(start, handleName)

        return
      }
    }


  }



  // nenhum domínio bateu e está é uma página 404 genérica?
  if ProxyRootConfig.NotFoundHandle != nil {
    ProxyRootConfig.NotFoundHandle(responseWriter, request)
  }
  timeMeasure( start, handleName )
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
