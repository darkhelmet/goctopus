package goctopus_test

import (
    G "github.com/darkhelmet/goctopus"
    . "launchpad.net/gocheck"
    "testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (td *TestSuite) TestConstructor(c *C) {
    o := G.New(make(chan bool))
    c.Assert(o.Len(), Equals, 1)
}

func (td *TestSuite) TestRequireThingsToBeChannels(c *C) {
    c.Assert(func() { G.New(1) }, PanicMatches, "only channels allowed, 1 at index 0 is not a channel")
}

func (td *TestSuite) TestRequireSendChannels(c *C) {
    c.Assert(func() { G.New(make(chan<- bool)) }, PanicMatches, "only recv channels allowed, channel at index 0 is a send-only channel")
}

func (td *TestSuite) TestRun(c *C) {
    c1, c2 := make(chan int, 1), make(chan bool, 1)
    c1 <- 5
    c2 <- true
    close(c1)
    close(c2)
    o := G.New(c1, c2)
    p := o.Run()

    for i := range p {
        switch v := i.(type) {
        case int:
            c.Assert(v, Equals, 5)
        case bool:
            c.Assert(v, Equals, true)
        default:
            c.Fatalf("unexpected value %#v", v)
        }
    }
}
