package server

import (
	"reflect"
	"testing"

	"github.com/jonathantorres/voy/internal/conf"
)

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
