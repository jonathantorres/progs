package conf

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jonathantorres/voy/internal/http"
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

func Load(confFile string) (*Conf, error) {
	file, err := openAndStripComments(confFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	file, err = parseIncludes(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = checkForSyntaxErrors(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conf, err := buildServerConf(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conf, nil
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
		s.parsePortOptions(opValue)
	case indexOption:
		s.parseIndexOptions(opValue)
	case errorLogOption:
		s.errorLog = opValue
	case accessLogOption:
		s.accessLog = opValue
	}

	// handle error pages
	if strings.Contains(opName, errorPageOption) {
		s.parseErrorPageOptions(opName, opValue)
	}
}

func (s *ServerConf) parsePortOptions(ports string) {
	portsStr := strings.Split(ports, ",")
	if len(portsStr) == 0 {
		return
	}
	s.ports = make([]int, 0)
	for _, p := range portsStr {
		pInt, err := strconv.Atoi(p)
		if err != nil {
			log.Println(err)
			return
		}
		s.ports = append(s.ports, pInt)
	}
}

func (s *ServerConf) parseIndexOptions(pages string) {
	pagesSli := strings.Split(pages, ",")
	if len(pagesSli) == 0 {
		return
	}
	s.indexPages = make([]string, 0)
	for _, p := range pagesSli {
		s.indexPages = append(s.indexPages, p)
	}
}

func (s *ServerConf) parseErrorPageOptions(errorType, page string) {
	eTypePieces := strings.Split(errorType, "_")
	// unlikely, but just in case
	if len(eTypePieces) == 1 {
		// TODO: maybe log this? or the function that checks the conf file should detect this?
		return
	}
	if s.errorPages == nil {
		s.errorPages = make([]ErrorPage, 0)
	}
	var errPage ErrorPage
	if len(eTypePieces) == 3 {
		// custom error page
		code, err := strconv.Atoi(eTypePieces[2])
		if err != nil {
			return
		}
		errPage = ErrorPage{
			code: code,
			page: strings.TrimSpace(page),
		}
	}
	if len(eTypePieces) == 2 {
		errPage = ErrorPage{
			code: http.StatusBadRequest,
			page: strings.TrimSpace(page),
		}
	}
	s.errorPages = append(s.errorPages, errPage)
}

func buildServerConf(file []byte) (*Conf, error) {
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
