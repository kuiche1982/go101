package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_plus(t *testing.T) {
	for i := 0; i < 10; i++ {
		v := i
		msg := fmt.Sprintf("test for %v + 10 + 10", v)
		t.Run(msg, func(c *testing.T) {
			res := plus(v, 10, 10)
			assert.Equal(t, res, v+20, msg)
		})
	}
}

func Test_plusExpandArray(t *testing.T) {
	ia := []int{1, 2, 3, 4}
	res := plus(ia...)
	assert.Equal(t, res, 10, "expanding array")
}
