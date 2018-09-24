package main

import (
    "sync"
    "fmt"
)

var wChan = make(chan int)
var rChan = make(chan int)
func write(v int) {
    wChan <- v
}
func read() int {
    return <-rChan
}
func worker() {
    var count int
    for {
        select {
        case v := <- wChan:
            count += v
        case rChan <- count:
        }
    }
}

func writer(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        write(1)
    }
}

func main() {
    go worker()
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go writer(&wg)
    }
    wg.Wait()

    fmt.Println(read())
}