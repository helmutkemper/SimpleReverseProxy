package SimpleReverseProxy

type LoadBalance int

var LoadBalances = [...]string{
	"",
	"round robin",
}

func (el LoadBalance) String() string {
	return LoadBalances[el]
}

const (
	LOAD_BALANCE_ROUND_ROBIN LoadBalance = iota + 1
)
