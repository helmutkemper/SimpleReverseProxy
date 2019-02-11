package SimpleReverseProxy

import "sync"

type NewRootConfigStt struct {
	l    sync.Mutex
	a    []ProxyRoute
	init bool
}

func (el *NewRootConfigStt) Clear() {
	defer el.l.Unlock()

	el.l.Lock()
	el.init = true
	el.a = make([]ProxyRoute, 0)
}
func (el *NewRootConfigStt) Append(v ProxyRoute) {
	defer el.l.Unlock()

	el.l.Lock()
	if el.init == false {
		el.init = true
		el.a = make([]ProxyRoute, 0)
	}

	el.a = append(el.a, v)
}
func (el *NewRootConfigStt) Len() int {
	defer el.l.Unlock()

	el.l.Lock()
	return len(el.a)
}
func (el *NewRootConfigStt) GetKey(i int) ProxyRoute {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a[i]
}
func (el *NewRootConfigStt) Get() []ProxyRoute {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a
}
func (el *NewRootConfigStt) Set(v []ProxyRoute) {
	defer el.l.Unlock()

	el.a = v

	el.l.Lock()
}
