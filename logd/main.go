package main

import (
	"flag"
	"fmt"
	"github.com/siddontang/go-log/log"
	"os"
)

var logFile = flag.String("logfile", "./logd.log", "file to log")
var logType = flag.String("logtype", "common", "log file handler type: common, rotating, timerotating")
var rotatingMaxSize = flag.Int("size", 1024*1024*1024, "rotating log file max file size")
var rotatingBackup = flag.Int("backup", 0, "rotating log file max count")
var timeRotatingWhen = flag.Int("when", 0, "timerotating when type:0 second, 1 minute, 2 hour, 3 day")
var timeRotatingInterval = flag.Int("interval", 1, "timerotating interval")
var addr = flag.String("addr", "127.0.0.1:11183", "server listen address")

func main() {
	flag.Parse()

	var h log.Handler
	var err error

	switch *logType {
	case "common":
		h, err = log.NewFileHandler(*logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND)
	case "rotating":
		h, err = log.NewRotatingFileHandler(*logFile, *rotatingMaxSize, *rotatingBackup)
	case "timerotating":
		h, err = log.NewTimeRotatingFileHandler(*logFile, int8(*timeRotatingWhen), *timeRotatingInterval)
	default:
		fmt.Printf("invalid log type %s\n", *logType)
		return
	}

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	s, err := newServer(*addr, h)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	s.Run()
}
