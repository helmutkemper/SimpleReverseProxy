package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

type ContainerRemoveDataIn struct {
  Id      string
}

func ContainerRemove(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerStartDataIn
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
    output.ToOutput(0, err, []int{}, w)
    return
  }

  if err := cli.ContainerRemove(ctx, inDataLStt.Id, types.ContainerRemoveOptions{}); err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  output.ToOutput( 0, err, []int{}, w )
}
