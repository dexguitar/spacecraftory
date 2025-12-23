package redis

import (
	"context"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/dexguitar/spacecraftory/platform/pkg/cache"
)

func (c *client) TxPipeline(ctx context.Context, fn func(cache.TxPipeliner) error) error {
	return c.withConn(ctx, func(ctx context.Context, conn redigo.Conn) error {
		if err := conn.Send("MULTI"); err != nil {
			return err
		}

		p := &txPipeliner{conn: conn}
		if err := fn(p); err != nil {
			_, connErr := conn.Do("DISCARD")
			if connErr != nil {
				return connErr
			}
			return err
		}

		_, err := conn.Do("EXEC")
		return err
	})
}

type txPipeliner struct {
	conn redigo.Conn
}

func (p *txPipeliner) HashSet(key string, values any) error {
	return p.conn.Send("HSET", redigo.Args{key}.AddFlat(values)...)
}

func (p *txPipeliner) Expire(key string, expiration time.Duration) error {
	return p.conn.Send("EXPIRE", key, int(expiration.Seconds()))
}

func (p *txPipeliner) SAdd(key, value string) error {
	return p.conn.Send("SADD", key, value)
}
