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
	User          string
	Group         string
	DefaultServer *ServerConf
	Vhosts        []ServerConf
	Workers       int
}

type ServerConf struct {
	Names      []string
	Root       string
	Ports      []int
	IndexPages []string
	ErrorPages []ErrorPage
	ErrorLog   string
	AccessLog  string
}

type ErrorPage struct {
	Code int
	Page string
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
	if c.DefaultServer == nil {
		c.DefaultServer = &ServerConf{}
	}
	switch opName {
	case userOption:
		c.User = opValue
	case groupOption:
		c.Group = opValue
	case workersOption:
		w, _ := strconv.Atoi(opValue)
		c.Workers = w
	default:
		c.DefaultServer.addOption(opName, opValue)
	}
}

func (c *Conf) addVhost(vhost ServerConf) {
	if c.Vhosts == nil {
		c.Vhosts = make([]ServerConf, 0, 10)
	}
	c.Vhosts = append(c.Vhosts, vhost)
}

func (s *ServerConf) addOption(opName string, opValue string) {
	switch opName {
	case nameOption:
		s.parseNameOptions(opValue)
	case rootOption:
		s.Root = opValue
	case portOption:
		s.parsePortOptions(opValue)
	case indexOption:
		s.parseIndexOptions(opValue)
	case errorLogOption:
		s.ErrorLog = opValue
	case accessLogOption:
		s.AccessLog = opValue
	}

	// handle error pages
	if strings.Contains(opName, errorPageOption) {
		s.parseErrorPageOptions(opName, opValue)
	}
}

func (s *ServerConf) parseNameOptions(serverNames string) {
	names := strings.Split(serverNames, ",")
	if len(names) == 0 {
		return
	}
	s.Names = make([]string, 0)
	for _, n := range names {
		// don't include duplicated names
		containsName := false
		for _, n2 := range s.Names {
			if n == n2 {
				containsName = true
			}
		}
		if !containsName {
			s.Names = append(s.Names, n)
		}
	}
}

func (s *ServerConf) parsePortOptions(ports string) {
	portsStr := strings.Split(ports, ",")
	if len(portsStr) == 0 {
		return
	}
	s.Ports = make([]int, 0)
	for _, p := range portsStr {
		pInt, err := strconv.Atoi(p)
		if err != nil {
			log.Println(err)
			return
		}
		s.Ports = append(s.Ports, pInt)
	}
}

func (s *ServerConf) parseIndexOptions(pages string) {
	pagesSli := strings.Split(pages, ",")
	if len(pagesSli) == 0 {
		return
	}
	s.IndexPages = make([]string, 0)
	for _, p := range pagesSli {
		s.IndexPages = append(s.IndexPages, p)
	}
}

func (s *ServerConf) parseErrorPageOptions(errorType, page string) {
	eTypePieces := strings.Split(errorType, "_")
	// unlikely, but just in case
	if len(eTypePieces) == 1 {
		// TODO: maybe log this? or the function that checks the conf file should detect this?
		return
	}
	if s.ErrorPages == nil {
		s.ErrorPages = make([]ErrorPage, 0)
	}
	var errPage ErrorPage
	if len(eTypePieces) == 3 {
		// custom error page
		code, err := strconv.Atoi(eTypePieces[2])
		if err != nil {
			return
		}
		errPage = ErrorPage{
			Code: code,
			Page: strings.TrimSpace(page),
		}
	}
	if len(eTypePieces) == 2 {
		errPage = ErrorPage{
			Code: http.StatusBadRequest,
			Page: strings.TrimSpace(page),
		}
	}
	s.ErrorPages = append(s.ErrorPages, errPage)
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
