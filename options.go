package pool

import "time"

type Options struct {
	initialCap  int           // 初始化数量
	maxCap      int           // 最大值
	idleTimeout time.Duration // 空闲超时时间
	internal    time.Duration // 健康检查时间间隔
	waitTimeout time.Duration // 获取连接池等待时间
}

type Option func(*Options)

func WithInitialCap(size int) Option {
	return func(o *Options) {
		o.initialCap = size
	}
}

func WithMaxCap(size int) Option {
	return func(o *Options) {
		o.maxCap = size
	}
}

func WithIdleTimeout(time time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = time
	}
}

func WithInternal(time time.Duration) Option {
	return func(o *Options) {
		o.internal = time
	}
}

func WithWaitTimeout(time time.Duration) Option {
	return func(o *Options) {
		o.waitTimeout = time
	}
}
