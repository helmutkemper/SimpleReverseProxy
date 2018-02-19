package marketPlaceProcy

import (
  "net/http"
  "regexp"
  "net/url"
  "time"
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

  // a ideia era ter mais controle sobre o relógio, mas, terminou ficando inacabada
  now := time.Now()

  start := time.Now()

  var handleName string

  //defer não funcionou direito
  //defer timeMeasure( start, handleName )

  request.ExpRegMatches = make( map[string]string )
  queryString := make( map[string][]string )

  // Trata o domínio e o separa
  // fixme: rever isto. isto é necessário? deixa o sistema mais lento
  matched, err := regexp.MatchString(ProxyRootConfig.DomainExpReg, r.Host)
  if err != nil {
    // há um erro grave na expreg do domínio
    log.Debugf( "The regular expression in charge of identifying the domain data has a serious error and the reverse proxy system can not continue. ExpReg: '/%v/' Error: %v", ProxyRootConfig.DomainExpReg, err.Error() )
    log.Criticalf( "The regular expression in charge of identifying the domain data has a serious error and the reverse proxy system can not continue. Error: %v", err.Error() )
    return
  }

  if matched == true {
    re := regexp.MustCompile(ProxyRootConfig.DomainExpReg)

    request.SubDomain = re.ReplaceAllString(r.Host,"${subDomain}")
    request.Domain = re.ReplaceAllString(r.Host, "${domain}")
    request.Port = re.ReplaceAllString(r.Host, "${port}")
  } else {
    // a equação de domínio não bateu
    log.Warnf( "Regular domain expression did not hit domain %v", r.Host )
    return
  }

  // trata a query string
  // fixme: isto é necessário aqui? deixa o sistema mais lento
  queryString, err = url.ParseQuery(r.URL.RawQuery)
  if err != nil {
    // há um erro na query string
    log.Infof( "The query string passed by the user does not appear to be in the correct format. Query String: %v Host: %v%v", r.URL.RawQuery, r.Host, r.URL.Path )
  }

  request.QueryString = queryString

  // Trata todas as rotas
  for keyRoute, route := range ProxyRootConfig.Routes {

    handleName = route.Name

    if route.Domain.SubDomain != "" {
      route.Domain.SubDomain += "."
    }

    if route.Domain.Port != "" {
      route.Domain.Port = ":" + route.Domain.Port
    }

    if r.Host != route.Domain.SubDomain + route.Domain.Domain + route.Domain.Port {
      continue
    }

    // O domínio foi encontrado
    if route.Path.ExpReg != "" && ( route.Path.Method == "" || route.Path.Method == r.Method ) {

      matched, err = regexp.MatchString(route.Path.ExpReg, r.URL.Path)
      if matched == true {
        re := regexp.MustCompile(route.Path.ExpReg)
        for k, v := range re.SubexpNames() {
          if k == 0 || v == "" {
            continue
          }

          request.ExpRegMatches[v] = re.ReplaceAllString(r.URL.Path, `${`+v+`}`)
        }

        if ProxyRootConfig.Routes[ keyRoute ].Handle.Handle != nil {
          ProxyRootConfig.Routes[ keyRoute ].Handle.Handle(responseWriter, request)
          ProxyRootConfig.Routes[ keyRoute ].Handle.TotalTime += time.Since( start ) * time.Nanosecond
          ProxyRootConfig.Routes[ keyRoute ].Handle.UsedSuccessfully += 1
          timeMeasure( start, handleName )
          return
        }

      } else {
        continue
      }

    } else if ( route.Path.Method == "" || route.Path.Method == r.Method ) && ( route.Path.Path == r.URL.Path || route.Path.Path == "" ) {

      if ProxyRootConfig.Routes[ keyRoute ].Handle.Handle != nil {
        ProxyRootConfig.Routes[ keyRoute ].Handle.Handle(responseWriter, request)
        ProxyRootConfig.Routes[ keyRoute ].Handle.TotalTime += time.Since( start ) * time.Nanosecond
        ProxyRootConfig.Routes[ keyRoute ].Handle.UsedSuccessfully += 1
        timeMeasure( start, handleName )
        return
      }

      // O domínio foi encontrado, porém, o path dentro do domínio não
    } else {
      continue
    }

    if route.ProxyEnable == false {
      ProxyRootConfig.Routes[ keyRoute ].Handle.Handle(responseWriter, request)
      ProxyRootConfig.Routes[ keyRoute ].Handle.TotalTime += time.Since( start ) * time.Nanosecond
      ProxyRootConfig.Routes[ keyRoute ].Handle.UsedSuccessfully += 1
      timeMeasure( start, handleName )
      return
    }

    loopCounter := 0

    for {
      passEnabled := false
      keyUrlToUse := 0
      externalServerUrl := ""
      passNextRoute := false
      // Procura pela próxima rota para uso que esteja habilitada
      for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers{
        if ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopOk == false && ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].Enabled == true && ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopError == false {
          passNextRoute = true
          passEnabled = true
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopOk = true
          keyUrlToUse = urlKey
          break
        }
      }

      // A próxima rota não foi encontrada
      if passNextRoute == false {
        // Limpa todas as indicações de próxima rota
        for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers{
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopOk = false
        }

        // Procura por uma rota habilitada e que não houve um erro na tentativa anterior
        for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers {
          if ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].Enabled == true && ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopError == false {
            ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopOk = true
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
          for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers {
            ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopError = false
          }

          // Procura por uma rota habilitada mesmo que tenha tido erro na etapa anterior
          // Uma rota desabilitada teve vários erros consecutivos, por isto, foi desabilitada temporariamente
          for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers {
            if ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].Enabled == true {
              ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].LastLoopOk = true
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
        for urlKey := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers{
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ urlKey ].Enabled = true
        }

        //aconteceu um erro grave, todas as rotas falharam com erros consecutivos e foram habilitadas a força para tentar de qualquer modo
        log.Warnf( "All %v domain routes are disabled by error and the system is trying all routes anyway.", r.Host )

        loopCounter += 1
        continue
      }

      externalServerUrl  = ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Url

      containerUrl, err := url.Parse(externalServerUrl)
      if err != nil {
        // Avisar que houve erro no parser
        log.Criticalf( "The route '%v - %v' of the domain '%v' is wrong. Error: %v", ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host, err.Error() )
        loopCounter += 1

        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorCounter += 1
        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter += 1
        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].LastLoopError = true

        if ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {

          // avisar que rota foi removida
          log.Criticalf( "The route '%v - %v' of the domain '%v' is wrong and has been disabled indefinitely until it is corrected by the admin.", ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host )

          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Enabled = false
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Forever = true
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].DisabledSince = now
        }

        // Houveram erros excessivos e o processo foi abortado
        if loopCounter >= ProxyRootConfig.MaxLoopTry {

          // Página de erro específica do domínio
          if ProxyRootConfig.Routes[ keyRoute ].Domain.ErrorHandle != nil {
            ProxyRootConfig.Routes[ keyRoute ].Domain.ErrorHandle(responseWriter, request)

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
        log.Warnf( "The route '%v - %v' of the domain '%v' returned an error. Error: %v", ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host, transport.Error.Error() )
        loopCounter += 1

        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorCounter += 1
        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter += 1
        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].LastLoopError = true

        if ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter >= ProxyRootConfig.ConsecutiveErrorsToDisable {
          // avisar que rota foi removida
          log.Warnf( "The route '%v - %v' of the domain '%v' returned many consecutive errors and was temporarily disabled.", ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Name, externalServerUrl, r.Host )

          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].Enabled = false
          ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].DisabledSince = now
        }

        // Houveram erros excessivos e o processo foi abortado
        if loopCounter >= ProxyRootConfig.MaxLoopTry {

          log.Criticalf( "The '%v' domain returned more %v consecutive errors and the error page was displayed to the user.", r.Host, ProxyRootConfig.MaxLoopTry )

          // Página de erro específica do domínio
          if ProxyRootConfig.Routes[ keyRoute ].Domain.ErrorHandle != nil {
            ProxyRootConfig.Routes[ keyRoute ].Domain.ErrorHandle(responseWriter, request)

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

      ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].ErrorConsecutiveCounter = 0
      ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].UsedSuccessfully += 1
      ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrlToUse ].TotalTime += time.Since( start ) * time.Nanosecond

      // LastLoopError evita um loop infinito em rotas com erro de resposta
      for keyUrl := range ProxyRootConfig.Routes[ keyRoute ].ProxyServers{
        ProxyRootConfig.Routes[ keyRoute ].ProxyServers[ keyUrl ].LastLoopError = false
      }

      timeMeasure( start, handleName )
      return
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
    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: sessionName, Value: sessionId(), Expires: expiration}
    http.SetCookie(w, &cookie)
  }

  cookie, _ = r.Cookie(sessionName)
  fmt.Printf("cookie: %q\n", cookie)*/
}
