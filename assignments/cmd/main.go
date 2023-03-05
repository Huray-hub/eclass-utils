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

	"github.com/Huray-hub/eclass-utils/assignments/assignment"
	"github.com/Huray-hub/eclass-utils/assignments/calendar"
	"github.com/Huray-hub/eclass-utils/assignments/cmd/flags"
	"github.com/Huray-hub/eclass-utils/assignments/cmd/output"
	"github.com/Huray-hub/eclass-utils/assignments/config"
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

	opts, creds, err := config.Import()
	if err != nil {
		log.Fatal(err.Error())
	}

	flags.Read(opts, creds)

	err = config.Ensure(opts, creds)
	if err != nil {
		log.Fatal(err.Error())
	}

	service, err := assignment.NewService(ctx, opts, *creds, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	assignments, err := service.FetchAssignments(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = output.PrintAssignments(assignments, opts.PlainText)
	if err != nil {
		log.Fatal(err.Error())
	}

	if opts.ExportICS {
		path, err := calendar.Export(assignments, opts.BaseDomain)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("stored in\n%v\n", path)
	}
}

func handleSignals(signalCh <-chan os.Signal, cancelFunc context.CancelFunc) {
	for signal := range signalCh {
		switch signal {
		case syscall.SIGTERM:
			fmt.Println(" signal:", signal.String())
			cancelFunc()
			os.Exit(0)
		case syscall.SIGINT:
			fmt.Println(" signal:", signal.String())
			cancelFunc()
			os.Exit(0)
		case syscall.SIGQUIT:
			fmt.Println(" signal:", signal.String())
			cancelFunc()
			os.Exit(0)
		default:
			fmt.Println(" unhandled/unknown signal")
		}
	}
}
