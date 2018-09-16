package main

import (
    "math/rand"
    "fmt"
    "time"
)

func randomGen(count int) chan int {
    ch := make(chan int)
    go func() {
        defer fmt.Println("randomGen exit...")
        // defer close(ch)
        for ; count > 0 ; count-- {
            v := int(rand.Intn(100))
            ch <- v
        }
    }()
    return ch
}
func main() {
    ch1 := randomGen(10)
    ch2 := randomGen(20)
    for i := 0; i < 50; i++ {
        time.Sleep(time.Millisecond * 300)
        select {
        case <- ch1:
            fmt.Println(i, " recv from ch1")
        case <- ch2:
            fmt.Println(i, " recv from ch2")
        default:
            fmt.Println("       ", i, " heihei...")
        }
    }
}