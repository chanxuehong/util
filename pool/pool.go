// 1.3 之前的 对象池 简单实现
package pool

import (
	"errors"
)

type Pool struct {
	newFunc func() interface{}
	buffer  chan interface{}
}

func New(New func() interface{}, size int) (*Pool, error) {
	if New == nil {
		return nil, errors.New("New == nil")
	}
	if size < 1 {
		return nil, errors.New("size < 1")
	}

	p := &Pool{
		newFunc: New,
		buffer:  make(chan interface{}, size),
	}
	return p, nil
}

func (p *Pool) Get() interface{} {
	select {
	case v := <-p.buffer:
		return v
	default:
		return p.newFunc()
	}
}

func (p *Pool) Put(x interface{}) {
	select {
	case p.buffer <- x:
	default:
	}
}
