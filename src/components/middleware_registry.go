package components

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	EXEC_CYCLE_SLEEP_DURATION = 500
)

type Handler = func(url string, header *Header, body string) (string, error)

type MiddlewareRegistry struct {
	name     string
	registry map[string]Handler
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		name:     "default",
		registry: map[string]Handler{},
	}
}

func (mr *MiddlewareRegistry) Add(path string, handler Handler) {
	mr.registry[path] = handler
}

func (mr *MiddlewareRegistry) Exec(lb *LoadBalancer, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		port := lb.Port
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Printf("Error trying to use port %s for TCP..", port)
			return
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error trying to receive next call on connection.")
				continue
			}

			received_request, _ := NewRequestFromListener(&conn)
			lb.FindService()

			handler := mr.getHandler(received_request)
			if handler == nil {
				received_request.SendResponse("<p> 404 NOT FOUND </p>")
				continue
			}

			response, err := handler(lb.targetHost, received_request.Header, received_request.Body)

			if err != nil {
				log.Println(err)
				received_request.SendResponse("<p> 404 NOT FOUND </p>")
				continue
			}

			received_request.SendResponse(response)
			time.Sleep(EXEC_CYCLE_SLEEP_DURATION * time.Millisecond)
		}
	}()
}

func (mr *MiddlewareRegistry) getHandler(req *Request) Handler {
	path := req.Header.Path
	handler := mr.registry[path]
	return handler
}
