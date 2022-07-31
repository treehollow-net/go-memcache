package go_memcache

import (
	"sync"
	"time"
)

type MemCache struct {
	sync.RWMutex // for control concurrency
	Items        map[string]*Item
	GCInterval   int32 // Garbage Collection
}

func New(GCInterval int32) *MemCache {
	m := &MemCache{
		Items:      make(map[string]*Item),
		GCInterval: GCInterval,
	}
	go m.GC()
	return m
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

func (m *MemCache) GC() {
	for {
		select {
		case <-time.After(time.Duration(m.GCInterval) * time.Second):
			var keys []string

			m.Lock()
			for key, item := range m.Items {
				if item.IsExpired() {
					keys = append(keys, key)
				}
			}

			for _, key := range keys {
				m.Remove(key)
			}
			m.Unlock()
		}
	}
}
