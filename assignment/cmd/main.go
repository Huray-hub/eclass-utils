package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Huray-hub/eclass-utils/assignment"
	"github.com/Huray-hub/eclass-utils/assignment/calendar"
	"github.com/Huray-hub/eclass-utils/assignment/cmd/internal/flags"
	"github.com/Huray-hub/eclass-utils/assignment/cmd/internal/output"
	"github.com/Huray-hub/eclass-utils/assignment/config"
)

func init() {
	homeCache, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	path := filepath.Join(homeCache, "eclass-utils")
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(
		filepath.Join(path, "assignments.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.SetOutput(file)
}

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancelFunc := context.WithCancel(context.Background())

	go handleSignals(signalCh, cancelFunc)

	cfg, err := config.ImportDefault()
	if err != nil {
		log.Fatal(err)
	}

	flags.Read(cfg)

	err = config.Ensure(cfg)
	if err != nil {
		log.Fatal(err)
	}

	service, err := assignment.NewService(ctx, *cfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	opts := cfg.Options
	assignments, err := service.FetchAssignments(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = output.PrintAssignments(assignments, opts.PlainText)
	if err != nil {
		log.Fatal(err)
	}

	if opts.ExportICS {
		path, err := calendar.Export(assignments, opts.BaseDomain)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("stored in\n%v\n", path)
	}
}

func handleSignals(signalCh <-chan os.Signal, cancelFunc context.CancelFunc) {
	for signal := range signalCh {
		switch signal {
		case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
			fmt.Println(" signal received:", signal.String())
			cancelFunc()
			os.Exit(0)
		default:
			fmt.Println(" unhandled/unknown signal")
		}
	}
}
