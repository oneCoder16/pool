# pool
golang 连接池
# 连接池接口
```
type Pool interface {
	Init()
	Get() Conn
	Put(Conn) error
	Close()
}
```
# 相关问题
## 连接如何管理

1. 通过配置管理连接池
```
  type Options struct {
    initialCap  int           // 初始化数量
    maxCap      int           // 最大值
    idleTimeout time.Duration // 空闲超时时间，健康检查程序将手动释放连接
    internal    time.Duration // 健康检查时间间隔
    waitTimeout time.Duration // 获取连接池等待时间
  }
```
2. 使用管道存储连接
## 连接创建过程
1. 初始化时根据 `initialCap` 配置生成连接
2. 获取连接时，超过 `waitTimeout` 时间后将手动创建连接
## 连接销毁过程
1. 健康检查检测连接超过空闲超时时间或 server 端已经将连接关闭，则将连接关闭掉
## 连接如何进行复用
1. 程序使用完后，需要手动指定 `Put` 方法将连接放回连接池中
## 如何做到并发安全
1. 使用管道数据结构实现
