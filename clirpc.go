package clirpc

import (
  "time"
  "fmt"
  "os/exec"
)
const (
  DefPort = "58085"
)

type RawSession struct {
	vif_id,
	username,
	mac,
	port,
	svid,
	cvid,
	session_id,
	ip_addr,
	mtu,
	ingress_cir,
	egress_cir,
	rx_pkts,
	tx_pkts,
	rx_bytes,
	tx_bytes,
	uptime string
}
type Rcli interface {
  FindUserSession(user string) (RawSession,error)
  DiscUserSession(user string) error
}

type Listener struct {
	Sleep time.Duration
}

func (l *Listener) GetUser(line []byte, ack *bool) (err error) {
  user:= string(line)
  cmd := exec.Command("./sh_user.sh", user)
	out, err := cmd.CombinedOutput()

  if err != nil {
     fmt.Println(err.Error())
     return
  }
 	if len(out) == 0 {
		return
	}

	return
}

