package logger

import (
	"fmt"
	"strconv"
)

type LogCaller struct {
	name  string
	color int
	err   bool
}

func (caller LogCaller) Out(output string) {
	if !caller.err {
		fmt.Printf("\033[%sm[%s]\033[0m %s\n", strconv.Itoa(caller.color), caller.name, output)
	} else {
		fmt.Printf("\033[1m\033[91m[Error]\033[0m\033[1m Error occurred in \033[%sm[%s]\033[0m\033[1m:\033[0m\n	%s\n", strconv.Itoa(caller.color), caller.name, output)
	}
}

func New(name string, color int, err bool) *LogCaller {
	p := new(LogCaller)
	p.name = name
	p.color = color
	p.err = err
	return p
}
