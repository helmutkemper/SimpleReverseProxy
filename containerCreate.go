package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/api/types/network"
)

type ContainerCreateStt struct {
  Docker    ContainerDockerConfig
  Host      ContainerHostConfig
  Network   ContainerNetworkConfig
  Name      string
}

func ContainerCreate(contextAStt context.Context, clientAStt client.Client, setupLAStt []ContainerCreateStt) (error, []container.ContainerCreateCreatedBody) {
  var returnLAStt []container.ContainerCreateCreatedBody

  for _, createDataLStt := range setupLAStt {
    dockerLStt := container.Config(createDataLStt.Docker)
    hostLStt := container.HostConfig(createDataLStt.Host)
    netWorkLStt := network.NetworkingConfig(createDataLStt.Network)

    createRespLStt, err := clientAStt.ContainerCreate(contextAStt, &dockerLStt, &hostLStt, &netWorkLStt, createDataLStt.Name)
    if err != nil {
      return err, []container.ContainerCreateCreatedBody{}
    }

    returnLAStt = append(returnLAStt, createRespLStt)
  }

  return nil, returnLAStt
}
