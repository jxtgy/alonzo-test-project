package main

import (
	"fmt"
)

func main() {/*
//123
	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Dial err")
		return
	}
	defer redisConn.Close()

	cnt := 100

	key := "alonzo_test_hash"
	array := RenderInts()

	fmt.Println("len:", len(array))
	s := time.Now()

	_, err = redis.Ints(redisConn.Do("hmget", key, array))
	if err != nil {
		fmt.Println("hmGet error")
		return
	}

	fmt.Println(time.Since(s))
*/
	m:= map[int]string{}
	m[1]="a"
	for k := range m {
		fmt.Println(k)
	}


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
