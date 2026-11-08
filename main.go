package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func handleCommand(args []string) string {
	cmd := strings.ToUpper(args[0])
	argsCount := len(args) - 1

	if isArityValid(cmd, argsCount) {
		return error(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd))
	}

	switch cmd {
	case "PING":
		if argsCount > 0 {
			return bulkString(args[1])
		}
		return simpleString("PONG")
	case "ECHO":
		return bulkString(args[1])
	case "COMMAND":
		return simpleString("OK")
	}

	return error(fmt.Sprintf("ERR unknown command '%s'", cmd))
}

var arity = map[string][2]int{
	"PING": {0, 1},
	"ECHO": {1, 1},
}

func isArityValid(cmd string, argsCount int) bool {

	value, exists := arity[cmd]

	if !exists {
		return false
	}
	if !(argsCount >= value[0] && argsCount <= value[1]) {
		return true
	}
	return false
}

func simpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}
func error(s string) string {
	return fmt.Sprintf("-%s\r\n", s)
}
func integer(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}
func bulkString(s string) string {
	if len(s) == 0 {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		args := parseArgs(line)
		response := handleCommand(args)
		fmt.Print(response)
	}
}

func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	for _, ch := range line {
		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true
		case ch == '"' && inQuotes:
			inQuotes = false
		case ch == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
