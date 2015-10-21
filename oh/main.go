package main

import "os"

func main() {
	app := CreateCliApp()
	app.Run(os.Args)
}
