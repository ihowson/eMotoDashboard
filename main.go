package main

import (
	"context"
	"log"

	"github.com/ihowson/eMotoDashboard/bike"
	"github.com/ihowson/eMotoDashboard/gui"
)

func main() {
	// TODO: logging with rotation https://github.com/natefinch/lumberjack
	// TODO: adopt https://github.com/rs/zerolog
	// also https://stackoverflow.com/questions/36139061/best-way-to-roll-log-file-in-go

	ctx := context.Background()

	model, bike := bike.Build()

	go func() {
		err := MetricsServer(ctx, bike)
		if err != nil && err != context.DeadlineExceeded {
			log.Fatalf("MetricsServer: %v", err)
		}
	}()

	gui := gui.MotoGUI{
		Model: model,
		Bike:  bike,
	}
	gui.Run()

	// TODO: send cancel context into system and each component
	// cancel := ca.Run()
	// defer cancel()
}
