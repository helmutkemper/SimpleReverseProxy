package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/volume"
)

type VolumesCreateDataIn struct{
  Options       volume.VolumesCreateBody    `json:"options"`
}

/*
  {
    "options": {
      "Name": "data"
    }
  }
  {
    "options": {
      "Name": "cfg"
    }
  }
 */
func VolumesCreate(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt VolumesCreateDataIn
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

  createOut, err := cli.VolumeCreate(ctx, inDataLStt.Options)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  output.ToOutput(1, nil, []types.Volume{ createOut }, w)
}
