package cashe

import (
	"sync"
	"task/pkg/model"
)

type Cashe struct {
	sync.RWMutex
	items map[string]model.Order
}

func New() *Cashe {
	orders := make(map[string]model.Order)

	return &Cashe{items: orders}
}

func (c *Cashe) Set(key string, value model.Order) {

	c.Lock()
	defer c.Unlock()
	c.items[key] = value
}

func (c *Cashe) Get(key string) (model.Order, bool) {

	c.RLock()
	defer c.RUnlock()

	if item, ok := c.items[key]; ok {
		return item, true
	}

	return model.Order{}, false
}
