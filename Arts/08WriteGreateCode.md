https://scene-si.org/2018/07/24/writing-great-go-code/

Writing Great Go Code
After writing two books on the subject (API Foundations in Go and 12 Factor Applications with Docker and Go) and years of writing various Go microservices, I do want to put some thoughts down as to what it takes to write great Go code.

But first of all, let me explain this very plainly to all you who are reading this. Writing great code is subjective. You might have completely different ideas as to what is considered great code, and we might only agree on some points. On the other hand, neither of us may be wrong in such a case, we’re just coming from two different perspective viewpoints, and may have chosen to solve our engineering problems in a different way.

Packages
This is very important, and you might strongly disagree - if you’re writing Go microservices, you may keep all your code in one package. There are of course very strong and valid viewpoints to the opposite, some of which are:

Keep a separate package only for your defined types,
Maintain a service layer which is transport agnostic,
Maintain a repository layer in addition to your service layer
If you’re counting, the minimum package count for a microservice is 1. If you have a large microservice, with websocket and HTTP gateways, you may end up with a minimum of 5 packages (types, repository, service, websocket and http packages).

The simple microservice doesn’t really care about abstracting business logic away from the data storage layer (repository), or even from the transport layer (websocket, http). You write your code, it digests data and spits out your responses. However, adding more packages in the mix, solves a few issues. For example, if you’re familiar with SOLID principles, the ’S’ stands for “single responsibility”. If we break everything down into packages, these become their responsibilities:

types - declares structs and possibly some mutators of these structs,
repository - it’s a data storage layer that deals with storing and reading structs,
service - would be the implementation of business logic that wraps repositories,
http, websocket, … - the transport layers, which all invoke the service layer
Of course, depending on your use case, it may make sense to break these down even further, for example you could have types/request and types/response that would separate some structs better. This way you can have request.Message and response.Message instead of MessageRequest and MessageResponse. It may make more sense if those have been separated from the start.

But, to push home the original point - don’t feel bad if you’re using only some of these package declarations. Big software like Docker uses only a types package under it’s server package, and that’s all it really needs. The other packages it uses (like a global errors package), may just as well be a third party package.

It’s also worth noting that it’s much easier to share the structures and functions you’re working on, if you’re living in the same package. If you had structs that depend on each-other, spliting them up into two or more different packages might lead you to encounter the diamond dependency problem. The solution for that one is relatively obvious - rewrite your code to be stand-alone, or make all your code live in the same package.

So which? Both ways work for me. If I’m being fully pedantic about it, splitting it up into more packages makes it cumbersome to add new code, as you’ll likely have to modify all of them to add a single API call. Jumping between packages may be a bit of a cognitive overhead, if it’s not very clear how you’ve laid them out. Intuition can only take you so far, so in many cases you’ll easier navigate the project, if it only has one or two packages.

You definitely don’t want to go with many small packages either (aka “Tiny Package Syndrome”).

Errors
Errors, if descriptive, may be the only tool that a developer might have to review issues in production. This is why it’s very important to either handle errors gracefully, or to pass them all the way to a layer of your application which will take the error and appropriately log it, if it cannot be handled. These are some of the characteristics that the standard library error type lacks:

the error doesn’t include a stack trace,
you can’t stack errors,
errors are pre-instantiated
It is however very simple to work around these issues by using a third party errors package - my favorite, pkg/errors. There are also other third party error packages available, but this one is written by Dave Cheney (a gopher god, take your time to read his blog), and is very much the de-facto standard when it comes to upgrading your error handling. His post “Don’t just check errors, handle them gracefully” is a must-read.

Errors with stack traces
The pkg/errors package will add context (a stack trace) to a newly created error, when invoking errors.New. When printing the error, it may look like this:

users_test.go:34: testing error Hello world
        github.com/crusttech/crust/rbac_test.TestUsers
                /go/src/github.com/crusttech/crust/rbac/users_test.go:34
        testing.tRunner
                /usr/local/go/src/testing/testing.go:777
        runtime.goexit
                /usr/local/go/src/runtime/asm_amd64.s:2361
Considering, that the complete error message was “Hello world”, printing that little bit of context with a %+v with fmt.Printf or similar - acomplishes a great thing when it comes to finding the root issue of your error(s). You know exactly where the error was created (being the operative word). Of course, when it comes to the standard library, the errors package, and the native error type - don’t provide a stack trace. With pkg/errors however, it’s easy to add one. For example:

resp, err := u.Client.Post(fmt.Sprintf(resourcesCreate, resourceID), body)
if err != nil {
        return errors.Wrap(err, "request failed")
}
In this example, the pkg/errors package adds context to a standard library error - being, your own error message ("request failed") and a stack trace to go with it. By calling errors.Wrap, the stack trace is added at this point, you so may track down the error to exactly this line.

Stacked errors
Your filesystem, database, or something else might throw relatively non-descript errors. For example, MySQL may throw this kind of forced error:

ERROR 1146 (42S02): Table 'test.no_such_table' doesn't exist
This isn’t very nice to handle. However, you can use errors.Wrap(err, "databaseError") to stack a new error on top. With this, you can better handle "databaseError" for example. The pkg/errors package will keep the actual error cause behind a causer interface:

type causer interface {
       Cause() error
}
So, errors are stacked and no context is lost. Which brings me to a side-note, that the MySQL error is a typed error, and includes more information behind it than just the error string. This means that it’s possible to handle it better too:

if driverErr, ok := err.(*mysql.MySQLError); ok {
    if driverErr.Number == mysqlerr.ER_ACCESS_DENIED_ERROR {
        // Handle the permission-denied error
    }
}
This example is gracefully borrowed from this stackoverflow thread.

Errors are pre-instantiated
What actually is an error? Very simply, an error is an struct, that implements the following interface:

type error interface {
	Error() string
}
In the case of net/http, the package exposes several error types as variables, as shown in the documentation. Adding the stack traces at this point isn’t possible (Go doesn’t allow code execution for global var declarations, only type declarations are possible). Secondly, if the stdlib would add the stack trace to the error - it wouldn’t point at the location where the error is returned, but where it’s declared as a variable (globals).

This means, that you still need to enforce calling something like return errors.WithStack(ErrNotSupported) everywhere later in your package code. It’s not a terrible pain, but unfortunately you can’t just import pkg/errors and have all your existing errors come with stack traces. It does require some by-hand invocation, if you’re not already using errors.New to instantiate your errors.

Logging
The obvious continuation of that is logging, or - more appropriately, structured logging. Here again a number of packages are available, like sirupsen/logrus or my favorite apex/log. These packages also support sending the logs to remote machines or services, where, hopefully, you have good tooling to monitor them.

When it comes to the standard log package, an option which I don’t see used very often is creating a custom logger, and passing flags like log.Lshortfile or log.LUTC to it, to again, get that little bit of context, and make your life easier - especially if you’re dealing with servers in different timezones.

const (
        Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
        Ltime                         // the time in the local time zone: 01:23:23
        Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
        Llongfile                     // full file name and line number: /a/b/c/d.go:23
        Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
        LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
        LstdFlags     = Ldate | Ltime // initial values for the standard logger
)
Even if you’re not creating a custom logger, you can modify the default with a SetFlags call (playground link):

package main

import (
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Hello, playground")
}
And the result is:

2009/11/10 23:00:00 main.go:9: Hello, playground
Don’t you like knowing where you printed something? Tracking down some code gets that little bit easier.

Interfaces
If you’re writing interfaces, name the parameters in the interface as well. Think about the following snippet:

type Mover interface {
	Move(context.Context, string, string) error
}
Do you know what the parameters here represent? Just naming the arguments in interfaces makes it completely clear.

type Mover interface {
	Move(context.Context, source string, destination string)
}
Often as well, I see some interface which produces a concrete type as the return value. A far under-used practice is to declare interfaces in a way, where the result is populated in the receiver, based on some known struct or interface arguments. This may be one of the most powerful interfaces in Go:

type Filler interface {
	Fill(r *http.Request) error
}

func (s *YourStruct) Fill(r *http.Request) error {
	// here you write your code...
}
It’s much more likely, that one or many structs would implement that interface. In comparison:

type RequestParser interface {
	Parse(r *http.Request) (*types.ServiceRequest, error)
}
… this interface returns a concrete type (and not an interface). Usually such code ends up littering your code base with interfaces that each have one implementation only and which are unusable outside your application / package structure.

Bonus trick
If you want to ensure at compile time, that one of your structs conforms and fully implements an interface (or many interfaces), you can do that like this:

var _ io.Reader = &YourStruct{}
var _ fmt.Stringer = &YourStruct{}
If you’re missing some of the functions required by those interfaces, the compiler will scream at you. The _ character hints that the variable isn’t used (a throw-away), so there are no side effects, the compiler even fully optimizes those lines away from the final binary.

Empty interfaces
This might be the more controversial opinion in contrast to the above - but I feel that using interface{} is sometimes very valid. In the example of HTTP API responses, the final step is usually json encoding, which takes an interface argument:

func (enc *Encoder) Encode(v interface{}) error
So, in effect, you can completely avoid typing out your API responses into concrete types. I don’t mean to suggest that you should do this for everything, but there are certain cases where you can fully omit typed responses in your API, or at least typed in the sense of a concrete type declaration. The first example that comes to mind is using an anonymous struct.

body := struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
}{username, roles}
First off, it’s impossible to return this kind of struct from your function, other than returning an interface{}. And obviously, the json encoder takes anything, so passing this along as-is makes perfect sense (to me). While the tendency is to declare concrete types, it is one layer of indirection which may not be needed for what you’re doing. It’s also comfy for functions that contain some logic and may return various forms of anonymous structs.

Correction: anonymous structs aren’t impossible to return, it’s just very cumbersome to do it: playground

Thanks to @Ikearens at Discord Gophers #golang channel
And the second use case, which is a big one, is database-driven API design. I’ve written about this subject before, and I should note that it’s very much possible to implement an API being completely database driven. That also means that adding and changing fields is done in the database only, without adding additional indirection in the form of an ORM. Obviously, you’ll still need to declare your types to insert the data in the database, but reading from it may be completely optional.

// getThread fetches comments by data, order by ID
func (api *API) getThread(params *CommentListThread) (comments []interface{}, err error) {
	// calculate pagination parameters
	start := params.PageNumber * params.PageSize
	length := params.PageSize
	query := fmt.Sprintf("select * from comments where news_id=? and self_id=? and visible=1 and deleted=0 order by id %s limit %d, %d", params.Order, start, length)
	err = api.db.Select(&comments, query, params.NewsID, params.SelfID)
	return
}
Similarly, your application might act as a reverse proxy or just uses a schema-less database store. The intent in those cases is just passing data along.

A big caveat (and this is where you need to type the structs out), is that modifying the interface values from Go isn’t an easy task. You’ll have to cast them to various things like maps, slices or structs in order to even access some of the returned data on the Go side. If you can’t keep your structures immutable and just pass them along from DB (or other back-end service) to the json encoder, then obviously this pattern is not for you. In fact, yes, it’s easy to argue that in the end, no such empty-interface code should exist. That being said, sometimes it’s exactly what you need - when you don’t want to know anything about the payloads.

Generated code
Generated code is code too. Commit that shit. You want to be able to issue go get -u [project], right? If you’re generating mocks for testing, if you’re generating protoc/gRPC code, any kind of code gen you might likely have, commit that shit. In case of conflicts you can always throw it away and just regenerate it.

The only possibly exception is commiting something like a public_html folder with assets that you would package with rakyll/statik. Unless somebody wants to tell me that code generated by gomock will pollute your git history with megabytes of data on every commit? It doesn’t.

Closing statements
Another good read worth noting about go best and worst practices should be Idiomatic Go. If you’re not familiar with it, give it a read - it’s a good tie-in to this article.

I’d like to leave you here with an old Jeff Atwood post - The Best Code is No Code At All and a memorable closing quote from it:

If you love writing code– really, truly love to write code– you’ll love it enough to write as little of it as possible.

But, do write those unit tests. Eventually.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.

