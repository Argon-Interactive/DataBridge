package console

import (
	cfg "DataBridge/config"
	"DataBridge/server"
	"fmt"
)


var helpStr = []string{
	"start: Starts the server", 
	"stop: Stops the server",
	"rester: Restarts the server",
	"status: Logs the servers status",
	"exit: Exits the program",
}

func printHelp() { for _, str := range helpStr { fmt.Println(str) } }

func printStatus() { 
	if server.IsRunning() {
		fmt.Printf("Server is running on port: %s\n", cfg.GetConfig().ServerPort)
		return
	}
	fmt.Println("Server is not running")
}
