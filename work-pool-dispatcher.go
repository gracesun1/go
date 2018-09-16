package main

import (
    "sync"
    "time"
    "fmt"
)

type Server struct {
    dispatcher *Dispatcher
}

func client(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    req := fmt.Sprintf("req-%d", id)
    resp := server.Serve(req)
    fmt.Println("client ", id, req, " --- ", resp)
    // _ = resp
}

var server Server
func main() {
    // start server
    const WORKER_LIMIT = DB_LIMIT
    const JOB_LIMIT = 0
    dispatcher := &Dispatcher{
        FreeWorker: make(chan *Worker, WORKER_LIMIT),
        JobQueue: make(chan *Job, JOB_LIMIT),
    }
    go dispatcher.Dispatch()
    server.dispatcher = dispatcher
    var db = &DB{}
    for i := 0; i < WORKER_LIMIT; i++ {
        woker := &Worker{
            db: db,
            queue: make(chan *Job),
        }
        go woker.doWork(dispatcher)
    }
    // start client
    var wg sync.WaitGroup
    for i := 0; i < 100000; i++ {
        wg.Add(1)
        go client(i, &wg)
    }
    wg.Wait()
}

func (s *Server) Serve(request string) string {
    // db.Query()
    // create a query task(Job)
    job := &Job{
        request: request,
        done: make(chan struct{}),
    }
    // sent task to queue
    s.dispatcher.AddJob(job)
    // wait for task finish
    <- job.done
    // read task result, return it to user
    return job.result
}

// work-pool
type Job struct {
    request string
    done chan struct{}
    result string
}
type Worker struct {
    db *DB
    queue chan *Job
}
func (w *Worker) doWork(d *Dispatcher) {
    for {
        d.FreeWorker <- w
        job := <- w.queue
        resp := w.db.Query(job.request)
        job.result = resp
        // close(job.done)
        job.done <- struct{}{}
    }
}
type Dispatcher struct {
    FreeWorker chan *Worker
    JobQueue chan *Job
}
func (d *Dispatcher) AddJob(job *Job) {
    d.JobQueue <- job
}
func (d *Dispatcher) Dispatch() {
    for {
        // find a free worker
        worker := <- d.FreeWorker
        // fetch a job
        job := <- d.JobQueue
        // send the job to worker
        worker.queue <- job
        fmt.Println("                       ", len(d.FreeWorker))
    }
}

// DB: limited concurrency
const DB_LIMIT = 20
type DB struct {
    sync.Mutex
    count int
}
func (db *DB) Query(request string) string{
    db.Lock()
    db.count++
    if db.count > DB_LIMIT {
        panic(db.count)
    }
    db.Unlock()

    // do some...
    time.Sleep(time.Millisecond)

    db.Lock()
    db.count--
    db.Unlock()
    return "RESP: " + request
}