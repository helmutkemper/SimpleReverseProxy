package marketPlaceProcy

import (
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "context"
)

func NetworkList( contextAStt context.Context, clientAStt client.Client, setupLStt types.NetworkListOptions ) (error, []types.NetworkResource) {

  listOut, err := clientAStt.NetworkList(contextAStt, setupLStt)
  return err, listOut
}
