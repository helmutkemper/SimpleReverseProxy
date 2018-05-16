package marketPlaceProcy

import (
  "github.com/helmutkemper/moby/api/types"
  "github.com/helmutkemper/moby/client"
  "context"
)

func NetworkList( contextAStt context.Context, clientAStt client.Client, setupLStt types.NetworkListOptions ) (error, []types.NetworkResource) {

  listOut, err := clientAStt.NetworkList(contextAStt, setupLStt)
  return err, listOut
}
