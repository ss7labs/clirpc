package main

import (
	"clirpc"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type customPrompt struct {
	p string
	h string
}

type Remotes struct {
	rmt map[string]string
	rs  *clirpc.RawSession
}

func (r *Remotes) Init() {
	r.rmt = make(map[string]string)
	r.rmt["55-bras"] = "10.19.176.55"
	r.rmt["bng-a46"] = "10.20.14.15"
	r.rmt["bng-a33"] = "10.20.221.219"
	r.rmt["bng-a34"] = "10.168.42.35"
	r.rmt["bng-a27"] = "10.20.55.237"
}

func (r *Remotes) showUser(c *ishell.Context) {
	var args []string
	args = c.Args
	if err := checkUserMac(args[0]); err != nil {
		c.Err(err)
		return
	}
	r.RemoteCall("show", args[0])
	if r.rs == nil {
		return
	}

	blocked := false
	if r.rs.IngressCir == "-" {
		blocked = true
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	boldGreen := color.New(color.FgGreen, color.Bold).SprintFunc()

	if blocked {
		boldRed := color.New(color.FgRed, color.Bold).SprintFunc()
		c.Printf("%s %s %s %s%s %s%s %s\n", r.rs.Username, yellow(r.rs.Mac), r.rs.IpAddr, boldGreen("SVID:"), r.rs.Svid, boldGreen("CVID:"), r.rs.Cvid, boldRed("BLOCKED"))
	} else {
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		wtB := color.New(color.FgWhite, color.Bold).SprintFunc()
		wtHi := color.New(color.FgHiWhite, color.Bold).SprintFunc()
		up := strings.Split(r.rs.IngressCir, ";")
		upTT := up[1]
		dn := strings.Split(r.rs.EgressCir, ";")
		outFmt := "%s %s %s %s %s%s %s%s %s%s%s %s%s%s %s%s\n"
		c.Printf(outFmt, magenta(r.rs.Host), wtHi(r.rs.Username), yellow(r.rs.Mac), r.rs.IpAddr, boldGreen("SVID:"), r.rs.Svid, boldGreen("CVID:"), r.rs.Cvid, boldGreen("UP:"), wtB(up[0])+";", cyan(upTT+"(TT)"), boldGreen("DN:"), wtB(dn[0])+";", cyan(dn[1]+"(TT)"), yellow("MTU:"), wtHi(r.rs.Mtu))
	}
	r.rs = nil
}

func (r *Remotes) RemoteCall(cmd, user string) {
	var err error
	for h, ip := range r.rmt {
		if cmd == "show" {
			err = r.shUser(ip, user)
			if err == nil {
				r.rs.Host = h
				break
			}
		}
		if cmd == "disc" {
			err = discUser(ip, user)
			if err == nil {
				break
			}
		}
	}
}
func checkHostAvail(host string) error {
	var rcvd int
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return err
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 1
	pinger.OnFinish = func(stats *ping.Statistics) {
		/*
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		*/
		log.Trace().Msg(stats.Addr + " packets received " + strconv.Itoa(stats.PacketsRecv))
		rcvd = stats.PacketsRecv
	}

	err = pinger.Run()
	if err != nil {
		return err
	}
	//fmt.Println("PING ", err, pinger.IPAddr(),rcvd)
	if rcvd == 0 {
		err = errors.New("Host not available")
		return err
	}
	return nil
}
func (r *Remotes) shUser(ip, user string) error {
	err := checkHostAvail(ip)

	if err != nil {
		return err
	}
	srvAddr := ip + ":" + clirpc.DefPort
	client, err := rpc.Dial("tcp", srvAddr)
	if err != nil {
		log.Trace().Msg(srvAddr + " Dial failed")
		return err
	}
	defer client.Close()
	reply := new(clirpc.RawSession)
	var line []byte
	line = []byte(user)
	showCall := client.Go("Listener.GetUser", line, reply, nil)
	<-showCall.Done
	if showCall.Error != nil {
		return showCall.Error
	}
	r.rs = reply
	return nil
}

func discUser(ip, user string) error {
	srvAddr := ip + ":" + clirpc.DefPort
	client, err := rpc.Dial("tcp", srvAddr)
	if err != nil {
		return err
	}
	defer client.Close()
	reply := false
	var line []byte
	line = []byte(user)
	call := client.Go("Listener.DiscUser", line, &reply, nil)
	<-call.Done
	if call.Error != nil {
		return call.Error
	}
	return nil
}

func (r *Remotes) discUser(c *ishell.Context) {
	var args []string
	args = c.Args
	if err := checkUserMac(args[0]); err != nil {
		c.Err(err)
		return
	}
	r.RemoteCall("disc", args[0])
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	//tcp connections will be permanent
	rmt := &Remotes{}
	rmt.Init()

	cp := &customPrompt{}
	cp.createPrompt()

	cfg := &readline.Config{Prompt: cp.p}
	shell := ishell.NewWithConfig(cfg)

	addCommands(shell, rmt)

	if len(os.Args) > 1 && os.Args[1] == "--trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	// when started with "exit" as first argument, assume non-interactive execution
	if len(os.Args) > 1 && os.Args[1] == "exit" {
		shell.Process(os.Args[2:]...)
	} else if len(os.Args) > 2 && os.Args[2] == "exit" {
		shell.Process(os.Args[3:]...)
	} else {
		shell.Println(cp.h)
		// start shell
		shell.Run()
		// teardown
		shell.Close()
	}
}

func addCommands(sh *ishell.Shell, r *Remotes) {
	showCmd := &ishell.Cmd{
		Name: "show",
		Help: "Show running system information",
	}
	showCmd.AddCmd(&ishell.Cmd{
		Name: "mac",
		Help: "Show user info by mac address",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("no login/username"))
				return
			}
			r.showUser(c)
		},
	})
	showCmd.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "Show user info by login/username",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("no login/username"))
				return
			}
			r.showUser(c)
		},
	})
	showCmd.AddCmd(&ishell.Cmd{
		Name: "version",
		Help: "Show software version",
	})
	sh.AddCmd(showCmd)

	discCmd := &ishell.Cmd{
		Name: "disconnect",
		Help: "Disconnect session",
	}
	discCmd.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "Disconnect session by login/username",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("no login/username"))
				return
			}
			r.discUser(c)
		},
	})
	discCmd.AddCmd(&ishell.Cmd{
		Name: "mac",
		Help: "Disconnect session by mac address",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("no mac address"))
				return
			}
			r.discUser(c)
		},
	})
	sh.AddCmd(discCmd)
}

var unamePatt *regexp.Regexp = regexp.MustCompile("^[0-9]{6}@gts$|^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})|(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")

func checkUserMac(str string) error {
	if !unamePatt.MatchString(str) {
		return errors.New("user should be like as 490036@gts or 3D:F2:C9:A6:B3:4F")
	}
	return nil
}

func (cp *customPrompt) createPrompt() {
	info := color.New(color.FgBlack, color.BgWhite).SprintFunc()
	cp.h = fmt.Sprintf("Type %s for commands list", info("help"))

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	username := "agts"
	user, err := user.Current()
	if err == nil {
		username = user.Username
	}

	d := color.New(color.FgCyan, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	prompTmpl := d("[") + "%s" + red("@") + "%s" + d("]")
	cp.p = fmt.Sprintf(prompTmpl, username, hostname)
}
