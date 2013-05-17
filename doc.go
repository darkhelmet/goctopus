/*
The goctopus package multiplexes many channels onto one channel.

After building an Octopus with an arbitrary number of channels, you
can Run the Octopus and it will give you back a new receive only channel
of empty interface values. You can receive on this channel until it closes, at
which points all values from all the other channels have been processed.

    import "github.com/darkhelmet/goctopus"

    c1 := make(chan int, 1)
    c1 <- 5
    close(c1)

    c2 := make(chan bool, 1)
    c2 <- true
    close(c2)

    c := goctopus.New(c1, c2).Run()
    for i := range c {
        switch v := i.(type) {
        case int:
            fmt.Printf("got int %d", v)
        case bool:
            fmt.Printf("got bool %t", v)
        }
    }

This will print

    got int 5
    got bool true

*/
package goctopus
