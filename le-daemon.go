package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"github.com/bsphere/le_go"
)

var (
	le  *le_go.Logger
	cmd *exec.Cmd
)

func msg(level, message string) {
	if level == "error" {
		fmt.Fprintln(os.Stderr, message)
	} else {
		fmt.Fprintln(os.Stdout, message)
	}
	if le != nil {
		le.Println(level, message)
	}
}

func checkError(err error) {
  if err != nil {
  	msg("error", fmt.Sprintf("Error: %s", err))
    os.Exit(3)
  }
}

func dumpToLe(retch chan bool, in io.ReadCloser, level string) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		msg(level, scanner.Text())
	}
	retch <- true
}

func main() {
	var tokenflag string

	flag.StringVar(&tokenflag, "token", "", "-token=<logentries_token>")

	flag.Parse()

	if tokenflag == "" || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: le-daemon -token=<logentries_token> [program with args to execute] ")
		flag.PrintDefaults()
		os.Exit(2)
	}

	var err error

	le, err = le_go.Connect(tokenflag)
	checkError(err)

	defer le.Close()

	cmd = exec.Command(flag.Args()[0], flag.Args()[1:]...)

	stdout, err := cmd.StdoutPipe()
  checkError(err)
  
  stderr, err := cmd.StderrPipe()
  checkError(err)

  retch := make(chan bool)

  go dumpToLe(retch, stdout, "info")
  go dumpToLe(retch, stderr, "error")

  msg("info", fmt.Sprintf("Started%v", flag.Args()))

  err = cmd.Run()
  checkError(err)

  <-retch
  <-retch

	exitCode := 1
	level := "error"
	exited := fmt.Sprintf("Exited(%s)", cmd.ProcessState.String())

	if (cmd.ProcessState.Success()) {
		level = "info"
		exitCode = 0
	}

	msg(level, exited)

	os.Exit(exitCode)
}
