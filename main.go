package main

import (
	"DataBridge/config"
	"DataBridge/console"
)

func main() {

	cfg.InitConfig(".env")

	console.RunInterface()

}
