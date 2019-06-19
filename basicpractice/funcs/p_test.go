package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_exit(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		p()
	})

}

func Test_plus(t *testing.T) {
	for i := 0; i < 10; i++ {
		v := i
		msg := fmt.Sprintf("test for %v + 10", v)
		t.Run(msg, func(c *testing.T) {
			res := plus(v, 10)
			assert.Equal(t, res, v+10, msg)
		})
	}
}

func Test_mvf(t *testing.T) {
	for i := 0; i < 10; i++ {
		v := i
		msg := fmt.Sprintf("test for %v + 10", v)
		t.Run(msg, func(c *testing.T) {
			ori, dbl := mrvalue(v)
			assert.Equal(t, v, ori, msg)
			assert.Equal(t, v*2, dbl, msg)
		})
	}
}
