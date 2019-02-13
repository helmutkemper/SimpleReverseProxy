package SimpleReverseProxy

type LogLevel int

var LogLevels = [...]string{
	"",
	"trace",
	"debug",
	"info",
	"warn",
	"error",
	"critical",
}

func (el LogLevel) String() string {
	return LogLevels[el]
}

const (
	LOG_LEVEL_TRACE LogLevel = iota + 1
	LOG_LEVEL_DEBUG
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_CRITICAL
)
