package clirpc

import (
	_ "fmt"
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
