package marketPlaceProcy

import (
  "time"
  "strings"
  "encoding/json"
  "golang.org/x/net/context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
  "errors"
)

type containerStatsDecode struct {
  // Common stats
  Read          time.Time                     `json:"read"`
  PreRead       time.Time                     `json:"preread"`

  // Linux specific stats, not populated on Windows.
  PidsStats     types.PidsStats               `json:"pids_stats,omitempty"`
  BlkioStats    types.BlkioStats              `json:"blkio_stats,omitempty"`

  // Windows specific stats, not populated on Linux.
  NumProcs      uint32                        `json:"num_procs"`
  StorageStats  types.StorageStats            `json:"storage_stats,omitempty"`

  // Shared stats
  CPUStats      types.CPUStats                `json:"cpu_stats,omitempty"`
  PreCPUStats   types.CPUStats                `json:"precpu_stats,omitempty"`
  MemoryStats   types.MemoryStats             `json:"memory_stats,omitempty"`

  Networks      map[string]types.NetworkStats `json:"networks,omitempty"`
}

type containerStatsOut struct {
  MemPercent float64
  PreviousCPU uint64
  PreviousSystem uint64
  CpuPercent float64
  BlkRead uint64
  BlkWrite uint64
  Mem uint64
  MemLimit uint64
  PidsStatsCurrent uint64
  NetRx float64
  NetTx float64
}

type containerDockerStats struct {
  ToDecode        types.ContainerStats
  statsDecode     containerStatsDecode
  Stats           containerStatsOut
}

func ( el *containerDockerStats ) Decode() error {
  if el.ToDecode.Body == nil {
    return errors.New("docker don't send any data to decode")
  }
  decoder := json.NewDecoder(el.ToDecode.Body)
  err := decoder.Decode(&el.statsDecode)
  if err != nil {
    return err
  }

  el.Stats.MemPercent = 0.0
  if el.statsDecode.MemoryStats.Limit != 0 {
    el.Stats.MemPercent = float64(el.statsDecode.MemoryStats.Usage) /
      float64(el.statsDecode.MemoryStats.Limit) * 100.0
  }
  el.Stats.PreviousCPU = el.statsDecode.PreCPUStats.CPUUsage.TotalUsage
  el.Stats.PreviousSystem = el.statsDecode.PreCPUStats.SystemUsage
  el.Stats.CpuPercent = el.calculateCPUPercentUnix()
  el.Stats.BlkRead, el.Stats.BlkWrite = el.calculateBlockIO()
  el.Stats.Mem = el.statsDecode.MemoryStats.Usage
  el.Stats.MemLimit = el.statsDecode.MemoryStats.Limit
  el.Stats.PidsStatsCurrent = el.statsDecode.PidsStats.Current
  el.Stats.NetRx, el.Stats.NetTx = el.calculateNetwork()

  return nil
}

func ( el *containerDockerStats ) calculateCPUPercentUnix() float64 {
  var (
    cpuPercent = 0.0
    // calculate the change for the cpu usage of the container in
    // between readings
    cpuDelta = float64(el.statsDecode.CPUStats.CPUUsage.TotalUsage) -
      float64(el.Stats.PreviousCPU)
    // calculate the change for the entire system between readings
    systemDelta = float64(el.statsDecode.CPUStats.SystemUsage) -
      float64(el.Stats.PreviousSystem)
  )

  if systemDelta > 0.0 && cpuDelta > 0.0 {
    cpuPercent = (cpuDelta / systemDelta) *
      float64(len(el.statsDecode.CPUStats.CPUUsage.PercpuUsage)) * 100.0
  }
  return cpuPercent
}

func ( el *containerDockerStats ) calculateBlockIO() (blkRead uint64, blkWrite uint64) {
  for _, bioEntry := range el.statsDecode.BlkioStats.IoServiceBytesRecursive {
    switch strings.ToLower(bioEntry.Op) {
    case "read":
      blkRead = blkRead + bioEntry.Value
    case "write":
      blkWrite = blkWrite + bioEntry.Value
    }
  }
  return
}

func ( el *containerDockerStats ) calculateNetwork() (float64, float64) {
  var rx, tx float64

  for _, v := range el.statsDecode.Networks {
    rx += float64(v.RxBytes)
    tx += float64(v.TxBytes)
  }
  return rx, tx
}

type ContainerStartByIdDataIn struct {
  Id        string
}

func ContainerStatsById(w ProxyResponseWriter, r *ProxyRequest) {
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

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  decode := containerDockerStats{}
  decode.ToDecode, err = cli.ContainerStats(ctx, inDataLStt.Id, false)
  err = decode.Decode()
  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  output.ToOutput(1, err, []containerStatsOut{ decode.Stats }, w)
}
