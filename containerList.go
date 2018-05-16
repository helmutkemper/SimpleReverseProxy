package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

func ContainerList(contextAStt context.Context, clientAStt client.Client, setupLAStt types.ContainerListOptions) (error, []types.Container){
  containersLStt, err := clientAStt.ContainerList(contextAStt, setupLAStt)

  return err, containersLStt
}
