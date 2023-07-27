package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	timeout := time.Second * 10
	retryTimeout := time.Millisecond * 250
	flag.DurationVar(&timeout, "timeout", timeout, "connection timeout to server")
	flag.DurationVar(&retryTimeout, "retry_timeout", retryTimeout, "timeout between attempts to connect")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "You must pass exactly 2 params: host and port\n")
		os.Exit(2)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	addr := fmt.Sprintf("%s:%s", host, port)

	var err error
	var conn net.Conn
	var good bool
	deadline := time.After(timeout)

	for !good {
		select {
		case <-deadline:
			fmt.Fprintf(os.Stderr, "Can't connect to %s: %v\n", addr, err)
			os.Exit(1)
		default:
			conn, err = net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't connect to server... Retry after %v\n", retryTimeout)
				time.Sleep(retryTimeout)
				continue
			}

			good = true
		}
	}
	defer conn.Close()

	fmt.Printf("Connect to %s\n", conn.LocalAddr())

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("Exit\n")
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Can't read string: %v\n", err)
			break
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't send %s: %v\n", input, err)
			continue
		}
	}
}
