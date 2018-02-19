package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
)

type ContainerStopDataIn struct {
  Id      string
}

func ContainerStopById(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerStopDataIn
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

  if err := cli.ContainerStop(ctx, inDataLStt.Id, nil); err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  output.ToOutput( 0, err, []int{}, w )
}
