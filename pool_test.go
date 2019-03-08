package magic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {

	assert := assert.New(t)

	p := NewPool(2)
	m0 := p.Get()
	assert.NotNil(m0)

	m1 := p.Get()
	assert.NotNil(m1)

	c := make(chan *Cookie, 1)
	go func() {
		m := p.Get()
		c <- m
	}()

	// test that pool is waiting
	select {
	case m := <-c:
		assert.NotNil(m)
		t.Fatal("MUST NOT get a pointer from full pool")
	case <-time.After(200 * time.Millisecond):
		assert.True(true)
	}

	// test that set will signal and release the wating gorutine
	p.Set(m0)
	select {
	case m := <-c:
		assert.NotNil(m)
	case <-time.After(1 * time.Second):
		t.Fatal("MUST get a pointer from non empty pool")
	}

	// test that you can't set more items than get
	assert.NotPanics(func() { p.Set(m1) })
	assert.Panics(func() { p.Set(m1) })

	// there is nothing to test, just call it
	// however just this found a bug (locked mutex) in the code
	assert.NotPanics(func() { p.CloseOne(m0) })
	assert.NotPanics(func() { p.CloseAll() })

	// you can't close more items than in pool
	m, err := New()
	assert.NoError(err)
	assert.Panics(func() { p.CloseOne(m) })

}
