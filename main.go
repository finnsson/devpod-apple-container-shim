package main

import (
	"fmt"
	"os"

	"github.com/finnsson/devpod-apple-container-shim/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: devpod-apple-container-shim <command> [container-id]\n")
		fmt.Fprintf(os.Stderr, "Commands: find, command, run, start, stop, delete, arch, logs\n")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	// Most commands take container-id as second arg
	containerID := ""
	if len(os.Args) > 2 {
		containerID = os.Args[2]
	}

	var err error
	switch subcommand {
	case "find":
		err = cmd.Find(containerID)
	case "command":
		err = cmd.Command(containerID)
	case "run":
		err = cmd.Run(containerID)
	case "start":
		err = cmd.Start(containerID)
	case "stop":
		err = cmd.Stop(containerID)
	case "delete":
		err = cmd.Delete(containerID)
	case "arch":
		err = cmd.Arch()
	case "logs":
		err = cmd.Logs(containerID)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", subcommand)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
