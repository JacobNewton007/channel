package main

import (
	"fmt"
	"github.com/JacobNewton007/channel/channel"
)

func main() {

	c := channel.NewChannel(10)
	c.Send(2)
	c.Send(3)
	c.Send(4)
	ch, ok := c.Receive()
	fmt.Println(ch, ok)
	ch, ok = c.Receive()
	fmt.Println(ch, ok)
	ch, ok = c.Receive()
	fmt.Println(ch, ok)
	c.Close()

}
