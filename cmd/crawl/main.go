package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"crawler/internal/logger"

	"crawler"
)

func main() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c

		os.Exit(0)
	}()

	logger := logger.NewStderr()

	if err := run(logger); err != nil {
		logger.Errorf("can't run crawler: %v", err)
		os.Exit(1)
	}
}

func run(logger logger.Logger) error {
	flag.Parse()

	var httpClient http.Client

	crawler := crawler.Crawler{
		Logger:     logger,
		HTTPGetter: &httpClient,
	}

	for _, arg := range flag.Args() {
		startingURL, err := url.Parse(arg)
		if err != nil {
			return err
		}

		crawler.Crawl(startingURL, func(url *url.URL) {
			fmt.Println(url)
		})
	}

	return nil
}
