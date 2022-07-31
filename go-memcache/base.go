package go_memcache

import "sync"

type MemCache struct {
	sync.RWMutex // for control concurrency
	Items        map[string]*Item
	GCInterval   int32 // Garbage Collection
}

func (m *MemCache) Remove(key string) {
	m.Lock()
	delete(m.Items, key)
	m.Unlock()
}

func (m *MemCache) Get(key string) (isOk bool, value string) {
	m.RLock()
	var item *Item
	item, isOk = m.Items[key]
	m.RUnlock()
	if isOk {
		if !isExpired(item.CreatedAt, item.Expiration) {
			return true, item.Value
		}
		m.Lock()
		item, isOk = m.Items[key]
		if isOk && isExpired(item.CreatedAt, item.Expiration) {
			delete(m.Items, key)
		}
		m.Unlock()
	}
	return false, ""
}

func (m *MemCache) Set(key string, value string, expiration int32) {
	m.Lock()
	m.Items[key] = &Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
		CreatedAt:  nowTimeStamp(),
	}
	m.Unlock()
}

func (m *MemCache) Contains(key string) bool {
	m.RLock()
	_, isOk := m.Items[key]
	m.RUnlock()
	return isOk
}
