package main

import (
	"fmt"
	"testing"
	"time"
	"vx"
)

type user struct {
	Name         any `vx:"name=name, type=string, required, minLength=3"`
	Age          any `vx:"name=age, type=float64, required"`
	Location     any `vx:"name=location, type=[]string"`
	AssocOrBonus any `vx:"type=map[string]string"`
}

func BenchmarkValidateStruct(b *testing.B) {
	u := user{
		Name:         "shikhar",
		Age:          20,
		Location:     []string{"India"},
		AssocOrBonus: map[string]string{"key": "value"},
	}

	for i := 0; i < b.N; i++ {
		vx.ValidateStruct(u)
	}
}

func main() {
	res := testing.Benchmark(BenchmarkValidateStruct)
	fmt.Println(res)
	fmt.Println("Op per sec:", time.Duration(1000000000).Abs().Nanoseconds()/res.NsPerOp(), "op/s")
}
