package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCommentsAreStripped(t *testing.T) {
	cases := []struct {
		confFile     string
		confFileWant string
	}{
		{"testdata/conf1.txt", "testdata/conf1_want.txt"},
		{"testdata/conf2.txt", "testdata/conf2_want.txt"},
	}

	for _, c := range cases {
		file, err := openAndStripComments(c.confFile)
		if err != nil {
			t.Errorf("error reading conf file %s", err)
		}
		want, err := ioutil.ReadFile(c.confFileWant)
		if err != nil {
			t.Errorf("error reading conf_want file %s", err)
		}
		if !bytes.Equal(file, want) {
			t.Errorf("bytes in file are not equal")
		}
	}
}

func TestIncludedFilesAreParsed(t *testing.T) {
	t.Skipf("todo")
}

func TestNoSyntaxErrorsAreFound(t *testing.T) {
	t.Skipf("todo")
}

func TestSyntaxErrorsAreFound(t *testing.T) {
	t.Skipf("todo")
}

func TestServerConfIsBuilt(t *testing.T) {
	wantConf := Conf{
		User:    "www-data",
		Group:   "www-data",
		Workers: 5,
		DefaultServer: &ServerConf{
			Names:      []string{"localhost"},
			Root:       "/var/www/localhost",
			Ports:      []int{80, 443},
			IndexPages: []string{"index.html", "index.htm"},
			ErrorPages: []ErrorPage{
				{
					Code: 400,
					Page: "error.html",
				},
				{
					Code: 404,
					Page: "404.html",
				},
			},
			ErrorLog:  "/etc/log/httpd/errors.log",
			AccessLog: "/etc/log/httpd/access.log",
		},
		Vhosts: []ServerConf{
			{
				Names:      []string{"mydomain.com"},
				Root:       "/var/www/mydomain.com/public",
				Ports:      []int{8081},
				IndexPages: []string{"index.html"},
				ErrorPages: []ErrorPage{
					{
						Code: 400,
						Page: "error.html",
					},
				},
				ErrorLog:  "/etc/log/httpd/mydomain.com.log",
				AccessLog: "/etc/log/httpd/mydomain.com.log",
			},
			{
				Names:      []string{"example.com", "www.example.com"},
				Root:       "/var/www/mydomain.com/public",
				Ports:      []int{8081},
				IndexPages: []string{"index.html"},
				ErrorPages: []ErrorPage{
					{
						Code: 400,
						Page: "error.html",
					},
				},
				ErrorLog:  "/etc/log/httpd/mydomain.com.log",
				AccessLog: "/etc/log/httpd/mydomain.com.log",
			},
			{
				Names:      []string{"test.com", "www.test.com"},
				Root:       "/var/www/mydomain.com/public",
				Ports:      []int{8081},
				IndexPages: []string{"index.html"},
				ErrorPages: []ErrorPage{
					{
						Code: 400,
						Page: "error.html",
					},
				},
				ErrorLog:  "/etc/log/httpd/mydomain.com.log",
				AccessLog: "/etc/log/httpd/mydomain.com.log",
			},
		},
	}
	cases := []struct {
		confFile string
		want     Conf
	}{
		{"testdata/conf3.txt", wantConf},
	}

	for _, c := range cases {
		file, err := ioutil.ReadFile(c.confFile)
		if err != nil {
			t.Errorf("%s", err)
		}
		conf, err := buildServerConf(file)
		if err != nil {
			t.Errorf("%s", err)
		}
		if !reflect.DeepEqual(*conf, c.want) {
			t.Errorf("the configurations are not equal")
		}
	}
}
