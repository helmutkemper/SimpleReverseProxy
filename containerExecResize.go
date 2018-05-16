package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

type ContainerExecResizeDataIn struct{
  Id          string
  Options     types.ResizeOptions
}

func ContainerExecResize(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerExecResizeDataIn
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

  err = cli.ContainerExecResize(ctx, inDataLStt.Id, inDataLStt.Options)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  output.ToOutput( 1, err, []int{}, w )
}
