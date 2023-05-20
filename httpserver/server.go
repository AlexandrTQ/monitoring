package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"monitoring/service/logging"
)

type server struct {
	*http.Server
}

func New(addr string) (*server, error) {
	appRouter, err := newAppRouter()
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		Addr:    addr,
		Handler: appRouter,
	}

	return &server{&srv}, nil
}

func (srv *server) Start() {
	go func() {
		logging.Debugf("Server start at: %s", srv.Addr)

		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			if err != nil {
				logging.Errorf("Error on Server listen stop - %s", err.Error())
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	err := srv.Shutdown(context.Background())
	if err != nil {
		logging.Errorf("Error on Server shutdown - %s", err.Error())
		panic(err)
	}

	logging.Debugf("Server gracefully shutdown")
}
