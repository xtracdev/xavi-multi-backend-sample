//session implements a context aware plugin that can add a session id
package session

import (
	"context"
	"github.com/xtracdev/xavi/plugin"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type sessionKey int

const SessionKey sessionKey = 111

func NewSessionWrapper(args ...interface{}) plugin.Wrapper {
	return new(SessionWrapper)
}

var mutex sync.Mutex

var seed = rand.NewSource(time.Now().UnixNano())
var gen = rand.New(seed)

type SessionWrapper struct{}

func (lw SessionWrapper) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		sessionId := gen.Intn(999999999)
		mutex.Unlock()

		newR := r.WithContext(context.WithValue(r.Context(), SessionKey, sessionId))

		h.ServeHTTP(w, newR)
	})
}
