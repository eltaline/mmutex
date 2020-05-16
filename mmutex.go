package mmutex

import (
	"math/rand"
	"sync"
	"time"
)

type Mutex struct {
	locks       map[interface{}]interface{}
	m           *sync.RWMutex
	lockRetries  int
	lockDelay    float64
	stdtDelay   float64
	lockFactor  float64
	lockJitter  float64
}

func (m *Mutex) IsLock(key interface{}) bool {

	m.m.RLock()

	if _, ok := m.locks[key]; ok {
		m.m.RULock()
		return true
	} else {
		m.m.RULock()
		return false
	}

	return false

}

func (m *Mutex) TryLock(key interface{}) bool {

	for i := 0; i < m.lockRetries; i++ {

		m.m.Lock()

		if _, ok := m.locks[key]; ok {
			m.m.Unlock()
			time.Sleep(m.moff(i))
		} else {
			m.locks[key] = struct{}{}
			m.m.Unlock()
			return true
		}

	}

	return false

}

func (m *Mutex) UnLock(key interface{}) {

	m.m.Lock()
	delete(m.locks, key)
	m.m.Unlock()

}

func (m *Mutex) moff(retries int) time.Duration {

	if retries == 0 {
		return time.Duration(m.stdtDelay) * time.Nanosecond
	}

	moff, max := m.stdtDelay, m.lockDelay
	for moff < max && retries > 0 {
		moff *= m.lockFactor
		retries--
	}

	if moff > max {
		moff = max
	}

	moff *= 1 + m.lockJitter*(rand.Float64()*2-1)
	if moff < 0 {
		return 0
	}

	return time.Duration(moff) * time.Nanosecond

}

func NewMMutex() *Mutex {

	return &Mutex{
		locks:         make(map[interface{}]interface{}),
		m:             &sync.RWMutex{},
		lockRetries:   450,
		lockDelay:     10000000,
		stdtDelay:     10000,
		lockFactor:    1.1,
		lockJitter:    0.2,
	}

}
