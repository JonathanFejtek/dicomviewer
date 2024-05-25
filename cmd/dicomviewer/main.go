package main

import (
	"dicomviewer/http"
	"flag"
)

type cliArgs struct {
	serverPort string
}

func parseCLIargs() cliArgs {
	portPtr := flag.String("port", http.DefaultPort, "listen port for server")

	flag.Parse()
	return cliArgs{
		serverPort: *portPtr,
	}
}

func main() {
	args := parseCLIargs()

	service := http.NewServer(
		http.UsePort(
			args.serverPort,
		),
	)

	service.ListenAndServe()
}
