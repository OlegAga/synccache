package synccache

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Cache struct {
	sync.RWMutex

	items      map[string]*Item
	lastChange time.Time
}

type Item struct {
	Value interface{}
	TTL   time.Duration

	Created       time.Time
	Updated       time.Time
	Accessed      time.Time
	AccessedCount int64
}

func New(cleanupDuration time.Duration,
	saveDuration time.Duration, filePersist string) CacheI {
	c := &Cache{
		items: map[string]*Item{},
	}
	go cleaner(c, cleanupDuration)
	go saver(c, saveDuration, filePersist)
	return c
}

func cleaner(c *Cache, d time.Duration) {
	t := time.Tick(d)
	for _ = range t {
		c.RemoveExpired()
	}
}

func (c *Cache) RemoveExpired() {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.items {
		if v.isExpired() {
			delete(c.items, k)
		}
	}
}

func (item *Item) isExpired() bool {
	result := false
	if item.TTL != 0 && item.Created.Add(item.TTL).Sub(time.Now()) < 0 {
		result = true
	}
	return result
}

func saver(c *Cache, d time.Duration, f string) {
	t := time.Tick(d)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-t:
			err := c.Save(f)
			if err != nil {
				log.Fatalf("Error save to file %v", err)
			}
		case <-s:
			err := c.Save(f)
			if err != nil {
				log.Fatalf("Error save to file before exit %v", err)
			}
			os.Exit(1)
		}
	}
}

func (c *Cache) Save(f string) error {
	_F, err := os.Create(f)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(_F)
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Value)
	}
	err = enc.Encode(c.items)
	if err != nil {
		return err
	}
	return _F.Close()
}

func (c *Cache) Load(f string) error {
	_F, err := os.Open(f)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(_F)
	oldItems := map[string]*Item{}
	err = dec.Decode(&oldItems)
	if err != nil {
		return err
	} else {
		c.Lock()
		defer c.Unlock()
		for k, v := range oldItems {
			cur, ok := c.items[k]
			if (!ok || cur.isExpired()) &&
				!v.isExpired() {
				c.items[k] = v
			}
		}
	}
	return _F.Close()
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, errors.New("Cache have not for get key: " + key)
	}
	if item.isExpired() {
		delete(c.items, key)
		return nil, errors.New("Cache have not for get valid key: " + key)
	}
	item.Accessed = time.Now()
	item.AccessedCount += 1
	return item.Value, nil
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()
	_, ok := c.items[key]
	if ok {
		return errors.New("Cache have already for set key: " + key)
	} else {
		item := Item{}
		item.Value = value
		item.TTL = ttl
		timenow := time.Now()
		item.Created = timenow
		c.items[key] = &item
		c.lastChange = timenow
	}
	return nil
}

func (c *Cache) Update(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	_, ok := c.items[key]
	if !ok {
		return errors.New("Cache have not for update key: " + key)
	} else {
		c.items[key].Value = value
		timenow := time.Now()
		c.items[key].Updated = timenow
		c.lastChange = timenow
	}
	return nil
}

func (c *Cache) Remove(k string) error {
	c.Lock()
	defer c.Unlock()
	delete(c.items, k)
	return nil
}

func (c *Cache) Keys() string {
	c.RLock()
	defer c.RUnlock()
	var keys []string
	for k, _ := range c.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func (c *Cache) LastChange() time.Time {
	return c.lastChange
}
