package g

import (
	"log"
	"runtime"
)

const (
	VERSION    = "1.0.5"
	LogFile    = "nqm.log"
	Running    = 0
	NotRunning = 1
)

var NQMRunningVersion = "<UNDEFINED>"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
