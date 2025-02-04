package mssql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

type WaitForOpts struct {
	Signal       context.Context
	WatchDogMs   int
	InitialDelay int
	MaxDelay     int
	Rate         float64
	Immediate    bool
}

func waitFor(exp func() error) error {
	// TODO
	return waitForWithOpts(exp, WaitForOpts{})
}

func waitForWithOpts(exp func() error, opts WaitForOpts) error {
	if 1.0 > opts.Rate || opts.Rate > 10 {
		return fmt.Errorf("rate for exp backoff must be in range 1 <= RATE <= 10, got %v", opts.Rate)
	}
	if opts.InitialDelay < 1 {
		return fmt.Errorf("initial delay value must be a positive integer, got %d", opts.InitialDelay)
	}
	if opts.MaxDelay < opts.InitialDelay {
		return fmt.Errorf("max delay value must be a positive integer bigger than initialDelay, got %d (initialDelay = %d)", opts.MaxDelay, opts.InitialDelay)
	}
	if opts.WatchDogMs < opts.MaxDelay {
		return fmt.Errorf("watchDogMs value must be a positive integer bigger than maxDelay, got %d (maxDelay = %d)", opts.WatchDogMs, opts.MaxDelay)
	}

	currentDelay := time.Duration(opts.InitialDelay) * time.Millisecond
	if !opts.Immediate {
		time.Sleep(currentDelay)
		currentDelay = time.Duration(math.Min(float64(currentDelay)*opts.Rate, float64(opts.MaxDelay))) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(opts.Signal, time.Duration(opts.WatchDogMs)*time.Millisecond)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.New("TimeoutError")
		default:
			if err := exp(); err == nil {
				return nil
			}
			time.Sleep(currentDelay)
			currentDelay = time.Duration(math.Min(float64(currentDelay)*opts.Rate, float64(opts.MaxDelay))) * time.Millisecond
		}
	}
}
