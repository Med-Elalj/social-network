package hub

import "social-network/sn/ws"

var HUB = ws.NewHub()

func init() {
	go HUB.Run()
}
