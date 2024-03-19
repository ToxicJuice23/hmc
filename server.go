// TODO PUSH TO PROD!!!!

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Not enough arguments.\nUsage: %s [middleman's address:port]\n", os.Args[0])
		os.Exit(1)
	}

	middleman, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when listening for connections: %s\n", err.Error())
		os.Exit(1)
	}

	middleman.Write([]byte("host"))
	buf := make([]byte, 100)
	n, err := middleman.Read(buf)
	if n < 1 {
		fmt.Fprintf(os.Stderr, "Middleman shut us off.\nError: %s\n", err.Error())
		os.Exit(1)
	}

	// start shell
	cmd := exec.Command("bash")
	stdout, _ := cmd.StdoutPipe()
	stdin, _ := cmd.StdinPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	fmt.Println("Started bash session.")
	go func() {
		for {
			buf := make([]byte, 1024)
			_, err := middleman.Read(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Lost connection with middleman.\n")
				middleman.Close()
				os.Exit(1)
			}
			if bytes.Contains(buf, []byte("\003")) {
				cmd.Process.Kill()
				// start shell again
				cmd = exec.Command("bash")
				stdout, _ = cmd.StdoutPipe()
				stdin, _ = cmd.StdinPipe()
				stderr, _ = cmd.StderrPipe()
				cmd.Start()
				fmt.Println("Restarted bash session.")
				continue
			}
			stdin.Write(buf)
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			stderr.Read(buf)
			_, err := middleman.Write(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Lost connection with middleman.\n")
				middleman.Close()
				os.Exit(1)
			}
		}
	}()

	for {
		buf := make([]byte, 1024)
		stdout.Read(buf)
		_, err := middleman.Write(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lost connection with middleman.\n")
			middleman.Close()
			os.Exit(1)
		}
	}
}
