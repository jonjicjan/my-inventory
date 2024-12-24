package main

import "log"

func main() {
	app := App{}
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}
	app.Run("localhost:10000")
}
