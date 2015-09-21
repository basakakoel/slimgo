package session

import (
	"crypto/cipher"
	"net/http"
	"sync"
)

var sessionProviderCookie = &SessionProvider{}

type SessionStoreCookie struct {
	id     string
	values map[interface{}]interface{}
	lock   sync.RWMutex
}

type SessionProviderCookie struct {
	conf  *configCookie
	age   int
	block cipher.Block
}

type configCookie struct {
	SecurityKey  string
	BlockKey     string
	SecurityName string
	CookieName   string
	Maxage       int
}

func (this *SessionStoreCookie) Set(key, value interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.values[key] = value
	return nil
}

func (this *SessionStoreCookie) Get(key interface{}) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.values[key]; ok {
		return v
	} else {
		return nil
	}
}

func (this *SessionStoreCookie) Delete(key interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.values, key)
	return nil
}

func (this *SessionStoreCookie) GetSessionId() string {
	return this.id
}

func (this *SessionStoreCookie) Release(responseWriter http.ResponseWriter) {

}
