package marketPlaceProcy

import (
  "runtime"
  "reflect"
)

type funcMapType map[string]interface{}
func( el *funcMapType )Add( fn interface{} ){
  (*el)[ runtime.FuncForPC( reflect.ValueOf( fn ).Pointer() ).Name() ] = fn
}

var FuncMap funcMapType

func init(){
  FuncMap = make( funcMapType )
}