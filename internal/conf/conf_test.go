package conf

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
	defServer := ServerConf{
		name:       "localhost",
		root:       "/var/www/localhost",
		ports:      nil,
		indexPages: nil,
		errorPages: nil,
		errorLog:   "/etc/log/voy/errors.log",
		accessLog:  "/etc/log/voy/access.log",
	}
	vhosts := []ServerConf{
		{
			name:       "mydomain.com",
			root:       "/var/www/mydomain.com/public",
			ports:      nil,
			indexPages: nil,
			errorPages: nil,
			errorLog:   "/etc/log/voy/mydomain.com.log",
			accessLog:  "/etc/log/voy/mydomain.com.log",
		},
	}
	wantConf := Conf{
		user:          "www-data",
		group:         "www-data",
		defaultServer: &defServer,
		vhosts:        vhosts,
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
