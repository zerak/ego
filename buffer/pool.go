package buffer

import (
	"sync"
	"time"
)

const (
	MaxGCTime = time.Second
)

type Pool struct {
	pool       *sync.Pool
	collection chan *Buffer
}

func (p *Pool) Get() (buf *Buffer, fromPool bool) {
	bufObj := p.pool.Get()
	if bufObj == nil {
		buf = NewBuffer()
	} else {
		fromPool = true
		buf = bufObj.(*Buffer)
		buf.Reset()
	}
	return
}

func (p *Pool) Put(buf *Buffer, fromPool bool) {
	if fromPool && buf != nil {
		select {
		case p.collection <- buf:
		default:
			// put buffer fail
			// this case will use temporary obj
			// and may cause memory fragmentation
		}
	}
}

func (p *Pool) run() {
	for buf := range p.collection {
		if buf.GC(MaxGCTime) {
			p.pool.Put(buf)
		}
	}
}

func NewPool(initSize, maxSize int) *Pool {
	if maxSize < initSize {
		maxSize = initSize
	}
	p := &Pool{
		pool:       &sync.Pool{New: func() interface{} { return NewBuffer() }},
		collection: make(chan *Buffer, maxSize),
	}
	for i := 0; i < initSize; i++ {
		p.pool.Put(p.pool.New())
	}
	go p.run()
	return p
}
