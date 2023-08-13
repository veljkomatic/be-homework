package subscriber

import (
	"context"
	"strings"
	"sync"
)

// Subscriber is responsible for subscribing, unsubscribing and testing addresses
type Subscriber interface {
	// Subscribe subscribes to address
	Subscribe(context context.Context, address string) error
	// UnSubscribe unsubscribes from address
	UnSubscribe(context context.Context, address string) error
	// Test tests if address is subscribed
	Test(context context.Context, address string) (bool, error)
}

var _ Subscriber = (*subscriber)(nil)

// subscription represents subscription to address
// now it is just a stub, but in future it could be more complex
// having more fields and complex logic around it
type subscription struct {
	Type  string `json:"type"`
	Event string `json:"event"`
	// TODO: add more fields
}

type subscriber struct {
	storage map[string]*subscription
	mutex   sync.RWMutex
}

func NewSubscriber() Subscriber {
	return &subscriber{
		storage: make(map[string]*subscription),
	}
}

func (s *subscriber) Subscribe(context context.Context, address string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// make sure that address is always lower case
	lowerCaseAddress := strings.ToLower(address)
	s.storage[lowerCaseAddress] = &subscription{}
	return nil
}

func (s *subscriber) UnSubscribe(context context.Context, address string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// make sure that address is always lower case
	lowerCaseAddress := strings.ToLower(address)
	delete(s.storage, lowerCaseAddress)
	return nil
}

func (s *subscriber) Test(context context.Context, address string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	// make sure that address is always lower case
	lowerCaseAddress := strings.ToLower(address)
	_, exists := s.storage[lowerCaseAddress]
	return exists, nil
}
