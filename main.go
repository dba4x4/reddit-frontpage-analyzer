package main

import (
	"runtime"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/analyzer"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/server"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	util.InitConfig()
	go server.Start()
	analyzer.Start()
}
