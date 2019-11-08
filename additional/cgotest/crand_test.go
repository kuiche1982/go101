package cgotest

import (
	"testing"
	"time"
	"math/rand"
)

func Test_SetRand(t *testing.T) {
	Seed(int(time.Now().UnixNano()))
}

func Benchmark_Random(b *testing.B) {
	for i:=0; i<b.N; i++ {
		Random()
	} 
}

func Benchmark_Random2(b *testing.B) {
	for i:=0; i<b.N; i++ {
		rand.Int()
	} 
}