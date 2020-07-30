package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	file, err := findFile(req.URL)
	if err != nil {
		writeErrorResponse(res, http.StatusNotFound, err.Error())
		return
	}
	if _, err = io.Copy(res, file); err != nil {
		writeErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}
}

func findFile(url *url.URL) (*os.File, error) {
	filepath := url.Path
	if filepath[0] == '/' {
		filepath = filepath[1:]
	}
	if filepath == "" {
		filepath = "index.html"
	} else if strings.HasSuffix(filepath, "/") {
		filepath = filepath+"index.html"
	}
	// TODO: validate the file extension
	// the extension should be supported
	// return error if the extension is not supported
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func writeErrorResponse(res http.ResponseWriter, statusCode int, msg string) {
	res.Header().Set("Content-type", "text/html")
	res.WriteHeader(statusCode)
	fmt.Fprintf(res, "%s", msg)
}

func setDefaultHeaders(res http.ResponseWriter) {
	res.Header().Set("Server", nameAndVersion)
	res.Header().Set("Connection", "close")
}
