package SimpleReverseProxy

import (
"github.com/helmutkemper/ntp"
"time"
"strconv"
)

type networkTimeStt struct {
  Loc *time.Location
  Dif time.Duration
}

var NetworkTime networkTimeStt

func(el *networkTimeStt)StringDay() string {
  var now = time.Now()
  now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC).In(el.Loc).Add(el.Dif)

  y := strconv.FormatInt( int64( now.Year() ), 10 )
  m := strconv.FormatInt( int64( now.Month() ), 10 )
  if len( m ) == 1 {
    m = "0" + m
  }
  d := strconv.FormatInt( int64( now.Day() ), 10 )
  if len( d ) == 1 {
    d = "0" + d
  }

  return y + "-" + m + "-" + d
}
func(el *networkTimeStt)Now() time.Time {
  return time.Now().In(el.Loc).Add(el.Dif)
}

func(el *networkTimeStt) Since(t time.Time) time.Duration {
  return el.Now().Sub(t)
}

func init() {
  NetworkTime = networkTimeStt{}

  // fixme: constante no código America/Sao_Paulo
  NetworkTime.Loc, _ = time.LoadLocation(`America/Sao_Paulo`)

  response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
  if err != nil {
    NetworkTime.Dif = 0
  } else {
    NetworkTime.Dif = response.ClockOffset
  }
}
