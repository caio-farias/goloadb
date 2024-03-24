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

func NewRequest(ctx *Context) (*Request, error) {
	conn, err := net.Dial("tcp", ctx.Url)
	if err != nil {
		log.Printf("Error trying to connect with %s", ctx.Url)
		return nil, err
	}

	req := &Request{
		conn:     &conn,
		Header:   ctx.Header,
		response: false,
	}

	return req, nil
}

func NewRequestFromListener(conn *net.Conn) (*Request, error) {
	(*conn).SetDeadline(time.Now().Add(utils.GetTimeInSeconds(10)))
	readb := make([]byte, src.BUFFER_LENGTH)

	_, read_err := (*conn).Read(readb[0:])
	if read_err != nil {
		log.Print("Failed to read from connection:", read_err)
		return nil, read_err
	}

	if src.VERBOSE {
		log.Printf("## Received request on connection -> %s \n", (*conn).RemoteAddr().String())
	}

	content := string(readb)
	req := &Request{
		conn:            conn,
		contentReceived: content,
	}

	var header_content, body, rest string

	utils.Unpack(strings.Split(content, utils.END_HEADER_PATTERN), &header_content, &body, &rest)
	req.Body = body + rest
	req.Header = NewHeader(header_content, false)

	return req, nil
}

func (req *Request) AwaitResponse() (string, error) {
	log.Printf("Awaiting response from client -> %s \n", (*req.conn).RemoteAddr().String())
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

func (req *Request) String() string {
	return string(req.Header.String() + req.Body)
}

func (req *Request) Send(body string) (int, error) {
	req.Body = body
	if src.VERBOSE {
		log.Printf("## Message sent to client -> %s \n\n", (*req.conn).RemoteAddr().String())

	}
	req_str := req.String() + utils.END_OF_HEADER_LINE
	res, err := (*req.conn).Write([]byte(req_str))
	if err != nil {
		utils.PrintAndSleep(DEFAULT_SLEEP_TIME, "Could not send message"+(*req.conn).RemoteAddr().Network())
		return res, err
	}
	return res, nil
}

func (req *Request) SendNow() (int, error) {
	if src.VERBOSE {
		log.Printf("## Message sent to client -> %s \n\n", (*req.conn).RemoteAddr().String())
	}

	res, err := (*req.conn).Write([]byte(req.String()))
	if err != nil {
		utils.PrintAndSleep(DEFAULT_SLEEP_TIME, "Could not send message"+(*req.conn).RemoteAddr().Network())
		return res, err
	}
	return res, nil
}

func Ping(ctx *Context) (int, error) {
	req, err := NewRequest(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	n, err := req.Send("")
	(*req.conn).Close()
	return n, err
}

func (req *Request) GetPath() string {
	return req.Header.Path
}

func (req *Request) GetMethod() string {
	return req.Header.Method
}
