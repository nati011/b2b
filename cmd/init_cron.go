package main

import (
	"errors"

	"b2b.nati011.github.com/internal/adapter/primary/cron"
	"github.com/go-co-op/gocron/v2"
)

var (
	ErrFailedToBuildCrons = errors.New("¯\\_(ツ)_/¯, failed to init cron jobs")
)

func InitCron(s gocron.Scheduler) {
	err := cron.BuildCrons(s)
	if err != nil {
		panic(ErrFailedToBuildCrons)
	}
}
