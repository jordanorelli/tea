package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	addr := "127.0.0.1:54321"

	switch len(os.Args) {
	case 0:
		fmt.Fprintln(os.Stderr, "how did you even get here")
		os.Exit(1)
	case 1:
		break
	case 2:
		addr = os.Args[1]
	case 3:
		addr = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	default:
		fmt.Fprintln(os.Stderr, "too many args")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("listening on: %s\n", addr)
	http.Serve(lis, &server{counters: make(map[string]int)})
}
