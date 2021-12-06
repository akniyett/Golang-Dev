package http

import (
	"context"
	"6/internal/store"
	"log"
	"net/http"
	lru "github.com/hashicorp/golang-lru"
	"6/internal/message_broker"
	"time"
	"github.com/go-chi/chi"

)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 30 * time.Second
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	cache 		*lru.TwoQueueCache
	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	accessoriesResource := NewAccessoryResource(s.store, s.broker, s.cache)
	r.Mount("/accessories", accessoriesResource.Routes())

	clothingsResource := NewCLothingResource(s.store, s.broker, s.cache)
	r.Mount("/clothings", clothingsResource.Routes())

	usersResource := NewUserResource(s.store, s.broker, s.cache)
	r.Mount("/users", usersResource.Routes())

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() 

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {

	<-s.idleConnsCh
}