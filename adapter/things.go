package adapter

import (
	"encoding/json"
	"errors"
	"github.com/xtracdev/xavi-multi-backend-sample/session"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/plugin/timing"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

var mutex sync.Mutex

func callThingBackend(thing string, h http.Handler, r *http.Request) string {
	recorder := httptest.NewRecorder()
	mutex.Lock()
	defer mutex.Unlock()
	h.ServeHTTP(recorder, r)
	return recorder.Body.String()
}

//HandleThings provides a handler that responds with data from the thing1 and thing2 backends.
var HandleThings plugin.MultiBackendHandlerFunc = func(m plugin.BackendHandlerMap, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	sid, ok := ctx.Value(session.SessionKey).(int)
	if ok {
		println("-----> session:", sid)
	}

	c := make(chan string)

	thing1Handler, ok := m["thing1"]
	if !ok {
		http.Error(w, "No backend named thing1 in context", http.StatusInternalServerError)
		return
	}

	thing2Handler, ok := m["thing2"]
	if !ok {
		http.Error(w, "No backend named thing2 in context", http.StatusInternalServerError)
		return
	}

	end2endTimer := timing.TimerFromContext(ctx)
	cont := end2endTimer.StartContributor("backend stuff")
	go func() { c <- callThingBackend("thing one", thing1Handler, r) }()
	go func() { c <- callThingBackend("thing two", thing2Handler, r) }()

	var results []string
	timeout := time.After(150 * time.Millisecond)
	for i := 0; i < 2; i++ {
		select {
		case result := <-c:
			results = append(results, result)
			cont.End(nil)
		case <-timeout:
			cont.End(errors.New("timeout error"))
			http.Error(w, "Timeout", http.StatusInternalServerError)
			return
		}
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(results)
	if err != nil {
		http.Error(w, "Error encoding results", http.StatusInternalServerError)
	}

}

func HandleThingsFactory(bhMap plugin.BackendHandlerMap) *plugin.MultiBackendAdapter {
	return &plugin.MultiBackendAdapter{
		BackendHandlerCtx: bhMap,
		Handler:           HandleThings,
	}
}
