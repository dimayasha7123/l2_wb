package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "You must pass exactly 1 param: port\n")
		os.Exit(2)
	}

	port := flag.Arg(0)
	addr := fmt.Sprintf(":%s", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen %s: %v\n", addr, err)
		os.Exit(1)
	}
	fmt.Printf("Start listening at %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't accept conn: %v\n", err)
			continue
		}

		fmt.Printf("Accept conn at %s\n", conn.RemoteAddr())

		go func(conn net.Conn) {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				fmt.Printf("%s: %s\n", conn.RemoteAddr(), scanner.Text())
			}

			fmt.Printf("Conn %s was closed\n", conn.RemoteAddr())
		}(conn)
	}
}
