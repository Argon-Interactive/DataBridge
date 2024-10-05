package main

import (
	"DataBridge/packages/config"
	"DataBridge/packages/console"
)

func main() {

	cfg.InitConfig(".env")

	console.RunInterface()

}
