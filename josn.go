package marketPlaceProcy

import (
  "net/http"
  "encoding/json"
  log "github.com/helmutkemper/seelog"
)

type MetaJSonOutStt struct{
  TotalCount  int64               `json:"TotalCount"`
  Error       string              `json:"Error"`
}

type JSonOutStt struct{
  Meta              MetaJSonOutStt      `json:"Meta"`
  Objects           interface{}         `json:"Objects"`
  geoJSonHasOutput  bool                `json:"-"`
}

func( el *JSonOutStt ) ToOutput( totalCountAInt int64, errorAErr error, dataATfc interface{}, w ProxyResponseWriter ) {
  var errorString = ""

  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  if errorAErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    errorString = errorAErr.Error()
    totalCountAInt = 0
  } else {
    w.WriteHeader(http.StatusOK)
  }

  el.Meta = MetaJSonOutStt{
    Error: errorString,
    TotalCount: totalCountAInt,
  }

  if errorAErr != nil {
    el.Objects = []int{}
  } else {
    switch dataATfc.(type) {
    default:
      el.Objects = dataATfc
    }
  }

  if err := json.NewEncoder( w ).Encode( el ); err != nil {
    log.Warn( err )
  }
}
