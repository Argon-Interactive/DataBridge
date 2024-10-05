package server

import (
	"DataBridge/config"
	"DataBridge/database"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

var echoCtx* echo.Echo
var serverState int32

func RunServer() {
	if IsRunning() { fmt.Println("Server is already running. Ignoring."); return }
	serverState = 0
	dbmgr.InitDB()
	echoCtx = echo.New()
	initServer()
	go func() { 
		setRunning(true)
		fmt.Printf("Server is running on port: %s\n", cfg.GetConfig().ServerPort)
		echoCtx.Start(":" + cfg.GetConfig().ServerPort);
		fmt.Println("Server has shut down.")
		setRunning(false)
	}()
}

func KillServer() {
	if !IsRunning() { return }
	ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
	defer cancel()
	echoCtx.Shutdown(ctx)
	dbmgr.CloseDB()
}

func IsRunning() bool { atomic.LoadInt32(&serverState); return serverState == 1 }
func setRunning(b bool) { var i int32; i = 0; if(b) { i = 1 }; atomic.StoreInt32(&serverState, i) } 
