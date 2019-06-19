https://scene-si.org/2017/12/21/introduction-to-reflection/

Introduction to Reflection
Reflection is the capability of a programming language to inspect data structures during runtime. When it comes to Go, this means that reflection can be used to traverse over public struct fields, retrieve information about tags of individual fields, and possibly other slightly more dangerous things.

You might know, that several packages in the Go standard library use reflection for their purposes. The example being cited most often is the implementation for encoding/json, which is commonly used to parse JSON documents into structures, and to encode the structure values back into JSON as needed.

I want to give you a slightly different example, a structure for a message payload from a chat bot that I’m currently working on:

type Message struct {
	ID         uint64    `db:"id"`
	Channel    string    `db:"channel"`
	UserName   string    `db:"user_name"`
	UserID     string    `db:"user_id"`
	UserAvatar string    `db:"user_avatar"`
	Message    string    `db:"message"`
	RawMessage string    `db:"message_raw"`
	MessageID  string    `db:"message_id"`
	Stamp      time.Time `db:"stamp"`
}
I want to log this message into the database, using a named SQL query. Basically this means that I have to generate a SQL query like this one:

insert into messages set id=:id, channel=:channel, ...
And then simply use jmoiron.sqlx to perform a db.NamedExec(query, message). Of course, depending on the number of structures that you have, writing out these queries might be prohibitive, time consuming, and error prone if you’re still refactoring your database schemata.

Wouldn’t it be cool if we could generate the query from a passed structure?

We can actually import the reflect package and do just that. There are a number of ORM packages which resort to similar methods as the one that I’m about to walk you through.

Resolving a struct to a reflect.Value
In order to traverse over structure fields, we will need to produce a reflect.Value instance first. This instance holds many functions that we can then use for traversal. Producing a reflect.Value is relatively straightforward.

message_value := reflect.ValueOf(message)
In order to iterate over the struct fields, we need to call message_value.NumField() to get the count of the fields in a struct. In case we pass a pointer to reflect.ValueOf, and then attempt to call NumField(), we will end up with a panic;

panic: reflect: call of reflect.Value.NumField on ptr Value
In order to resolve this, we can check the message_value.Kind() for a pointer result, and then get the actual value, resolving the pointer:

if message_value.Kind() == reflect.Ptr {
	message_value = message_value.Elem()
}
Calling message_value.NumField() will now report the correct number of fields within this struct. We can use this value to iterate over each one and get the needed names and values that we require.

Reading field details
There is a significant number of things we can get from struct fields, but what we’re interested into is reading the tag information from field declarations. Where it gets a bit fuzzy is that reflect.Value is dealing with actual values which are set in the field, but if you’re after things like the name of the field (UserName for example), or the associated tags, you will need a reflect.Type.

Let’s say we wanted to create a list of all the fields in the struct, with the name of the field, the specific tag for “db”, and the actual value set in the field, the code will look like this:

message_fields := make([]struct {
	Name  string
	Tag   string
	Value interface{}
}, message_value.NumField())

for i := 0; i < len(message_fields); i++ {
	fieldValue := message_value.Field(i)
	fieldType := message_value.Type().Field(i)
	message_fields[i].Name = fieldType.Name
	message_fields[i].Value = fieldValue.Interface()
	message_fields[i].Tag = fieldType.Tag.Get("db")
}
As you see, the literal value of each field is retrieved with a call to reflect.Value.Interface(), while the type information like the name of the field and the specific tag is retrieved from a reflect.Type. You can run the full example on go playground.

Putting it all together
As we need less than what we already made, the trick to producing the sql query which we need is just throwing away some code and just using the tag values, and parsing them slightly:

func insert(table string, data interface{}) string {
	message_value := reflect.ValueOf(data)
	if message_value.Kind() == reflect.Ptr {
		message_value = message_value.Elem()
	}

	message_fields := make([]string, message_value.NumField())

	for i := 0; i < len(message_fields); i++ {
		fieldType := message_value.Type().Field(i)
		message_fields[i] = fieldType.Tag.Get("db")
	}

	sql := "insert into " + table + " set"
	for _, tagFull := range message_fields {
		if tagFull != "" && tagFull != "-" {
			tag := strings.Split(tagFull, ",")
			sql = sql + " " + tag[0] + "=:" + tag[0] + ","
		}
	}
	return sql[:len(sql)-1]
}
And here’s the final go playground code.

There are some specific notes about this:

For our example, we are not traversing further into the struct, so it doesn’t matter if a struct field is a pointer or not in our case. The reflect.Type information works regardless of the actual value.
If you would need to traverse fields further, be sure to resolve pointer values, just like we initially do with message_value, check for a pointer and get that Elem()
There are several great packages that provide some functionality with reflection. If you want to see more real world examples, take a look at codegangsta/inject and fatih/structs.
Obligatory reading: The Laws Of Reflection, by Rob Pike.
And a final word of warning: If you’re coming from languages like PHP or JavaScript: the lack of type safety there is ridiculous. You might want to use reflection to get around some of these things, but you should exercise judgement in doing so. You chose to use Go for a reason, and if you’re going to be resorting to reflection to do some things that you know from these languages, you’ll inevitably cancel out your own reasons why you’re moving to Go in the first place.

In most cases during my history with Go, resorting to reflection is somewhat of a crutch. Even if this article example is a bit contrived. You’re very much able to type out those few struct fields into a string slice or the final query itself, avoiding the use of reflection alltogether. In fact, if whatever API you’re writing ever reaches any kind of real traffic, the reflection will be the first low hanging fruit to get rid of.

Case and point: Even for JSON, there’s the package json-iterator/go, which is a drop in replacement for the standard library encoding/json package. In this implementation the reliance on reflection was greatly reduced, resulting in significant speed gains. While it’s impossible for the reflection to be refactored away in the light of such uncertainty, you might gravitate towards gRPC and Protobuf if that is becoming an issue for you. In that case, full type safety is achieved with code generation.

Exercise wisdom.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.