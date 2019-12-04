dig
A Better main()
Now that we know how the dig container works, let’s use it to build a better main.

func BuildContainer() *dig.Container {
  container := dig.New()

  container.Provide(NewConfig)
  container.Provide(ConnectDatabase)
  container.Provide(NewPersonRepository)
  container.Provide(NewPersonService)
  container.Provide(NewServer)

  return container
}

func main() {
  container := BuildContainer()

  err := container.Invoke(func(server *Server) {
    server.Run()
  })

  if err != nil {
    panic(err)
  }
}
The only thing we haven’t seen before here is the error return value from Invoke. If any provider used by Invoke returns an error, our call to Invoke will halt and that error will be returned.

Even though this example is small, it should be easy to see some of the benefits of this approach over our “standard” main. These benefits become even more obvious as our application grows larger.

One of the most important benefits is the decoupling of the creation of our components from the creation of their dependencies. Say, for example, that our PersonRepository now needs access to the Config. All we have to do is change our NewPersonRepository constructor to include the Config as an argument. Nothing else in our code changes.

Other large benefits are lack of global state, lack of calls to init (dependencies are created lazily when needed and only created once, obviating the need for error-prone init setup) and ease of testing for individual components. Imagine creating your container in your tests and asking for a fully-build object to test. Or, create an object with mock implementations of all dependencies. All of these are much easier with the DI approach.


https://github.com/uber-go/dig



A reflection based dependency injection toolkit for Go.

Good for:
Powering an application framework, e.g. Fx.
Resolving the object graph during process startup.
Bad for:
Using in place of an application framework, e.g. Fx.
Resolving dependencies after the process has already started.
Exposing to user-land code as a Service Locator.
Installation
We recommend consuming SemVer major version 1 using your dependency manager of choice.

$ glide get 'go.uber.org/dig#^1'
$ dep ensure -add "go.uber.org/dig@v1"
$ go get 'go.uber.org/dig@v1'
Stability
This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0.


