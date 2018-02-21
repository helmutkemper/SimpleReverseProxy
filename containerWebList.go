package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

func ContainerWebList(w ProxyResponseWriter, r *ProxyRequest) {
  //w.Header().Set("Access-Control-Allow-Origin", "http://htmlpanel.localhost:8888")
  output := JSonOutStt{}

  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
  }

  containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{ All: true })
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
  }

  output.ToOutput( len(containers), err, containers, w )
}
