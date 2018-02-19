package marketPlaceProcy

import (
  "github.com/docker/docker/client"
  "context"
  "net/http"
)

func NewEnvClient() (error, *client.Client, context.Context) {
  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    return err, nil, nil
  }
  return nil, cli, ctx
}

func NewClient(host, version string, clientHttp *http.Client, headersHttp map[string]string) (error, *client.Client, context.Context) {
  ctx := context.Background()
  cli, err := client.NewClient(host, version, clientHttp, headersHttp)
  if err != nil {
    return err, nil, nil
  }
  return nil, cli, ctx
}