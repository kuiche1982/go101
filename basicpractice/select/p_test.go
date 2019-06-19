package main

import (
	"testing"
)

func Test_exit(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		p()
	})

}

func Test_timeout(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		timeout()
	})

}

func Test_nonblocking(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		nonblocking()
	})

}

func Test_timerfun(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		timerfun()
	})

}

func Test_tickerfun(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		tickerfun()
	})

}
