package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var defaultPort = 9090
var port = flag.Int("p", defaultPort, "server port")
var showVersion = flag.Bool("v", false, "print server version")

func main() {
	flag.Parse()
	if *showVersion {
		printVersion()
	}
	fmt.Printf("fserve running on port %d\n", *port)

	addr := fmt.Sprintf("localhost:%d", *port)
	handler := ServerHandler{}
	log.Fatal(http.ListenAndServe(addr, &handler))
}

type ServerHandler struct{}

func (handler *ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	setDefaultHeaders(res)
	file, err := serveFile(req.URL)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintf(res, "error based response")
	}
	fmt.Fprintf(res, file)
}

func serveFile(url *url.URL) (*os.File, error) {
	filename := url.Path
	if filename[0] == '/' {
		filename = filename[1:]
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func setDefaultHeaders(res http.ResponseWriter) {
	res.Header().Set("Server", nameAndVersion)
	res.Header().Set("Connection", "close")
}
