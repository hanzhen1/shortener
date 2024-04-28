package sequence

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 基于redis实现一个发号器

type Redis struct {
	//redis连接
	conn redis.Redis
}

func NewRedis(redisAddr string) Sequence {
	return &Redis{
		conn: *redis.New(redisAddr),
	}
}

// Next 取下个号
func (r *Redis) Next() (res uint64, err error) {
	//使用redis实现发号器
	//incr 将 key 中储存的数字值增一 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作
	//Pipelined 方法，它会在函数退出时调用 Exec
	var incr *redis.IntCmd
	err = r.conn.PipelinedCtx(context.Background(), func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(context.Background(), "Sequence:Number transmitter ID")
		return nil
	})
	if err != nil {
		panic(err)
	}
	// 在pipeline执行后获取到结果
	return uint64(incr.Val()), nil
}
