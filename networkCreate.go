package marketPlaceProcy

import (
  "github.com/helmutkemper/moby/api/types"
  "github.com/helmutkemper/moby/client"
  "context"
)

type NetworkCreateStt struct {
  Name      string
  Options   types.NetworkCreate
}

func NetworkCreate( contextAStt context.Context, clientAStt client.Client, setupLAStt []NetworkCreateStt ) (error, []types.NetworkCreateResponse) {
  var returnLAStt []types.NetworkCreateResponse

  for _, netDataLStt := range setupLAStt {
    responseLStt, err := clientAStt.NetworkCreate(contextAStt, netDataLStt.Name, netDataLStt.Options)
    if err != nil {
      return err, []types.NetworkCreateResponse{}
    }
    returnLAStt = append( returnLAStt, responseLStt )
  }

  return nil, returnLAStt
}
