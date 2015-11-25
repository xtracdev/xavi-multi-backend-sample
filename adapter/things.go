package adapter

import (
	"encoding/json"
	"github.com/xtracdev/xavi/plugin"
	"net/http"
	"net/http/httptest"
	"time"
)


func callThingBackend(h http.Handler, r *http.Request) string {
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, r)
	return recorder.Body.String()
}


//HandleThings provides a handler that responds with data from the thing1 and thing2 backends.
var HandleThings plugin.MultiBackendHandlerFunc = func(m plugin.BackendHandlerMap, w http.ResponseWriter, r *http.Request) {
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

	go func() { c <- callThingBackend(thing1Handler,r) }()
	go func() { c <- callThingBackend(thing2Handler,r) }()

	var results []string
	timeout := time.After(150 * time.Millisecond)
	for i := 0; i < 2; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
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
		Ctx:     bhMap,
		Handler: HandleThings,
	}
}
