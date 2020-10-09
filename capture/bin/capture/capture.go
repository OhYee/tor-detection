package main

import (
	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/tor-detection/capture/lib/log"
	"github.com/OhYee/tor-detection/capture/lib/socks5"
)

func main() {
	socks := socks5.NewSocks5Server()
	err := socks.Start("0.0.0.0", 6666)
	if err != nil {
		log.Error.Println(errors.ShowStack(err))
	}
	for {
	}
}
