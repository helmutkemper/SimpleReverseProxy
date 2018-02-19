package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

type NetworkCreateDataIn struct{
  Name          string                      `json:"name"`
  Options       types.NetworkCreate         `json:"options"`
}

/*
{
  "name": "mongo",
  "options": {
    "Driver": "overlay",
    "Attachable": true,
  }
}
{
  "name": "mongos",
  "options": {
    "Driver": "overlay",
    "Attachable": true,
  }
}
*/
func NetworkWebCreate(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt NetworkCreateDataIn
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

  createOut, err := cli.NetworkCreate(ctx, inDataLStt.Name, inDataLStt.Options)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  output.ToOutput(1, nil, []types.NetworkCreateResponse{ createOut }, w)
}
