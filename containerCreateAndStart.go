package marketPlaceProcy

import (
  "bytes"
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/api/types/network"
)

type ContainerDockerConfig   container.Config
type ContainerHostConfig     container.HostConfig
type ContainerNetworkConfig  network.NetworkingConfig

type CreateAndStartDataIn struct{
  Docker        ContainerDockerConfig
  Host          ContainerHostConfig
  Network       ContainerNetworkConfig
  Name          string
  WaitForErrors bool
}

func ContainerCreateAndStart(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt CreateAndStartDataIn
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

  docker  := container.Config(inDataLStt.Docker)
  host    := container.HostConfig(inDataLStt.Host)
  network := network.NetworkingConfig(inDataLStt.Network)

  createResp, err := cli.ContainerCreate(ctx, &docker, &host, &network, inDataLStt.Name)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  if err := cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  buf := new(bytes.Buffer)
  if inDataLStt.WaitForErrors == true {
    cli.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)

    out, err := cli.ContainerLogs(ctx, createResp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
    if err != nil {
      output.ToOutput(0, err, []int{}, w)
      return
    }

    if out != nil {
      buf.ReadFrom(out)
    }
  }

  output.ToOutput( 1, err, []struct{
    Created   container.ContainerCreateCreatedBody
    OutPut    string
  }{{ createResp, buf.String() }}, w )
}
