package main

import (
	// "./deploysaurus"
	"github.com/ys/deploysaurus/deploysaurus"
)
import _ "github.com/joho/godotenv/autoload"

func main() {

	//Launch Build processor
	// Deploysaurus == old name for Deployer
	events := make(chan deploysaurus.Event, 1024)
	for w := 1; w <= 3; w++ {
		go deploysaurus.Deploysaurus(w, events)
	}
	//Launch Server
	deploysaurus.LaunchServer(events)
}
