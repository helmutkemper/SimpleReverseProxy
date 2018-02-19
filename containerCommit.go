package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

type ContainerCommitStt struct {
  Id        string
  Options   types.ContainerCommitOptions
}

func ContainerCommit(contextAStt context.Context, clientAStt client.Client, setupLAStt []ContainerCommitStt) (error, []types.IDResponse) {
  var returnLAStt []types.IDResponse

  for _, contDataLStt := range setupLAStt {
    commitResp, err := clientAStt.ContainerCommit(contextAStt, contDataLStt.Id, contDataLStt.Options)
    if err != nil {
      return err, []types.IDResponse{}
    }

    returnLAStt = append( returnLAStt, commitResp )
  }

  return nil, returnLAStt
}