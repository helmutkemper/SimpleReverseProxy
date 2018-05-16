package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

type NetworkListDataIn struct{
  Options       types.NetworkListOptions    `json:"options"`
}

func NetworkWebList(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt NetworkListDataIn
  var err, _ = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  listOut, err := cli.NetworkList(ctx, inDataLStt.Options)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  output.ToOutput(len(listOut), nil, listOut, w)
}
