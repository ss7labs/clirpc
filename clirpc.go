package clirpc

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	DefPort  = "58085"
	showRcli = "/usr/local/bin/showuser"
	discRcli = "/usr/local/bin/discuser"
	rcli     = "ip netns exec tr /usr/sbin/rcli"
)

type RawSession struct {
	VifId,
	Username,
	Mac,
	Port,
	Svid,
	Cvid,
	SessionId,
	IpAddr,
	Mtu,
	IngressCir,
	EgressCir,
	RxPkts,
	TxPkts,
	RxBytes,
	TxBytes,
	Uptime,
	Host string
}
type Rcli interface {
	FindUserSession(user string) (RawSession, error)
	DiscUserSession(user string) error
}

type Listener struct {
	Sleep time.Duration
}

func (l *Listener) GetUser(line []byte, reply *RawSession) (err error) {
	user := string(line)
	args := []string{user}
	cmd := exec.Command(showRcli, args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("GetUser Command", err.Error(), string(out))
		return
	}
	if len(out) == 0 {
		return
	}
	s := strings.Split(string(out), "\t")
	reply.VifId = s[0]
	reply.Username = s[1]
	reply.Mac = s[2]
	reply.Port = s[3]
	reply.Svid = s[4]
	reply.Cvid = s[5]
	reply.SessionId = s[6]
	reply.IpAddr = s[7]
	reply.Mtu = s[8]
	reply.IngressCir = s[9]
	reply.EgressCir = s[10]
	reply.RxPkts = s[11]
	reply.TxPkts = s[12]
	reply.RxBytes = s[13]
	reply.TxBytes = s[14]
	reply.Uptime = s[15]
	return
}

func (l *Listener) DiscUser(line []byte, ack *bool) (err error) {
	user := string(line)
	args := []string{user}
	cmd := exec.Command(discRcli, args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Disc Command", err.Error(), string(out))
		return
	}
	if len(out) == 0 {
		return
	}
	*ack = true
	return
}
