package mssql

import "fmt"

type VerboseLogger struct {
	verbose bool
}

func (l VerboseLogger) Log(s ...any) {
	if l.verbose {
		fmt.Println(s...)
	}
}
