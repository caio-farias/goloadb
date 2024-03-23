package components

import (
	"log"
	"net"
	"server/src"
	"server/src/utils"
	"strings"
	"time"
)

const (
	DefaultHttp = "HTTP/1.1"
)

type Request struct {
	conn            *net.Conn
	response        bool
	contentReceived string
	Header          *Header
	Body            string
}

func NewRequest(url string, header *Header) (*Request, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Printf("Error trying to use server %s for TCP..", url)
		return nil, err
	}

	req := &Request{
		conn:     &conn,
		Header:   header,
		response: false,
	}

	req.parse()
	return req, nil
}

func (req *Request) Await() (string, error) {
	readb := make([]byte, src.BUFFER_LENGTH)
	_, read_err := (*req.conn).Read(readb[0:])

	if read_err != nil {
		log.Print("Failed to read from connection:", read_err)
		return "", read_err
	}

	if src.VERBOSE {
		log.Printf(">> Received response from client -> %s \n", (*req.conn).RemoteAddr().String())
	}
	(*req.conn).Close()
	return string(readb), nil
}

func NewRequestFromListener(conn *net.Conn) (*Request, error) {
	content, _ := readFromConnection(conn)
	req := &Request{
		conn:            conn,
		contentReceived: content,
	}

	req.parse()
	return req, nil
}

func (req *Request) String() string {
	return string(req.Header.String() + req.Body)
}

func (req *Request) Send(body string) (int, error) {
	req.Body = body
	if src.VERBOSE {
		log.Printf("<<< Response sent to client %s \n\n", (*req.conn).RemoteAddr().String())

	}
	res, err := (*req.conn).Write([]byte(req.String() + utils.END_OF_HEADER_LINE))
	if err != nil {
		utils.PrintAndSleep(DEFAULT_SLEEP_TIME, "Could not send message")
		return res, err
	}
	(*req.conn).Close()
	return res, nil
}

func (req *Request) SendResponse(body string) {
	req.parse()
	response := &Request{
		Header: &Header{
			ProtocolVersion: req.Header.ProtocolVersion,
			HeaderContent: map[string]string{
				utils.DateKey: utils.GetTimeHere(),
			},
		},
		Body: body,
	}
	(*req.conn).Write([]byte(response.String()))
	if src.VERBOSE {
		log.Printf("<<< Response sent to client %s \n\n", (*req.conn).RemoteAddr().String())

	}
	(*req.conn).Close()
}

func (req *Request) SendRequest(header Header, body string) *Request {
	response := &Request{
		Header: &header,
		Body:   body,
	}
	(*req.conn).Write([]byte(response.String()))
	return response
}

func readFromConnection(conn *net.Conn) (string, error) {
	(*conn).SetDeadline(time.Now().Add(utils.GetTimeInSeconds(10)))
	readb := make([]byte, src.BUFFER_LENGTH)
	_, read_err := (*conn).Read(readb[0:])

	if read_err != nil {
		log.Print("Failed to read from connection:", read_err)
		return "", read_err
	}

	if src.VERBOSE {
		log.Printf(">>> Received request on connection -> %s \n", (*conn).RemoteAddr().String())
	}

	return string(readb), nil
}

func (h *Request) parse() *Request {
	var header, body string

	utils.Unpack(strings.Split(h.contentReceived, utils.END_HEADER_PATTERN), &header, &body)
	h.Header = h.readHeader(header)
	h.Body = body
	return h
}

func (req *Request) readHeader(content string) *Header {
	lines := strings.Split(content, utils.END_OF_HEADER_LINE)
	var protocol, version, path, method, protocol_version, status, statusMesage string
	if req.response {
		utils.Unpack(strings.Fields(lines[0]), &protocol_version, &status, &statusMesage, &statusMesage)
	} else {
		utils.Unpack(strings.Fields(lines[0]), &method, &path, &protocol_version)
	}
	utils.Unpack(strings.Split(protocol_version, "/"), &protocol, &version)

	header := &Header{
		Method:          method,
		Path:            path,
		Protocol:        protocol,
		ProtocolVersion: protocol_version,
	}

	header_map := make(map[string]string)
	var key, val string
	const key_ex = utils.HostKey + ":"

	for _, line := range lines[1:] {
		if strings.Contains(line, key_ex) {
			utils.Unpack(strings.Split(line, key_ex), &key, &val)
			continue
		}
		utils.Unpack(strings.Split(line, ": "), &key, &val)
		header_map[key] = val
	}

	header.HeaderContent = header_map
	return header
}
