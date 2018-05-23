package SimpleReverseProxy

import (
  "runtime"
  "reflect"
)

type funcMapType map[string]interface{}
func( el *funcMapType )Add( fn interface{} ){
  (*el)[ runtime.FuncForPC( reflect.ValueOf( fn ).Pointer() ).Name() ] = fn
}
func( el *funcMapType )List() []string {
  var ret = make( []string, len(*el) )
  var count = 0

  for k := range *el {
    ret[count] = k
    count += 1
  }

  return ret
}

var FuncMap funcMapType

func init(){
  FuncMap = make( funcMapType )
}