package main

import (
    "time"
    "math/rand"
    "fmt"
)

func main() {
    go bake()
    go packSmoker()
    // go packSmoker()
    // go packSmoker()
    time.Sleep(time.Hour)
}

type Bread string
var buff = make(chan Bread, 3)
func bake() {
    fmt.Println("baker start!")
    defer fmt.Println("baker stop!")
    for i := 0; i < 30; i++ {
        bakeTime := 500 + (int)(rand.Intn(200))
        time.Sleep(time.Duration(bakeTime) * time.Millisecond)
        b := Bread(fmt.Sprintf("bread-%d", i))
        fmt.Println("bake: ", b)

TRY:    for {
            select {
            case buff <- b:
                break TRY
            default:
                coolTime := 200
                time.Sleep(time.Duration(coolTime) * time.Millisecond)
                fmt.Println("bake cool: ")
            }
        }
    }
    close(buff)
}

func packSmoker() {
    fmt.Println("packSmoker enter!")
    defer fmt.Println("packSmoker leave!")
    for {
        select {
        case b, ok := <- buff:
            if ok {
                packTime := 500 + (int)(rand.Intn(600))
                time.Sleep(time.Duration(packTime) * time.Millisecond)
                fmt.Println("             pack: ", b)
            } else {
                return
            }
        default:
            fmt.Println("                            pack-smook: ")
            smookTime := 500
            time.Sleep(time.Duration(smookTime) * time.Millisecond)
        }
    }
}