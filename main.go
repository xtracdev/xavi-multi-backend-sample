package main

import (
	"os"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/multi-backend-sample/adapter"
)

func main() {
	runner.Run(os.Args[1:], func(){
		plugin.RegisterMultiBackendAdapterFactory("handle-things", adapter.HandleThingsFactory)
	})
}
