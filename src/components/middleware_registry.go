package components

import (
	"log"
	"net"
	"server/src/utils"
	"sync"
	"time"
)

const (
	EXEC_CYCLE_SLEEP_DURATION = 1
)

type Handler = func(ctx *Context) error

type MiddlewareRegistry struct {
	name             string
	lb               *LoadBalancer
	handler_registry map[string]Handler
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		name:             "default",
		handler_registry: map[string]Handler{},
	}
}

func (mr *MiddlewareRegistry) AddHandler(path string, handler Handler) {
	mr.handler_registry[path] = handler
}

func (mr *MiddlewareRegistry) Exec(lb *LoadBalancer, wg *sync.WaitGroup) {
	mr.lb = lb

	mr.AddHandler("/info", func(ctx *Context) error {
		ctx.Header = &Header{
			ProtocolVersion: DefaultHttp,
			Status:          "200",
			HeaderContent: map[string]string{
				utils.UserAgenteKey:  utils.Fake_User_Agent,
				utils.AcceptKey:      utils.AcceptAll,
				utils.ConnectionKey:  utils.KeepAlive,
				utils.ContentTypeKey: utils.ApplicationJson,
			},
		}
		ctx.Body = lb.GetDiscoveryInfo()
		return nil
	})

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

			received_request, err := NewRequestFromListener(&conn)
			if err != nil {
				conn.Close()
				time.Sleep(EXEC_CYCLE_SLEEP_DURATION * time.Millisecond)
				return
			}

			lb.FindService()
			ctx := &Context{
				Url:           lb.targetHost,
				Header:        received_request.Header,
				Body:          received_request.Body,
				DiscoveryInfo: lb.GetDiscoveryInfo(),
			}

			handler := mr.handler_registry[received_request.GetPath()]

			if handler != nil {
				err = handler(ctx)
				if err != nil {
					ctx.Header = &Header{
						ProtocolVersion: DefaultHttp,
						Status:          "500",
						HeaderContent: map[string]string{
							utils.ContentTypeKey: utils.TextHtml,
						},
					}
					ctx.Body = "<center><h1>ERROR</h1></center>"
				}
			} else {
				ctx.Header = &Header{
					ProtocolVersion: DefaultHttp,
					Status:          "404",
					HeaderContent: map[string]string{
						utils.ContentTypeKey: utils.TextHtml,
					},
				}
				ctx.Body = "<center><h1>404 NO FOUND</h1></center>"
			}

			received_request.Body = ctx.Body
			received_request.Header = ctx.Header
			received_request.SendNow()

			(*received_request.conn).Close()
			time.Sleep(EXEC_CYCLE_SLEEP_DURATION * time.Millisecond)
		}
	}()
}
