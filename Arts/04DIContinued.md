https://scene-si.org/2016/07/07/dependency-injection-continued/

Dependency Injection continued
The previous post on Dependency Injection in Go stirred up some heated responses on Reddit and Twitter. I’m writing this post to illustrate some good benefits of Dependency Injection in Go when applied to some benefit. I will also demonstrate an additional, declarative DI pattern of the Factory model presented in the previous installment on this subject.

Introduction
But first - if you’re a hardcore Go programmer, I want you to think about something. I want to clear the air, because we might not agree on something here. Dependency Injection is not a framework, it’s not a way of life, it’s just a way of passing stuff to where you need it.


Aaron Patterson
✔
@tenderlove
 I think if "dependency injection" was just referred to as "passing stuff in", nobody would care so much

1,042
2:08 AM - Jul 1, 2016
Twitter Ads info and privacy
645 people are talking about this
Passing stuff in is not a crime - it’s a basic pattern of programming. It’s only personal preference of what you need to do when your data model changes, and how much you know about how stuff gets passed in.

A good decision is based on knowledge. Strive to know as much as you can.

Templating
As I mentioned on the Reddit thread, one very benefitial way of using Dependency Injection is to invoke functions which don’t have a frozen interface. But Go is a strongly typed language, what are you doing with your functions that they are changing parameters so often?

Let me introduce you to ERB style templating, in the form of the package SlinSo/egon. Templates are written with tags that denote blocks for output, logic, parameters and used packages - and finally, they are compiled to actual Go code! A very simple template like this:

<%% import "strings" %%>
<%! name string %>
Hello <%= strings.TrimSpace(name) %>!
Becomes the strongly typed Go code that you are used to.

package app

import (
  "fmt"
  "html"
  "io"
  "strings"
)

func TestTemplate(w io.Writer, name string) error {
  io.WriteString(w, "\n")
  io.WriteString(w, "\nHello ")
  io.WriteString(w, html.EscapeString(fmt.Sprintf("%v", strings.TrimSpace(name))))
  io.WriteString(w, "!")
  return nil
}
And not only that, you can call functions from various packages and objects which you pass to the template.

For a very quick comparison against text/template or html/template packages, using egon allows you to:

Keep all the typed variables, get back compile warnings/errors
Write actual Go code (and not some pseudo language)
Write arithmetic/calculations/i18n code at the final layer
Invoke functions from packages like strings and other
The main pitfall about this template language is that you have to run egon to compile your templates to .go files, and then run the standard go build or go run commands. When you’re using the template packages in the go standar library - you don’t have the benefits of the go language itself, which are listed above.

Injection and templating
Okay, sure, you might not connect the dots yet with such a simple example. The templates itself have access to the full capabilities of the go language. A real world example (slightly shortened for brevity) is the following piece of template code - part of an edit form on one of our back-end products.

<%%
import (
        "app/api"
  . "app/api/structs"
        "app/common"
)
%%><%! err common.Error, repository Repository, command Command, parents []Command %>
If you read my article GoLang Tips and Tricks you will notice the import . "package" line above. I am importing all the declared public structures from “app/api/structs” package - notably Repository and Command structs.

At some point, I could just decide that I’d also like to get a RepositoryUser. Or even an inject.Injector to invoke a another template. Invoking other templates is actually a common pattern in web development - you can reuse components, widgets or even large parts of forms.

So, as we would add our additional parameter and recompile the template, the declaration of the template function would change. Invoking the template would stay as is: injector.Invoke(templates.IndexTemplate) - no changes in code are required, as long as the codegangsta/inject package can resolve the dependency. If you’re adding new data to pass to the template, you’re already modifying the model part of your MVC to put this data into the injector.

Okay, fine, but …
I know, templating on such a level is not very common with Go. You can read the Analyzing Go code with BigQuery for an insight of what people use - and even the standard library templating packages are not there. Apart from this very specific use case, I don’t have an example on hand where I’d like to use injector.Invoke. I’m just getting the most benefit out of using it here in this very specific way.

Factory method
If you remember, a Factory method is the pattern where you use a function to create your objects when you need them. While you can collect all your factory methods in a single struct, you might have a benefit of declaring them individually.

package main

import "log"
import _ "github.com/go-sql-driver/mysql"
import "github.com/jmoiron/sqlx"
import "github.com/namsral/flag"

var (
  dsn = flag.String("dsn", "api:api@tcp(db1:3306)/api?collation=utf8_general_ci&parseTime=true&loc=Local", "DSN for database connection")
)

type DatabaseFactory struct {
}
func (r *DatabaseFactory) GetDatabase(name string) (*sqlx.DB, error) {
  return sqlx.Open("mysql", *dsn)
}
func (r *DatabaseFactory) ReleaseDatabase(db *sqlx.DB) {
  db.Close()
}

type Command struct {
  Id      int
  Name    string
  Command string
}
type Repository struct {
  Id   int
  Name string
}

type ViewCommandEdit struct {
  *DatabaseFactory
  Name       string
  Command    Command
  Repository Repository
}

func main() {
  view := ViewCommandEdit{ Name: "admin-command-edit" }
  db, err := view.GetDatabase("default")
  if err != nil {
    log.Fatal("Can't connect to database");
  }
  defer view.ReleaseDatabase(db)
  // call something like:
  //   err = db.Get(&view.Command, "select * from command order by id limit 0,1;")
  //   err = db.Select(&view.Repository, "select * from repository where repository_id=?", view.Command.Id)
  // and pass the view to the template...
}
In the example above, I’m declaring a DatabaseFactory struct, which implements GetDatabase and ReleaseDatabase functions. In the ViewCommandEdit struct I’m embedding the struct, and adding fields which are important to my application. As an example I’m retrieving a database handle (I could use the name parameter to connect to a different database).

As the declaration of the object factory is explicit (embedding a struct), we’re not relying on any kind of generic object factory, or dependency injection package. You can embed multiple structs. In fact, if we could package all the needed types into ViewCommandEdit, and pass this object to the template - achieving the same result as you would with injector.Invoke.

Of course, we have a pattern of sub-templates, which still require us to work with injection, without resorting to interface{} and then casting back to a concrete type struct. It’s a trade-off.

Closing words
Dependency Injection seems to have a bad connotation in the Go community, the most common argument is that it is not idiomatic, and that it should be put on the pile of anti-patterns which were identified in other languages, and which in turn are also part of the reason why Go exists.

I couldn’t agree more with these sentiments. Nobody should want to turn Go into Java.

Whatever you do, just don’t (use) panic. Learn what Dependency Injection wants to acomplish, and then choose your own way. If anything I want you to know that it exists and show you some cases where it can be a powerful tool in “passing stuff in”.

These are just some things that scratch the surface of Golang practices. I’m writing about these and how to start writing your API in Go in my book, API foundations in Go, you should check it out. I’m also writing about and using docker through the book.

This article is part of a series, “Diving Deeper with Go”:

Golang Tips and Tricks
Advanced Go Tips and Tricks
Dependency Injection patterns in Go
While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.