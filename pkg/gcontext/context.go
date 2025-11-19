package gcontext

import (
	"sync"
)

type GContext interface {
	Set(key string, value any)
	Get(key string) any
	Del(key string)
}

var gContext *GlobalContext

type GlobalContext struct {
	flagset map[string]any
	mu      sync.RWMutex
}

func NewGlobalContext() GContext {
	gContext = &GlobalContext{
		flagset: make(map[string]any),
	}
	return gContext
}

func GetGlobalContext() GContext {
	return gContext
}

func (c *GlobalContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flagset[key] = value
}

func (c *GlobalContext) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flagset[key]
}

func (c *GlobalContext) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.flagset, key)
}
