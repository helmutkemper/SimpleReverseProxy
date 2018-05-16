package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

type ContainerExecCreateDataIn struct{
  Id          string
  Options     types.ExecConfig
}

func ContainerExecCreate(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerExecCreateDataIn
  var err, id = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  if id != "" {
    inDataLStt.Id = id
  }

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  execResp, err := cli.ContainerExecCreate(ctx, inDataLStt.Id, inDataLStt.Options)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  output.ToOutput( 1, err, []types.IDResponse{ execResp }, w )
}
