package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func main() {
	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Dial err")
		return
	}
	defer redisConn.Close()

	cnt := 1112

	key := "alonzo_test_hash"
	array := []int{}
	for n := 1; n < cnt; n++ {
		array = append(array, n)
		//redisConn.Do("hset", key, n, 1)
	}
	fmt.Println("len:", len(array))
	s := time.Now()

	_, err = redis.Ints(redisConn.Do("hmget", key, array))
	if err != nil {
		fmt.Println("hmget error")
		return
	}

	fmt.Println(time.Since(s))

	/*s1 := time.Now()
	for n := 1; n < cnt; n++ {
		redisConn.Send("hget", key, n)
	}

	if err = redisConn.Flush(); err != nil {
		fmt.Println("flush err")
	}

	for n := 1; n < cnt; n++ {
		redisConn.Receive()
	}
	fmt.Println(time.Since(s1))
	*/

	/*
		key := "alonzo_test"
		for n := 1; n < cnt; n++ {
			redisConn.Do("sadd", key, n)
		}
		s := time.Now()

		for n := 1; n < cnt; n++ {
			redisConn.Do("SISMEMBER", key, n)
		}

		fmt.Println(time.Since(s))

		for n := 1; n < cnt; n++ {
			redisConn.Send("SISMEMBER", n)
		}

		if err = redisConn.Flush(); err != nil {
			fmt.Println("flush err")
		}

		for n := 1; n < cnt; n++ {
			redisConn.Receive()
		}
		fmt.Println(time.Since(s))
	*/

}
