package main

import (
	"os"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi-multi-backend-sample/adapter"
	"github.com/xtracdev/xavi-multi-backend-sample/session"
)

func main() {
	runner.Run(os.Args[1:], func(){
		plugin.RegisterMultiBackendAdapterFactory("handle-things", adapter.HandleThingsFactory)
		plugin.RegisterWrapperFactory("SessionId", session.NewSessionWrapper)
	})
}
