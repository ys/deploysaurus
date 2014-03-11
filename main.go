package main

import (
	"github.com/ys/deploysaurus/deploysaurus"
)
import _ "github.com/joho/godotenv/autoload"

func main() {

	//Launch Build processor
	// Deploysaurus == old name for Deployer
	// 3 deploysaurus in the kitchen, that's one more than in Jurassic Park
	// We should be able to catch and eat the children
	events := make(chan deploysaurus.Event, 1024)
	for w := 1; w <= 3; w++ {
		go deploysaurus.Deploysaurus(w, events)
	}
	//Launch Server
	deploysaurus.LaunchServer(events)
}
