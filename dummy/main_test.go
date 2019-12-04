package main

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

)

var (
	systemTest *bool

	coverCounters = make(map[string][]uint32)
	coverBlocks   = make(map[string][]testing.CoverBlock)
)

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

func TestSystem(t *testing.T) {

	fmt.Println("in TestSystem:")
	if *systemTest {

		pkgs := []string{"kuitest/dummy"}

		myCover := testing.Cover{
			Mode:            "count",
			Counters:        coverCounters,
			Blocks:          coverBlocks,
			CoveredPackages: covered(pkgs...),
		}

		testing.RegisterCover(myCover)

		go func() {
			defer func() {
				if r := recover(); r!= nil {
					fmt.Printf("err: %v\n", r)
				}
			}()

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					fmt.Printf("cov: [%f] %+v\n", testing.Coverage(), myCover)
					// we can get the data in this way, so we can send the data to any collector :D
				}
			}
		}()

		main()
	}
}

func covered(pkgs ...string) string {
	return " in " + strings.Join(pkgs, ", ")
}
