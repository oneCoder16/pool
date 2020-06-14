package main

import (
	"errors"
	"fmt"
	"github.com/oneCoder16/pool"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type P struct {
	pools map[string]pool.Pool
	mu    sync.RWMutex
}

func (p *P) RegisterDbPool(name string, opts []pool.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()

	dbPool := pool.NewPool(opts)

	dbPool.Dial = func() pool.Conn {
		dbOpt := &Options{
			UserName: "root",
			Psw:      "19951206",
			Addr:     "127.0.0.1",
			Name:     "test_db",
		}

		db, err := OpenDB(dbOpt)
		if err != nil {
			fmt.Printf("%s\ng", err)
		}
		return pool.NewDbWrapConn(db)
	}
	dbPool.CheckConnAlive = func(conn pool.Conn) bool {
		return true
	}

	dbPool.Init()

	p.pools[name] = dbPool
}

func (p *P) Select(name string) (pc pool.Pool, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if pc, ok := p.pools[name]; ok {
		return pc, nil
	}

	return nil, errors.New("not found pool")
}

func main() {

	var poolOpts = []pool.Option{
		pool.WithInitialCap(2),
		pool.WithMaxCap(3),
		pool.WithIdleTimeout(5 * time.Second),
		pool.WithInternal(1 * time.Second),
		pool.WithWaitTimeout(1 * time.Second),
	}

	p := &P{
		pools: make(map[string]pool.Pool),
		mu:    sync.RWMutex{},
	}

	p.RegisterDbPool("test_db", poolOpts)
	pc, err := p.Select("test_db")

	if err != nil {
		fmt.Println(err)
	}

	conn1 := pc.Get()
	conn2 := pc.Get()
	conn3 := pc.Get()
	conn4 := pc.Get()

	pc.Put(conn1)
	pc.Put(conn2)
	pc.Put(conn3)
	pc.Put(conn4)
	fmt.Println("ok")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	<-quit
	fmt.Println("exit")
}
