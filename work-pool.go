package main

import (
    "sync"
    "time"
    "fmt"
)

type Server int

func client(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    req := fmt.Sprintf("req-%d", id)
    resp := server.Serve(req)
    // fmt.Println("client ", id, req, " --- ", resp)
    _ = resp
}

var server Server
func main() {
    // start server
    const WORKER_LIMIT = DB_LIMIT
    var db = &DB{}
    for i := 0; i < WORKER_LIMIT; i++ {
        woker := &Worker{
            db: db,
        }
        go woker.doWork()
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
    JobQueue <- job
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
var JobQueue = make(chan *Job)
type Worker struct {
    db *DB
}
func (w *Worker) doWork() {
    for {
        job := <- JobQueue
        resp := w.db.Query(job.request)
        job.result = resp
        // close(job.done)
        job.done <- struct{}{}
    }
}

// DB: limited concurrency
const DB_LIMIT = 50
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