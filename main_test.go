package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

const (
	serverPort = "7878"
	serverUrl  = "http://localhost:" + serverPort
)

var serverCmd *exec.Cmd
var tests = []struct {
	url          string
	responseCode int
}{
	{serverUrl + "/testdata/index.html", http.StatusOK},
	{serverUrl + "/testdata/test.html", http.StatusOK},
	{serverUrl + "/testdata/dummy.pdf", http.StatusOK},
	{serverUrl + "/testdata/robots.txt", http.StatusOK},
	{serverUrl + "/testdata/snake.png", http.StatusOK},
	{serverUrl + "/testdata/snake2.jpg", http.StatusOK},
}

func TestServerRequest(t *testing.T) {
	err := startServer()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	// make http request and check response
	for _, test := range tests {
		resp, err := http.Get(test.url)
		if err != nil {
			t.Fatalf("error making request: %s", err)
		}
		if resp.StatusCode != test.responseCode {
			t.Errorf("expected response code of %d, got %d", test.responseCode, resp.StatusCode)
		}
	}
	err = stopServer()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func startServer() error {
	compileCmd := exec.Command("go", "build")
	err := compileCmd.Run()
	if err != nil {
		return fmt.Errorf("problem building fserve: %s", err)
	}
	serverCmd = exec.Command("./fserve", "-p", serverPort)
	err = serverCmd.Start()
	if err != nil {
		return fmt.Errorf("problem starting fserve: %s", err)
	}
	time.Sleep(1 * time.Second)
	return nil
}

func stopServer() error {
	if serverCmd.Process != nil {
		err := serverCmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("problem killing fserve process: %s", err)
		}
		cleanCmd := exec.Command("go", "clean")
		err = cleanCmd.Run()
		if err != nil {
			return fmt.Errorf("problem cleaning files: %s", err)
		}
		return nil
	}
	return fmt.Errorf("problem stopping fserve: process is not running")
}
