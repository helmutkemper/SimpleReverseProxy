package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/api/types/network"
)

type ContainerCreateDataIn struct {
  Docker    ContainerDockerConfig
  Host      ContainerHostConfig
  Network   ContainerNetworkConfig
  Name      string
}

func ContainerWebCreate(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerCreateDataIn
  var err, _ = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  docker  := container.Config(inDataLStt.Docker)
  host    := container.HostConfig(inDataLStt.Host)
  netWork  := network.NetworkingConfig(inDataLStt.Network)

  createResp, err := cli.ContainerCreate(ctx, &docker, &host, &netWork, inDataLStt.Name)

  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  output.ToOutput(1, err, []container.ContainerCreateCreatedBody{ createResp }, w)
}