package pool

import (
	"errors"
	"sync"
	"time"
)

type Pool interface {
	Init()
	Get() Conn
	Put(Conn) error
	Close()
}

// 默认连接池
type pool struct {
	mu             sync.Mutex
	conns          chan Conn
	Dial           func() Conn
	CheckConnAlive func(Conn) bool
	opts           *Options
}

func NewPool(opts []Option) *pool {
	p := &pool{
		mu:   sync.Mutex{},
		opts: &Options{},
	}

	for _, o := range opts {
		o(p.opts)
	}

	p.conns = make(chan Conn, p.opts.maxCap)
	return p
}

func (p *pool) Init() {
	for i := 0; i < p.opts.initialCap; i++ {
		p.conns <- p.Dial()
	}

	// 注册健康检查协程
	p.RegisterChecker(p.opts.internal, p.Checker)
}

func (p *pool) Get() (conn Conn) {
	select {
	case conn := <-p.conns:
		conn.SetUseTime(time.Now())
		return conn
	case <-time.After(p.opts.waitTimeout):
		return p.Dial()
	}
}

func (p *pool) Put(conn Conn) error {
	if conn == nil {
		return errors.New("connection closed")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conns == nil {
		return conn.Close()
	}

	select {
	case p.conns <- conn:
		return nil
	default:
		// 连接池已满直接关闭
		return conn.Close()
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.Dial = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)

	for conn := range p.conns {
		conn.Close()
	}
}

func (p *pool) RegisterChecker(internal time.Duration, checker func(conn Conn) bool) {
	if internal <= 0 || checker == nil {
		return
	}

	go func() {

		for {

			time.Sleep(internal)
			length := len(p.conns)

			for i := 0; i < length; i++ {

				select {
				case conn := <-p.conns:
					if !checker(conn) {
						break
					} else {
						p.Put(conn)
					}
				default:
					break
				}

			}
		}

	}()
}

func (p *pool) Checker(conn Conn) bool {

	// check timeout
	if conn.GetUseTime().Add(p.opts.idleTimeout).Before(time.Now()) {
		return false
	}

	// check conn is alive or not
	if !p.CheckConnAlive(conn) {
		return false
	}

	return true
}
