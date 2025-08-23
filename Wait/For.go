package Wait

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Wait options
type Opts struct {
	PreDelay time.Duration
	MinDelay time.Duration
	MaxDelay time.Duration
	Rate     float64
}

// Wait for exp to return without an error, retrying its calls with
// an increasing exponential backoff
func For(ctx context.Context, exp func() error, opts Opts) error {
	if opts.MinDelay == 0 {
		opts.MinDelay = time.Duration(500) * time.Millisecond
	}
	if opts.MaxDelay == 0 {
		opts.MaxDelay = time.Duration(5000) * time.Millisecond
	}
	if opts.Rate == 0 {
		opts.Rate = 1.5
	}

	if 1.0 > opts.Rate || opts.Rate > 10 {
		return fmt.Errorf("rate for exp backoff must be in range 1 <= RATE <= 10, got %v", opts.Rate)
	}
	if opts.MinDelay < 1 {
		return fmt.Errorf("min delay value must be a positive integer, got %d", opts.MinDelay)
	}
	if opts.MaxDelay < opts.MinDelay {
		return fmt.Errorf("max delay value must be a positive integer bigger than initialDelay, got %d (minDelay = %d)", opts.MaxDelay, opts.MinDelay)
	}

	currentDelay := opts.MinDelay
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation timed out: %w", ctx.Err())
		default:
			if opts.PreDelay != 0 {
				time.Sleep(opts.PreDelay)
				opts.PreDelay = 0
			}
			if err := exp(); err == nil {
				return nil
			}
			time.Sleep(currentDelay)
			currentDelay = time.Duration(math.Min(float64(currentDelay)*opts.Rate, float64(opts.MaxDelay)))
		}
	}
}
