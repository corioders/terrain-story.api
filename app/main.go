package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/corioders/terrain-story.api/foundation"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/log"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run() error {
	logger := log.New(os.Stdout, "github.com/corioders/terrain-story.api")

	config, err := foundation.LoadConfig()
	if err != nil {
		return err
	}

	app := foundation.NewApplication(logger, config)

	router, err := newRouter(app)
	if err != nil {
		return err
	}

	server := http.Server{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Addr:         ":" + config.Web.Port,
		Handler:      router,
	}

	app.RegisterOnStop("Server", func() error {
		return server.Shutdown(context.Background())
	})

	errChan := make(chan error)

	go func() {
		fmt.Println(app.GetConfig())
		err = server.ListenAndServe()
		if err != nil {
			errChan <- errors.WithStack(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	select {
	case err = <-errChan:
		return err

	case <-signalChan:
		return app.Stop()
	}
}
