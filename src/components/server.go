package components

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type ServerHandler func(w http.ResponseWriter, r *http.Request) error

type Server struct {
	smux *http.ServeMux
	port string
	lb   *LoadBalancer
}

func NewServer(port string) *Server {
	return &Server{
		smux: http.NewServeMux(),
		port: ":" + port,
	}
}

func (s *Server) EnableLoadBalancing(lb *LoadBalancer, path string) {
	s.lb = lb
	lb.enableBalancing = true
	defer lb.EnableLoadBalancing()

	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) error {
		targetHost := s.lb.DiscoverService()
		client := http.Client{Transport: &http.Transport{
			DisableCompression: true,
		}}

		resBuffer, err := client.Get(targetHost)

		if err != nil {
			return err
		}

		defer resBuffer.Body.Close()

		resBytes := make([]byte, 1024)
		if _, err = resBuffer.Body.Read(resBytes); err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		w.Write(resBytes)
		return nil
	})
}

func (s *Server) GetDiscoveryInfo() string {
	return s.lb.GetDiscoveryInfo()
}

func (s *Server) HandleFunc(path string, fn ServerHandler) {
	if fn == nil {
		log.Println("Handler is nil.")
	}
	s.smux.HandleFunc(path, s.handlerWrapper(fn))
}

func (s *Server) handlerWrapper(handler ServerHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {

			switch svError := err.(type) {
			case *ServerError:
				json.NewEncoder(w).Encode(map[string]any{
					"status": svError.status,
					"error":  svError.mssg,
				})
			default:
				json.NewEncoder(w).Encode(map[string]any{
					"status": 500,
					"error":  "Something bad happened",
				})
			}
			log.Println(">> ERROR: ", err)

		}
	}
}

func (s *Server) Run() {
	if err := http.ListenAndServe(s.port, s.smux); err != nil {
		log.Fatalln(">> ERROR: ", err)
	}
}
