package SimpleReverseProxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helmutkemper/go-radix"
	log "github.com/helmutkemper/seelog"
	"github.com/helmutkemper/util"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const kLogFormat = "[%Level::%Date %Time] %Msg%n"

type ProxyConfig struct {
	readyToJSon bool `json:"-"`

	/*
	   Configuração do seelog
	   @see https://github.com/cihub/seelog
	*/
	SeeLogConfig string `json:"seeLogConfig"`

	/*
	   Expressão regular que identifica o domínio do site
	*/
	DomainExpReg string `json:"domainExpReg"`

	/*
	   Função de erro genérica, caso a função do domínio não seja definida
	*/
	ErrorHandle         ProxyHandlerFunc `json:"-"`
	ErrorHandleAsString string           `json:"ErrorHandle"`

	/*
	   Função de page not found genérica, caso a função do domínio não seja definida
	*/
	NotFoundHandle         ProxyHandlerFunc `json:"-"`
	NotFoundHandleAsString string           `json:"NotFoundHandle"`

	/*
	   Tamanho de caracteres do token de segurança
	*/
	UniqueIdLength int `json:"uniqueIdLength"`

	/*
	   URL do servidor principal
	*/
	ListenAndServe string `json:"listenAndServe"`

	/*
	   Quantidade máxima de loop quando todas as rotas do proxy falham
	*/
	MaxLoopTry int

	/*
	   Quantidades de erros consecutivos para desabilitar uma rota do proxy.
	   A ideia é que uma rota do proxy possa está dando erro temporário, assim, o código desabilita a rota por um tempo e
	   depois habilita de novo para testar se a mesma voltou.
	   Caso haja apenas uma instabilidade, a rota continua.
	*/
	ConsecutiveErrorsToDisable int64

	/*
	   Tempo para manter uma rota do proxy desabilitada antes de testar novamente
	*/
	TimeToKeepDisabled time.Duration

	/*
	   Há uma função em loop infinito e a cada x período de tempo, ela verifica se alguma rota está desabilitada e reabilita
	   caso o tempo de espera tenha sido excedido
	*/
	TimeToVerifyDisabled time.Duration

	/*
	   Rotas do servidor proxy
	*/
	Routes ProxyRouteAStt
}

/*
Esta função adiciona novas rotas ao proxy
{
    "name": "news",
    "domain": {
      "subDomain": "news",
      "domain": "localhost",
      "port": "8888"
    },
    "proxyEnable": true,
    "proxyServers": [
    {
          "name": "docker 1 - ok",
          "url": "http://localhost:2368"
    },
    {
          "name": "docker 2 - ok",
          "url": "http://localhost:2368"
    },
    {
      "name": "docker 3 - ok",
          "url": "http://localhost:2368"
    }
  ]
}
*/

// Esta função coloca a rota nova em 'ProxyNewRootConfig' e espera uma nova chamada em uma rota qualquer para que a
// nova rota tenha efeito. Isso é transparente para o usuário final, mas, a rota não pode entrar em vigor durante o
// processamento da rota anterior, ou o sistema trava, devido a mudança dos 'ponteiros'
func (el *ProxyConfig) RouteAdd(w ProxyResponseWriter, r *ProxyRequest) {
	var newRoute ProxyRoute
	var output = JSonOutStt{}

	if ProxyNewRootConfig.Len() != 0 {
		ProxyRootConfig.Routes.Set(ProxyNewRootConfig.Get())
	}

	err := json.NewDecoder(r.R.Body).Decode(&newRoute)

	if err != nil {
		output.ToOutput(0, err, []int{}, w)
		return
	}

	if newRoute.ProxyEnable == false {
		output.ToOutput(0, errors.New("this function only adds new routes that can be used in conjunction with the reverse proxy"), []int{}, w)
		return
	}

	if newRoute.ProxyServers.Len() == 0 {
		output.ToOutput(0, errors.New("this function must receive at least one route that can be used in conjunction with the reverse proxy"), []int{}, w)
		return
	}

	for _, route := range newRoute.ProxyServers.Get() {
		if route.Name == "" {
			output.ToOutput(0, errors.New("every route must have a name assigned to it"), []int{}, w)
			return
		}

		_, err := url.Parse(route.Url)
		if err != nil {
			output.ToOutput(0, errors.New("the route of name '"+route.Name+"' presented the following error: "+err.Error()), []int{}, w)
			return
		}
	}

	// Habilita todas as rotas, pois, o padrão do go é false
	for urlKey := range newRoute.ProxyServers.Get() {
		tmp := newRoute.ProxyServers.GetKey(urlKey)
		tmp.Enabled = true
		newRoute.ProxyServers.SetKey(urlKey, tmp)
	}

	ProxyNewRootConfig.Set(ProxyRootConfig.Routes.Get())
	ProxyNewRootConfig.Append(newRoute)

	el.RoutePrepare()

	output.ToOutput(ProxyNewRootConfig.Len(), nil, ProxyNewRootConfig, w)
}

func (el *ProxyConfig) AddRouteToProxyStt(newRoute ProxyRoute) error {
	var err error

	if ProxyNewRootConfig.Len() != 0 {
		ProxyRootConfig.Routes.Set(ProxyNewRootConfig.Get())
	}

	if newRoute.ProxyEnable == false {
		return errors.New("this function only adds new routes that can be used in conjunction with the reverse proxy")
	}

	if newRoute.ProxyServers.Len() == 0 {
		return errors.New("this function must receive at least one route that can be used in conjunction with the reverse proxy")
	}

	for _, route := range newRoute.ProxyServers.Get() {
		if route.Name == "" {
			return errors.New("every route must have a name assigned to it")
		}

		_, err = url.Parse(route.Url)
		if err != nil {
			return errors.New("the route of name '" + route.Name + "' presented the following error: " + err.Error())
		}
	}

	// Habilita todas as rotas, pois, o padrão do go é false
	for urlKey := range newRoute.ProxyServers.Get() {
		tmp := newRoute.ProxyServers.GetKey(urlKey)
		tmp.Enabled = true
		newRoute.ProxyServers.SetKey(urlKey, tmp)
	}

	// O index é usado como ponteiro para algumas funções e contadores
	newRoute.Index = ProxyRootConfig.Routes.Len()

	ProxyNewRootConfig.Set(ProxyRootConfig.Routes.Get())
	ProxyNewRootConfig.Append(newRoute)

	el.RoutePrepare()

	return nil
}

func (el *ProxyConfig) AddRouteFromFuncStt(newRoute ProxyRoute) error {
	var err error

	if ProxyNewRootConfig.Len() != 0 {
		ProxyRootConfig.Routes.Set(ProxyNewRootConfig.Get())
	}

	if newRoute.ProxyEnable == true {
		return errors.New("this function only adds new routes that can't be used in conjunction with the reverse proxy")
	}

	if newRoute.ProxyServers.Len() != 0 { //fixme: colocar um erro aqui
		return errors.New("fixme: colocar um erro aqui")
	}

	for _, route := range newRoute.ProxyServers.Get() {
		if route.Name == "" {
			return errors.New("every route must have a name assigned to it")
		}

		_, err = url.Parse(route.Url)
		if err != nil {
			return errors.New("the route of name '" + route.Name + "' presented the following error: " + err.Error())
		}
	}

	// Habilita todas as rotas, pois, o padrão do go é false
	for urlKey := range newRoute.ProxyServers.Get() {
		tmp := newRoute.ProxyServers.GetKey(urlKey)
		tmp.Enabled = true
		newRoute.ProxyServers.SetKey(urlKey, tmp)
	}

	// O index é usado como ponteiro para algumas funções e contadores
	newRoute.Index = ProxyRootConfig.Routes.Len()

	ProxyNewRootConfig.Set(ProxyRootConfig.Routes.Get())
	ProxyNewRootConfig.Append(newRoute)

	el.RoutePrepare()

	return nil
}

/*
Esta função elimina rotas do proxy
{
    "name": "name_of_route"
}
*/
func (el *ProxyConfig) RoutePrepare() {

	ProxyRadix = radix.New()

	for _, newRoute := range ProxyNewRootConfig.Get() {

		var separatorHost = ""
		var separatorPath = ""
		if !strings.HasSuffix(newRoute.Domain.Host, "/") {
			separatorHost = "/"
		}

		if !strings.HasPrefix(newRoute.Path.Path, "/") {
			separatorPath = "/"
		}

		if newRoute.ProxyEnable == true {

			ProxyRadix.Insert(newRoute.Domain.Host, newRoute)

		} else {

			if newRoute.Path.Method == "" {
				var list = []string{"GET", "POST", "DELETE", "PUT", "HEAD", "PATCH", "OPTIONS"}
				for _, v := range list {
					ProxyRadix.Insert(newRoute.Domain.Host+separatorHost+v+separatorPath+newRoute.Path.Path, newRoute)
				}
			} else {
				ProxyRadix.Insert(newRoute.Domain.Host+separatorHost+newRoute.Path.Method+separatorPath+newRoute.Path.Path, newRoute)
			}

		}
	}
}

func (el *ProxyConfig) RouteDelete(w ProxyResponseWriter, r *ProxyRequest) {
	// Esta função coloca a rota nova em 'ProxyNewRootConfig' e espera uma nova chamada em uma rota qualquer para que a
	// nova rota tenha efeito. Isso é transparente para o usuário final, mas, a rota não pode entrar em vigor durante o
	// processamento da rota anterior, ou o sistema trava, devido a mudança dos 'ponteiros'

	var newRoute ProxyRoute
	var output = JSonOutStt{}

	err := json.NewDecoder(r.R.Body).Decode(&newRoute)

	if err != nil {
		output.ToOutput(0, err, []int{}, w)
		return
	}

	var i int
	var l = ProxyRootConfig.Routes.Len()
	var nameFound = false
	for i = 0; i != l; i += 1 {
		if ProxyRootConfig.Routes.GetKey(i).Name == newRoute.Name {
			nameFound = true
			break
		}
	}

	if nameFound == true {
		if i == 0 {
			ProxyNewRootConfig.Set(ProxyRootConfig.Routes.Get()[1:])
		} else if i == ProxyRootConfig.Routes.Len()-1 {
			ProxyNewRootConfig.Set(ProxyRootConfig.Routes.Get()[:ProxyRootConfig.Routes.Len()-1])
		} else {
			dataFromProxyRootConfig := ProxyRootConfig.Routes.Get()
			ProxyNewRootConfig.Set(append(dataFromProxyRootConfig[0:i], dataFromProxyRootConfig[i+1:]...))
		}
	}

	if ProxyRootConfig.Routes.GetKey(i).ProxyEnable == false {
		output.ToOutput(0, errors.New("this function can only remove the routes used with the reverse proxy, not being able to remove other types of routes"), []int{}, w)
		return
	}

	b, e := json.Marshal(ProxyNewRootConfig)
	if e != nil {
		w.Write([]byte(e.Error()))
		return
	}
	w.Write(b)

	el.RoutePrepare()

	output.ToOutput(ProxyNewRootConfig.Len(), nil, ProxyNewRootConfig, w)
}

type ConfigStart struct {
	Log   Log
	Proxy ProxyEventConfig
}

type ProxyEventConfig struct {
}

type LogConfig struct {
	Path    string
	MaxRoll int
	MaxSize int
	Format  string
	Console Boolean
}

func (el *LogConfig) GetMaxRoll() string {
	return strconv.FormatInt(int64(el.MaxRoll), 10)
}

func (el *LogConfig) GetMaxSize() string {
	return strconv.FormatInt(int64(el.MaxSize), 10)
}

type Log struct {
	MinLevel LogLevel
	MaxLevel LogLevel

	// For pervasive information on states of all elementary constructs. Use 'Trace' for in-depth debugging to
	// find problem parts of a function, to check values of temporary variables, etc.
	Trace LogConfig

	// For general information on the application's work. Use 'Info' level in your code so that you could
	// leave it 'enabled' even in production. So it is a 'production log level'.
	Info LogConfig

	// For detailed system behavior reports and diagnostic messages to help to locate problems during
	// development.
	Debug LogConfig

	// For indicating small errors, strange situations, failures that are automatically handled in a safe
	// manner.
	Warn LogConfig

	// For severe failures that affects application's workflow, not fatal, however (without forcing app
	// shutdown).
	Error LogConfig

	// For producing final messages before application’s death. Note: critical messages force immediate flush
	// because in critical situation it is important to avoid log message losses if app crashes.
	Critical LogConfig
}

func (el *Log) Prepare() error {
	var err error

	if el.MinLevel == 0 {
		el.MinLevel = LOG_LEVEL_WARN
	}

	if el.MaxLevel == 0 {
		el.MaxLevel = LOG_LEVEL_CRITICAL
	}

	// Trace
	if el.Trace.Path == "" {
		el.Trace.Path = "./log/info.log"
	}

	if err = util.DirMake(el.Trace.Path); err != nil {
		return err
	}

	if el.Trace.Format == "" {
		el.Trace.Format = kLogFormat
	}

	if el.Trace.MaxRoll == 0 {
		el.Trace.MaxRoll = 2
	}

	if el.Trace.MaxSize == 0 {
		el.Trace.MaxSize = 50 * 1024 * 1024
	}

	if el.Trace.Console == BOOL_NOT_SET {
		el.Trace.Console = BOOL_TRUE
	}

	// Info
	if el.Info.Path == "" {
		el.Info.Path = "./log/info.log"
	}

	if err = util.DirMake(el.Info.Path); err != nil {
		return err
	}

	if el.Info.Format == "" {
		el.Info.Format = kLogFormat
	}

	if el.Info.MaxRoll == 0 {
		el.Info.MaxRoll = 2
	}

	if el.Info.MaxSize == 0 {
		el.Info.MaxSize = 50 * 1024 * 1024
	}

	if el.Info.Console == BOOL_NOT_SET {
		el.Info.Console = BOOL_TRUE
	}

	// Debug
	if el.Debug.Path == "" {
		el.Debug.Path = "./log/debug.log"
	}

	if err = util.DirMake(el.Debug.Path); err != nil {
		return err
	}

	if el.Debug.Format == "" {
		el.Debug.Format = kLogFormat
	}

	if el.Debug.MaxRoll == 0 {
		el.Debug.MaxRoll = 2
	}

	if el.Debug.MaxSize == 0 {
		el.Debug.MaxSize = 50 * 1024 * 1024
	}

	if el.Debug.Console == BOOL_NOT_SET {
		el.Debug.Console = BOOL_TRUE
	}

	// Warn
	if el.Warn.Path == "" {
		el.Warn.Path = "./log/warn.log"
	}

	if err = util.DirMake(el.Warn.Path); err != nil {
		return err
	}

	if el.Warn.Format == "" {
		el.Warn.Format = kLogFormat
	}

	if el.Warn.MaxRoll == 0 {
		el.Warn.MaxRoll = 2
	}

	if el.Warn.MaxSize == 0 {
		el.Warn.MaxSize = 200 * 1024 * 1024
	}

	if el.Warn.Console == BOOL_NOT_SET {
		el.Warn.Console = BOOL_TRUE
	}

	// Error
	if el.Error.Path == "" {
		el.Error.Path = "./log/warn.log"
	}

	if err = util.DirMake(el.Error.Path); err != nil {
		return err
	}

	if el.Error.Format == "" {
		el.Error.Format = kLogFormat
	}

	if el.Error.MaxRoll == 0 {
		el.Error.MaxRoll = 2
	}

	if el.Error.MaxSize == 0 {
		el.Error.MaxSize = 200 * 1024 * 1024
	}

	if el.Error.Console == BOOL_NOT_SET {
		el.Error.Console = BOOL_TRUE
	}

	// Critical
	if el.Critical.Path == "" {
		el.Critical.Path = "./log/warn.log"
	}

	if err = util.DirMake(el.Critical.Path); err != nil {
		return err
	}

	if el.Critical.Format == "" {
		el.Critical.Format = kLogFormat
	}

	if el.Critical.MaxRoll == 0 {
		el.Critical.MaxRoll = 2
	}

	if el.Critical.MaxSize == 0 {
		el.Critical.MaxSize = 200 * 1024 * 1024
	}

	if el.Critical.Console == BOOL_NOT_SET {
		el.Critical.Console = BOOL_TRUE
	}

	return nil
}

// Inicializa algumas variáveis
func (el *ProxyConfig) Prepare(config ConfigStart) {

	var traceConsole string
	if config.Log.Trace.Console == BOOL_TRUE {
		traceConsole = "<console/>"
	}

	var debugConsole string
	if config.Log.Debug.Console == BOOL_TRUE {
		debugConsole = "<console/>"
	}

	var infoConsole string
	if config.Log.Info.Console == BOOL_TRUE {
		infoConsole = "<console/>"
	}

	var warnConsole string
	if config.Log.Warn.Console == BOOL_TRUE {
		warnConsole = "<console/>"
	}

	var errorConsole string
	if config.Log.Error.Console == BOOL_TRUE {
		errorConsole = "<console/>"
	}

	var criticalConsole string
	if config.Log.Critical.Console == BOOL_TRUE {
		criticalConsole = "<console/>"
	}

	// Configura o log como arquivos com tamanho limitado. Um arquivo, info.log para coisas simples e um arquivo warn.log
	// para coisas que devem ser observadas pelo administrador
	if el.SeeLogConfig == "" {
		el.SeeLogConfig = `<seelog minlevel="warn" maxlevel="critical" type="sync">
  <outputs formatid="all">
    <filter levels="trace">
      <rollingfile type="size" filename="` + config.Log.Trace.Path + `" maxrolls="` + config.Log.Trace.GetMaxRoll() + `" maxsize="` + config.Log.Trace.GetMaxSize() + `" />` + traceConsole + `
    </filter>
    <filter levels="debug">
      <rollingfile type="size" filename="` + config.Log.Debug.Path + `" maxrolls="` + config.Log.Debug.GetMaxRoll() + `" maxsize="` + config.Log.Debug.GetMaxSize() + `" />` + debugConsole + `
    </filter>
    <filter levels="info">
      <rollingfile type="size" filename="` + config.Log.Info.Path + `" maxrolls="` + config.Log.Info.GetMaxRoll() + `" maxsize="` + config.Log.Info.GetMaxSize() + `" />` + infoConsole + `
    </filter>
    <filter levels="warn">
      <rollingfile type="size" filename="` + config.Log.Warn.Path + `" maxrolls="` + config.Log.Warn.GetMaxRoll() + `" maxsize="` + config.Log.Warn.GetMaxSize() + `" />` + warnConsole + `
    </filter>
    <filter levels="error">
      <rollingfile type="size" filename="` + config.Log.Error.Path + `" maxrolls="` + config.Log.Error.GetMaxRoll() + `" maxsize="` + config.Log.Error.GetMaxSize() + `" />` + errorConsole + `
    </filter>
    <filter levels="critical">
      <rollingfile type="size" filename="` + config.Log.Critical.Path + `" maxrolls="` + config.Log.Critical.GetMaxRoll() + `" maxsize="` + config.Log.Critical.GetMaxSize() + `" />` + criticalConsole + `
    </filter>
  </outputs>
  <formats>
    <format id="trace"    format="` + config.Log.Trace.Format + `"/>
    <format id="debug"    format="` + config.Log.Debug.Format + `"/>
    <format id="info"     format="` + config.Log.Info.Format + `"/>
    <format id="warn"     format="` + config.Log.Warn.Format + `"/>
    <format id="error"    format="` + config.Log.Error.Format + `"/>
    <format id="critical" format="` + config.Log.Critical.Format + `"/>
    <format id="all"      format="` + kLogFormat + `"/>
  </formats>
</seelog>`
	}

	logger, err := log.LoggerFromConfigAsBytes([]byte(el.SeeLogConfig))
	if err != nil {
		fmt.Printf("Erro na configuração do log. Error: %v", err.Error())
	}
	log.UseLogger(logger)

	// Define o tamanho do token como sendo 30 caracteres
	if el.UniqueIdLength == 0 {
		el.UniqueIdLength = 30
	}

	// Expressão regular do domínio
	if el.DomainExpReg == "" {
		el.DomainExpReg = `^(?P<subDomain>[a-zA-Z0-9]??|[a-zA-Z0-9]?[a-zA-Z0-9-.]*?[a-zA-Z0-9]*)[.]*(?P<domain>[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9]):*(?P<port>[0-9]*)$`
	}

	// Após 20 tentativas de acessar todos os containers da rota, uma mensagem de erro é exibida
	if el.MaxLoopTry == 0 {
		el.MaxLoopTry = 20
	}

	// Caso um container apresente mais de 10 erros consecutivos, o mesmo é desabilitado
	if el.ConsecutiveErrorsToDisable == 0 {
		el.ConsecutiveErrorsToDisable = 10
	}

	// Deixa um container desabilitado por 90 segundos após vários erros consecutivos
	if el.TimeToKeepDisabled == 0 {
		el.TimeToKeepDisabled = time.Second * 90
	}

	// Faz um teste a cada 30 segundos para saber se há containers desabilitados além do tempo
	if el.TimeToVerifyDisabled == 0 {
		el.TimeToVerifyDisabled = time.Second * 30
	}

	// Função de erro padrão do sistema
	if el.ErrorHandle == nil {
		el.ErrorHandle = el.ProxyError
	}

	// Função de page not found padrão do sistema
	if el.NotFoundHandle == nil {
		el.NotFoundHandle = el.ProxyNotFound
	}

	// Habilita todas as rotas do proxy, pois, o padrão do go é false
	for routesKey := range el.Routes.Get() {
		for urlKey := range el.Routes.GetKey(routesKey).ProxyServers.Get() {
			tmp := el.Routes.GetKey(routesKey).ProxyServers.GetKey(urlKey)
			tmp.Enabled = true
			//el.Routes.GetKey(routesKey).ProxyServers.SetKey(urlKey, tmp)
			route := el.Routes.GetKey(routesKey)
			route.ProxyServers.SetKey(urlKey, tmp)
			el.Routes.SetKey(routesKey, route)
		}
	}

	FuncMap[runtime.FuncForPC(reflect.ValueOf(el.ProxyError).Pointer()).Name()] = el.ProxyError
	FuncMap[runtime.FuncForPC(reflect.ValueOf(el.ProxyNotFound).Pointer()).Name()] = el.ProxyNotFound

	for routesKey := range el.Routes.Get() {
		FuncMap[runtime.FuncForPC(reflect.ValueOf(el.Routes.GetKey(routesKey).Domain.NotFoundHandle).Pointer()).Name()] = el.Routes.GetKey(routesKey).Domain.NotFoundHandle
		FuncMap[runtime.FuncForPC(reflect.ValueOf(el.Routes.GetKey(routesKey).Domain.ErrorHandle).Pointer()).Name()] = el.Routes.GetKey(routesKey).Domain.ErrorHandle
		FuncMap[runtime.FuncForPC(reflect.ValueOf(el.Routes.GetKey(routesKey).Handle.Handle).Pointer()).Name()] = el.Routes.GetKey(routesKey).Handle.Handle
	}

	el.readyToJSon = true

	el.RoutePrepare()
}

func (el *ProxyConfig) ProxyError(w ProxyResponseWriter, r *ProxyRequest) {
	w.Write([]byte(`<html><header><style>body{height:100%; position:relative}div{margin:auto;height: 100%;width: 100%;position:fixed;top:0;bottom:0;left:0;right:0;background:blue;}div.center{margin:auto;height: 70%;width: 70%;}</style></header><body><div><div style="color:#ffff;" class="center"><p style="text-align: center; background-color: #888888;">There is something very wrong!</p><p>&nbsp;</p>The address is correct, but no server has responded correctly. The system administrator will be informed about this.<p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!</div></div></body></html>`))
}

func (el *ProxyConfig) ProxyNotFound(w ProxyResponseWriter, r *ProxyRequest) {
	w.Write([]byte(`<html><header><style>body{height:100%; position:relative}div{margin:auto;height: 100%;width: 100%;position:fixed;top:0;bottom:0;left:0;right:0;background:blue;}div.center{margin:auto;height: 70%;width: 70%;}</style></header><body><div><div style="color:#ffff;" class="center"><p style="text-align: center; background-color: #888888;">Page Not Found!</p><p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!<p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!</div></div></body></html>`))
}

// Verifica se há urls do proxy desabilitadas e as habilita depois de um tempo
// A ideia é que o servidor possa está fora do ar por um tempo, por isto, ele remove a rota por algum tempo, para evitar
// chamadas desnecessárias ao servidor
func (el *ProxyConfig) VerifyDisabled() {
	for {
		for routesKey := range el.Routes.Get() {
			for urlKey := range el.Routes.GetKey(routesKey).ProxyServers.Get() {
				tmp := el.Routes.GetKey(routesKey).ProxyServers.GetKey(urlKey)
				if time.Since(tmp.DisabledSince) >= el.TimeToKeepDisabled && tmp.Enabled == false && tmp.Forever == false {
					tmp.ErrorConsecutiveCounter = 0
					tmp.Enabled = true
					//el.Routes.GetKey(routesKey).ProxyServers.SetKey(urlKey, tmp)
					route := el.Routes.GetKey(routesKey)
					route.ProxyServers.SetKey(urlKey, tmp)
					el.Routes.SetKey(routesKey, route)
				}
			}
		}

		time.Sleep(el.TimeToVerifyDisabled)
	}
}

func (el *ProxyConfig) ProxyRoutes(w ProxyResponseWriter, r *ProxyRequest) {

	byteJSon, err := json.Marshal(ProxyRootConfig.Routes)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(byteJSon)
}
func (el *ProxyConfig) MarshalJSON() ([]byte, error) {
	if el.readyToJSon == false {
		return []byte{}, errors.New("call prepare() before json.Marshal() function")
	}
	return json.Marshal(&struct {
		SeeLogConfig               string        `json:"seeLogConfig"`
		DomainExpReg               string        `json:"domainExpReg"`
		ErrorHandleAsString        string        `json:"ErrorHandle"`
		NotFoundHandleAsString     string        `json:"NotFoundHandle"`
		UniqueIdLength             int           `json:"uniqueIdLength"`
		ListenAndServe             string        `json:"listenAndServe"`
		MaxLoopTry                 int           `json:"maxLoopTry"`
		ConsecutiveErrorsToDisable int64         `json:"consecutiveErrorsToDisable"`
		TimeToKeepDisabled         time.Duration `json:"timeToKeepDisabled"`
		TimeToVerifyDisabled       time.Duration `json:"timeToVerifyDisabled"`
		Routes                     []ProxyRoute  `json:"routes"`
	}{
		SeeLogConfig:               el.SeeLogConfig,
		DomainExpReg:               el.DomainExpReg,
		ErrorHandleAsString:        runtime.FuncForPC(reflect.ValueOf(el.ErrorHandle).Pointer()).Name(),
		NotFoundHandleAsString:     runtime.FuncForPC(reflect.ValueOf(el.NotFoundHandle).Pointer()).Name(),
		UniqueIdLength:             el.UniqueIdLength,
		ListenAndServe:             el.ListenAndServe,
		MaxLoopTry:                 el.MaxLoopTry,
		ConsecutiveErrorsToDisable: el.ConsecutiveErrorsToDisable,
		TimeToKeepDisabled:         el.TimeToKeepDisabled,
		TimeToVerifyDisabled:       el.TimeToVerifyDisabled,
		Routes:                     el.Routes.Get(),
	})
}
func (el *ProxyConfig) UnmarshalJSON(data []byte) error {
	type tmpStt struct {
		SeeLogConfig               string        `json:"seeLogConfig"`
		DomainExpReg               string        `json:"domainExpReg"`
		ErrorHandleAsString        string        `json:"ErrorHandle"`
		NotFoundHandleAsString     string        `json:"NotFoundHandle"`
		UniqueIdLength             int           `json:"uniqueIdLength"`
		ListenAndServe             string        `json:"listenAndServe"`
		MaxLoopTry                 int           `json:"maxLoopTry"`
		ConsecutiveErrorsToDisable int64         `json:"consecutiveErrorsToDisable"`
		TimeToKeepDisabled         time.Duration `json:"timeToKeepDisabled"`
		TimeToVerifyDisabled       time.Duration `json:"timeToVerifyDisabled"`
		Routes                     []ProxyRoute  `json:"routes"`
	}
	var tmp = tmpStt{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	el.SeeLogConfig = tmp.SeeLogConfig
	el.DomainExpReg = tmp.DomainExpReg

	if tmp.ErrorHandleAsString != "" {
		el.ErrorHandle = FuncMap[tmp.ErrorHandleAsString].(ProxyHandlerFunc)
	}

	if tmp.NotFoundHandleAsString != "" {
		el.NotFoundHandle = FuncMap[tmp.NotFoundHandleAsString].(ProxyHandlerFunc)
	}

	el.UniqueIdLength = tmp.UniqueIdLength
	el.ListenAndServe = tmp.ListenAndServe
	el.MaxLoopTry = tmp.MaxLoopTry
	el.ConsecutiveErrorsToDisable = tmp.ConsecutiveErrorsToDisable
	el.TimeToKeepDisabled = tmp.TimeToKeepDisabled
	el.TimeToVerifyDisabled = tmp.TimeToVerifyDisabled
	el.Routes.Set(tmp.Routes)

	return nil
}

func NewStartConfig() ConfigStart {
	return ConfigStart{}
}

func init() {
	ProxyRootConfig.Routes.Set(ProxyNewRootConfig.Get())
}
