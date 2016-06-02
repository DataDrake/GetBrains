package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"github.com/DataDrake/GetBrains/core"
)

func usage() {
	fmt.Fprintln(os.Stderr,"USAGE: getbrains COMMAND [TOOL]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = func () {usage()}
	flag.Parse()
	args := flag.Args()

	var dist string
	var cmd string
	tool := "all"

	//Handle Args
	switch len(args) {
	case 2:
		tool = args[1]
		fallthrough
	case 1:
		cmd = args[0]
	default:
		usage()
		os.Exit(1)
	}

	//Handle OS
	switch runtime.GOOS {
	case "darwin":
		dist = "mac"
	case "freebsd":
		//Sorry, but this is the closest for you folks
		dist = "linux"
	default:
		dist = runtime.GOOS
	}

	//Handle Command
	switch cmd {
	case "add": fallthrough
	case "install":

	case "update": fallthrough
	case "upgrade":

	case "remove": fallthrough
	case "uninstall":

	case "info":
		core.PrintRelease(tool,dist)
	case "download":
		core.DownloadTool(tool,dist)
	case "cleanup":
		core.DownloadCleanup(tool)
	default:
		fmt.Fprintf(os.Stderr,"ERROR: Command \"%s\" is not valid\n",cmd)
		os.Exit(1)
	}

	os.Exit(0)
}
