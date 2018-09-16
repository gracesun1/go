package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Mail struct {
    sender int
    receiver int
    content string
}
var user1_send_mailbox = make(chan *Mail)
var user1_recv_mailbox = make(chan *Mail)

var user2_send_mailbox = make(chan *Mail)
var user2_recv_mailbox = make(chan *Mail)

type PostMan struct {

}
func (p *PostMan) doWork() {
    for {
        select {
        case mail := <- user1_send_mailbox:
            go p.post(mail)
        case mail := <- user2_send_mailbox:
            go p.post(mail)
        }
    }
}
func (p *PostMan) post(mail *Mail) {
    switch mail.receiver {
    case 1:
        user1_recv_mailbox <- mail
    case 2:
        user2_recv_mailbox <- mail
    default:
        // drop it!
    }
}

func user1() {
    for {
        // check recv mailbox
        select {
        case mail := <- user1_recv_mailbox:
            fmt.Println("       user1", *mail)
        default:
            // no mail comming..., do nothing
        }
        // check and send mail
        r := rand.Intn(100)
        switch {
        case r <= 70:
            // do nothing
        case r <= 80:
            sendMail(user1_send_mailbox, 1, 1, "hehe")
        default:
            sendMail(user1_send_mailbox, 1, 2, "hehe2")
        }
        time.Sleep(time.Second)
    }
}

func user2() {
    for {
        // check recv mailbox
        select {
        case mail := <- user2_recv_mailbox:
            fmt.Println("       user2", *mail)
        default:
            // no mail comming..., do nothing
        }
        // check and send mail
        r := rand.Intn(100)
        switch {
        case r <= 80:
            // do nothing
        case r <= 90:
            sendMail(user2_send_mailbox, 2, 1, "haha")
        default:
            sendMail(user2_send_mailbox, 2, 2, "haha2")
        }
        time.Sleep(time.Second)
    }
}

func sendMail(mb chan *Mail, sender int, receiver int, content string) {
    mail := &Mail {
        sender: sender,
        receiver: receiver,
        content: content,
    }
    fmt.Printf("%#v\n", *mail)
    mb <- mail
}

func main() {
    rand.Seed(time.Now().Unix())
    var p = &PostMan{}
    go p.doWork()
    go user1()
    go user2()

    time.Sleep(time.Hour)
}

// type User struct {
//     name string
//     inbox []string
//     outbox []string
// }
//
// type Worker struct {}
//
// var ch1 = make(chan string, 3)
// var ch2 = make(chan string, 3)
// func main() {
//
//     a := User{ name: "xia"}
//     a.outbox =    append(a.outbox, "beijing")
//
//     b := User{ name: "han"}
//     b.outbox =    append(b.outbox, "hongkong")
//
//     go Collect(a)
//
//
// }
//
// func Collect(u User){
//     for i,v :=range(u.outbox){
//         ch1 <- v
//     }
// }
//
// func Distribute(){
//
// }

