https://scene-si.org/2016/06/13/advanced-go-tips-and-tricks/
Advanced Go Tips And Tricks
We covered a few useful tricks in the previous installment of this blog series, and now I’m going to show you a few more which might make your life even simpler.

Nested structures when parsing JSON data
Contrary to what you might believe, there’s no need to declare structures individually, you may also declare and use them in nested form. When dealing with more complex JSON documents, this has a number of advantages.

Let’s consider this simple JSON file:

{
  "id": 1,
  "name": "Tit Petric",
  "address": {
    "street": "Viska cesta 49c",
    "zip": "1000",
    "city": "Ljubljana",
    "country": "Slovenia"
  }
}
The type for this JSON document can be declared as:

type Person struct {
  Id      int    `json:"id"`
  Name    string `json:"name"`
  Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
  } `json:"address"`
}
By using this form of declaration, you clearly and explicitly define the hierarchy of the JSON document you are parsing. Keep in mind that the declaration doesn’t need to define all fields, but just the ones you will be using. In skipping some fields, you are making the JSON parsing more resilient to changes in the JSON document.

Anonymous structs
The Address property in the previous section is called an anonymous struct. You may declare the complete structure as anonymous by explicitly assigning it to a variable:

person := struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Address struct {
    City string `json:"city"`
    Country string `json:"country"`
  } `json:"address"`
}{}
And you can use this variable the same you would a Person{} from the previous example.

err = json.Unmarshal(contents, &person)
Using embedding to your advantage
As mentioned in the previous installment, you can embed two structs together by ommiting the name of the embedded struct. You can use this to your advantage with, for example, sync.Mutex, sync.RWMutex or sync.WaitGroup. You can actually embed many structs, so your structure may perform the functions of them all.

An example from a project I’m working on is using two embedded structs:

type RunQueue struct {
  sync.RWMutex
  sync.WaitGroup
  // ...
  flagIsDone  bool
}
func (r *RunQueue) Close() {
  r.Lock()
  defer r.Unlock()
  r.flagIsDone = true
}
func (r *RunQueue) IsDone() bool {
  r.RLock()
  defer r.RUnlock()
  return r.flagIsDone
}
Leveraging sync.RWMutex
The declaration of the RunQueue struct above leverages the sync.RWMutex to provide synchronous access to the object from many goroutines. A goroutine may use Close to finish the execution of the goroutine queue. Each worker in the queue would call IsDone to check if the queue is still active.

Leveraging sync.WaitGroup
The RunQueue struct leverages a sync.WaitGroup to provide queue clean up and statistics like time spent. While I can’t provide all the code, the basic use is like this:

func (r *RunQueue) Runner() {
  fmt.Printf("Starting %d runners\n", runtime.NumCPU())
  for idx := 1; idx <= runtime.NumCPU(); idx++ {
    go r.runOnce(idx)
  }
}
func NewRunQueue(jobs []Command) RunQueue {
  q := RunQueue{}
  for idx, job := range jobs {
    if job.SelfId == 0 {
      q.Dispatch(&jobs[idx])
    }
  }
  q.Add(len(q.jobs)) // sync.WaitGroup
  return q
}
runnerQueue := NewRunQueue(commands)
go runnerQueue.Finisher()
go runnerQueue.Runner()
runnerQueue.Wait() // sync.WaitGroup
The main idea of the program I’m building is that it starts runtime.NumCPU() runners, which handle execution of a fixed length of commands. The WaitGroup comes into play very simply:

NewRunQueue calls *wg.Add(count of jobs)
Individual jobs are processed with RunQueue.runOnce, they call *wg.Done()
RunnerQueue.Wait() (*wg.Wait()) will wait until all jobs are processed
Limiting goroutine parallelization
At one point I’ve struggled to create a queue manager, which would parallelize workloads to a fixed limit of parallel goroutines. My idea was to register a slot manager, which would provide a pool of running tasks. If no pool slot is available, I’d sleep for a few seconds before again trying to get a slot. It was frustrating.

Just look at the loop from the Runner function above:

for idx := 1; idx <= runtime.NumCPU(); idx++ {
  go r.runOnce(idx)
}
This is an elegant way to limit parallelization to N routines. There is no need to bother youself with some kind of routine allocation pool structs. The runOnce function should only do a few things:

Listen for new jobs in an infinite loop, read jobs from a channel
Perform the job without new goroutines
The reason to read the jobs from a channel is that the read from a channel is a blocking operation. The function will just wait there until a new job appears on the channel it’s reading from.

func (r *RunQueue) runOnce(idx int) {
  for {
    queueJob, ok := <-r.runQueue
    if !ok {
      return
    }
    // run tasks
[...]
The job needs to be executed without a goroutine, or with nested *WaitGroup.Wait() call. The reason for this should be obvious - as soon as you start a new goroutine, it gets executed in parallel and the runOnce function reads the next job from the queue. This means that the limitation of how many tasks are running in parallel would not be enforced.

While I have you here...
It would be great if you buy one of my books:

API Foundations in Go
12 Factor Apps with Docker and Go
I promise you'll learn a lot more if you buy one. Buying a copy supports me writing more about similar topics. Say thank you and buy my books.
Feel free to send me an email if you want to book my time for consultancy/freelance services. I'm great at APIs, Go, Docker, VueJS and scaling services, among many other things.

