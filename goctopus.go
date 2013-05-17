// +build go1.1
package goctopus

import (
    "fmt"
    "reflect"
)

type Octopus struct {
    cases []reflect.SelectCase
    c     chan interface{}
}

// Get how many arms the octopus has.
func (o *Octopus) Len() int {
    return len(o.cases)
}

// Run the octopus, get a new channel that everything comes out of.
func (o *Octopus) Run() <-chan interface{} {
    o.c = make(chan interface{})
    go o.pump(o.c)
    return o.c
}

func (o *Octopus) pump(c chan interface{}) {
    for {
        if len(o.cases) == 0 {
            break
        }

        i, v, ok := reflect.Select(o.cases)

        if ok {
            o.c <- v.Interface()
        } else {
            // Delete this channel
            o.cases = append(o.cases[:i], o.cases[i+1:]...)
        }
    }
    close(o.c)
}

func verifyChannel(index int, v reflect.Value) {
    if reflect.Chan != v.Kind() {
        panic(fmt.Errorf("only channels allowed, %#v at index %d is not a channel", v.Interface(), index))
    }

    switch v.Type().ChanDir() {
    case reflect.SendDir:
        panic(fmt.Errorf("only recv channels allowed, channel at index %d is a send-only channel", index))
    default:
        // All good
    }
}

// Build a new octopus.
func New(channels ...interface{}) *Octopus {
    cases := make([]reflect.SelectCase, 0, len(channels))
    for index, c := range channels {
        value := reflect.ValueOf(c)
        verifyChannel(index, value)
        scase := reflect.SelectCase{
            Dir:  reflect.SelectRecv,
            Chan: value,
        }
        cases = append(cases, scase)
    }
    return &Octopus{cases: cases}
}
