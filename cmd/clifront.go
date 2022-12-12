package main

import (
	"errors"
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"os"
	"os/user"
  "regexp"
  "net/rpc"
)

type customPrompt struct {
	p string
	h string
}

type Remotes struct {
  conns []*rpc.Client
}
func (r *Remotes) Connect(){
}

func main() {
//tcp connections will be permanent
  rmt := &Remotes{}
  rmt.Connect()

  cp := &customPrompt{}
	cp.createPrompt()

	cfg := &readline.Config{Prompt: cp.p}
	shell := ishell.NewWithConfig(cfg)

	addCommands(shell,rmt)

  // when started with "exit" as first argument, assume non-interactive execution
	if len(os.Args) > 1 && os.Args[1] == "exit" {
		shell.Process(os.Args[2:]...)
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
var unamePatt *regexp.Regexp = regexp.MustCompile("^[0-9]{6}@gts$|^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$")

func checkUserMac(str string) error {
  if !unamePatt.MatchString(str) {
    return errors.New("user should be like as 490036@gts or 3D:F2:C9:A6:B3:4F")
  }
  return nil
}

func (r *Remotes) showUser(c *ishell.Context){
  var args []string
  args = c.Args
  if err := checkUserMac(args[0]); err != nil {
		c.Err(err)
		return 
	}
}

func (r *Remotes) discUser(c *ishell.Context){
  var args []string
  args = c.Args
  if err := checkUserMac(args[0]); err != nil {
		c.Err(err)
		return 
	}
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
