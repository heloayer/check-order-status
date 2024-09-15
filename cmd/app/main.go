package main

import (
	"github.com/gin-gonic/gin"
	"github.com/heloayer/check-order-status/config"
	"github.com/heloayer/check-order-status/core/external_api/providers"
	"github.com/heloayer/check-order-status/internal/resources/http/handler"
	"github.com/heloayer/check-order-status/internal/service"
	"github.com/heloayer/check-order-status/internal/service/order"
	"github.com/heloayer/check-order-status/service/poller"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("app failed to run")
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load configuration")
		return err
	}

	providerList := []providers.Provider{
		providers.NewJSONProvider(cfg.JsonAPIConfig),
		providers.NewMockAPIProvider(cfg.MockAPIConfig),
	}

	orderService := order.NewOrderService(providerList)
	services := service.New(orderService)
	router := gin.Default()
	server := handler.NewServer(cfg, router, *services)

	p := poller.NewPoller(providerList, cfg.PollInterval)

	go p.Start()
	defer p.Stop()

	log.Info().Msg("starting server on :8080")
	go func() {
		if err := http.ListenAndServe(":8080", server.Router); err != nil {
			log.Error().Err(err).Msg("failed to run server")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("shutting down server...")
	return nil
}
