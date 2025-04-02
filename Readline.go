package initdockerdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Readline struct {
	hadOutput bool
	reader    *bufio.Reader
}

func NewReadline() Readline {
	return Readline{
		hadOutput: false,
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (rl *Readline) Question(question string, defaultAnswer string) (string, error) {
	rl.hadOutput = true
	fmt.Printf("%s (%s): ", question, defaultAnswer)
	answer, err := rl.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}
	return answer, nil
}
