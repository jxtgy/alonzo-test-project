package utils_test

import (
	"fmt"
	"going/utils"
	"testing"
)

func BenchmarkItoIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := utils.ItoIP(174343044)
		if s == "" {
			b.Errorf("fail")
		}
	}
}

func ExampleGetLocalIP() {
	fmt.Println("big endian local ip:", utils.GetLocalIP())
	fmt.Println("little endian local ip:", utils.GetLittleEndianLocalIP())
	//output:
	//11
}

func ExampleGCD() {
	gcd := utils.GCD(30, 10, 5)
	fmt.Println(gcd)
	//output:
	//5
}
