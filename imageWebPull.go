package marketPlaceProcy

import (
  "bytes"
  "encoding/json"
  "encoding/base64"
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
)

//{"username":"username","password":"password"}
type ImagePullRegistry struct {
  UserName                  string                        `json:"username"`
  PassWord                  string                        `json:"password"`
}
type ImagePullOptions struct {
  All                       bool
  Registry                  ImagePullRegistry
  RegistryAuth              string // RegistryAuth is the base64 encoded credentials for the registry
  PrivilegeFunc             types.RequestPrivilegeFunc    `json:"-"`
  Platform                  string
}
type PullDataIn struct{
  ImageName                 string                        `json:"imagename"`
  Options                   ImagePullOptions                   `json:"options"`
}
/*
{
  "imagename": "alpine:latest",
  "options": {
    "Registry": {
      "username": "username",
      "password": "password"
    }
  }
}
*/
func ImageWebPull(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt PullDataIn
  var err, _ = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  encodedJSON, err := json.Marshal(inDataLStt.Options.Registry)
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  inDataLStt.Options.RegistryAuth = base64.URLEncoding.EncodeToString(encodedJSON)

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  imageOut, err := cli.ImagePull(ctx, inDataLStt.ImageName, types.ImagePullOptions{
    All: inDataLStt.Options.All,
    RegistryAuth: inDataLStt.Options.RegistryAuth,
    PrivilegeFunc: inDataLStt.Options.PrivilegeFunc,
    Platform: inDataLStt.Options.Platform,
  })
  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  buf := new(bytes.Buffer)
  buf.ReadFrom(imageOut)

  output.ToOutput(1, nil, []string{ buf.String() }, w)
}