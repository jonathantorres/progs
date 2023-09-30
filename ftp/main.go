package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

const confJSON = `{
	"root": "/Users/jonathantorres/gftp_test",
	"users": [
		{
			"username": "jt",
			"password": "test",
			"root": "/jt"
		},
		{
			"username": "test",
			"password": "test",
			"root": "/test"
		}
	]
}`

const (
	ControlPort                    = 21
	DefaultHost                    = "localhost"
	defaultCmdSize                 = 512
	TransferTypeAscii TransferType = "A"
	TransferTypeImage TransferType = "I"
)

type TransferType string

type Server struct {
	Host string
	Port int
	Conf *ServerConf
}

type ServerConf struct {
	Root  string
	Users []*User
}

// the current active session
type Session struct {
	user         *User
	server       *Server
	tType        TransferType
	passMode     bool
	controlConn  net.Conn
	dataConn     net.Conn
	dataConnPort uint16
	dataConnChan chan struct{}
	cwd          string
}

// the current user logged in for this session
type User struct {
	Username string
	Password string
	Root     string
}

func main() {
	conf, err := loadConf()
	if err != nil {
		fmt.Fprintf(os.Stderr, "server conf error: %s\n", err)
		os.Exit(1)
	}
	s := &Server{
		Host: DefaultHost,
		Port: 9010, // for testing purposes
		Conf: conf,
	}
	err = s.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		os.Exit(1)
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept error:  %s\n", err)
			continue
		}
		go s.handleClient(conn)
	}
	return nil
}

func (s *Server) handleClient(conn net.Conn) {
	err := sendResponse(conn, 220, "") // welcome message
	if err != nil {
		fmt.Fprintf(os.Stderr, "error response: %s\n", err)
		return
	}
	session := &Session{
		controlConn: conn,
		server:      s,
	}
	session.start()
}

func (s *Session) start() {
	for {
		clientCmd := make([]byte, defaultCmdSize)
		_, err := s.controlConn.Read(clientCmd)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stderr, "connection finished by client %s\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "error read: %s\n", err)
				sendResponse(s.controlConn, 500, "")
			}
			s.controlConn.Close()
			break
		}
		err = s.handleCommand(clientCmd)
		if err != nil {
			sendResponse(s.controlConn, 500, "")
			continue
		}
	}
}

func (s *Session) openDataConn(port uint16) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.server.Host, port))
	if err != nil {
		return err
	}
	s.dataConnPort = port
	s.dataConnChan = make(chan struct{})
	go func() {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "data conn: accept error:  %s\n", err)
			return
		}
		go s.handleDataTransfer(conn, l)
	}()
	return nil
}

func (s *Session) handleDataTransfer(conn net.Conn, l net.Listener) {
	s.dataConn = conn

	var sig struct{}
	// send signal to command that the connection is ready
	s.dataConnChan <- sig

	// wait until the command finishes, then close the connection
	<-s.dataConnChan

	s.dataConn = nil
	s.dataConnPort = 0
	s.dataConnChan = nil
	defer conn.Close()
	defer l.Close()
}

func (s *Session) handleCommand(clientCmd []byte) error {
	clientCmdStr := trimCommandLine(clientCmd)
	cmd := ""
	cmdParams := ""
	foundFirstSpace := false
	for _, r := range clientCmdStr {
		if !foundFirstSpace && r == ' ' {
			foundFirstSpace = true
			continue
		}
		if foundFirstSpace {
			cmdParams += string(r)
		} else {
			cmd += string(r)
		}
	}
	if cmdParams == "" {
		return s.execCommand(cmd, "")
	} else {
		return s.execCommand(cmd, cmdParams)
	}
	return sendResponse(s.controlConn, 500, "")
}

func (s *Session) execCommand(cmd string, cmdArgs string) error {
	var err error = nil
	fmt.Fprintf(os.Stdout, "cmd: %s\n", cmd)
	switch cmd {
	case CommandUser:
		err = runCommandUser(s, cmdArgs)
	case CommandPassword:
		err = runCommandPassword(s, cmdArgs)
	case CommandPrintDir:
		err = runCommandPrintDir(s)
	case CommandChangeDir:
		err = runCommandChangeDir(s, cmdArgs)
	case CommandType:
		err = runCommandType(s, cmdArgs)
	case CommandPassive:
		err = runCommandPasv(s)
	case CommandList:
		err = runCommandList(s, cmdArgs)
	case CommandRetrieve:
		err = runCommandRetrieve(s, cmdArgs)
	case CommandAcceptAndStore:
		err = runCommandAcceptAndStore(s, cmdArgs)
	case CommandSystemType:
		err = runCommandSystemType(s)
	case CommandChangeParent, CommandChangeToParentDir:
		err = runCommandChangeParent(s)
	case CommandMakeDir, CommandMakeADir:
		err = runCommandMakeDir(s, cmdArgs)
	case CommandDelete:
		err = runCommandDelete(s, cmdArgs)
	default:
		err = runUninmplemented(s)
	}
	return err
}

// list of commands
const (
	CommandAbort             = "ABOR"
	CommandAccount           = "ACCT"
	CommandAuthData          = "ADAT"
	CommandAllo              = "ALLO"
	CommandAppend            = "APPE"
	CommandAuth              = "AUTH"
	CommandAvail             = "AVBL"
	CommandClear             = "CCC"
	CommandChangeParent      = "CDUP"
	CommandConf              = "CONF"
	CommandCsId              = "CSID"
	CommandChangeDir         = "CWD"
	CommandDelete            = "DELE"
	CommandDirSize           = "DSIZ"
	CommandPrivProtected     = "ENC"
	CommandExtAddrPort       = "EPRT"
	CommandExtPassMode       = "EPSV"
	CommandFeatLis           = "FEAT"
	CommandHelp              = "HELP"
	CommandHost              = "HOST"
	CommandLang              = "LANG"
	CommandList              = "LIST"
	CommandLongAddrPort      = "LPRT"
	CommandLongPassMode      = "LPSV"
	CommandLastModTime       = "MDTM"
	CommandModCreatTime      = "MFCT"
	CommandModFact           = "MFF"
	CommandModLastModTime    = "MFMT"
	CommandInteProtect       = "MIC"
	CommandMakeDir           = "MKD"
	CommandListDir           = "MLSD"
	CommandObjData           = "MLST"
	CommandMode              = "MODE"
	CommandFileNames         = "NLST"
	CommandNoOp              = "NOOP"
	CommandOptions           = "OPTS"
	CommandPassword          = "PASS"
	CommandPassive           = "PASV"
	CommandBufSizeProt       = "PBSZ"
	CommandPort              = "PORT"
	CommandDataChanProtLvl   = "PROT"
	CommandPrintDir          = "PWD"
	CommandQuit              = "QUIT"
	CommandReinit            = "REIN"
	CommandRestart           = "REST"
	CommandRetrieve          = "RETR"
	CommandRemoveDir         = "RMD"
	CommandRemoveDirTree     = "RMDA"
	CommandRenameFrom        = "RNFR"
	CommandRenameTo          = "RNTO"
	CommandSite              = "SITE"
	CommandFileSize          = "SIZE"
	CommandMountFile         = "SMNT"
	CommandSinglePortPassive = "SPSV"
	CommandServerStatus      = "STAT"
	CommandAcceptAndStore    = "STOR"
	CommandStoreFile         = "STOU"
	CommandFileStruct        = "STRU"
	CommandSystemType        = "SYST"
	CommandThumbnail         = "THMB"
	CommandType              = "TYPE"
	CommandUser              = "USER"
	CommandChangeToParentDir = "XCUP"
	CommandMakeADir          = "XMKD"
	CommandPrintCurDir       = "XPWD"
	CommandRemoveTheDir      = "XRMD"
	CommandSendMail          = "XSEM"
	CommandSendTerm          = "XSEN"
)

// server status codes
var statusCodes = map[uint16]string{
	110: "Restart marker replay",
	120: "Service ready in a few minutes",
	125: "Data connection already open",
	150: "File status okay, about to open data connection",
	200: "Command Ok",
	202: "Command not implemented",
	211: "System status",
	212: "Directory Status",
	213: "File Status",
	214: "Help message",
	215: "NAME system type",
	220: "Service Ready",
	221: "Closing control connection",
	225: "Data connection open, no transfer in progress",
	226: "Closing data connection. File action ok",
	227: "Entering passive mode",
	228: "Entering long passive mode",
	229: "Entering extended passive mode",
	230: "User logged in, proceed. Logged out if appropriate",
	231: "User logged out, service terminated.",
	232: "Logout command noted",
	234: "Authentication mechanism accepted",
	250: "Requested file action ok, completed",
	257: "Path created",
	331: "Username okay, need password.",
	332: "Need account for login.",
	350: "Requested file action pending more information",
	400: "Command not accepted, please try again",
	421: "Service not available, closing control connection",
	425: "Can't open data connection",
	426: "Connection closed, transfer aborted",
	430: "Invalid username or password",
	434: "Requested host unavailable",
	450: "Requested file action not taken",
	451: "Requested action aborted. Local error in processing",
	452: "Requested action not taken. Insufficient storage space in system. File unavailable",
	500: "Unknown error",
	501: "Syntax error in parameters or arguments. ",
	502: "Command not implemented",
	503: "Bad sequence of commands",
	504: "Command not implemented for that parameter.",
	530: "Not logged in.",
	532: "Need account for storing files.",
	534: "Could Not Connect to Server - Policy Requires SSL",
	550: "File not found, error encountered",
	551: "Requested action aborted. Page type unknown.",
	552: "Requested file action aborted. Exceeded storage allocation",
	553: "Requested action not taken. File name not allowed.",
	631: "Integrity protected reply",
	632: "Confidentiality and integrity protected reply",
	633: "Confidentiality protected reply",
}

func GetStatusCodeMessage(statusCode uint16) (string, error) {
	for code, statusMsg := range statusCodes {
		if code == statusCode {
			return statusMsg, nil
		}
	}
	return "", errors.New("status code not found")
}

// command functions
func runCommandUser(session *Session, username string) error {
	userFound := false
	for _, u := range session.server.Conf.Users {
		if u.Username == username {
			userFound = true
			session.user = u
			break
		}
	}
	if userFound {
		return sendResponse(session.controlConn, 331, "")
	}
	return sendResponse(session.controlConn, 430, "")
}

func runCommandPassword(session *Session, pass string) error {
	passFound := false
	for _, u := range session.server.Conf.Users {
		if u.Username == session.user.Username && u.Password == pass {
			passFound = true
			session.user = u
			break
		}
	}
	if passFound {
		// change to home directory
		err := os.Chdir(session.server.Conf.Root + session.user.Root)
		if err != nil {
			return sendResponse(session.controlConn, 550, "")
		}
		return sendResponse(session.controlConn, 230, "")
	}
	return sendResponse(session.controlConn, 430, "")
}

func runCommandPrintDir(session *Session) error {
	return sendResponse(session.controlConn, 257, "\"/"+session.cwd+"\" is current directory\n")
}

func runCommandChangeDir(session *Session, dir string) error {
	cwd := session.cwd
	if dir[0] == '/' {
		// moving to relative path
		cwd = dir[1:]
	} else {
		if cwd != "" {
			cwd += "/" + dir
		} else {
			cwd = dir
		}
	}
	err := os.Chdir(session.server.Conf.Root + session.user.Root + "/" + cwd)
	if err != nil {
		return sendResponse(session.controlConn, 550, "")
	}
	session.cwd = cwd
	return sendResponse(session.controlConn, 250, "CWD successful. \"/"+dir+"\" is current directory\n")
}

func runCommandType(session *Session, typ string) error {
	selectedTransferType := TransferType(typ)
	if selectedTransferType == TransferTypeAscii || selectedTransferType == TransferTypeImage {
		session.tType = TransferType(typ)
		return sendResponse(session.controlConn, 200, "Transfer type Ok")
	}
	return sendResponse(session.controlConn, 504, "")
}

func runCommandPasv(session *Session) error {
	session.passMode = true
	addr, err := findOpenAddr()
	if err != nil {
		return sendResponse(session.controlConn, 425, "")
	}
	respParts := make([]string, 0)
	for i := 0; i < len(addr.IP); i++ {
		respParts = append(respParts, strconv.Itoa(int(addr.IP[i])))
	}

	var p uint16 = uint16(addr.Port)
	var p1 uint8 = uint8(p >> 8)
	var p2 uint8 = uint8(p)
	respParts = append(respParts, strconv.Itoa(int(p1)))
	respParts = append(respParts, strconv.Itoa(int(p2)))
	respMsg := strings.Join(respParts, ",")

	if err = session.openDataConn(p); err != nil {
		return sendResponse(session.controlConn, 425, "")
	}
	return sendResponse(session.controlConn, 227, respMsg)
}

func runCommandList(session *Session, file string) error {
	// wait until the data connection is ready for sending/receiving data
	<-session.dataConnChan

	path := session.server.Conf.Root + session.user.Root + "/" + session.cwd
	if file != "" {
		path += "/" + file
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed listing directory: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	dirFiles := make([]string, 0)
	for _, f := range files {
		line := getFileLine(f)
		dirFiles = append(dirFiles, line)
	}
	dirData := strings.Join(dirFiles, "\n")
	_, err = session.dataConn.Write([]byte(dirData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed writing data: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	var sig struct{}
	session.dataConnChan <- sig
	return sendResponse(session.controlConn, 200, "")
}

func runCommandRetrieve(session *Session, filename string) error {
	<-session.dataConnChan
	path := session.server.Conf.Root + session.user.Root + "/" + session.cwd + "/" + filename
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	_, err = io.Copy(session.dataConn, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error transferring file: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	file.Close()
	var sig struct{}
	session.dataConnChan <- sig
	return sendResponse(session.controlConn, 200, "")
}

func runCommandAcceptAndStore(session *Session, filename string) error {
	<-session.dataConnChan
	path := session.server.Conf.Root + session.user.Root + "/" + session.cwd + "/" + filename
	fileData, err := ioutil.ReadAll(session.dataConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error receiving file: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	_, err = file.Write(fileData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing bytes to new file: %s\n", err)
		return sendResponse(session.controlConn, 450, "")
	}
	file.Close()
	var sig struct{}
	session.dataConnChan <- sig
	return sendResponse(session.controlConn, 200, "")
}

func runCommandSystemType(session *Session) error {
	return sendResponse(session.controlConn, 215, "UNIX Type: L8")
}

func runCommandChangeParent(session *Session) error {
	cwd := session.cwd
	pieces := strings.Split(cwd, "/")
	if len(pieces) <= 1 {
		cwd = ""
	} else {
		cwd = strings.Join(pieces[:len(pieces)-1], "/")
	}
	err := os.Chdir(session.server.Conf.Root + session.user.Root + "/" + cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err chdir: %s\n", err)
		return sendResponse(session.controlConn, 550, "")
	}
	session.cwd = cwd
	base := path.Base(cwd)

	return sendResponse(session.controlConn, 200, "CDUP successful. \"/"+base+"\" is current directory\n")
}

func runCommandMakeDir(session *Session, dirName string) error {
	cwd := session.cwd
	err := os.Mkdir(session.server.Conf.Root+session.user.Root+"/"+cwd+"/"+dirName, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err mkdir: %s\n", err)
		return sendResponse(session.controlConn, 550, "")
	}
	return sendResponse(session.controlConn, 200, fmt.Sprintf("Directory %s created", dirName))
}

func runCommandDelete(session *Session, filename string) error {
	cwd := session.cwd
	err := os.Remove(session.server.Conf.Root + session.user.Root + "/" + cwd + "/" + filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err remove file: %s\n", err)
		return sendResponse(session.controlConn, 550, "")
	}
	return sendResponse(session.controlConn, 200, fmt.Sprintf("File %s deleted", filename))
}

func runUninmplemented(session *Session) error {
	return sendResponse(session.controlConn, 502, "")
}

func sendResponse(conn net.Conn, statusCode uint16, respMsg string) error {
	codeMsg, err := GetStatusCodeMessage(statusCode)
	var code uint16
	if err != nil {
		code = 500
		respMsg = err.Error()
	} else {
		code = statusCode
		if respMsg == "" {
			respMsg = codeMsg
		}
	}
	respMsg = fmt.Sprintf("%d %s\n", code, respMsg)
	_, err = conn.Write([]byte(respMsg))
	if err != nil {
		return err
	}
	return nil
}

func findOpenAddr() (*net.TCPAddr, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return nil, errors.New("tcp address could not be resolved")
	}
	defer l.Close()
	return addr, nil
}

func trimCommandLine(clientCmd []byte) string {
	trimmedCommand := ""
	for _, b := range clientCmd {
		if rune(b) != 0x00 && rune(b) != '\r' && rune(b) != '\n' {
			trimmedCommand += string(b)
		}
	}
	return trimmedCommand
}

func loadConf() (*ServerConf, error) {
	conf := &ServerConf{}
	err := json.Unmarshal([]byte(confJSON), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getFileLine(file os.FileInfo) string {
	mode := file.Mode().String()
	size := file.Size()
	modTime := fmt.Sprintf("%s %d %d", file.ModTime().Month().String(), file.ModTime().Day(), file.ModTime().Year())
	return fmt.Sprintf("%s %d %s %s", mode, size, modTime, file.Name())
}
