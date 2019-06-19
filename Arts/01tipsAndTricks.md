https://scene-si.org/2016/06/01/golang-tips-and-tricks/

Golang tips and tricks
It’s been a while since I’ve started programming with Go, and I managed to pick up certain tricks along the way, which I’ll present you with. Use them wisely and you’ll be on your way to creating that beautiful product or service you’ve been itching to start.

Importing packages
When you start to use a language, you notice how to include code. While Go doesn’t include files, it does include packages. C has #include, Java has import with wildcards, PHP has include, and Go has import with package name. But there’s additional options to this import which usually people skip over.

import . “package”
The . (dot) in front of the package name imports the package into your namespace. If you’d import the “fmt” package this way, you could just call Println and other functions without using the fmt. prefix.

import “./package”
This is a relative import. It works for me in some cases, but it’s better to be avoided. It doesn’t work with vendoring (sad story) and it was moved into the Unplanned milestone on GitHub.

You can replace the “dot” with an absolute path of your projects. If your main.go lives inside the app folder, you can, legally, use import "app/package" to work around this.

Circular dependencies
It’s the belief of Golang developers that circular dependencies are design flaws. Well, they are, but I’m guessing that if you have package alice that depends on package bob for some things, and then you have package bob that depends on package alice for some things,… it shouldn’t be necessary to join packages alice and bob into one package to resolve the circular dependency.

For the most part, it’s not. Keep a very basic programming principle in mind:

Write programs that do one thing and do it well.

There are two ways to solve this issue, by keeping most of your things in tact. I don’t exactly know if this is a pattern or not, but it served me well. When you have packages alice and bob that depend on each-other, they really, usually, just need access to the same structs, and modify them as needed.

Modifier pattern
A way to solve this unfortunately, is not with code, in the sense that there’s not a one-liner that would solve this. For example, I’d have two structs, with different selection of fields.

package alice
type Alice struct {
    Id int
    Name string
    IsAdmin bool
}
and

package bob
type Bob struct {
    Name string
}
In the above example, I’ll use the Bob struct when getting data from a POST request to an API endpoint. I’ll then copy this data into the Alice struct, and fill out additional properties.

A simple way to still have this kind of workflow is to separate all the structs in an individual package. This way package alice and package bob will create and modify structs from package charlie. No conflict.

Copying values between semi-compatible structs (not casting)
In the above example, you’d likely like to copy the data from a Bob struct into Alice. As the structures don’t have the same underlying fields and types, you can’t cast one into another. But you can implement this copying yourself.

import "github.com/fatih/structs"

func Fill(dest interface{}, src interface{}) {
    mSrc := structs.Map(src)
    mDest := structs.Map(dest)
    for key, val := range mSrc {
        if _, ok := mDest[key]; ok {
            structs.New(dest).Field(key).Set(val)
        }
    }
}
You’d use it very simply like this:

a := Alice{}
b := Bob{ Name: "Bobby" }
Fill(&a, b)
And you’d get a.Name filled out.

Inheritance for additional type safety
To make the above example safer, you can embed one struct into another (composition).

package alice
import "app/bob"
type Alice struct {
    Id int
    Bob // embed Bob struct
    IsAdmin bool
}
This method is safer as you’re ensuring that there can be no type conflict between properties shared between Bob and Alice object types.

Merging slices
If you’re working with slices, which you mostly are, you’ll inevitably find yourself in a situation where you’ll have to merge two slices into one. The novice way of doing this is something like the following:

alice := []string{"foo", "bar"}
bob := []string{"verdana", "tahoma", "arial"}
for _, val := range bob {
    alice = append(alice, val)
}
A slightly more expert way of doing this is:

alice = append(alice, bob...)
The ellipsis or variadic argument expands bob in this case to all it’s members, effectively achieving the same result, without that loop. You can read more about it in the official documentation: Appending to and copying slices

You could use it to wrap logic around log.Printf, which for example, trims and appends a newline at the end of the format string.

import "log"
import "strings"

func Log(format string, v ...interface{}) {
    format = strings.TrimSuffix(format, "\n") + "\n"
    log.Printf(format, v...)
}
While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.


