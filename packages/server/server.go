package server

import (
	"DataBridge/packages/config"
	"DataBridge/packages/database"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

var echoCtx* echo.Echo
var status int32

func RunHTTPSServer() {
	if IsRunning() { fmt.Println("Server is already running. Ignoring."); return }
	setStatus(Stoped)
	dbmgr.InitDB()
	echoCtx = echo.New()
	initServer()
	go func() { 
		setStatus(RunningHTTPS)
		fmt.Printf("Server is running on port: %s\n", cfg.GetConfig().ServerPort)
		echoCtx.StartTLS(":" + cfg.GetConfig().ServerPort, "server.crt", "server.key");
		fmt.Println("Server has shut down.")
		setStatus(Stoped)
	}()
}

func RunHTTPServer() {
	if IsRunning() { fmt.Println("Server is already running. Ignoring."); return }
	setStatus(Stoped)
	dbmgr.InitDB()
	echoCtx = echo.New()
	initServer()
	go func() { 
		setStatus(RunningHTTP)
		fmt.Printf("Server is running on port: %s\n", cfg.GetConfig().ServerPort)
		echoCtx.Start(":" + cfg.GetConfig().ServerPort);
		fmt.Println("Server has shut down.")
		setStatus(Stoped)
	}()
}

func KillServer() {
	if !IsRunning() { return }
	ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
	defer cancel()
	echoCtx.Shutdown(ctx)
	dbmgr.CloseDB()
}

type serverStatus int32

const (
	Stoped serverStatus = iota
	RunningHTTP
	RunningHTTPS
)

func GetStatus() serverStatus { atomic.LoadInt32(&status); return serverStatus(status) }
func IsRunning() bool { return GetStatus() != Stoped }
func setStatus(s serverStatus) { atomic.StoreInt32(&status, int32(s)) } 
