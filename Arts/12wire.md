https://github.com/google/wire

Google recently released their own DI container called wire. It avoids runtime reflection by building the container using code generation. I would recommend using it rather than dig

https://blog.drewolson.org/go-dependency-injection-with-wire

Go Dependency Injection with Wire
Several months ago I wrote a blog post about dependency injection in go. In the time since that post was written, Google has released a fantastic new dependency injection container for go called wire. I much prefer wire to other containers in the go ecosystem. This post will explain why.

A (Very) Brief Primer on Dependency Injection
Dependency injection (DI) is a style of writing code such that the dependencies of a particular object (or, in go, a struct) are provided at the time the object is initialized.

Imagine we have a simple system that will take a list of URLs, perform an HTTP GET against each of the URLs, and finally concatenate the results of these requests together.

Here’s some code that solves this problem without using the DI pattern.

package main

import (
	"bytes"
	"fmt"
)

type Logger struct{}

func (logger *Logger) Log(message string) {
	fmt.Println(message)
}

type HttpClient struct {
	logger *Logger
}

func (client *HttpClient) Get(url string) string {
	client.logger.Log("Getting " + url)

	// make an HTTP request
	return "my response from " + url
}

func NewHttpClient() *HttpClient {
	logger := &Logger{}
	return &HttpClient{logger}
}

type ConcatService struct {
	logger *Logger
	client *HttpClient
}

func (service *ConcatService) GetAll(urls ...string) string {
	service.logger.Log("Running GetAll")

	var result bytes.Buffer

	for _, url := range urls {
		result.WriteString(service.client.Get(url))
	}

	return result.String()
}

func NewConcatService() *ConcatService {
	logger := &Logger{}
	client := NewHttpClient()

	return &ConcatService{logger, client}
}

func main() {
	service := NewConcatService()

	result := service.GetAll(
		"http://example.com",
		"https://drewolson.org",
	)

	fmt.Println(result)
}
Note that both the HttpClient and ConcatService structs take it upon themselves to build everything they need in their own initializers. Both HttpClient and ConcatService build their own Logger, though we may in practice only want (or need) a single Logger instance.

Our main function, on the other hand, looks nice and simple. We initialize only the ConcatService because it is the only entry point from the top of our program.

And here’s the same code written in a style that uses DI.

package main

import (
	"bytes"
	"fmt"
)

type Logger struct{}

func (logger *Logger) Log(message string) {
	fmt.Println(message)
}

type HttpClient struct {
	logger *Logger
}

func (client *HttpClient) Get(url string) string {
	client.logger.Log("Getting " + url)

	// make an HTTP request
	return "my response from " + url
}

func NewHttpClient(logger *Logger) *HttpClient {
	return &HttpClient{logger}
}

type ConcatService struct {
	logger *Logger
	client *HttpClient
}

func (service *ConcatService) GetAll(urls ...string) string {
	service.logger.Log("Running GetAll")

	var result bytes.Buffer

	for _, url := range urls {
		result.WriteString(service.client.Get(url))
	}

	return result.String()
}

func NewConcatService(logger *Logger, client *HttpClient) *ConcatService {
	return &ConcatService{logger, client}
}

func main() {
	logger := &Logger{}
	client := NewHttpClient(logger)
	service := NewConcatService(logger, client)

	result := service.GetAll(
		"http://example.com",
		"https://drewolson.org",
	)

	fmt.Println(result)
}
Now we only initialize a single Logger instance and share it across both our HttpClient and ConcatService structs. We can explicitly choose when to create new instances of our dependencies and when to reuse the same instance (creating, in effect, a singleton within our dependency graph).

Also, our structs no longer have the responsibility for building their dependencies. This means our structs are less tightly coupled to their dependencies. If our dependencies are interfaces, we may choose to provide different implementations of them depending on the context in which we are running the program. For example, we may want to provide a real HttpClient while running our program in production but a fake client when running our tests (to prevent our tests from hitting the network). This flexibility is very powerful.

However, our main has become more complicated. It now needs to be aware of every relationship in the dependency graph between our structs. This will become more burdensome and brittle as our project grows in size.

DI Containers
A Dependency injection container is designed to solve the problem of manually wiring up dependencies while allowing us to retain all the benefits of programming in a DI style. This is where wire comes into play.

Let’s modify the example above to use wire to build our dependency graph.

First, we’ll pull the code that builds the ConcatService out to a new file called container.go.

package main

func CreateConcatService() *ConcatService {
	logger := &Logger{}
	client := NewHttpClient(logger)
	return NewConcatService(logger, client)
}
We can now use our new CreateConcatService function from main.

func main() {
	service := CreateConcatService()

	result := service.GetAll(
		"http://example.com",
		"https://drewolson.org",
	)

	fmt.Println(result)
}
With these changes in place, we’ll update our container.go file to use wire.

//+build wireinject

package main

import (
	"github.com/google/wire"
)

func CreateConcatService() *ConcatService {
	panic(wire.Build(
		Logger{},
		NewHttpClient,
		NewConcatService,
	))
}
We’ve made two very important updates to this file:

We’ve added a build constraint on the first line. This means this file will only be considered when wire is building the project and not when we’re building the project in our normal workflows.
We’ve replaced the implementation of CreateConcatService with a magic-looking placeholder body. The body panics, but inside of the panic it calls the function wire.Build. We pass to wire.Build all of the constructors in our application. These arguments can be functions and structs.
With these changes in place, we can now use the wire CLI to automatically generate the correct implementation of our CreateConcatService function.

$ go get github.com/google/wire/cmd/wire
$ wire
You should see something like this output (assuming the name of your current module is example):

$ wire
example: wrote ~/example/wire_gen.go
Let’s peek inside this wire_gen.go file.

// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

// Injectors from container.go:

func CreateConcatService() *ConcatService {
	logger := &Logger{}
	httpClient := NewHttpClient(logger)
	concatService := NewConcatService(logger, httpClient)
	return concatService
}
Wire has implemented our CreateConcatService function for us. In fact, it has implemented the exact function we would have written ourselves. Note that this file has a build constraint saying this file will be considered any time we aren’t running the wire command. Also note that it conveniently includes a go:generate comment. This means in the future we can regenerate files created by wire simply by running go generate ./....

How did wire know how to write this function? When we ran the wire command, it looked for any files tagged with the //+build wireinject build constraint. Inside of these files, it looked for any function that had a call to wire.Build inside of its body. Finally, it looked at the return type of that function and all of the arguments provided to that wire.Build call. It now knew what type we ultimately wanted the function to return and all the building blocks of our dependency graph, along with their inputs and outputs. Given all of this, wire wrote the function exactly as we would have written it.

Benefits of Wire
One very important benefit of wire when compared to other DI containers in the go ecosystem is explicitness. Because wire uses code generation, the resulting container code is obvious and readable. We can look inside of the wire_gen.go file for any package and see exactly what is being built. It also means that the container is “just go code” and gets all the benefits of compile time type safety that all other go code enjoys.

There are a number of other features in wire that we haven’t even touched on. wire allows easy handling of error values within your dependency graph as well as clean up functions (for things like database connections and files).

Finally, once the code has been generated, wire containers can be simply left alone. If you don’t change your dependency graph, developers within your application never even need to know that wire is being used at all. In fact, the generated wire_gen.go file doesn’t have any dependency on wire. Even when updates need to happen, wire is a good go citizen and leverages go:generate to make the process more intuitive.

Wrapping Up
There have been countless discussions about the benefits using a dependency injection approach when building software. Even so, the go community has been highly skeptical of DI containers in the past. Often this skepticism boils down the notion that the benefits of a container do not outweigh its costs.

The authors of wire have clearly heard these complaints and tackled them head-on. The resulting code is not “magic”, it is just plain go code. The resulting code does not fail at runtime, it can leverage compile time type safety. The resulting code is readable – it’s literally what you would have typed yourself.

Given these design decisions, wire seems like a perfect fit for the go community. I sincerely hope wire becomes a de facto choice for any large go application.



