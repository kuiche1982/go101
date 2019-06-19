https://scene-si.org/2017/09/27/things-to-know-about-http-in-go/


Things to know about HTTP in Go
Go has a very strong standard library, and one of the most used parts of it is the net/http package, which provides structures that make development of HTTP servers and clients very straightforward. There are a few edge cases, where a deeper understanding of the http and related packages is very welcome.

Most RESTful HTTP API requests don’t have to deal with many things. People generally need to read query variables from GET requests, or Form variables from POST requests, and in some cases read the POST body in order to save a file or get a JSON payload. People resort to the available net/http structures that expose the relevant fields and functions for this purpose.

net/http/httputil
Unfortunately, not a lot of people know about httputil, and it was a long while before I saw it in the wild. The package provides an useful utility function, DumpRequest. No doubt, some time in your past you might have written code that would dump the request method (POST/GET), the request r.URL.Path, or even something like this:

for name, headers := range r.Header {
	for _, h := range headers {
		fmt.Printf("%v: %v\n", name, h)
	}
}
Don’t worry. We’ve all been there.

If you want to debug your HTTP requests, all you really need to do is import the net/http/httputil package, and invoke DumpRequest with the parameter *http.Request and a boolean to specify if you want to dump the request body as well. The function returns a []byte, error, and you could use it like this:

dump := func(r *http.Request) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}
The function call will dump your request method, URI with query parameters, headers and request body if you have one. It should look like this:

POST /api/v3/projects/1234567/notices?key=FIXME HTTP/1.1
Host: 10.1.1.2:3000
Accept-Encoding: gzip
Content-Length: 617
Content-Type: application/json
User-Agent: Go-http-client/1.1

POST_REQUEST_BODY_HERE
"Things to know about HTTP in Go: Use httputil.DumpRequest for debugging #golang" via @TitPetric

Click to Tweet
Request body is an io.ReadCloser
With ErrorHub we’re dealing with some compressed HTTP payloads. This means that the payload needs some decoding. With functional javascript or PHP you would chain several function calls that would do for example, gzip and base64 decoding, before you would decode the JSON payload.

Go is much more elegant in this respect. You can of course follow bad practice from other languages and call something like ioutil.ReadAll on the request body to get the contents, which you then in turn pass through several functions to get the decoded result.

Or, you can realize that http.Request.Body is an io.ReadCloser. This means that you can decode your payload as a stream. For example, the Sentry client raven-go sends the JSON payload first compressed and then base64 encoded.

The decoder for the payload looks something like this:

base64decoder := base64.NewDecoder(base64.StdEncoding, r.Body)
gz, err := zlib.NewReader(base64decoder)
if err != nil {
        return err
}
defer gz.Close()

decoder := json.NewDecoder(gz)
var t SentryV6Notice
err = decoder.Decode(&t)
if err != nil {
	return err
}
r.Body.Close()
// ...
The interface io.ReadCloser also satisfies io.Reader. With this we can first create a base64 decoder, which we then pass into the zlib.NewReader to create a zlib decoder/reader, and finally pass that one into json.NewDecoder, which we can use to decode the payload.

Dealing with io.ReadCloser and io.Reader is straightforward. The above implementation has the benefit of being quite efficient in speed/memory use, in comparison with working with []byte or string variables. Those will inevitably use a more memory and perform worse than above.

"Things to know about HTTP in Go: request.Body is an io.ReadCloser #golang" via @TitPetric

Click to Tweet
Handler and HandlerFunc
There are two types that declare a signature for HTTP handlers. The most common type, http.HandlerFunc is a type alias for func(http.ResponseWriter,*http.Request). You can pass this into http.HandleFunc(). The less common type, http.Handler is an interface, which should implement a ServeHTTP(http.ResponseWriter,*http.Request).

As an example of the second form, there is a http.FileServer() function provided. This is the example use as-is from the documentation, which will serve the files in your /tmp directory:

http.Handle("/", http.FileServer(http.Dir("/tmp")))
Since you may not want or need to provide a signature for http.Handle but a http.HandleFunc, you can “wrap” the code like this:

func FileServer(path string) http.HandlerFunc {
	server := http.FileServer(http.Dir(path))
	return func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}
}

http.HandleFunc("/", FileServer("/tmp"))
The example itself isn’t great - the unwrapped function provided doesn’t add on any functionality. But let’s consider the following use case. VueJS apps can use browser history API to simulate pageloads. This means that when you click on /about, the javascript takes care of rendering the web page. Other frameworks like React and Angular work in the same way. To reliably support an user pressing refresh in the browser, any non-existent page on the server should return the contents of /index.html.

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(assetPath string, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(assetPath, "index.html")
		requestedPage := path.Join(assetPath, r.URL.Path)
		_, err := os.Stat(requestedPage)
		if err != nil {
			// serve index if page doesn't exist
			http.ServeFile(w, r, indexPage)
			return
		}
		serve.ServeHTTP(w, r)
	}
}
This function uses os.Stat() to figure out if a file doesn’t exist, and in that case serve index.html. Everything is wrapped into a http.HandlerFunc. This is an example straight from the vue-serve project.

"Go #protip - you can wrap other http handlers to extend functionality #golang" via @TitPetric

Click to Tweet
The ResponseWriter
The http.ResponseWriter is an io.Writer. This means that you can encode your JSON directly to the underlying writer. There are again many examples encoding JSON first, and then calling Write() on the ResponseWriter. If you want to chain some gzip compression onto this, going with writers should be your best bet.

type JSON struct {
        value interface{}
}

func (self *JSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        encoder := json.NewEncoder(w)
        err := encoder.Encode(self.value)
        if err != nil {
                http.Error(w, err.Error(), 503)
        }
}
Similarly, we could use io.Copy to read out files. I mean, would you really like to read files which might be several GB in size, before you write them out to the ResponseWriter? Thought so.

Bonus tip: The handler httputil.ReverseProxy implements a copyBuffer function which is more suited to proxying requests. It has a similar signature to io.Copy.

"Know about httputil.ReverseProxy? Leverage ResponseWriter as io.Writer :) #golang" via @TitPetric

Click to Tweet
Testing HTTP handlers
If you want to automate some tests, as you should, people usually resort to external tooling and validators to test their API responses. Go has everything you need to write tests without those. Using net/http/httptest you can create a server, which will return data from your http.Handler without actually using sockets. Let’s try to see how to do that:

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type JSON struct {
	value interface{}
}

func (self *JSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(self.value)
	if err != nil {
		http.Error(w, err.Error(), 503)
	}
}

func main() {
	handler := &JSON{"hello world"}
	server := httptest.NewServer(handler)
	defer server.Close()

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := http.Get(server.URL)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	log.Printf("received: %d %s\n", resp.StatusCode, body)
}
The complete example is on go playground. A testing http server is created with httptest.NewServer which takes http.Handler as the parameter. You can request things from the server by referencing server.URL in the call to http.Get. You can test pretty much anything you write without extra tooling.

"Start testing your HTTP handlers without external tooling, use net/http/httputil package #golang" via @TitPetric

Click to Tweet
Conclusion
HTTP, as much as it can be very simple with POST/GET requests, has a much wider scope. Knowing some of the best practices and options which you have when developing servers should help you along towards better implementations. While it’s certainly possible to create quick microservices, you should take some care in creating better microservices as well.

If you’re into improving your apps, Go or something else, check out and subscribe to ErrorHub. We’re writing an error catching service that will let you know if your app has errors that you need to handle. If you want to learn more about Go and how to do stuff with it, check out the book list below.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.

