package server

import (
	"fmt"
	"io/ioutil"
	"log"
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

func TestMain(m *testing.M) {
	cmd, err := startServer("testdata/voy.conf")
	if err != nil {
		log.Fatalf("server could not be started: %s\n", err)
	}
	m.Run()
	err = stopServer(cmd)
	if err != nil {
		log.Fatalf("server could not be stopped: %s\n", err)
	}
}

func TestSimpleGetRequest(t *testing.T) {
	res, err := http.Get("http://localhost:8081")
	if err != nil {
		t.Fatalf("error sending GET request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
}

func TestSimplePostRequest(t *testing.T) {
	postData := url.Values{}
	postData.Set("foo", "bar")
	postData.Set("baz", "zass")
	postData.Set("number", "one")
	postData.Set("name", "John")
	postData.Set("state", "Florida")

	res, err := http.PostForm("http://localhost:8081", postData)
	if err != nil {
		t.Fatalf("error sending POST request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
}

func TestLargePostRequest(t *testing.T) {
	f, err := os.Open("internal/server/testdata/post_data.txt")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	res, err := http.Post("http://localhost:8081", "text/plain", f)
	if err != nil {
		t.Fatalf("error sending POST request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
}

func TestLargeRequestLine(t *testing.T) {
	url := `http://localhost:8081/vod/FetchOrderDetails?itemId=233756825167&transactionId=1921811535013&mkevt=1&mkpid=0&emsid=e11401.m43700.l49689&mkcid=7&ch=osgood&euid=cd3dbb358e3b4633b21e32c5e7b1ded2&bu=43783229363&exe=98631&ext=232562&some1=43783229363&test1=1234567789898232&tryid=12938129381293812938&console=912839812938123&logid=nqt%3DAAAAEAAAACAgAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAABAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAgAAAQAAAAAAAACAAAAgAAAAAgAAAAAAAAAAAAAAAgA**%26nqc%3DAAAAEAAAACAgAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAABAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAgAAAQAAAAAAAACAAAAgAAAAAgAAAAAAAAAAAAAAAgA**%26mdbreftime%3D1622479417918%26es%3D0%26ec%3D1&osub=-1~1&crd=20210531095122&segname=11401&sojTags=ch%3Dch%2Cbu%3Dbu%2Cnqt%3Dnqt%2Cnqc%3Dnqc%2Cmdbreftime%3Dmdbreftime%2Ces%3Des%2Cec%3Dec%2Cexe%3Dexe%2Cext%3Dext%2Cexe%3Dexe%2Cext%3Dext%2Cosub%3Dosub%2Ccrd%3Dcrd%2Csegname%3Dsegname%2Cchnl%3Dmkcid`
	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("error sending GET request: %s\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected a 200 response, got: %d\n", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s\n", err)
	}
	if string(b) != "Hello, world" {
		t.Fatalf("incorrect response body, got %s but want %s\n", string(b), "Hello, world")
	}
	if res.ContentLength != int64(len("Hello, world")) {
		t.Fatalf("content length of response does not match, got %d want %d\n", res.ContentLength, len("Hello, world"))
	}
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
