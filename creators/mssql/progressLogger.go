package mssql

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/huh/spinner"
	"golang.org/x/term"
)

type ProgressLogger struct {
	verbose    bool
	cancelFunc context.CancelFunc
	isTerminal bool
	wg         sync.WaitGroup
}

func NewProgressLogger(verbose bool) ProgressLogger {
	return ProgressLogger{
		verbose:    verbose,
		isTerminal: term.IsTerminal(int(os.Stdout.Fd())),
	}
}

func (l *ProgressLogger) LogState(state string) {
	if l.verbose {
		fmt.Println(state)
		return
	}

	if !l.isTerminal {
		return
	}

	// Cancel previous spinner if running
	if l.cancelFunc != nil {
		l.cancelFunc()
		l.wg.Wait() // Wait for previous spinner to finish
	}

	// Start new spinner in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	l.cancelFunc = cancel

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		spinner.New().
			Title(state).
			Context(ctx).
			Run()
	}()
}

func (l *ProgressLogger) LogVerbose(s ...any) {
	if l.verbose {
		fmt.Println(s...)
	}
}

func (l *ProgressLogger) Done() {
	if l.cancelFunc != nil {
		l.cancelFunc()
		l.wg.Wait() // Wait for spinner goroutine to finish
		l.cancelFunc = nil
	}
}
