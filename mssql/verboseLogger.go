package mssql

import "fmt"

type VerboseLogger struct {
	verbose bool
}

func (l VerboseLogger) Log(s string) {
	if l.verbose {
		fmt.Printf("verbose: %s\n", s)
	}
}
