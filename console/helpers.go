package console

import (
	cfg "DataBridge/config"
	"DataBridge/server"
	"fmt"
)


var helpStr = []string{
	"run: Starts the HTTP server", 
	"run-s: Starts the HTTPS server", 
	"stop: Stops the server",
	"restart: Restarts the server",
	"status: Logs the servers status",
	"exit: Exits the program",
}

func printHelp() { for _, str := range helpStr { fmt.Println(str) } }

func printStatus() { 
	if server.IsRunning() {
		var protocol string
		if server.GetStatus() == server.RunningHTTPS { protocol = "HTTPS" } else { protocol = "HTTP" }
		fmt.Printf("%s Server is running on port: %s\n", protocol, cfg.GetConfig().ServerPort)
		return
	}
	fmt.Println("Server is not running")
}

func restart() {
	if !server.IsRunning() { fmt.Println("Server is not running"); return; }
	state := server.GetStatus()
	server.KillServer()
	switch state {
	case server.RunningHTTP: server.RunHTTPServer()
	case server.RunningHTTPS: server.RunHTTPSServer()
	}
}

func test() {
}
