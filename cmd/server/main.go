package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rodinv/pow/internal/book"
	"github.com/rodinv/pow/internal/delivery"
	"github.com/rodinv/pow/internal/hashcash"
	"github.com/rodinv/pow/internal/server"
	"github.com/rodinv/pow/internal/storage"
	"github.com/rodinv/pow/pkg/config"
)

func main() {
	// init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// init config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("initializing config")
		return
	}

	// init hashcash
	hashStore := storage.New()
	pow, err := hashcash.New(24, hashStore)
	if err != nil {
		log.Fatal().Err(err).Msg("initializing hashcash")
		return
	}

	// init handlers
	h := delivery.New(book.New(), pow)

	// init server and handlers
	s := server.New(cfg.Host, cfg.Port)
	s.AddHandler("challenge", h.Challenge)
	s.AddHandler("get_quote", h.GetQuote)

	// run and wait
	err = s.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	defer s.Stop()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}
