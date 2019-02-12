package SimpleReverseProxy

import "sync"

type ProxyRouteAStt struct {
	l    sync.Mutex
	a    []ProxyRoute
	init bool
}

func (el *ProxyRouteAStt) Clear() {
	defer el.l.Unlock()

	el.l.Lock()
	el.init = true
	el.a = make([]ProxyRoute, 0)
}
func (el *ProxyRouteAStt) Append(v ProxyRoute) {
	defer el.l.Unlock()

	el.l.Lock()
	if el.init == false {
		el.init = true
		el.a = make([]ProxyRoute, 0)
	}

	el.a = append(el.a, v)
}
func (el *ProxyRouteAStt) Len() int {
	defer el.l.Unlock()

	el.l.Lock()
	return len(el.a)
}
func (el *ProxyRouteAStt) GetKey(i int) ProxyRoute {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a[i]
}
func (el *ProxyRouteAStt) Get() []ProxyRoute {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a
}
func (el *ProxyRouteAStt) Set(v []ProxyRoute) {
	defer el.l.Unlock()

	el.a = v

	el.l.Lock()
}
