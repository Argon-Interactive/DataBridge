package console

import (
	"DataBridge/packages/server"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func RunInterface() {
	fmt.Println("Welcome to the DataBridge!")
	quit := make(chan struct{})
	input := make(chan string)

	go signalLisiner(quit, input)
	go inputLisiner(quit, input)

	fmt.Printf(">> ")
	for {
		command := strings.ToLower(<-input)
		switch command {
		case "r": fallthrough
		case "run": server.RunHTTPServer()
		case "r-s": fallthrough
		case "run-s": server.RunHTTPSServer()
		case "stop": server.KillServer()
		case "restart": restart()
		case "help": printHelp()
		case "status": printStatus()
		case "clr": fmt.Println("\033[H\033[2J")
		case "q": fallthrough
		case "quit": close(quit); server.KillServer(); return
		default: fmt.Println("Unknown command. Type \"help\" for a list of commands.")
		}
		time.Sleep(2 * time.Millisecond)
		fmt.Printf(">> ")
	}
}

func signalLisiner(q chan(struct{}), i chan(string)) {
	sigchan := make(chan(os.Signal), 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	select{
	case <-q: return;
	case <-sigchan: fmt.Println(""); i <- "quit"; return
	}
}

func inputLisiner(q chan(struct{}), i chan(string)) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-q: return
		default: if scanner.Scan() { i <- scanner.Text() }
		}
	}
}
