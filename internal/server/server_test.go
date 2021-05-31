package server

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jonathantorres/voy/internal/conf"
)

func TestSimpleGetRequest(t *testing.T) {
	cmd, err := startServer("testdata/voy.conf")
	if err != nil {
		t.Fatalf("server could not be started: %s\n", err)
	}
	res, err := http.Get("http://localhost:8081")
	if err != nil {
		t.Fatalf("error sending GET request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
	err = stopServer(cmd)
	if err != nil {
		t.Fatalf("server could not be stopped: %s\n", err)
	}
}

func TestSimplePostRequest(t *testing.T) {
	cmd, err := startServer("testdata/voy.conf")
	if err != nil {
		t.Fatalf("server could not be started: %s\n", err)
	}
	postData := url.Values{}
	postData.Set("foo", "bar")
	postData.Set("baz", "zass")
	postData.Set("number", "one")
	postData.Set("name", "John")
	postData.Set("state", "Florida")

	res, err := http.PostForm("http://localhost:8081", postData)
	if err != nil {
		t.Fatalf("error sending GET request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
	err = stopServer(cmd)
	if err != nil {
		t.Fatalf("server could not be stopped: %s\n", err)
	}
}

func TestLargePostRequest(t *testing.T) {
	// TODO: send a POST request with a large amount of data
}

func TestGetPortsToListen(t *testing.T) {
	tests := []struct {
		c    *conf.Conf
		want []int
	}{
		{
			&conf.Conf{
				User:  "foo",
				Group: "bar",
				DefaultServer: &conf.ServerConf{
					Names: []string{"one", "two", "three"},
					Root:  "/foo/bar",
					Ports: []int{80, 8080, 8081},
				},
			},
			[]int{80, 8080, 8081},
		},
		{
			&conf.Conf{
				User:  "baz",
				Group: "bazzer",
				DefaultServer: &conf.ServerConf{
					Names: []string{"bee", "sting", "print"},
					Root:  "/tmp/server",
					Ports: []int{80, 443, 80, 80, 9090, 9090, 9091},
				},
			},
			[]int{80, 443, 9090, 9091},
		},
	}

	for _, test := range tests {
		ports, err := getPortsToListen(test.c)
		if err != nil {
			t.Errorf("getPortsToListen() returned an error: %s\n", err)
		}
		if !reflect.DeepEqual(ports, test.want) {
			t.Errorf("getPortsToListen() returned the wrong ports, got %v but want %v\n", ports, test.want)
		}
	}
}

func TestGetPortsToListenWithNoPorts(t *testing.T) {
	c := &conf.Conf{
		User:  "foo",
		Group: "bar",
		DefaultServer: &conf.ServerConf{
			Names: []string{"one", "two", "three"},
			Root:  "/foo/bar",
			Ports: []int{},
		},
	}
	if _, err := getPortsToListen(c); err == nil {
		t.Errorf("getPortsToListen() should return an error\n")
	}
}

func startServer(confPath string) (*exec.Cmd, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if strings.Contains(cwd, "internal/server") {
		err := os.Chdir("../../")
		if err != nil {
			return nil, err
		}
	}
	compCmd := exec.Command("go", "build")
	err = compCmd.Run()
	if err != nil {
		return nil, fmt.Errorf("problem building voy: %s", err)
	}
	serverCmd := exec.Command("./voy", "-c", confPath)
	err = serverCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("problem starting voy: %s", err)
	}
	time.Sleep(1 * time.Second) // wait a little bit so that everything is ready
	return serverCmd, nil
}

func stopServer(serverCmd *exec.Cmd) error {
	if serverCmd == nil {
		return fmt.Errorf("problem stopping voy: process is not running")
	}
	if serverCmd.Process != nil {
		err := serverCmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("problem killing voy process: %s", err)
		}
		cleanCmd := exec.Command("go", "clean")
		err = cleanCmd.Run()
		if err != nil {
			return fmt.Errorf("problem cleaning files: %s", err)
		}
		return nil
	}
	return fmt.Errorf("problem stopping voy: process is not running")
}
