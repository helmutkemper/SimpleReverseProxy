package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

func ContainerStopAllRunning(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  for _, containerLStt := range containers {
    if err := cli.ContainerStop(ctx, containerLStt.ID, nil); err != nil {
      output.ToOutput(0, err, []int{}, w)
      return
    }
  }

  output.ToOutput(0, err, []int{}, w)
}