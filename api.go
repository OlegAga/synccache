package synccache

import (
	"time"
)

// Cache is the contract for all of the cache backends that are supported by
// this package
type CacheI interface {

	// Get returns single item from the backend if the requested item is not
	// found, returns NotFound err
	Get(key string) (interface{}, error)

	// Set sets a single item to the backend
	Set(key string, value interface{}, ttl time.Duration) error

	// Set sets a single item to the backend
	Update(key string, value interface{}) error

	// Delete deletes single item from backend
	Remove(key string) error

	// Get keys
	Keys() string

	// Persist db
	Save(file string) error

	// Ressurect db
	Load(file string) error

	// Clean up db
	RemoveExpired()

	// Time last change
	LastChange() time.Time
}
