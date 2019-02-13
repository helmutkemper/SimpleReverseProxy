package SimpleReverseProxy

type Boolean int

var Booleans = [...]string{
	"false",
	"",
	"true",
}

func (el Boolean) String() string {
	return Booleans[el]
}

const (
	BOOL_FALSE Boolean = iota - 1
	BOOL_NOT_SET
	BOOL_TRUE
)
