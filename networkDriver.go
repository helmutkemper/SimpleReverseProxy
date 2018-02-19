package marketPlaceProcy

type NetworkDriver int

var networkDrivers = [...]string {
  "bridge",
  "overlay",
}

func (el NetworkDriver) String() string {
  return networkDrivers[el]
}

const (
  BRIDGE NetworkDriver = iota
  OVERLAY
)
