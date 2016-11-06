package main

import (
	"runtime"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/analyzer"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/server"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	util.InitConfig()
	go server.Start()
	analyzer.Start()
}
