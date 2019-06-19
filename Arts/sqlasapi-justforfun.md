https://scene-si.org/2018/02/07/sql-as-an-api/

SQL as an API
If you haven’t been living under a rock, you’ll know that recently there is an uptake in popularity of “Functions as a service”. In the open source community, the OpenFaaS project from Alex Ellis has received good traction, and recently Amazon Lambda announced support for Go. These systems allow you to scale with demand, and execute your CLI programs via API.

Motivation behind Lambda/FaaS
Let’s call this for what it is - the whole “serverless” movement is marketing for cloud stacks like AWS, that allow you to hand over any server management to them, for, hopefully, a fraction of your income. In concrete terms, this means that AWS and similar solutions take your application, and take steps to run and scale it based on demand in their data centers.

But you probably know all this already.

But did you know, this functionality already existed as CGI, since, according to Wikipedia, 1993, and formally in the form of a RFC since 1997? Everything old is new again. The intent of CGI (Comon Gateway Interface) is:

In computing, Common Gateway Interface (CGI) offers a standard protocol for web servers to execute programs that execute like Console applications (also called Command-line interface programs) running on a server that generates web pages dynamically.

Source: Wikipedia

So in terms of Go, the simplest FaaS/CGI server can be written in about 10 lines of code. The Go stdlib already includes net/http/cgi which does all the heavy lifting. Including PHP CGI can be as short as these few lines of code:

func phpHandler() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        handler := new(cgi.Handler)
        handler.Dir = "/var/www/"
        handler.Path = handler.Dir + "api"
        args := []string{r.RequestURI}
        handler.Args = append(handler.Args, args...)
        fmt.Printf("%v", handler.Args)
        handler.ServeHTTP(w, r)
    }
}
Using it is straightforward: http.HandleFunc("/api/", phpHandler()). Of course, why exactly you would want to use it is beyond me. Since the dawn of time, CGI has been plagued with performance issues, the biggest of which is the pressure it exerts on the system. For each request, an os.Exec call is made, which isn’t really the most system friendly call. In fact, you’d likely want to keep this close to 0, if you’re serving any kind of real traffic.

This is why CGI evolved into FastCGI.

FastCGI is a variation on the earlier Common Gateway Interface (CGI); FastCGI’s main aim is to reduce the overhead associated with interfacing the web server and CGI programs, allowing a server to handle more web page requests at once.

Source: Wikipedia

While I won’t be implementing a FastCGI server (there’s also net/http/fcgi in the stdlib), I did want to illustrate the performance pitfalls in going towards such an implementation. Of course, when you’re running your programs on AWS, you tend not to care about this, since they have the server capacity to scale out your request volume.

Working around CGI
If there’s something that I’ve learned over the years, it’s the fact that most services are data driven. This means, there’s some form of database available which is holding the mission critical data. According to one Quora answer, Twitter uses at least 8 different forms of a databases, from MySQL, Cassandra and Redis, to other more esoteric ones.

In fact, most of my work revolves around writing APIs that would mostly read stuff from the database and return it as JSON via a REST call. While these queries often can’t be expressed in only one SQL statement, there’s a significant number of times when they absolutely can. How about instead of writing a CGI program that would do something, we would write an SQL script that would do something without an overhead of os.Exec?

Challenge accepted.

SQL as an API
While I’m not taking this into the absolute behemoth that SQL as an API can be, I do want to create something at least remotely usable. I want to create API calls by creating an .sql file on disk, and I want this API call to take any parameters defined in the HTTP request query. This means that we can filter SQL query results, by the parameters we pass to the API. I have chosen MySQL and sqlx for this task.

I’ve recently created several chat bots for Twitch, Slack, Youtube, Discord and it seems I’ll be working on a Telegram one soon. Specifically their intent is similar - connect to a number of channels, log messages, tally up some statistics and respond to commands or questions. There was a need to expose this data via APIs for a web front-end written in Vue.js; while not all the API calls could be implemented with SQL, a significant number of them could. For example:

listing all channels
listing channels by channel ID
These two calls are relatively similar and easy to implement. Specifically I created just two files, that provide this information:

api/channels.sql (responds to /api/channels)

select * from channels
api/channelByID.sql (responds to /api/channelByID?id=...)

select * from channels where id=:id
As you see, there’s really not much to creating a new API call which you can express as an SQL query. I’ve tried to design the system so that it’s possible that you’ll create a api/name.sql file, and this file is immediately reachable as /api/name. Any HTTP query parameters are converted into a map[string]interface{} and passed to the SQL query as bound parameters. The SQL driver takes care of escaping those.

I also took care of formatting the errors. If you can’t reach your database, or the .sql file doesn’t exist for a certain api call, there’s an error returned which looks like this:

{
    "error": {
        "message": "open api/messages.sql: no such file or directory"
    }
}
Using URL Query parameters for the SQL query
Getting the query parameters in Go is done by issuing a Query() call on the *url.URL structure contained in the request object. The call returns an url.Values object, which is a type alias of map[string][]string.

We need to convert this object to pass it into an sqlx statement. We need to create a a map[string]interface{}. This is because the sqlx function we need to call, accepts this form of parameters for the query (sqlx.NamedStmt.Queryx). Let’s convert them and issue the query:

params := make(map[string]interface{})
urlQuery := r.URL.Query()
for name, param := range urlQuery {
        params[name] = param[0]
}

stmt, err := db.PrepareNamed(string(query))
if err != nil {
        return err
}
rows, err := stmt.Queryx(params)
if err != nil {
        return err
}
What we’re left with is the rows variable, which we can iterate to receive each individual row. We’ll need to add them to a slice, which we’ll encode to JSON in the final steps of the API call.

for rows.Next() {
        row := make(map[string]interface{})
        err = rows.MapScan(row)
        if err != nil {
                return err
        }
Here’s where things become more interesting. The values contained in each row should be converted to something that the JSON encoder understands.

Because the underlying type is []uint8, we would need to convert this to strings first. If we don’t convert to string, a JSON representation of this structure will be automatically base64 encoded. Since query responses can be represented with map[string]string, and uint8 is a type alias to byte, we opt into this conversion.

rowStrings := make(map[string]string)
for name, val := range row {
        switch tval := val.(type) {
        case []uint8:
                ba := make([]byte, len(tval))
                for i, v := range tval {
                        ba[i] = byte(v)
                }
                rowStrings[name] = string(ba)
        default:
                return fmt.Errorf("Unknown column type %s %#v", name, spew.Sdump(val))
        }
}
Here we have a rowStrings object representing each returned SQL row, which can now be encoded into JSON without issues. All we need to do is append these into a result, encode it and return the encoded value. The full (and relatively short) code is available on titpetric/sql-as-a-service.

Caveat emptor
While this approach has the distinct benefit of an API layer which is database driven, there are a number of scenarios that need to be considered in order to make it suitable for a wider array of usage. For example:

Ordering result sets

We can’t really order result sets in this way. This is because we can’t bind a query value into an order by parameter, as SQL doesn’t allow this. Sanitizing this also becomes a problem, and you can’t even use a function as a replacement Completely not possible: ORDER BY :column IF(:order = 'asc', 'asc', 'desc').

Parameters

In order to create some sort of pagination rules, MySQL provides a LIMIT $offset, $length clause. While you can supply these as query parameters, we again have the issue of not being able to bind them in this place, or have a good way to clamp their values. However, the error returned is something along the lines of “Undeclared variable: …“.

Multiple SQL queries

Generally, we could use multiple SQL queries to produce a single result. That would need to be enabled in the driver however, and is generally disabled and frowned upon because it’s one of the main causes for SQL injection attacks. In a fair and just world, however, something like this would be possible:

set @start=:start;
set @length=:length;
select * from channels order by id desc limit @start, @length;
Alas, it’s not; the above doesn’t work even remotely. If you try it with a MySQL client it’s going to error out. While the variables are created, they can’t be used in either the order by or limit clauses.

So, say goodbye to pagination?

A work-around which I’ve seen in the wild, most notably Twitch and Tumblr APIs, allows you to pass a since or previousID value. This will result in a query like this:

select * from messages where id < :previousID order by id desc limit 0, 20
This allows you to iterate through a table based on a predefined page size. It does require either K-sortable primary keys (sonyflake / snowflake IDs), or a secondary sequential column which you can use for sorting.

Functions

SQL servers are no dumb beasts, in fact they are super powerful. One of these sources of power is the ability to create SQL functions, or, procedures. With it you can create an implementation of complex business logic. While typically in the domain of a DBA, or at least a database oriented programmer, it’s easy enough to read up a bit on it and try out some things yourself.

It’s a good way to work around the multiple-SQL-query limitation of the clients, but it requires you to keep all your database logic on the database side. This is outside of the comfort zone of most programmers that I know. Expect basically to handle bunches of cases with IF’s.

Conclusion
If you really have simple queries that produce the result set that you like, you can achieve a lot with the SQL as a service approach. Besides being friendly to your systems resources, you can also provide a way to add functionality from other programmers that don’t necessarily know Go, but can write up a decent SQL query.

With the increasing logic demands for what an API endpoint should return however, there is a need for a scripting language which would approach or overcome the limitations of PL/SQL and it’s various flavors implemented in various RDBMS.

Or you can always hook up a Javascript VM like dop251/goja into your app, and take that front-end programmer on your team for a ride he/she might never forget. There are also LUA vms implemented in pure go available, if you prefer something smaller than the whole ES5.1+ “runtime”.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.