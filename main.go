package main

import (
	"fmt"
	"github.com/JacobNewton007/channel/channel"
)

func main() {

	c := channel.NewChannel(10)
	c.Send(2)
	c.Send(3)
	c.Close()
	ch, ok := c.Receive()
	if ok {
		fmt.Println(ch)
	}

}
