// go1.3 之前的 对象池 简单实现, go1.3 有了自己的对象池 sync.Pool
package pool

type Pool struct {
	newFunc func() interface{}
	buffer  chan interface{}
}

func New(New func() interface{}, size int) *Pool {
	return &Pool{
		newFunc: New,
		buffer:  make(chan interface{}, size),
	}
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
