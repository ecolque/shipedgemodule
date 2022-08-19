package core

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ecolque/shipedgemodule/wb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Config struct {
	Port      string
	JWTSecret string
	PQUrl     string
}

type Server interface {
	Config() *Config
	PQ() *gorm.DB
	Hub() *wb.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	pq     *gorm.DB
	hub    *wb.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}
func (b *Broker) PQ() *gorm.DB {
	return b.pq
}

func (b *Broker) Hub() *wb.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config, db *gorm.DB) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("Secret JWT is required")
	}
	if config.PQUrl == "" {
		return nil, errors.New("Database utl is required")
	}
	return &Broker{
		config: config,
		router: mux.NewRouter(),
		pq:     db,
		hub:    wb.NewHub(),
	}, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	handler := cors.Default().Handler(b.router)
	srv := &http.Server{
		Addr:    b.config.Port,
		Handler: handler,
	}
	go b.hub.Run()
	binder(b, b.router)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
