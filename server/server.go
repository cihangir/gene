package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cenkalti/backoff"
	"github.com/facebookgo/freeport"
	"github.com/rcrowley/go-tigertonic"
)

// Initer is for initializing the api endpoints for a module
type Initer interface {
	Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux
}

type Server struct {
	Mux    *tigertonic.TrieServeMux
	Port   int
	Server *tigertonic.Server
}

func New(api Initer) *Server {
	s := &Server{
		Mux: tigertonic.NewTrieServeMux(),
	}

	s.Mux = api.Init(s.Mux)
	return s
}

func (s *Server) Listen() error {
	server, err := s.crateServer()
	if err != nil {
		return err
	}

	s.Server = server

	fmt.Println("started listening on: ", s.Server.Server.Addr)
	// todo move those to a more generic place
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Println(<-ch)
	s.Server.Close()
	fmt.Println("closed")

	return nil
}

func (s *Server) crateServer() (*tigertonic.Server, error) {
	ticker := backoff.NewTicker(
		backoff.NewExponentialBackOff(),
	)

	var err error
	var server *tigertonic.Server
	for _ = range ticker.C {
		if server, err = s.startListening(); err != nil {
			log.Printf("err while creating the server: %s, will try again", err.Error())
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) startListening() (*tigertonic.Server, error) {
	var err error
	if s.Port == 0 {
		s.Port, err = freeport.Get()
		if err != nil {
			return nil, err
		}
	}

	server := tigertonic.NewServer(
		fmt.Sprintf("0.0.0.0:%d", s.Port),
		s.Mux,
	)

	go func() {
		err := server.ListenAndServe()
		if nil != err {
			log.Fatalf("listen err:", err.Error())
		}
	}()

	return server, nil
}
