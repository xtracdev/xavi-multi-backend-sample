package main

import (
	"github.com/xtracdev/xavi-multi-backend-sample/adapter"
	"github.com/xtracdev/xavi-multi-backend-sample/session"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/plugin/recovery"
	"github.com/xtracdev/xavi/plugin/timing"
	"github.com/xtracdev/xavi/runner"
	"os"
)

func main() {
	runner.Run(os.Args[1:], func() {
		plugin.RegisterMultiBackendAdapterFactory("handle-things", adapter.HandleThingsFactory)
		plugin.RegisterWrapperFactory("SessionId", session.NewSessionWrapper)
		plugin.RegisterWrapperFactory("Timing", timing.NewTimingWrapper)
		plugin.RegisterWrapperFactory("Recovery", recovery.NewRecoveryWrapper)
	})
}
