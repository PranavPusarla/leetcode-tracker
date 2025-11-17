package main

import "flag"

import (
	"github.com/computer-geek64/leetcode-tracker/config"
	"github.com/computer-geek64/leetcode-tracker/server"
)

func main() {
	var listenAddress string
	var configFilepath string
	flag.StringVar(&listenAddress, "l", "127.0.0.1", "Address to bind HTTP server to")
	flag.StringVar(&listenAddress, "listen", "127.0.0.1:8000", "Address to bind HTTP server to")
	flag.StringVar(&configFilepath, "c", "config/config.yaml", "Path to configuration file")
	flag.StringVar(&configFilepath, "config", "config/config.yaml", "Path to configuration file")
	flag.Parse()

	var conf = config.Load(configFilepath)
	var httpServer = server.NewServer(conf)
	httpServer.Run(listenAddress)
}
