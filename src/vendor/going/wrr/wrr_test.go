package wrr_test

import (
	"flag"
	"fmt"
	"going/wrr"
)

var m = map[string]uint32{
	"a": 5,
	"b": 1,
	"c": 1,
}

var count = flag.Int("count", 7, "wrr choose count")

func Example() {
	w1 := wrr.NewSWRR(m)
	w2 := wrr.NewWRR(m)

	fmt.Println("hello swrr")
	for i := 0; i < *count; i++ {
		fmt.Println(w1.Next())
	}
	fmt.Println("hello wrr")
	for i := 0; i < *count; i++ {
		fmt.Println(w2.Next())
	}
	//output:
	//hello swrr
	//a
	//a
	//b
	//a
	//c
	//a
	//a
	//hello wrr
	//a
	//a
	//a
	//a
	//a
	//b
	//c
}
