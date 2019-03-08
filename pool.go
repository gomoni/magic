package magic

// package pool - gorutine safe bounded pool of magic cookies

import (
	"sync"
)

// pool - pool of magic cookies, similar to sync/Pool
// 2. does not deallocate
// 3. have upper limit, if reached, gorutine is paused until resource got available again
type Pool struct {
	pool   []*Cookie
	limit  int
	getted int
	cond   *sync.Cond
}

func NewPool(limit int) *Pool {
	m := &sync.Mutex{}
	c := sync.NewCond(m)
	return &Pool{
		pool:   make([]*Cookie, 0),
		limit:  limit,
		getted: 0,
		cond:   c,
	}
}

// Get - get new magic pointer or the one from pool
//		 blocks if there are not enough pointers
//		 panics if magic creation fails
func (p *Pool) Get() *Cookie {
	p.cond.L.Lock()
	defer p.cond.L.Unlock() // XXX: in a test we catch panic, but keep the mutex locked
	//      so subsequent calls are locked indefinitelly,
	//		use defer to ensure mutex is unlocked even in panic situation
	for p.getted >= p.limit {
		p.cond.Wait()
	}
	var c *Cookie
	if len(p.pool) == 0 {
		var err error
		c, err = New()
		if err != nil {
			panic(err)
		}
	} else {
		c, p.pool = p.pool[0], p.pool[1:]
	}
	p.getted += 1
	return c
}

// Set - return magic pointer to the pool
//		 panics if there are more items in the pool
func (p *Pool) Set(m *Cookie) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock() // XXX: in a test we catch panic, but keep the mutex locked
	//      so subsequent calls are locked indefinitelly,
	//		use defer to ensure mutex is unlocked even in panic situation
	if len(p.pool) >= p.limit-1 {
		panic("Cannot set more items into pool than a limit, check the get/set API usage")
	}
	p.pool = append(p.pool, m)
	p.getted -= 1
	p.cond.Signal()
}

// closeAll - close all magic pointers and broadcast
func (p *Pool) CloseAll() {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	for _, m := range p.pool {
		m.Close()
	}
	p.pool = make([]*Cookie, 0)
	p.getted = 0
	p.cond.Broadcast()
}

// closeOne - drop the magic cookie and signal
// calls Close everytime
// panics if there are no getted cookies
func (p *Pool) CloseOne(m *Cookie) {
	m.Close()
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	if p.getted == 0 {
		panic("Cannot close more items than get'ed, check usage of get/set/closeOne API")
	}
	p.getted -= 1
	p.cond.Signal()
}
