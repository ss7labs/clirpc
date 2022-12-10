package main

import (
  "fmt"
  "github.com/abiosoft/ishell/v2"
  "github.com/abiosoft/readline"
  "os"
  "os/user"
  "github.com/fatih/color"
)
type customPrompt struct {
  p string
  h string
}

func main() {
  cp := &customPrompt{}
  cp.createPrompt()

  cfg := &readline.Config{Prompt: cp.p}
 	shell := ishell.NewWithConfig(cfg)
	shell.Println(cp.h)

  // start shell
	shell.Run()
	// teardown
	shell.Close()
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
  prompTmpl := d("[")+"%s"+red("@")+"%s"+d("]")
  cp.p = fmt.Sprintf(prompTmpl,username,hostname)
}
