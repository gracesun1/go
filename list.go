package main

import (
    "fmt"
)

type LinkNode struct {
    next *LinkNode
    value int
}

func (head *LinkNode) Append(v int) *LinkNode {
    if head == nil {
        return &LinkNode{
            next: nil,
            value: v,
        }
    }
    var tail *LinkNode
    // find the tail
    for tail = head; tail.next != nil; tail = tail.next {    }
    tail.next = &LinkNode{
        next: nil,
        value: v,
    }
    return head
}

func (head *LinkNode) Remove(v int) *LinkNode {
    newHead := head
    // find node match v, remove it
    for pre, cur := (*LinkNode)(nil), head; cur != nil; {
        if cur.value == v {
            if pre == nil {
                // remove the head node
                newHead = head.next
            } else {
                // remove the middle/last node
                pre.next = cur.next
            }
            break
        }
        // pre, cur move forwards
        pre, cur = cur, cur.next
    }
    return newHead
}

func (head *LinkNode) Show(s string) {
    fmt.Println("============", s, "=============")
    i := 0
    for p := head; p != nil; p = p.next {
        fmt.Printf("[%d] value=%d\n", i, p.value)
        i++
    }
}

func main() {
    h := (*LinkNode)(nil).Append(1)
    h = h.Append(100)
    h = h.Append(2)
    h = h.Append(3)
    h.Show("showing 1,100,2,3")
    h = h.Remove(-1)
    h.Show("showing remove -1: 1,100,2,3")
    h = h.Remove(2)
    h.Show("showing remove 2: 1,100,3")
    h = h.Remove(1)
    h.Show("showing remove 1: 100,3")
    h = h.Remove(3)
    h.Show("showing remove 3: 100,")

    h = nil
    h.Remove(-1)
    h.Show("showing remove -1: ")
}