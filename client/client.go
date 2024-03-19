package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Not enough arguments.\nUsage: %s [server's address:port]\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured while connecting to the server.\n%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Dialing: 127.0.0.1:8080...")

	conn.Write([]byte("client"))
	n, err := conn.Read(make([]byte, 100))
	if n < 1 {
		fmt.Fprintf(os.Stderr, "Server shut us off.\nError: %s\n", err.Error())
		os.Exit(1)
	}

	bufSize := 2048
	go func() {
		for {
			buf := make([]byte, bufSize)
			_, err := conn.Read(buf)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on recv().\n")
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "%s", buf)
		}
	}()

	go func() {
		for {
			<-c
			char := byte(3)
			msg := make([]byte, 1)
			msg = append(msg, char)
			conn.Write(msg)
		}
	}()

	for {
		buf2 := make([]byte, bufSize)
		fmt.Print("> ")
		os.Stdin.Read(buf2)
		_, err := conn.Write(buf2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on write().\n")
			os.Exit(1)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
