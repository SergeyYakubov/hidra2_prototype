package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/sergeyyakubov/dcomp/dcomp/cli"
	"github.com/sergeyyakubov/dcomp/dcomp/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

var (
	flHelp = flag.Bool("help", false, "Print usage")
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dComp"); ret {
		return
	}

	flag.Parse()

	if *flHelp || flag.NArg() == 0 {
		flag.Usage()
		cli.PrintAllCommands()
		return
	}

	if flag.Arg(0) == "daemon" {
		daemon.Start(flag.Args()[1:])
	} else {
		if err := cli.SetDaemonConfiguration(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := cli.DoCommand(flag.Arg(0), flag.Args()[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
