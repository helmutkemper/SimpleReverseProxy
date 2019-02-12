package SimpleReverseProxy

import "sync"

type ProxyUrlArrStt struct {
	l    sync.Mutex
	a    []ProxyUrl
	init bool
}

func (el *ProxyUrlArrStt) Clear() {
	defer el.l.Unlock()

	el.l.Lock()
	el.init = true
	el.a = make([]ProxyUrl, 0)
}
func (el *ProxyUrlArrStt) Append(v ProxyUrl) {
	defer el.l.Unlock()

	el.l.Lock()
	if el.init == false {
		el.init = true
		el.a = make([]ProxyUrl, 0)
	}

	el.a = append(el.a, v)
}
func (el *ProxyUrlArrStt) Len() int {
	defer el.l.Unlock()

	el.l.Lock()
	return len(el.a)
}
func (el ProxyUrlArrStt) GetKey(i int) ProxyUrl {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a[i]
}
func (el ProxyUrlArrStt) Get() []ProxyUrl {
	defer el.l.Unlock()

	el.l.Lock()
	return el.a
}
func (el *ProxyUrlArrStt) Set(v []ProxyUrl) {
	defer el.l.Unlock()

	el.a = v

	el.l.Lock()
}
func (el *ProxyUrlArrStt) SetKey(k int, v ProxyUrl) {
	defer el.l.Unlock()

	el.a[k] = v

	el.l.Lock()
}
