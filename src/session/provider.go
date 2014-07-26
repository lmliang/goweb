package session

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestory(sid string) error
	SessionGC(maxLifeTime int)
}

var providers = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session;Register provider is nil")
	}

	if _, had := providers[name]; had {
		panic("Session:Register called twice for provider " + name)
	}

	providers[name] = provider
}
