package mssql

import (
	"context"
	"fmt"
	"math"
	"time"
)

type WaitForOpts struct {
	InitialDelay int
	MaxDelay     int
	Rate         float64
	Immediate    bool
}

func waitFor(ctx context.Context, exp func() error) error {
	return waitForWithOpts(ctx, exp, WaitForOpts{})
}

func waitForWithOpts(ctx context.Context, exp func() error, opts WaitForOpts) error {
	if opts.InitialDelay == 0 {
		opts.InitialDelay = 200
	}
	if opts.MaxDelay == 0 {
		opts.MaxDelay = 5000
	}
	if opts.Rate == 0 {
		opts.Rate = 1.5
	}

	if 1.0 > opts.Rate || opts.Rate > 10 {
		return fmt.Errorf("rate for exp backoff must be in range 1 <= RATE <= 10, got %v", opts.Rate)
	}
	if opts.InitialDelay < 1 {
		return fmt.Errorf("initial delay value must be a positive integer, got %d", opts.InitialDelay)
	}
	if opts.MaxDelay < opts.InitialDelay {
		return fmt.Errorf("max delay value must be a positive integer bigger than initialDelay, got %d (initialDelay = %d)", opts.MaxDelay, opts.InitialDelay)
	}

	currentDelay := time.Duration(opts.InitialDelay) * time.Millisecond
	if !opts.Immediate {
		time.Sleep(currentDelay)
		currentDelay = time.Duration(math.Min(float64(currentDelay)*opts.Rate, float64(opts.MaxDelay))) * time.Millisecond
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation timed out: %w", ctx.Err())
		default:
			if err := exp(); err == nil {
				return nil
			}
			time.Sleep(currentDelay)
			currentDelay = time.Duration(math.Min(float64(currentDelay)*opts.Rate, float64(opts.MaxDelay))) * time.Millisecond
		}
	}
}
