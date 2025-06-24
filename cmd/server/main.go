package main

import "github.com/jonilsonds9/goexpert-modulo-7-apis/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
