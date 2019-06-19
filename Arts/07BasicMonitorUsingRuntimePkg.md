https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/

Basic monitoring of Go apps with the runtime package
You might be wondering - especially if you’re just beginning to work with Go, how you might add monitoring to your microservice application. As people with some sort of track record will tell you - monitoring is hard. And what I’m telling you - at least basic monitoring doesn’t need to be. You don’t need to set up a Prometheus cluster to get reporting for your simple applications, in fact, you might not even need external services to add a simple printout of your apps statistics.

But which properties of our app are we interested in? The Go runtime package contains functions that interact with Go’s runtime system - like the scheduler and memory manager. This means that we can access some of the applications internals:

Goroutines
A goroutine is a very lightweight thread that the Go scheduler manages for us. A typical issue that may occur in any code is called “leaking goroutines”. The reasons for that issue vary anywhere from forgetting to set default http timeouts, sql timeouts, general lack of support for context package cancellations, sending data on closed channels, and so on. When this occurs, a goroutine may stay alive indefinitely and never release the resources that it uses.

A very basic function that we might be interested in is runtime.NumGoroutine() int, which returns the number of goroutines that are currently existing. Just by printing this number and inspecting it over some period of time, we can be reasonably sure that we might be leaking goroutines, and then investigating those issues.

Memory usage
Memory usage issues are common in the Go world. While most people tend to use pointers effectively (much more effectively than anything written in Node.js anyway), an often experienced issue in regards to performance is related to memory allocation. To demonstrate a simple, but inefficient way to reverse a string:

package main

import (
	"strings"
	"testing"
)

func BenchmarkStringReverseBad(b *testing.B) {
	b.ReportAllocs()

	input := "A pessimist sees the difficulty in every opportunity; an optimist sees the opportunity in every difficulty."

	for i := 0; i < b.N; i++ {
		words := strings.Split(input, " ")
		wordsReverse := make([]string, 0)
		for {
			word := words[len(words)-1:][0]
			wordsReverse = append(wordsReverse, word)
			words = words[:len(words)-1]
			if len(words) == 0 {
				break
			}
		}
		output := strings.Join(wordsReverse, " ")
		if output != "difficulty. every in opportunity the sees optimist an opportunity; every in difficulty the sees pessimist A" {
			b.Error("Unexpected result: " + output)
		}
	}
}

func BenchmarkStringReverseBetter(b *testing.B) {
	b.ReportAllocs()

	input := "A pessimist sees the difficulty in every opportunity; an optimist sees the opportunity in every difficulty."

	for i := 0; i < b.N; i++ {
		words := strings.Split(input, " ")
		for i := 0; i < len(words)/2; i++ {
			words[len(words)-1-i], words[i] = words[i], words[len(words)-1-i]
		}
		output := strings.Join(words, " ")
		if output != "difficulty. every in opportunity the sees optimist an opportunity; every in difficulty the sees pessimist A" {
			b.Error("Unexpected result: " + output)
		}
	}
}
view rawmain_test.go hosted with ❤ by GitHub
The bad function has unnecessary allocations, namely:

we create an empty slice to hold the resulting string,
we append to the slice (append allocates memory as needed, but not optimally)
The benchmarks and the related output because of the call to b.reportAllocs() paints an accurate picture:

BenchmarkStringReverseBad-4              1413 ns/op             976 B/op          8 allocs/op
BenchmarkStringReverseBetter-4            775 ns/op             480 B/op          3 allocs/op
Another aspect of memory allocations due to the virtual-memory implemented in Go, are garbage collection pauses, or GC for short. A common phrase used in regards to GC pauses is a “stop the world”, noting that your application will completely stop responding during a GC pause. The google team continually improve the performance of GC, but the symptoms of poor memory management by inexperienced developers will remain a problem in the future.

The runtime package exposes runtime.ReadMemStats(m *MemStats) that fills a MemStats object. There are a number of fields in that struct that might serve as a good indicator of poor memory allocation strategies and related performance issues.

Alloc - currently allocated number of bytes on the heap,
TotalAlloc - cumulative max bytes allocated on the heap (will not decrease),
Sys - total memory obtained from the OS,
Mallocs and Frees - number of allocations, deallocations, and live objects (mallocs - frees),
PauseTotalNs - total GC pauses since the app has started,
NumGC - number of completed GC cycles
Approach
So, we started with the premise that we don’t want to use an external service for providing simple app monitoring. It’s my aim just to print the collected metrics to the console every once in a while. We should spin up a goroutine that will fetch this data every X seconds, and just print it to the console.

package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}
To use it, you can just call it from main with something like go NewMonitor(300) and it will print out your apps metrics every 5 minutes. You can then inspect these every once in a while either from the console or historical logs, to see how your application behaves. Any performance impact from adding it to your app is minimal.

{"Alloc":1143448,"TotalAlloc":1143448,"Sys":5605624,"Mallocs":8718,"Frees":301,"LiveObjects":8417,"PauseTotalNs":0,"NumGC":0,"NumGoroutine":6}
{"Alloc":1144504,"TotalAlloc":1144504,"Sys":5605624,"Mallocs":8727,"Frees":301,"LiveObjects":8426,"PauseTotalNs":0,"NumGC":0,"NumGoroutine":5}
...
I think having this output in the console is an useful insight that will let you know if you might be hitting on some problems in the near future.

Using expvar
Go actually comes with two built-ins that help us with monitoring our apps in production. One of those built-ins is the package expvar. The package provides a standardized interface to public variables, such as operation counters in servers. These variables will be then available, by default, on /debug/vars. Let’s put our metrics into the expvar store.

… After a few minutes, as soon as I registered the HTTP handler for expvar, I realized that the full MemStats struct is already available on it. That’s great!

In addition to adding the HTTP handler, this package registers the following variables:

cmdline os.Args
memstats runtime.Memstats
The package is sometimes only imported for the side effect of registering its HTTP handler and the above variables. To use it this way, link this package into your program:

import _ "expvar"

As the metrics are now already exported, you only need to point your monitoring system at your app and import the memstats output there. We still don’t have the goroutine count, I realize, but that’s easy to add. Import the expvar package and add the following lines:

// The next line goes at the start of NewMonitor()
var goroutines = expvar.NewInt("num_goroutine")
// The next line goes after the runtime.NumGoroutine() call
goroutines.Set(int64(m.NumGoroutine))
The field “num_goroutine” is now available in the /debug/vars output, next to the full memory statistics.

Going beyond basic monitoring
Another powerful addition to the Go stdlib is the net/http/pprof package. The package has many functions, but the main intent is to provide runtime profiling data for the go pprof tool, which is bundled in the go toolchain. With it, you can further inspect the operation of your app in production. You can check out one of my previous articles if you want to learn more about pprof and code optimisation:

Benchmarking Go programs,
Benchmarking Go programs, part 2
And, if you want continuous profiling of Go programs, there’s a Google service for that, StackDriver Profiler. But, if you want to run monitoring on your own infrastructure for whatever reason, Prometheus might be your best bet. Enter your email below if that’s something you’d like to read about.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.