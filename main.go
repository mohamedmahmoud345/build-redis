package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var store = map[string]string{}

func handleCommand(args []string) string {
	cmd := strings.ToUpper(args[0])
	argsCount := len(args) - 1
	if isArityInvalid(cmd, argsCount) {
		return errorReply(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd))
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
	case "SET":
		store[args[1]] = args[2]
		return simpleString("OK")
	case "GET":
		v, ok := store[args[1]]
		if !ok {
			return nullBulk()
		}
		return bulkString(v)
	case "DBSIZE":
		return integer(len(store))
	}

	return errorReply(fmt.Sprintf("ERR unknown command '%s'", cmd))
}

var arity = map[string][2]int{
	"PING": {0, 1},
	"ECHO": {1, 1},
	"GET":  {1, 1},
	"SET":  {2, 2},
}

func isArityInvalid(cmd string, argsCount int) bool {

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
func errorReply(s string) string {
	return fmt.Sprintf("-%s\r\n", s)
}
func integer(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}
func bulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}
func nullBulk() string {
	return "$-1\r\n"
}
func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		args, err := readCommand(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		fmt.Print(handleCommand(args))
	}
}

// *3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n

func readCommand(r *bufio.Reader) ([]string, error) {
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(line, "*") {
		n, _ := strconv.Atoi(line[1:])
		args := make([]string, 0, n)
		for i := 0; i < n; i++ {
			hdr, _ := readLine(r)
			size, _ := strconv.Atoi(hdr[1:])
			buf := make([]byte, size)
			io.ReadFull(r, buf)
			readLine(r)
			args = append(args, string(buf))
		}
		return args, nil
	}

	return parseArgs(line), nil
}

func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		if len(s) > 0 && err == io.EOF {
			return strings.TrimRight(s, "\r\n"), nil
		}
		return "", err
	}
	return strings.TrimRight(s, "\r\n"), nil
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
