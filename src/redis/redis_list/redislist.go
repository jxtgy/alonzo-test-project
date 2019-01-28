package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {
	redisConn := redisPool.Pool.Get()
	defer redisConn.Close()

	cnt := 100000
	prefix := "alonzo_test"
	key := ""
	for n := 1; n < cnt; n++ {
		redis.Int(redisConn.Do("sadd", key, n))
	}

	for n := 1; n < cnt; n++ {
		memberKey := fmt.Sprintf("%s", n)
		redis.Int(redisConn.Do("SISMEMBER", memberKey))
	}

	for n := 1; n < cnt; n++ {
		memberKey := fmt.Sprintf("%s", n)
		redisConn.Send("SISMEMBER", memberKey)
	}

	if err = redisConn.Flush(); err != nil {
		fmt.Println("flush err")
		return -1
	}

	for n := 1; n < cnt; n++ {
		redis.Int(redisConn.Receive())
	}
	fmt.Println("end")
}
