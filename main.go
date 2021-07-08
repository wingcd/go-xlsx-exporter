package main

import (
	"flag"
	"fmt"
)

var (
	VERSION = "0.0.1"
)

func main() {
	var v = flag.Bool("version", false, "获取版本")
	if *v {
		fmt.Println(VERSION)
	}
}
