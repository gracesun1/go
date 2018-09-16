package main

import (
    "fmt"
    "sync"
    "time"
    "errors"
)

const TOKEN = "3rd-token"
func ApiBy3rd() (token string) {
    fmt.Println("ApiBy3rd called, 1")
    time.Sleep(time.Second * 20)
    fmt.Println("ApiBy3rd called, 2")
    return TOKEN
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 200000; i++ {
        wg.Add(1)
        go clientCall(&wg)
    }
    wg.Wait()
}

func clientCall(wg *sync.WaitGroup) {
    defer wg.Done()

    var token string
    var err error = errors.New("init")
    var totalSpent time.Duration
    getToken := func() error {
        tmStart := time.Now()
        token, err = server()
        tmSpent := time.Since(tmStart)
        if tmSpent > time.Millisecond * 1000 {
            panic(tmSpent)
        }
        totalSpent += tmSpent
        return err
    }

    var retry int
    for retry = 0; retry < 500 && err != nil; retry++ {
        err = getToken()
        time.Sleep(time.Millisecond * 50)
    }
    if err != nil || token != TOKEN {
        panic(err)
    }
    fmt.Println(retry, token, totalSpent/time.Duration(retry))
}

var tokenCache string // this is the file in server to save token got
var lock = sync.Mutex{}
var tokenGetFlag = false
var errRetry = errors.New("retry")
func server() (string, error) {
    // check token cache
    lock.Lock()
    defer lock.Unlock()
    if tokenCache != "" {
        token := tokenCache
        return token, nil
    }

    // no token, we got it
    // check if we are trying get token
    if tokenGetFlag == false {
        tokenGetFlag = true
        go func() {
            token := ApiBy3rd()
            lock.Lock()
            defer lock.Unlock()
            tokenCache = token
            tokenGetFlag = false
        }()
    }
    return "", errRetry
}