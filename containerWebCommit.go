package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

type ContainerCommitDataIn struct {
  Id        string
  Options   types.ContainerCommitOptions
}

func ContainerWebCommit(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerCommitDataIn
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

  commitResp, err := cli.ContainerCommit(ctx, inDataLStt.Id, inDataLStt.Options)
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  output.ToOutput(1, err, []types.IDResponse{ commitResp }, w)
}