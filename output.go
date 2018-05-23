package SimpleReverseProxy

import (
  "net/http"
  "encoding/json"
  "log"
)

type MetaJSonOutStt struct{
  Limit             int                   `json:"Limit"`
  Next              string                `json:"Next"`
  Offset            int64                 `json:"Offset"`
  Previous          string                `json:"Previous"`
  TotalCount        int                   `json:"TotalCount"`
  Success           bool                  `json:"Success"`
  Error             interface{}           `json:"Error"`
}

type JSonOutStt struct{
  Meta              MetaJSonOutStt        `json:"Meta"`
  Objects           interface{}           `json:"Objects"`
  geoJSonHasOutput  bool                  `json:"-"`
}

func( JSonOutAStt *JSonOutStt ) ToOutput( totalCountAInt int, errorAErr interface{}, dataATfc interface{}, w http.ResponseWriter ) {
  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  if errorAErr != nil {
    w.WriteHeader(http.StatusInternalServerError)

    JSonOutAStt.Meta = MetaJSonOutStt{
      Error: errorAErr,
      Limit: 0,
      Next:  "",
      Offset:0,
      Previous: "",
      Success: false,
      TotalCount: 0,
    }
    JSonOutAStt.Objects = make( []int, 0 )
  } else {
    if totalCountAInt == 0 {
      w.WriteHeader(http.StatusOK)

      JSonOutAStt.Meta = MetaJSonOutStt{
        Error: "",
        Limit: 0,
        Next:  "",
        Offset:0,
        Previous: "",
        Success: true,
        TotalCount: totalCountAInt,
      }
      JSonOutAStt.Objects = make( []int, 0 )
    } else {
      JSonOutAStt.Meta = MetaJSonOutStt{
        Error: "",
        Limit: 20,
        Next:  "",
        Offset:0,
        Previous: "",
        Success: true,
        TotalCount: totalCountAInt,
      }

      JSonOutAStt.Objects = dataATfc
    }
  }

  if err := json.NewEncoder( w ).Encode( JSonOutAStt ); err != nil {
    log.Fatal( err )
  }
}