package marketPlaceProcy

import (
  "time"
  log "github.com/helmutkemper/seelog"
)

func timeMeasure( start time.Time, name string ) {
  elapsed := time.Since(start)
  log.Infof("%s: %s", name, elapsed)
}
