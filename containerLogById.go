package marketPlaceProcy

import (
  "bytes"
  "golang.org/x/net/context"
  "github.com/helmutkemper/moby/client"
  "github.com/helmutkemper/moby/api/types"
  "fmt"
)

type ContainerLogsByIdDataIn struct {
  Id        string
  Options   types.ContainerLogsOptions
}

func ContainerLogsById(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerLogsByIdDataIn
  var err, id = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  if id != "" {
    inDataLStt.Id = id
  }

  if inDataLStt.Options.ShowStderr == false && inDataLStt.Options.ShowStdout == false {
    inDataLStt.Options.ShowStderr = true
    inDataLStt.Options.ShowStdout = true
  }

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  out, err := cli.ContainerLogs(ctx, inDataLStt.Id, inDataLStt.Options)
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  buf := new(bytes.Buffer)
  if out != nil {
    buf.ReadFrom(out)
  }

  output.ToOutput(1, err, []string{ fmt.Sprintf("%v", buf.String() ) }, w)
}
