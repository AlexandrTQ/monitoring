package main

import "monitoring/service/cmd"

func main() {
	println("Start")
	cmd.Execute()
	println("End")
}
