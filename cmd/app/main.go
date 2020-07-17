package main

import (
	"context"
	"flag"
	"fmt"
	net "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mustafadubul/product/internal/service"
	"github.com/rs/zerolog"

	"github.com/mustafadubul/product/internal/handler/http"
	"github.com/mustafadubul/product/internal/repository/sqlite"
)

func main() {

	filePath := flag.String("path", "", "path to SQLite file")
	host := flag.String("host", "0.0.0.0:8080", "server port")
	debug := flag.Bool("debug", false, "sets log level to debug")
	inMemory := flag.Bool("db_in_memory", false, "choose to have the Database in memory")

	migrate := flag.Bool("migrate", false, "migrate database")

	flag.Parse()

	db, err := sqlite.Open(*inMemory, *filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	l := zerolog.New(os.Stdout).With().Logger()

	repo := sqlite.New(db)
	defer repo.Close()

	if *migrate {
		repo.Migrate()
	}

	svc := service.New(&l, repo)

	server := net.Server{
		Addr:    *host,
		Handler: http.NewHandler(&l, svc).Setup(),
	}

	go func() {
		l.Info().Str("host", *host).Msg("Listening...")
		if err := server.ListenAndServe(); err != net.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(2)
		} else {
			l.Info().Msg("Server closed!")
		}
	}()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	l.Info().Str("signal", sig.String()).Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Error().Err(err).Msg("Unable to shut down server")
	}

	l.Info().Msg("Server stopped")
	os.Exit(0)
}
