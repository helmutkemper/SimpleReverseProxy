package marketPlaceProcy

import (
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
  "time"
)

var containerStatsInterval time.Duration
var containerStatsLimit int
var containerStatsData map[string][]containerStatsOut

func ContainerWebStatsLogSetLimit( limit int ){
  containerStatsLimit = limit
}

func ContainerStatsLog() error {
  var idList = make( map[string]bool )

  cli, err := client.NewEnvClient()
  if err != nil {
    return err
  }

  containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
  if err != nil {
    return err
  }

  for _, containerData := range containers {
    idList[ containerData.ID ] = true

    if len( containerStatsData[ containerData.ID ] ) == 0 || len( containerStatsData[ containerData.ID ] ) != containerStatsLimit {
      containerStatsData[ containerData.ID ] = make( []containerStatsOut, containerStatsLimit )

      for i := 0; i != containerStatsLimit; i += 1 {
        containerStatsData[ containerData.ID ][ i ] = containerStatsOut{}
      }
    }

    decode := containerDockerStats{}
    decode.ToDecode, err = cli.ContainerStats(context.Background(), containerData.ID, false)
    if err != nil {
      return err
    }

    err = decode.Decode()
    if err != nil {
      return err
    }

    containerStatsData[ containerData.ID ] = containerStatsData[ containerData.ID ][1:]
    containerStatsData[ containerData.ID ] = append( containerStatsData[ containerData.ID ], decode.Stats )
  }

  // When a container is removed, its data remains and must be removed
  for idFound := range containerStatsData {
    if idList[ idFound ] == false {
      delete( idList, idFound )
    }
  }

  return nil
}

func ContainerStatsLogStart() {
  for {
    ContainerStatsLog()
    time.Sleep( containerStatsInterval )
  }
}

func ContainerWebStatsLog(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}
  output.ToOutput( len( containerStatsData ), nil, containerStatsData, w )
}

func ContainerWebStatsLogById(w ProxyResponseWriter, r *ProxyRequest) {
  output := JSonOutStt{}

  var inDataLStt ContainerStartByIdDataIn
  var err, id = Input(r, &inDataLStt)

  if err != nil {
    output.ToOutput( 0, err, []int{}, w )
    return
  }

  if id != "" {
    inDataLStt.Id = id
  }

  output.ToOutput( len( containerStatsData[ inDataLStt.Id ] ), nil, containerStatsData[ inDataLStt.Id ], w )
}

func init(){
  // This value directly influences the length of the graph kendo ui spark line
  containerStatsLimit = 99
  containerStatsInterval = time.Second * 2
  containerStatsData = make( map[string][]containerStatsOut )

  go ContainerStatsLogStart()
}
