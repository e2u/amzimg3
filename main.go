package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"e2u.io/amzimg3/lib"
)

var (
	allowSourceFile string
	port            uint
	address         string
	baseDir         string
	logfile         string
)

func init() {

	flag.StringVar(&allowSourceFile, "allow", "/opt/amzimg3/etc/allow_sources.txt", "allow srouce list file")
	flag.StringVar(&baseDir, "data", "/var/data", "cache image storage directory")
	flag.StringVar(&address, "address", "0.0.0.0", "listen address")
	flag.StringVar(&logfile, "logfile", "/var/log/amzimg3/stdout.log", "logfile")
	flag.UintVar(&port, "port", 8085, "listen port")
	flag.Parse()

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0666)
	if err != nil {
		panic(err)
	}
	log.SetFlags(0)
	log.SetOutput(f)
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

	fmt.Println("server started...")

	<-sigchan
}
