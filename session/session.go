package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

//session操作基础接口
type SessionStore interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	GetSessionId() string
	Release(w http.ResponseWriter) error
	ClearAll() error
}

type SessionProvider interface {
	SessionInit(gcTime int64, config string) error
	SessionGet(id string) (SessionStore, error)
	SessionExist(id string) bool
	SessionDestory(id string) error
	SessionGC()
}

var regProviders = make(map[string]SessionProvider)

//reg a session provider
func RegisterSessionProvider(key string, p SessionProvider) {
	if p == nil {
		panic("Session provider is null!")
	}
	if _, ok := regProviders[key]; ok {
		panic("Session provider " + key + " is exist!")
	}
	regProviders[key] = p
}

//conf of manager
type managerConf struct {
	CookieKey       string
	EnableSetCookie bool
	ProviderConf    string
	Domain          string
	CookieAge       int64
	SessionAge      int64
	SessionGcTime   int64
	SessionIdLength int
}

type Manager struct {
	provider SessionProvider
	conf     *managerConf
}

//new mana
func NewManager(pname, conf string) (*Manager, error) {
	provider, ok := regProviders[pname]
	if !ok {
		return nil, errors.New(pname + " not supported or not reg.")
	}
	cf := new(managerConf)
	cf.EnableSetCookie = true
	err := json.Unmarshal(conf, &cf)
	if err != nil {
		return nil, err
	}
	err = provider.SessionInit(cf.SessionAge, cf.ProviderConf)
	if err != nil {
		return nil, err
	}
	if cf.SessionIdLength == 0 {
		cf.SessionIdLength = 16
	}
	if cf.SessionAge == 0 {
		cf.SessionAge = cf.SessionGcTime
	}
	if cf.CookieKey == "" {
		cf.CookieKey = "Slimgosess"
	}
	return &Manager{
		provider: provider,
		conf:     cf,
	}, nil
}

func (this *Manager) SessionStart(request *http.Request, reponseWriter http.ResponseWriter) (sessStore SessionStore, errRtn error) {
	sessionId, err := request.Cookie(this.conf.CookieKey)
	if err != nil || sessionId.Value == "" {
		sessId, errIn := this.newSessionId()
		if errIn != nil {
			return nil, err
		}
		sessStore, errRtn = this.provider.SessionGet(sessId)
		sessionId := &http.Cookie{
			Name:     this.conf.CookieKey,
			Value:    url.QueryEscape(sessId),
			Path:     "/",
			HttpOnly: true,
			Domain:   this.conf.Domain,
		}
		if this.conf.CookieAge > 0 {
			sessionId.MaxAge = this.conf.CookieAge
			sessionId.Expires = time.Now().Add(time.Duration(this.conf.CookieAge) * time.Second)
		}
		if this.conf.EnableSetCookie {
			http.SetCookie(reponseWriter, sessionId)
		}
		request.AddCookie(sessionId)
	} else {
		sessId, errIn := url.QueryUnescape(sessionId.Value)
		if errIn != nil {
			return nil, errIn
		}
		if this.provider.SessionExist(sessId) {
			sessStore, errRtn = this.provider.SessionGet(sessId)
		} else {
			sessId, errIn = this.newSessionId()
			if errIn != nil {
				return nil, err
			}
			sessStore, errRtn = this.provider.SessionGet(sessId)
			sessionId := &http.Cookie{
				Name:     this.conf.CookieKey,
				Value:    url.QueryEscape(sessId),
				Path:     "/",
				HttpOnly: true,
				Domain:   this.conf.Domain,
			}
			if this.conf.CookieAge > 0 {
				sessionId.MaxAge = this.conf.CookieAge
				sessionId.Expires = time.Now().Add(time.Duration(this.conf.CookieAge) * time.Second)
			}
			if this.conf.EnableSetCookie {
				http.SetCookie(reponseWriter, sessionId)
			}
			request.AddCookie(sessionId)
		}
	}
}

func (this *Manager) newSessionId() (string, error) {
	b := make([]byte, this.conf.SessionIdLength)
	n, err := rand.Read(b)
	if err != nil || n != len(b) {
		return "", "new session id failed"
	}
	return hex.EncodeToString(b), nil
}

func (this *Manager) SessionDestory(request *http.Request, reponseWriter http.ResponseWriter) error {
	cookie, err := request.Cookie(this.conf.CookieKey)
	if err != nil {
		return err
	}
	err = this.provider.SessionDestory(cookie.Value)
	if err != nil {
		return err
	}
	cookieClr := &http.Cookie{
		Name:     this.conf.CookieKey,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
	http.SetCookie(reponseWriter, cookieClr)
	return nil
}

func (this *Manager) GetSessionStore(id string) (SessionStore, error) {
	ss, err := this.provider.SessionGet(id)
	if err != nil {
		return nil, err
	} else {
		return ss, nil
	}
}

func (this *Manager) GC() {
	this.provider.SessionGC()
	time.AfterFunc(time.Duration(this.conf.SessionGcTime)*time.Second, func() { this.GC() })
}
