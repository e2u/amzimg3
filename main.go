package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"e2u.io/amzimg3/lib"
)

var allowSourceFile string
var port uint
var address string
var baseDir string

func init() {
	flag.StringVar(&allowSourceFile, "allow", "/opt/amzimg3/etc/allow_sources.txt", "allow srouce list file")
	flag.StringVar(&baseDir, "dir", "/var/data", "cache image storage directory")
	flag.StringVar(&address, "address", "0.0.0.0", "listen address")
	flag.UintVar(&port, "port", 8085, "listen port")
	flag.Parse()
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		server := lib.NewServer()
		server.Address = address
		server.Port = port
		server.BaseDir = baseDir
		server.Start()
	}()

	go func() {
		lib.AllowRemoteSource = lib.NewAllowSourceByFile(allowSourceFile)
	}()

	log.Println("server started...")

	<-sigchan
}
