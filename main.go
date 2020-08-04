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
var useLogFile = flag.Bool("l", false, "save logs in a file (fserve.log)")
var fileLogger *log.Logger = nil

func main() {
	flag.Parse()
	if *showVersion {
		printVersion()
	}
	fmt.Printf("fserve running on port %d\n", *port)
	if *useLogFile {
		registerLogger();
	}

	addr := fmt.Sprintf("localhost:%d", *port)
	handler := ServerHandler{}
	err := http.ListenAndServe(addr, &handler)
	if err != nil {
		if *useLogFile {
			fileLogger.Printf("server error: %s", err)
		}
		log.Fatal(err)
	}
}

type ServerHandler struct{}

func (handler *ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	setDefaultHeaders(res)
	file, err := findFile(req.URL)
	if err != nil {
		writeErrorResponse(res, req, http.StatusNotFound, err.Error())
		return
	}
	fileinfo, err := file.Stat()
	if err != nil {
		writeErrorResponse(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	extPieces := strings.Split(fileinfo.Name(), ".")
	ext := extPieces[len(extPieces)-1]
	fileType := contentTypes[ext]
	res.Header().Set("Content-type", fileType.contentType)

	if _, err = io.Copy(res, file); err != nil {
		writeErrorResponse(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("%s %s %s %d", req.Method, req.URL.Path, req.Proto, http.StatusOK)
	if *useLogFile {
		fileLogger.Printf("%s %s %s %d", req.Method, req.URL.Path, req.Proto, http.StatusOK)
	}
}

const (
	FileTypeText = iota
	FileTypeBinary
)

type FileType struct {
	contentType string
	fileType    uint
}

var contentTypes = map[string]FileType{
	"html": FileType{"text/html", FileTypeText},
	"htm": FileType{"text/html", FileTypeText},
	"css": FileType{"text/css", FileTypeText},
	"md": FileType{"text/markdown", FileTypeText},
	"txt": FileType{"text/plain", FileTypeText},
	"xml": FileType{"text/xml", FileTypeText},
	"js": FileType{"application/javascript", FileTypeText},
	"json": FileType{"application/json", FileTypeText},
	"pdf": FileType{"application/pdf", FileTypeBinary},
	"zip": FileType{"application/zip", FileTypeBinary},
	"bmp": FileType{"image/bmp", FileTypeBinary},
	"gif": FileType{"image/gif", FileTypeBinary},
	"jpg": FileType{"image/jpeg", FileTypeBinary},
	"jpeg": FileType{"image/jpeg", FileTypeBinary},
	"ico": FileType{"image/x-icon", FileTypeBinary},
	"png": FileType{"image/png", FileTypeBinary},
	"tiff": FileType{"image/tiff", FileTypeBinary},
	"svg": FileType{"image/svg", FileTypeText},
	"mp3": FileType{"audio/mp3", FileTypeBinary},
	"mp4": FileType{"audio/mp4", FileTypeBinary},
	// "mp4": FileType{"video/mp4", FileTypeBinary},
	"mpeg": FileType{"audio/mpeg", FileTypeBinary},
	// "mpeg": FileType{"video/mpeg", FileTypeBinary},
	"ogg": FileType{"audio/ogg", FileTypeBinary},
	// "ogg": FileType{"video/ogg", FileTypeBinary},
	"quicktime": FileType{"video/quicktime", FileTypeBinary},
	"ttf": FileType{"font/ttf", FileTypeBinary},
	"woff": FileType{"font/woff", FileTypeBinary},
	"woff2": FileType{"font/woff2", FileTypeBinary},
}

func validateExtension(ext string) bool {
	_, ok := contentTypes[ext];
	return ok
}

func findFile(url *url.URL) (*os.File, error) {
	filepath := url.Path
	if filepath[0] == '/' {
		filepath = filepath[1:]
	}
	if filepath == "" {
		filepath = "index.html"
	} else if strings.HasSuffix(filepath, "/") {
		filepath = filepath + "index.html"
	}
	extPieces := strings.Split(filepath, ".")
	ext := extPieces[len(extPieces)-1]
	if ok := validateExtension(ext); !ok {
		return nil, fmt.Errorf("file extension %s is not suppored", ext)
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func registerLogger() {
	logfile, err := os.Create("fserve.log")
	if err != nil {
		log.Printf("error creating log file: %s", err)
		return
	}
	fileLogger = log.New(logfile, "", log.LstdFlags)
}

func writeErrorResponse(res http.ResponseWriter, req *http.Request, statusCode int, msg string) {
	res.Header().Set("Content-type", "text/html")
	res.WriteHeader(statusCode)
	fmt.Fprintf(res, "%s", msg)
	log.Printf("%s %s %s %d", req.Method, req.URL.Path, req.Proto, statusCode)
	if *useLogFile {
		fileLogger.Printf("%s %s %s %d", req.Method, req.URL.Path, req.Proto, statusCode)
	}
}

func setDefaultHeaders(res http.ResponseWriter) {
	res.Header().Set("Server", nameAndVersion)
	res.Header().Set("Connection", "close")
}
