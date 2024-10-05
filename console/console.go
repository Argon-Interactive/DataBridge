package console

import (
	"DataBridge/server"
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
		case "run": server.RunHTTPServer()
		case "run-s": server.RunHTTPSServer()
		case "stop": server.KillServer()
		case "restart": restart()
		case "help": printHelp()
		case "status": printStatus()
		case "test": test()
		case "exit": close(quit); server.KillServer(); return
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
	case <-sigchan: i <- "exit"; return
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
