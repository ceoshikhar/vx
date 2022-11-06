package main

import (
	"fmt"
	"testing"
	"time"
	"vx"
)

type simple struct {
	SimpleA any `vx:"name=simple_a, type=string"`
	SimpleB any `vx:"name=simple_b, type=[]string"`
}

type complex struct {
	ComplexA simple
	ComplexB any `vx:"name=complex_b, type=map[string]string"`
}

type trippleNestedStruct struct {
	RootA   any `vx:"type=string, required"`
	NestedB complex
}

func BenchmarkValidateStruct(b *testing.B) {
	v := trippleNestedStruct{
		RootA: 23,
		NestedB: complex{
			ComplexA: simple{
				SimpleA: 69,
				SimpleB: []int{},
			},
			ComplexB: 23,
		},
	}

	for i := 0; i < b.N; i++ {
		vx.ValidateStruct(v)
	}
}

func main() {
	res := testing.Benchmark(BenchmarkValidateStruct)
	fmt.Println(res)
	fmt.Println("Op per sec:", time.Duration(1000000000).Abs().Nanoseconds()/res.NsPerOp(), "op/s")
}
