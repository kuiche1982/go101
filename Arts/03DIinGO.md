https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/

Dependency Injection patterns in Go
Dependency Injection is a powerful approach to providing your applications with required objects and configuration instead of declaring them in code. This injection is usually used to provide mock objects for unit testing, or to run the same application in multiple environments.

Resolving dependencies in a container object
We’ll be using the codegangsta/inject package to illustrate how dependency injection works. In our first example, we have an struct which we would like to fill out with the wanted data types.

type DependencyContainer struct {
  Db    *DatabaseClient `inject`
  Redis *RedisClient    `inject`
}
The DependencyContainer struct here is the object which would hold our database client and our redis client in this case. The inject tag specifies that these fields may be resolved using inject.Apply. A generic method for filling out your DI container may look something like this:

func getInjector() inject.Injector {
  injector := inject.New()
  injector.Map(&DatabaseClient{"Hello from DatabaseClient"})
  injector.Map(&RedisClient{"Hello from RedisClient"})
  return injector
}

func getContainer(container interface{}) error {
  injector := getInjector()
  return injector.Apply(container)
}
Resolving the dependencies is as easy as this:

c := DependencyContainer{}
err := getContainer(&c)
fmt.Printf("Database: %s\n", c.Db.Name)
fmt.Printf("Redis: %s\n", c.Redis.Name)
It resolves the pointers as it should:

Database: Hello from DatabaseClient
Redis: Hello from RedisClient
The pattern itself is simple, but it does require declaring the dependencies in a declarative way, which may be overkill for simple microservices based applications.

Injecting function arguments
The package we are using supports injector.Invoke, which resolves the arguments which the function itself requires. This means you just declare your function as you would, and you can use the injector.Invoke to call it and resolve their arguments.

func useBoth(db *DatabaseClient, redis *RedisClient) {
  fmt.Printf("[invoke] Database & Redis: %s & %s\n", db.Name, redis.Name)
}
func useRedis(redis *RedisClient) {
  fmt.Printf("[invoke] Redis: %s\n", redis.Name)
}
func useDatabase(db *DatabaseClient) {
  fmt.Printf("[invoke] Database: %s\n", db.Name)
}
Let’s try these out:

injector := getInjector()
injector.Invoke(useBoth)
injector.Invoke(useRedis)
injector.Invoke(useDatabase)
The functions are called and their parameters successfully resolved.

[invoke] Database & Redis: Hello from DatabaseClient & Hello from RedisClient
[invoke] Redis: Hello from RedisClient
[invoke] Database: Hello from DatabaseClient
An awesome use case for dependency injection here is refactoring: you can add, remove and reorder function arguments, while not breaking your application. As long as you can resolve individual arguments, you don’t have to fix every function call to reflect the changed arguments.

Using multiple objects with the same type
Inevitably, there comes a time when an application needs to utilize multiple connections to different, for example, database servers. That would require injection of specific database objects for individual connections.

Type alias
One way to approach this, which is completely valid with the inject package we use, is to declare a type alias. This way, we can provide connection context with the type declaration itself. For example:

type DatabaseClient struct {
  Name string
}
type AdminDatabaseClient DatabaseClient
type AnonDatabaseClient DatabaseClient

func getInjector() inject.Injector {
  injector := inject.New()
  injector.Map(&AdminDatabaseClient{"Admin"})
  injector.Map(&AnonDatabaseClient{"Anonymous"})
  return injector
}
func useBoth(db *AdminDatabaseClient, db2 *AnonDatabaseClient) {
  fmt.Printf("[invoke] Databases: %s & %s\n", db.Name, db2.Name)
}
func useAdmin(db *AdminDatabaseClient) {
  fmt.Printf("[invoke] Admin: %s\n", db.Name)
}
func useAnon(db *AnonDatabaseClient) {
  fmt.Printf("[invoke] Anon: %s\n", db.Name)
}

func main() {
  injector := getInjector()
  injector.Invoke(useBoth)
  injector.Invoke(useAnon)
  injector.Invoke(useAdmin)
}
This has a few advantages - one is that connections you use get declared explicitly, which means you can document them with code and use available tooling to analyze usage.

[invoke] Databases: Admin & Anonymous
[invoke] Anon: Anonymous
[invoke] Admin: Admin
Of course, this method has the pitfall that you have to declare your objects beforehand and there might be some time before they are used, if they even are used.

Lazy loading
For those of you not familiar with “Lazy loading”, it’s the practice in software development to instantiate whatever objects you need just before the moment when you need them. This avoids having to declare and instantiate all objects beforehand, regardless of if you need them or not. In many cases, this results in better efficiency of your application.

Object factory
There is a classic pattern to provide Lazy loading in other programming languages, and it applies to Go as well. I’ve chosen to provide the most simple version of an object factory. An object factory’s purpose is to return objects based on function calls. Depending on your requirements, the functions may return a new instance every time, or they may return a shared object (singleton), which you re-use in your code.

type ObjectFactory struct {
}

func (r *ObjectFactory) GetDatabase(name string) (*DatabaseClient, error) {
  switch name {
  case "admin":
    return &DatabaseClient{"Administrator"}, nil
  case "anon":
    return &DatabaseClient{"Anonymous"}, nil
  }
  return nil, fmt.Errorf("Unknown database definition: '%s'", name)
}

func getInjector() inject.Injector {
  injector := inject.New()
  injector.Map(&ObjectFactory{})
  return injector
}
func useObjectFactory(of *ObjectFactory) {
  names := []string{"admin", "anon", "missing"}
  for _, name := range names {
    db, err := of.GetDatabase(name)
    if err != nil {
      fmt.Printf("[db] got error: %s\n", err)
      continue
    }
    fmt.Printf("[db] got database '%s'\n", db.Name)
  }
}
Invoking useObjectFactory gives us the expected output:

[db] got database 'Administrator'
[db] got database 'Anonymous'
[db] got error: Unknown database definition: 'missing'
As discussed in the previous example, there are a number of pitfalls here, the lack of explicit typing being one of them. Of course, you may provide your type aliases in the object factory, but you will have a harder time finding out where this object is used with existing code analysis tools.

Lazy Get
Well, we did kind of roll our own object factory there, but we don’t have to. We can still have the best of both worlds - lazy loading and explicit declarations. We can use the inject package which provides it. But, in this example we will use the tsaikd/inject package which adds this functionality to the original by codegangsta.

import "github.com/tsaikd/inject"

type DatabaseClient struct {
  Name string
}
type AdminDatabaseClient DatabaseClient
type AnonDatabaseClient DatabaseClient

type ObjectFactory struct {
}
func (r ObjectFactory) NewAdminDatabaseClient() *AdminDatabaseClient {
  return &AdminDatabaseClient{"Administrator"};
}
func (r ObjectFactory) NewAnonDatabaseClient() *AnonDatabaseClient {
  return &AnonDatabaseClient{"Anonymous"};
}

func getInjector() inject.Injector {
  of := ObjectFactory{};
  injector := inject.New()
  injector.Provide(of.NewAdminDatabaseClient)
  injector.Provide(of.NewAnonDatabaseClient)
  return injector
}
func useBoth(db *AdminDatabaseClient, db2 *AnonDatabaseClient) {
  fmt.Printf("[invoke] Databases: %s & %s\n", db.Name, db2.Name)
}
func useAdmin(db *AdminDatabaseClient) {
  fmt.Printf("[invoke] Admin: %s\n", db.Name)
}
func useAnon(db *AnonDatabaseClient) {
  fmt.Printf("[invoke] Anon: %s\n", db.Name)
}
Using inject.Provide and passing it a function allows us to achieve lazy loading when individual objects are required. The inject package checks the return types of these provider functions to resolve objects which are required at the time of inject.Invoke function calls.

General thoughts
Dependency Injection is a strong workhorse pattern which enables you to write very expressive code, while still using the benefits of typed syntax. There is a trade off in the area of handling errors - all dependencies should resolve, and errors can only trigger panic calls. Depending on your needs, you may combine these patterns to your advantage.

As I hinted at the beginning of the article, your requirements for dependency injection might also be shaped towards the goal of unit testing. The used inject packages above support resolving interfaces, which give you everything you need to put your right foot forward when writing unit tests.

This article is part of a series, “Diving Deeper with Go”:

Golang Tips and Tricks
Advanced Go Tips and Tricks
While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.


