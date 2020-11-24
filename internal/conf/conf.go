package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

// This is supposed to validate, test and parse the configuration file

type Conf struct {
	user          string
	group         string
	defaultServer *ServerConf
	vhosts        []ServerConf
}

type ServerConf struct {
	name       string
	root       string
	ports      []int
	indexPages []string
	errorPages []ErrorPage
	errorLog   string
	accessLog  string
}

type ErrorPage struct {
	code int
	page string
}

// TODO: must come from a standard location
// or specified as a command line param
var confFile = "./voy.conf"

func Load() error {
	file, err := openAndStripComments(confFile)
	if err != nil {
		log.Println(err)
		return err
	}
	file, err = parseIncludes(file)
	if err != nil {
		log.Println(err)
		return err
	}
	err = checkForSyntaxErrors(file)
	if err != nil {
		log.Println(err)
		return err
	}
	conf, err := buildServerConf(file)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(conf)
	return nil
}

func (c *Conf) addOption(opName string, opValue string) {
	if c.defaultServer == nil {
		c.defaultServer = &ServerConf{}
	}
	switch opName {
	case userOption:
		c.user = opValue
	case groupOption:
		c.group = opValue
	default:
		c.defaultServer.addOption(opName, opValue)
	}
}

func (c *Conf) addVhost(vhost ServerConf) {
	if c.vhosts == nil {
		c.vhosts = make([]ServerConf, 0, 10)
	}
	c.vhosts = append(c.vhosts, vhost)
}

func (s *ServerConf) addOption(opName string, opValue string) {
	switch opName {
	case nameOption:
		s.name = opValue
	case rootOption:
		s.root = opValue
	case portOption:
		s.ports = parsePortOptions(opValue)
	case indexOption:
		s.indexPages = parseIndexOptions(opValue)
	case errorPageOption:
		// TODO: handle dinamic error page types (404, 501 etc. etc.)
		s.errorPages = parseErrorPageOptions(opValue)
	case errorLogOption:
		s.errorLog = opValue
	case accessLogOption:
		s.accessLog = opValue
	}
}

func buildServerConf(file []byte) (*Conf, error) {
	// TODO: build the Conf structure based on the correctly
	// loaded configuration file
	r := bytes.NewReader(file)
	scanner := bufio.NewScanner(r)
	insideVhost := false
	conf := &Conf{}
	var curVhost *ServerConf
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.ContainsRune(line, equalSign) {
			// this is a line with an option
			ops := bytes.Split(line, []byte{byte(equalSign)})
			opName := string(bytes.TrimSpace(ops[0]))
			opValue := string(bytes.TrimSpace(ops[1]))
			if insideVhost {
				// option for the current virtual host
				curVhost.addOption(opName, opValue)
			} else {
				// top level or global option
				conf.addOption(opName, opValue)
			}
		} else if bytes.Contains(line, []byte(vhostOption)) {
			// this is a line with a vhost command
			insideVhost = true
			curVhost = &ServerConf{}
		} else if bytes.ContainsRune(line, closingBracket) {
			// closing bracket for a vhost command
			if insideVhost {
				conf.addVhost(*curVhost)
				curVhost = nil
				insideVhost = false
			}
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conf, nil
}

func parseIncludes(file []byte) ([]byte, error) {
	// TODO: this should load and expand every single "include"
	// in the configuration file
	return file, nil
}

func checkForSyntaxErrors(file []byte) error {
	// TODO: check for syntax errors in the file:
	// - Symbols that are not recognized
	// - Options that are unknown
	// - Opened brackets that are not closed (and viceversa)
	return nil
}

func openAndStripComments(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()
	file := make([]byte, 0)
	scanner := bufio.NewScanner(f)
scan:
	for scanner.Scan() {
		line := scanner.Bytes()
		foundComment := false
		for i, b := range line {
			if rune(b) == commentSign {
				foundComment = true
			}
			if foundComment {
				if i != 0 {
					file = append(file, byte('\n'))
				}
				continue scan
			} else {
				file = append(file, b)
			}
		}
		file = append(file, byte('\n'))
	}
	err = scanner.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return file, nil
}

func parsePortOptions(ports string) []int {
	return nil
}

func parseIndexOptions(pages string) []string {
	return nil
}

func parseErrorPageOptions(pages string) []ErrorPage {
	return nil
}
