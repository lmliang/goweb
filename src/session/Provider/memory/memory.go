package memory

import (
	"container/list"
	"session"
	"sync"
	"time"
)

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

var pder = &Provider{list: list.New()}

type SessionStore struct {
	sid        string
	accessTime time.Time
	value      map[interface{}]interface{}
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	value, ok := st.value[key]
	if ok {
		return value
	}

	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	v := make(map[interface{}]interface{}, 0)
	sess := &SessionStore{sid: sid, accessTime: time.Now(), value: v}
	element := pder.list.PushBack(sess)
	pder.sessions[sid] = element
	return sess, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	sess, ok := pder.sessions[sid]
	if ok {
		return sess.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}

	return nil, nil
}

func (pder *Provider) SessionDestory(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxLifeTime int) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}

		if (element.Value.(*SessionStore).accessTime.Unix() + int64(maxLifeTime)) < time.Now().Unix() {
			delete(pder.sessions, element.Value.(*SessionStore).sid)
			pder.list.Remove(element)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).accessTime = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}
