package main

import "testing"

func Test_exit(t *testing.T) {
	for i := 0; i < 10; i++ {
		v := i
		t.Run("run exit", func(t *testing.T) {
			p(v)
		})
	}

}

func Test_q(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		q()
	})
}

func Test_w(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		w()
	})
}

func Test_whatAmI(t *testing.T) {
	t.Run("run exit", func(t *testing.T) {
		whatAmI(true)
		whatAmI(1)
		whatAmI("hey")
	})
}
